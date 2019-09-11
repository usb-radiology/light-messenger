package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery" // https://godoc.org/github.com/PuerkitoBio/goquery
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func getResponseHTMLDoc(t *testing.T, request *http.Request) *goquery.Document {
	request.Header.Set(HTMLHeaderContentType, HTMLHeaderContentTypeValueHTML)

	response := getResponse(t, request)
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	doc, errHTMLDoc := goquery.NewDocumentFromReader(response.Body)
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
		t.Fatalf("%+v", errJSONDecode)
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
