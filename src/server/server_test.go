package server

import (
	"database/sql"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

var statements []string

func setupTest(t *testing.T) (*httptest.Server, *sql.DB) {

	initConfig, err := configuration.LoadAndSetConfiguration(filepath.Join("..", "..", "config-sample.json"))
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	db, errDb := lmdatabase.GetDB(initConfig)
	if errDb != nil {
		t.Fatal(errDb)
	}

	// only initialize statements once since reading from disk is a relatively slow operation
	if len(statements) == 0 {
		strStatements, errReadStatements := lmdatabase.ReadStatementsFromSQL(filepath.Join("..", "..", "res", "integration_test_setup.sql"))
		if errReadStatements != nil {
			t.Fatalf("%+v", errors.WithStack(errReadStatements))
		}

		statements = *strStatements
		// log.Printf("%+v", statements)
	}

	_, errExecStatements := lmdatabase.ExecStatements(db, statements)

	if errExecStatements != nil {
		t.Fatalf("%+v", errors.WithStack(errExecStatements))
	}

	/*
		for _, execResult := range *execResults {
			log.Printf("%+v", execResult)
		}
	*/

	router := getRouter(initConfig, db)
	ts := httptest.NewServer(router)

	return ts, db
}

func tearDownTest(t *testing.T, server *httptest.Server, db *sql.DB) {
	server.Close()
	db.Close()
}
