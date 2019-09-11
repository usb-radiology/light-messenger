package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationNotificationCreateShouldReturnJSONForLowPriority(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department             = "abc"
		modality               = "x"
		priority               = "3"
		priorityNumber float64 = 3
		now                    = time.Now()
	)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])
	assertExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-info", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnJSONForMediumPriority(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department             = "abc"
		modality               = "x"
		priority               = "2"
		priorityNumber float64 = 2
		now                    = time.Now()
	)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])
	assertExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-warning", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnJSONForHighPriority(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department             = "abc"
		modality               = "x"
		priority               = "1"
		priorityNumber float64 = 1
		now                    = time.Now()
	)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])
	assertExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-danger", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnJSONForHighPriorityAndArduinoStatus(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department             = "abc"
		modality               = "x"
		priority               = "1"
		priorityNumber float64 = 1
		now                    = time.Now()
	)

	arduinoStatus := lmdatabase.ArduinoStatus{
		DepartmentID: department,
		StatusAt:     now.Unix() - 1,
	}

	errArduinoStatusInsert := lmdatabase.ArduinoStatusInsert(db, arduinoStatus)
	if errArduinoStatusInsert != nil {
		t.Fatalf("%+v", errArduinoStatusInsert)
	}

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.NotNil(t, responseBodyStrings["ArduinoStatus"])
	arduinoStatusStrings := responseBodyStrings["ArduinoStatus"].(map[string]interface{})
	assert.Equal(t, department, arduinoStatusStrings["DepartmentID"])
	assert.Equal(t, float64(arduinoStatus.StatusAt), arduinoStatusStrings["StatusAt"])
	assertExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-danger", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnHTMLForLowPriority(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department = "abc"
		modality   = "x"
		priority   = "3"
		now        = time.Now()
	)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	doc := getResponseHTMLDoc(t, request)

	// fmt.Print(doc.Html())

	titleLinkSelection := doc.Find("header .card-header-title a")
	assert.Equal(t, 1, titleLinkSelection.Length())

	headerTagsSelection := doc.Find("header div.tags").Children()
	assert.Equal(t, 3, headerTagsSelection.Length())

	headerTagsSelection.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			assert.Equal(t, "Offen", s.Text())
			assert.True(t, s.HasClass("is-info"))
		}
		if i == 1 {
			assertExpectedTimeIsLessThanActualTime(t, now, s.Text())
		}
		if i == 2 {
			assert.True(t, s.HasClass("is-delete"))
			assert.Equal(t, "/modality/"+modality+"/department/"+department+"/cancel", s.AttrOr("ic-post-to", ""))
		}
	})

	priorityButtonsSelection := doc.Find("div.card-content a.button")
	assert.Equal(t, 3, priorityButtonsSelection.Length())

	priorityButtonsSelection.Each(func(i int, s *goquery.Selection) {
		_, existsDisabledAttr := s.Attr("disabled")
		_, existsIcPostToAttr := s.Attr("ic-post-to")

		if i == 0 {
			assert.True(t, s.HasClass("is-info"))
			assert.True(t, existsDisabledAttr)
			assert.False(t, existsIcPostToAttr)
		}
		if i == 1 {
			assert.True(t, s.HasClass("is-warning"))
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
		if i == 2 {
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
	})

	cardFooterItemsSelection := doc.Find("div.card-footer-item div.column")
	assert.Equal(t, 2, cardFooterItemsSelection.Length())

	cardFooterItemsSelection.Each(func(i int, s *goquery.Selection) {
		if i == 0 {

			assert.Equal(t, "Arduino Status", s.Children().First().Text())
		}
		if i == 1 {
			assert.True(t, s.HasClass("has-text-danger"))
			arduinoStatusIconSelection := s.Children().First()
			assert.True(t, arduinoStatusIconSelection.HasClass("fa-ban"))
			assert.Equal(t, "Kein Signal vom Arduino", arduinoStatusIconSelection.AttrOr("title", ""))
		}
	})

	tearDownTest(t, server, db)
}

func assertExpectedTimeIsLessThanActualTime(t *testing.T, expectedTime time.Time, actualTimeStr string) {
	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", expectedTime.Format("2006-01-02 ")+actualTimeStr)
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, expectedTime.Unix(), actualTime.Unix())
}
