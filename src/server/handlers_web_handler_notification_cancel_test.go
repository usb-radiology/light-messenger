package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationNotificationCancelShouldReturnJSON(t *testing.T) {

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

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/cancel", nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)

	assert.Equal(t, department, responseBodyStrings["Department"].(string))
	assert.Equal(t, modality, responseBodyStrings["Modality"].(string))

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCancelShouldReturnHTML(t *testing.T) {

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

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/cancel", nil)

	// then
	doc := getResponseHTMLDoc(t, request)

	assertNotificationHTMLNoPriority(t, doc, modality, department)

	tearDownTest(t, server, db)
}

func TestIntegrationNotificationCancelShouldReturnJSONWhenNoNotificationToCancel(t *testing.T) {

	// given
	server, db := setupTest(t)

	var (
		department = "abc"
		modality   = "x"
		// priorityInt = 1
		// now         = time.Now()
		// priority               = "3"
		// priorityNumber float64 = 3
	)

	// testNotificationInsert(t, db, department, priorityInt, modality, now.Unix())

	// when
	request, _ := http.NewRequest("GET", server.URL+"/modality/"+modality+"/department/"+department+"/cancel", nil)

	// then
	responseBodyStrings := getResponseBodyStrings(t, request)
	// fmt.Printf("%+v", responseBodyStrings)

	assert.Equal(t, department, responseBodyStrings["Department"].(string))
	assert.Equal(t, modality, responseBodyStrings["Modality"].(string))

	tearDownTest(t, server, db)
}
