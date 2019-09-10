# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME_PREFIX=light-messenger
BINARY_NAME=$(BINARY_NAME_PREFIX).exec
BINARY_UNIX=$(BINARY_NAME_PREFIX)-unix.exec

VERSION=$(shell git describe --tags --dirty --always)
BUILD_TIME=$(shell date +%FT%T%z)

LDFLAGS=-ldflags "-X github.com/usb-radiology/light-messenger/src/version.Version=$(VERSION) -X github.com/usb-radiology/light-messenger/src/version.BuildTime=$(BUILD_TIME)"

all: test build
embed: 
	rm src/server/rice-box.go
	rice embed-go -v -i github.com/usb-radiology/light-messenger/src/server
build: embed
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
test: embed
	$(GOTEST) -v ./...
test-unit: embed
	$(GOTEST) -v -run Unit ./...
test-integration:  embed
	$(GOTEST) -v -run Integration ./...
clean: 
	rm src/server/rice-box.go
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run: build
	./$(BINARY_NAME)
deps:
	$(GOMOD) tidy
	
# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
