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
	rice embed-go -v -i github.com/usb-radiology/light-messenger/src/server
build: embed
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
test-unit: 
	$(GOTEST) -v -run Unit ./...
test-integration: 
	$(GOTEST) -v -run Integration ./...
clean: 
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
