package lmdatabase

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/configuration"
)

func setUp(t *testing.T) *sql.DB {
	initConfig, err := configuration.LoadAndSetConfiguration("../../config-sample.json")
	if err != nil {
		t.Fatal(err)
	}

	db, errDb := GetDB(initConfig)
	if errDb != nil {
		t.Fatal(errDb)
	}

	sqlFilePath := filepath.Join("..", "..", "res", "create_tables.sql")

	fileContents, errFileRead := ioutil.ReadFile(sqlFilePath)
	if err != nil {
		t.Fatal(errFileRead)
	}

	statements := strings.Split(string(fileContents), ";")

	execStatements(db, statements)

	return db
}

func TestIntegrationShouldUpdateTheArduinoStatus(t *testing.T) {

	// given
	db := setUp(t)

	departmentID := "abc"
	var now int64
	now = 1000

	arduinoStatus := ArduinoStatus{
		DepartmentID: departmentID,
		StatusAt:     now - 10,
	}

	// when
	errInsert := InsertStatus(db, arduinoStatus)
	if errInsert != nil {
		t.Fatal(errInsert)
	}

	// then
	result, errQuery := IsAlive(db, departmentID, now)

	if errQuery != nil {
		t.Fatal(errQuery)
	}

	if result == nil {
		t.Errorf("Did not retrieve any rows")
	}

	fmt.Printf("%+v\n", result)

	assert.Equal(t, result.DepartmentID, arduinoStatus.DepartmentID)
	assert.Equal(t, result.StatusAt, arduinoStatus.StatusAt)

}
