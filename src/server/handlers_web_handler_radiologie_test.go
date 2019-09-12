package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationRadiologieShouldReturnJSONWithNotificationsHTML(t *testing.T) {

	// given
	server, db := setupTest(t)

	duration, _ := time.ParseDuration("-1h")

	var (
		department = "aod"
		oneHourAgo = time.Now().Add(duration)
	)

	testNotificationInsert(t, db, department, 1, "x", oneHourAgo.Unix())      // oldest, but highest prio
	testNotificationInsert(t, db, department, 2, "y", oneHourAgo.Unix()+1000) // medium prio
	testNotificationInsert(t, db, department, 3, "z", oneHourAgo.Unix()+2000) // most recent lowest prio

	// when
	request, _ := http.NewRequest("GET", server.URL+"/radiologie/"+department, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	// fmt.Printf("%+v", responseBodyStrings)
	// fmt.Printf("%+v", responseBodyStrings["ProcessedNotifications"])

	assert.Equal(t, department, responseBodyStrings["Department"].(string))
	assert.Nil(t, responseBodyStrings["ArduinoStatus"])

	doc := getDocument(t, responseBodyStrings["Notifications"].(string))

	notificationsSelection := doc.Find("div.content")
	assert.Equal(t, 3, notificationsSelection.Length())

	notificationsSelection.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			assertNotificationDisplayHTML(t, s, "is-danger", "x", oneHourAgo)
		}
		if i == 1 {
			assertNotificationDisplayHTML(t, s, "is-warning", "y", oneHourAgo)
		}
		if i == 2 {
			assertNotificationDisplayHTML(t, s, "is-info", "z", oneHourAgo)
		}
	})

	tearDownTest(t, server, db)
}

func TestIntegrationRadiologieShouldReturnJSONWithArduinoStatus(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department = "aod"
		now        = time.Now()
	)

	arduinoStatus := testArduinoStatusInsert(t, db, department, now.Unix())

	// when
	request, _ := http.NewRequest("GET", server.URL+"/radiologie/"+department, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	// fmt.Printf("%+v", responseBodyStrings)
	// fmt.Printf("%+v", responseBodyStrings["ProcessedNotifications"])

	assert.Equal(t, department, responseBodyStrings["Department"].(string))
	assert.NotNil(t, responseBodyStrings["ArduinoStatus"])
	arduinoStatusStrings := responseBodyStrings["ArduinoStatus"].(map[string]interface{})
	assert.Equal(t, department, arduinoStatusStrings["DepartmentID"])
	assert.Equal(t, float64(arduinoStatus.StatusAt), arduinoStatusStrings["StatusAt"])

	tearDownTest(t, server, db)
}

func assertNotificationDisplayHTML(t *testing.T, notificationSelection *goquery.Selection, priorityClass string, modality string, expectedTime time.Time) {

	partsSelection := notificationSelection.Find("div.control")
	assert.Equal(t, 2, partsSelection.Length())

	partsSelection.Each(func(partsSelectionIndex int, partSelection *goquery.Selection) {
		if partsSelectionIndex == 0 {
			tagsSelection := partSelection.Find("span.tag")
			assert.Equal(t, 3, tagsSelection.Length())

			tagsSelection.Each(func(tagsSelectionIndex int, tagSelection *goquery.Selection) {
				if tagsSelectionIndex == 0 {
					assert.True(t, tagSelection.HasClass(priorityClass))
					assert.Equal(t, "Offen", tagSelection.Text())
				}
				if tagsSelectionIndex == 1 {
					assert.True(t, tagSelection.HasClass("is-dark"))
					assertNotificationCardExpectedTimeIsLessThanActualTime(t, expectedTime, tagSelection.Text())
				}
				if tagsSelectionIndex == 2 {
					assert.True(t, tagSelection.HasClass("is-info"))
					assert.Equal(t, modality, tagSelection.Text())
				}
			})
		}
		if partsSelectionIndex == 1 {
			tagsSelection := partSelection.Find("a.tag")
			assert.Equal(t, 1, tagsSelection.Length())
		}
	})
}
