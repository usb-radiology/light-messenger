package lmdatabase

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

var statements []string

func GetTestDB(t *testing.T, initConfig *configuration.Configuration) *sql.DB {
	db, errDb := GetDB(initConfig)
	if errDb != nil {
		t.Fatal(errDb)
	}

	// only initialize statements once since reading from disk is a relatively slow operation
	if len(statements) == 0 {
		strStatements, errReadStatements := ReadStatementsFromSQL(filepath.Join("..", "..", "res", "integration_test_setup.sql"))
		if errReadStatements != nil {
			t.Fatalf("%+v", errors.WithStack(errReadStatements))
		}

		statements = *strStatements
		// log.Printf("%+v", statements)
	}

	_, errExecStatements := ExecStatements(db, statements)

	if errExecStatements != nil {
		t.Fatalf("%+v", errors.WithStack(errExecStatements))
	}

	/*
		for _, execResult := range *execResults {
			log.Printf("%+v", execResult)
		}
	*/

	return db
}

func setupTest(t *testing.T) *sql.DB {
	initConfig, err := configuration.LoadAndSetConfiguration(filepath.Join("..", "..", "config-sample.json"))
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	return GetTestDB(t, initConfig)
}

func tearDownTest(t *testing.T, db *sql.DB) {
	db.Close()
}
