package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery" // https://godoc.org/github.com/PuerkitoBio/goquery
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationIndexShouldReturnLinksForMTRAsAndRadiologists(t *testing.T) {

	// given
	server, db := setupTest(t)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/", nil)

	// then
	doc := getResponseHTMLDoc(t, request)

	links := doc.Find("a.is-large").Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("href", "")
	})

	expectedLinks := []string{"/mtra/ct", "/mtra/mr", "/mtra/nuk", "/radiologie/aod", "/radiologie/ctd", "/radiologie/msk", "/radiologie/nr", "/radiologie/nuk"}
	assert.EqualValues(t, expectedLinks, links)

	tearDownTest(t, server, db)
}

func testNotificationInsert(t *testing.T, db *sql.DB, department string, priority int, modality string, when int64) {
	err := lmdatabase.NotificationInsert(db, department, priority, modality, when)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
}

func testArduinoStatusInsert(t *testing.T, db *sql.DB, department string, when int64) lmdatabase.ArduinoStatus {
	arduinoStatus := lmdatabase.ArduinoStatus{
		DepartmentID: department,
		StatusAt:     when,
	}

	errArduinoStatusInsert := lmdatabase.ArduinoStatusInsert(db, arduinoStatus)
	if errArduinoStatusInsert != nil {
		t.Fatalf("%+v", errArduinoStatusInsert)
	}

	return arduinoStatus
}

func getResponseHTMLDoc(t *testing.T, request *http.Request) *goquery.Document {
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueHTML)

	response := getResponse(t, request)
	doc, errHTMLDoc := goquery.NewDocumentFromResponse(response)
	if errHTMLDoc != nil {
		t.Fatalf("%+v", errors.WithStack(errHTMLDoc))
	}

	return doc
}

func getDocument(t *testing.T, html string) *goquery.Document {
	doc, errHTMLDoc := goquery.NewDocumentFromReader(strings.NewReader(html))
	if errHTMLDoc != nil {
		t.Fatalf("%+v", errors.WithStack(errHTMLDoc))
	}
	return doc
}

func getResponseBodyStrings(t *testing.T, request *http.Request) map[string]interface{} {
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueJSON)

	response := getResponse(t, request)
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody interface{}
	errJSONDecode := json.NewDecoder(response.Body).Decode(&responseBody)
	if errJSONDecode != nil {
		t.Fatalf("%+v", errors.WithStack(errJSONDecode))
	}

	// printJSON(t, responseBody)

	responseBodyStrings := responseBody.(map[string]interface{})
	return responseBodyStrings
}

func getResponse(t *testing.T, request *http.Request) *http.Response {
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}

	return response
}

func printJSON(t *testing.T, data interface{}) {
	s, errJSONPrint := json.MarshalIndent(data, "", "\t")
	if errJSONPrint != nil {
		t.Fatalf("%+v", errJSONPrint)
	}
	fmt.Println(string(s))
}

func assertNotificationHTMLMediumPriority(t *testing.T, doc *goquery.Document, modality string, department string, now time.Time) {

	titleLinkSelection := doc.Find("header .card-header-title a")
	assert.Equal(t, 1, titleLinkSelection.Length())

	headerTagsSelection := doc.Find("header div.tags").Children()
	assert.Equal(t, 3, headerTagsSelection.Length())

	headerTagsSelection.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			assert.Equal(t, "Offen", s.Text())
			assert.True(t, s.HasClass("is-warning"))
		}
		if i == 1 {
			assertNotificationCardExpectedTimeIsLessThanActualTime(t, now, s.Text())
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
			assert.True(t, existsDisabledAttr)
			assert.False(t, existsIcPostToAttr)
		}
		if i == 2 {
			assert.True(t, s.HasClass("is-danger"))
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
	})

	assertNotificationHTMLArdunioStatusNoSignal(t, doc)
}

func assertNotificationHTMLNoPriority(t *testing.T, doc *goquery.Document, modality string, department string) {

	titleLinkSelection := doc.Find("header .card-header-title a")
	assert.Equal(t, 1, titleLinkSelection.Length())

	headerTagsSelection := doc.Find("header div.tags").Children()
	assert.Equal(t, 0, headerTagsSelection.Length())

	priorityButtonsSelection := doc.Find("div.card-content a.button")
	assert.Equal(t, 3, priorityButtonsSelection.Length())

	priorityButtonsSelection.Each(func(i int, s *goquery.Selection) {
		_, existsDisabledAttr := s.Attr("disabled")
		_, existsIcPostToAttr := s.Attr("ic-post-to")

		if i == 0 {
			assert.True(t, s.HasClass("is-info"))
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
		if i == 1 {
			assert.True(t, s.HasClass("is-warning"))
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
		if i == 2 {
			assert.True(t, s.HasClass("is-danger"))
			assert.False(t, existsDisabledAttr)
			assert.True(t, existsIcPostToAttr)
		}
	})

	assertNotificationHTMLArdunioStatusNoSignal(t, doc)
}

func assertNotificationHTMLArdunioStatusNoSignal(t *testing.T, doc *goquery.Document) {
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
}

func assertNotificationCardExpectedTimeIsLessThanActualTime(t *testing.T, expectedTime time.Time, actualTimeStr string) {
	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", expectedTime.Format("2006-01-02 ")+actualTimeStr)
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, expectedTime.Unix(), actualTime.Unix())
}
