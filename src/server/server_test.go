package server

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

func setupTestServer(t *testing.T) *httptest.Server {

	initConfig, err := configuration.LoadAndSetConfiguration(filepath.Join("..", "..", "config-sample.json"))
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	router := getRouter(initConfig)
	ts := httptest.NewServer(router)

	return ts
}

func tearDownTestServer(t *testing.T, server *httptest.Server) {
	server.Close()
}
