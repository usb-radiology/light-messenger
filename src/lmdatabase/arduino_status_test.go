package lmdatabase

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationShouldInsertArduinoStatusWhenNoneExisting(t *testing.T) {

	// given
	db := setupTest(t)

	departmentID := "abc"
	var now int64
	now = 1000

	arduinoStatus := ArduinoStatus{
		DepartmentID: departmentID,
		StatusAt:     now - 299,
	}

	// when
	{
		errInsert := ArduinoStatusInsert(db, arduinoStatus)
		if errInsert != nil {
			t.Fatalf("%+v", errors.WithStack(errInsert))
		}
	}

	// then
	result, errQuery := ArduinoStatusQueryWithin5MinutesFromNow(db, departmentID, now)

	if errQuery != nil {
		t.Fatalf("%+v", errors.WithStack(errQuery))
	}

	if result == nil {
		t.Errorf("Did not retrieve any rows")
	}

	// fmt.Printf("%+v\n", result)
	assert.Equal(t, result.DepartmentID, arduinoStatus.DepartmentID)
	assert.Equal(t, result.StatusAt, arduinoStatus.StatusAt)

	tearDownTest(t, db)
}

func TestIntegrationShouldInsertArduinoStatusWhenExists(t *testing.T) {

	// given
	db := setupTest(t)

	departmentID := "abc"
	var now, update int64
	now = 1000
	update = 2000

	arduinoStatus := ArduinoStatus{
		DepartmentID: departmentID,
		StatusAt:     now - 299,
	}

	arduinoStatusUpdate := ArduinoStatus{
		DepartmentID: departmentID,
		StatusAt:     update,
	}

	{
		errInsert := ArduinoStatusInsert(db, arduinoStatus)
		if errInsert != nil {
			t.Fatalf("%+v", errors.WithStack(errInsert))
		}
	}

	// when
	{
		errUpdate := ArduinoStatusInsert(db, arduinoStatusUpdate)
		if errUpdate != nil {
			t.Fatalf("%+v", errors.WithStack(errUpdate))
		}
	}

	// then
	result, errQuery := ArduinoStatusQueryWithin5MinutesFromNow(db, departmentID, update)

	if errQuery != nil {
		t.Fatal(errQuery)
	}

	if result == nil {
		t.Errorf("Did not retrieve any rows")
	}

	// fmt.Printf("%+v\n", result)
	assert.Equal(t, result.DepartmentID, arduinoStatus.DepartmentID)
	assert.Equal(t, result.StatusAt, arduinoStatusUpdate.StatusAt)

	tearDownTest(t, db)
}

func TestIntegrationShouldNotRetrieveArduinoStatusWhenOlderThan5Minutes(t *testing.T) {

	// given
	db := setupTest(t)

	departmentID := "abc"
	var now int64
	now = 1000

	arduinoStatus := ArduinoStatus{
		DepartmentID: departmentID,
		StatusAt:     now - 300,
	}

	// when
	{
		errInsert := ArduinoStatusInsert(db, arduinoStatus)
		if errInsert != nil {
			t.Fatalf("%+v", errors.WithStack(errInsert))
		}
	}

	// then
	result, errQuery := ArduinoStatusQueryWithin5MinutesFromNow(db, departmentID, now)

	if errQuery != nil {
		t.Fatalf("%+v", errors.WithStack(errQuery))
	}

	if result != nil {
		t.Fail()
	}

	tearDownTest(t, db)
}
