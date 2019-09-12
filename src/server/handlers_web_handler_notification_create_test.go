package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assertNotificationCardExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
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
	assertNotificationCardExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
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
	assertNotificationCardExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
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

	arduinoStatus := testArduinoStatusInsert(t, db, department, now.Unix()-1)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.NotNil(t, responseBodyStrings["ArduinoStatus"])
	arduinoStatusStrings := responseBodyStrings["ArduinoStatus"].(map[string]interface{})
	assert.Equal(t, department, arduinoStatusStrings["DepartmentID"])
	assert.Equal(t, float64(arduinoStatus.StatusAt), arduinoStatusStrings["StatusAt"])
	assertNotificationCardExpectedTimeIsLessThanActualTime(t, now, responseBodyStrings["CreatedAt"].(string))
	assert.Equal(t, department, responseBodyStrings["Department"])
	assert.Equal(t, modality, responseBodyStrings["Modality"])
	assert.Equal(t, priority, responseBodyStrings["Priority"])
	assert.Equal(t, "is-danger", responseBodyStrings["PriorityName"])
	assert.Equal(t, priorityNumber, responseBodyStrings["PriorityNumber"])

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCreateShouldReturnHTMLForMediumPriority(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department = "abc"
		modality   = "x"
		priority   = "2"
		now        = time.Now()
	)

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/prio/"+priority, nil)

	// then
	doc := getResponseHTMLDoc(t, request)

	// fmt.Print(doc.Html())

	assertNotificationHTMLMediumPriority(t, doc, modality, department, now)

	tearDownTest(t, server, db)
}
