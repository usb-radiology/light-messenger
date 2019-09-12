package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationVisierungShouldReturnJSONForAllCards(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		modality    = "x"
		priorityInt = 2
		now         = time.Now()
		// priority               = "3"
		// priorityNumber float64 = 3
	)

	testNotificationInsert(t, db, "aod", priorityInt, modality, now.Unix())
	testNotificationInsert(t, db, "ctd", priorityInt, modality, now.Unix())
	testNotificationInsert(t, db, "msk", priorityInt, modality, now.Unix())
	testNotificationInsert(t, db, "nr", priorityInt, modality, now.Unix())
	testNotificationInsert(t, db, "nuk", priorityInt, modality, now.Unix())

	// when
	request, _ := http.NewRequest("GET", server.URL+"/mtra/"+modality, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	// fmt.Printf("%+v", responseBodyStrings)
	// fmt.Printf("%+v", responseBodyStrings["ProcessedNotifications"])

	assertNotificationHTMLMediumPriority(t, getDocument(t, responseBodyStrings["AOD"].(string)), modality, "aod", now)
	assertNotificationHTMLMediumPriority(t, getDocument(t, responseBodyStrings["CTD"].(string)), modality, "ctd", now)
	assertNotificationHTMLMediumPriority(t, getDocument(t, responseBodyStrings["MSK"].(string)), modality, "msk", now)
	assertNotificationHTMLMediumPriority(t, getDocument(t, responseBodyStrings["NR"].(string)), modality, "nr", now)
	assertNotificationHTMLMediumPriority(t, getDocument(t, responseBodyStrings["NUK_NUK"].(string)), modality, "nuk", now)
	assert.Empty(t, responseBodyStrings["ProcessedNotifications"])

	tearDownTest(t, server, db)
}

func TestIntegrationVisierungShouldReturnJSONWithProcessedNotifications(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		modality    = "x"
		priorityInt = 2
		now         = time.Now()
		cancelledAt = now.Unix() - 2
		// priority               = "3"
		// priorityNumber float64 = 3
	)

	testNotificationInsert(t, db, "aod", priorityInt, modality, now.Unix()-1000) // cancelled
	testNotificationInsert(t, db, "ctd", priorityInt, modality, now.Unix()-500)  // confirmed

	errNotificationCancel := lmdatabase.NotificationCancel(db, modality, "aod", cancelledAt)
	if errNotificationCancel != nil {
		t.Fatalf("%+v", errors.WithStack(errNotificationCancel))
	}

	notification, errNotificationGetByDepartmentAndModality := lmdatabase.NotificationGetOpenNotificationByDepartmentAndModality(db, "ctd", modality)
	if errNotificationGetByDepartmentAndModality != nil {
		t.Fatalf("%+v", errors.WithStack(errNotificationGetByDepartmentAndModality))
	}

	_, errNotificationConfirm := lmdatabase.NotificationConfirm(db, notification.NotificationID, now.Unix())
	if errNotificationConfirm != nil {
		t.Fatalf("%+v", errors.WithStack(errNotificationConfirm))
	}

	// when
	request, _ := http.NewRequest("GET", server.URL+"/mtra/"+modality, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)
	processedNotifications := responseBodyStrings["ProcessedNotifications"].([]interface{})

	// fmt.Printf("%+v", processedNotifications)
	assert.NotEmpty(t, processedNotifications)

	{
		var ctdNotification map[string]interface{}
		ctdNotification = processedNotifications[0].(map[string]interface{})
		assert.Equal(t, "ctd", ctdNotification["DepartmentID"].(string))
		assert.Equal(t, modality, ctdNotification["Modality"].(string))
		assert.Equal(t, float64(priorityInt), ctdNotification["Priority"].(float64))
		assert.Equal(t, float64(-1), ctdNotification["CancelledAt"].(float64))
		assert.Equal(t, float64(now.Unix()), ctdNotification["ConfirmedAt"].(float64))
		assert.Equal(t, float64(now.Unix()-500), ctdNotification["CreatedAt"].(float64))
	}

	{
		var aodNotification map[string]interface{}
		aodNotification = processedNotifications[1].(map[string]interface{})
		assert.Equal(t, "aod", aodNotification["DepartmentID"].(string))
		assert.Equal(t, modality, aodNotification["Modality"].(string))
		assert.Equal(t, float64(priorityInt), aodNotification["Priority"].(float64))
		assert.Equal(t, float64(cancelledAt), aodNotification["CancelledAt"].(float64))
		assert.Equal(t, float64(-1), aodNotification["ConfirmedAt"].(float64))
		assert.Equal(t, float64(now.Unix()-1000), aodNotification["CreatedAt"].(float64))
	}

	tearDownTest(t, server, db)
}
