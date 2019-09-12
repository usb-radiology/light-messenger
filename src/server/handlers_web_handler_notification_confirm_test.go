package server

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationNotificationConfirmShouldReturnHTTP200WhenNotificationExists(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department  = "abc"
		modality    = "x"
		priorityInt = 1
		now         = time.Now()
		// priority               = "3"
		// priorityNumber float64 = 3
	)

	testNotificationInsert(t, db, department, priorityInt, modality, now.Unix())
	insertedNotification := getNotification(t, db, department, modality)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/notification/"+department+"/"+insertedNotification.NotificationID, nil)

	// then
	response := getResponse(t, request)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	notification := getNotificationByID(t, db, insertedNotification.NotificationID)
	assert.LessOrEqual(t, now.Unix(), notification.ConfirmedAt)

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationConfirmShouldReturnHTTP200WhenNoNotificationExists(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department = "abc"
		// modality    = "x"
		// priorityInt = 1
		// now = time.Now()
		// priority               = "3"
		// priorityNumber float64 = 3
	)

	// testNotificationInsert(t, db, department, priorityInt, modality, now.Unix())
	// insertedNotification := getNotification(t, db, department, modality)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/notification/"+department+"/xxx", nil)

	// then
	response := getResponse(t, request)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	tearDownTest(t, server, db)
}

func getNotification(t *testing.T, db *sql.DB, department string, modality string) *lmdatabase.Notification {
	notification, err := lmdatabase.NotificationGetOpenNotificationByDepartmentAndModality(db, department, modality)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	return notification
}

func getNotificationByID(t *testing.T, db *sql.DB, notificationID string) *lmdatabase.Notification {
	notification, err := lmdatabase.NotificationGetByID(db, notificationID)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	return notification
}
