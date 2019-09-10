package server

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationArduinoStatusShouldLogAtGivenTime(t *testing.T) {

	// given
	server, db := setupTest(t)

	departmentID := "abc"
	now := time.Now().Unix()

	// when

	response, errHTTPGet := http.Get(server.URL + "/nce-rest/arduino-status/" + departmentID + "-status")
	if errHTTPGet != nil {
		t.Fatalf("%+v", errors.WithStack(errHTTPGet))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	result, errQuery := lmdatabase.ArduinoStatusQueryWithin5MinutesFromNow(db, departmentID, now)
	if errQuery != nil {
		t.Fatalf("%+v", errors.WithStack(errQuery))
	}

	if result == nil {
		t.Errorf("Did not retrieve any rows")
	}

	assert.Equal(t, departmentID, result.DepartmentID)
	assert.LessOrEqual(t, now, result.StatusAt)

	tearDownTest(t, server, db)
}

func TestIntegrationArduinoGetOpenNotificationsShouldGetHighPriorityNotificationWhenMultipleExist(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		departmentID = "abc"
		modality1    = "x"
		modality2    = "y"
		modality3    = "z"
		now          = time.Now().Unix()
	)

	notificationInsert(t, db, departmentID, 1, modality1, now-10)
	notificationInsert(t, db, departmentID, 2, modality2, now-5)
	notificationInsert(t, db, departmentID, 3, modality3, now) // lowest priority came in last

	// when

	response, errHTTPGet := http.Get(server.URL + "/nce-rest/arduino-status/" + departmentID + "-open-notifications")
	if errHTTPGet != nil {
		t.Fatalf("%+v", errors.WithStack(errHTTPGet))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadResponse := ioutil.ReadAll(response.Body)
	if errReadResponse != nil {
		t.Fatalf("%+v", errors.WithStack(errReadResponse))
	}
	bodyString := string(body)

	assert.Equal(t, ";1;HIGH;", bodyString)

	tearDownTest(t, server, db)
}

func TestIntegrationArduinoGetOpenNotificationsShouldGetMediumPriorityNotificationWhenMultipleExist(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		departmentID = "abc"
		modality2    = "y"
		modality3    = "z"
		now          = time.Now().Unix()
	)

	notificationInsert(t, db, departmentID, 2, modality2, now-5)
	notificationInsert(t, db, departmentID, 3, modality3, now) // lowest priority came in last

	// when

	response, errHTTPGet := http.Get(server.URL + "/nce-rest/arduino-status/" + departmentID + "-open-notifications")
	if errHTTPGet != nil {
		t.Fatalf("%+v", errors.WithStack(errHTTPGet))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadResponse := ioutil.ReadAll(response.Body)
	if errReadResponse != nil {
		t.Fatalf("%+v", errors.WithStack(errReadResponse))
	}
	bodyString := string(body)

	assert.Equal(t, ";1;MEDIUM;", bodyString)

	tearDownTest(t, server, db)
}

func TestIntegrationArduinoGetOpenNotificationsShouldGetLowPriorityNotificationWhenMultipleExist(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		departmentID = "abc"
		modality2    = "y"
		modality3    = "z"
		now          = time.Now().Unix()
	)

	notificationInsert(t, db, departmentID, 3, modality2, now)
	notificationInsert(t, db, departmentID, 3, modality3, now)

	// when

	response, errHTTPGet := http.Get(server.URL + "/nce-rest/arduino-status/" + departmentID + "-open-notifications")
	if errHTTPGet != nil {
		t.Fatalf("%+v", errors.WithStack(errHTTPGet))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadResponse := ioutil.ReadAll(response.Body)
	if errReadResponse != nil {
		t.Fatalf("%+v", errors.WithStack(errReadResponse))
	}
	bodyString := string(body)

	assert.Equal(t, ";1;LOW;", bodyString)

	tearDownTest(t, server, db)
}

func TestIntegrationArduinoGetOpenNotificationsShouldGet0WhenNoNotificationsExist(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		departmentID = "abc"
	)

	// when

	response, errHTTPGet := http.Get(server.URL + "/nce-rest/arduino-status/" + departmentID + "-open-notifications")
	if errHTTPGet != nil {
		t.Fatalf("%+v", errors.WithStack(errHTTPGet))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, errReadResponse := ioutil.ReadAll(response.Body)
	if errReadResponse != nil {
		t.Fatalf("%+v", errors.WithStack(errReadResponse))
	}
	bodyString := string(body)

	assert.Equal(t, ";0;", bodyString)

	tearDownTest(t, server, db)
}

func notificationInsert(t *testing.T, db *sql.DB, departmentID string, priority int, modality string, createdAt int64) {
	errNotificationInsert := lmdatabase.NotificationInsert(db, departmentID, priority, modality, createdAt)
	if errNotificationInsert != nil {
		t.Fatalf("%+v", errNotificationInsert)
	}
}
