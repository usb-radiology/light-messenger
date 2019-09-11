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

func testNotificationInsert(t *testing.T, db *sql.DB, department string, priority int, modality string, when int64) {
	err := lmdatabase.NotificationInsert(db, department, priority, modality, when)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
}
