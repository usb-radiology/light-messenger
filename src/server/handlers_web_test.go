package server

import (
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

	response, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	defer response.Body.Close()

	// then
	assert.Equal(t, http.StatusOK, response.StatusCode)

	doc, errHTMLDoc := goquery.NewDocumentFromReader(response.Body)
	if errHTMLDoc != nil {
		t.Fatalf("%+v", errors.WithStack(errHTMLDoc))
	}

	links := doc.Find("a.is-large").Map(func(i int, s *goquery.Selection) string {
		return s.AttrOr("href", "")
	})

	expectedLinks := []string{"/mtra/ct", "/mtra/mr", "/mtra/nuk", "/radiologie/aod", "/radiologie/ctd", "/radiologie/msk", "/radiologie/nr", "/radiologie/nuk"}
	assert.EqualValues(t, expectedLinks, links)

	tearDownTest(t, server, db)
}
