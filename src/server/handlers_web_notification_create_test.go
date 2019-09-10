package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
)

func TestIntegrationNotificationCreateShouldReturnHTMLForLowPriority(t *testing.T) {

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
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueJSON)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody interface{}
	errJSONDecode := json.NewDecoder(response.Body).Decode(&responseBody)
	if errJSONDecode != nil {
		t.Fatalf("%+v", errJSONDecode)
	}

	// printJSON(t, responseBody)

	responseBodyStrings := responseBody.(map[string]interface{})

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])

	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+(responseBodyStrings["CreatedAt"].(string)))
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, now.Unix(), actualTime.Unix())

	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-info", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	/*
		doc, errHTMLDoc := goquery.NewDocumentFromReader(response.Body)
		if errHTMLDoc != nil {
			t.Fatalf("%+v", errors.WithStack(errHTMLDoc))
		}

		fmt.Print(doc.Html())

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
				actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+(responseBodyStrings["CreatedAt"].(string)))
				if errTimeParse != nil {
					t.Fatalf("%+v", errTimeParse)
				}
			}
			if i == 2 {
				assert.True(t, s.HasClass("is-delete"))
				assert.Equal(t, "/modality/"+modality+"/department/"+department+"/cancel", s.AttrOr("ic-post-to", ""))
			}
		})
	*/

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnHTMLForMediumPriority(t *testing.T) {

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
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueJSON)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody interface{}
	errJSONDecode := json.NewDecoder(response.Body).Decode(&responseBody)
	if errJSONDecode != nil {
		t.Fatalf("%+v", errJSONDecode)
	}

	// printJSON(t, responseBody)

	responseBodyStrings := responseBody.(map[string]interface{})

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])

	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+(responseBodyStrings["CreatedAt"].(string)))
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, now.Unix(), actualTime.Unix())

	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-warning", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnHTMLForHighPriority(t *testing.T) {

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
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueJSON)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody interface{}
	errJSONDecode := json.NewDecoder(response.Body).Decode(&responseBody)
	if errJSONDecode != nil {
		t.Fatalf("%+v", errJSONDecode)
	}

	// printJSON(t, responseBody)

	responseBodyStrings := responseBody.(map[string]interface{})

	assert.Nil(t, responseBodyStrings["ArduinoStatus"])

	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+(responseBodyStrings["CreatedAt"].(string)))
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, now.Unix(), actualTime.Unix())

	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-danger", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnHTMLForHighPriorityAndArduinoStatus(t *testing.T) {

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
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueJSON)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var responseBody interface{}
	errJSONDecode := json.NewDecoder(response.Body).Decode(&responseBody)
	if errJSONDecode != nil {
		t.Fatalf("%+v", errJSONDecode)
	}

	// printJSON(t, responseBody)

	responseBodyStrings := responseBody.(map[string]interface{})

	assert.NotNil(t, responseBodyStrings["ArduinoStatus"])
	arduinoStatusStrings := responseBodyStrings["ArduinoStatus"].(map[string]interface{})
	assert.Equal(t, department, arduinoStatusStrings["DepartmentID"])
	assert.Equal(t, float64(arduinoStatus.StatusAt), arduinoStatusStrings["StatusAt"])

	actualTime, errTimeParse := time.Parse("2006-01-02 15:04:05", now.Format("2006-01-02 ")+(responseBodyStrings["CreatedAt"].(string)))
	if errTimeParse != nil {
		t.Fatalf("%+v", errTimeParse)
	}
	// fmt.Printf("%d %d %s", now.Unix(), actualTime.Unix(), s.Text())
	assert.LessOrEqual(t, now.Unix(), actualTime.Unix())

	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-danger", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func printJSON(t *testing.T, data interface{}) {
	s, errJSONPrint := json.MarshalIndent(data, "", "\t")
	if errJSONPrint != nil {
		t.Fatalf("%+v", errJSONPrint)
	}
	fmt.Println(string(s))
}
