package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseCalendarShouldYieldItems(t *testing.T) {
	initCalFetcherVars()

	var files = []struct {
		location   string
		totalItems int
	}{
		{"../../../../testfiles/errorparse.html", 4},
		{"../../../../testfiles/normalparse.html", 5},
		{"../../../../testfiles/outdatedparse.html", 3},
	}

	for _, pt := range files {
		b, _ := ioutil.ReadFile(pt.location)
		items, _ := parseCalendar(&b)

		if len(*items) != pt.totalItems {
			t.Fatalf("Amount returned items does not match! Expected %d but got %d.", pt.totalItems, len(*items))
		}

	}
}

func TestParseCalenderWithWrongInputShouldYieldError(t *testing.T) {
	initCalFetcherVars()

	nonsenseCalSource := []byte("the regex will not trigger on this string")
	_, err := parseCalendar(&nonsenseCalSource)

	if err == nil {
		t.Fatal("Nonsense calendar string should have produced an error!")
	}
}

func TestFetchCalenderPageWithBrokenUrlShouldYieldError(t *testing.T) {
	initCalFetcherVars()

	_, err := fetchCalenderPage("http://localhost:60606")

	if err == nil {
		t.Fatal("Fetching a broken URL should have produced an error!")
	}
}

func TestRenderCalendarShouldYieldCorrectOutput(t *testing.T) {
	initCalFetcherVars()

	var iCals = []struct {
		location string
		items    []*CalItem
	}{
		{"../../../../testfiles/test.ical", []*CalItem{GetTestItem1(), GetTestItem2()}},
		{"../../../../testfiles/empty.ical", []*CalItem{}},
	}

	for _, ct := range iCals {
		var result bytes.Buffer
		renderCalendar(&ct.items, &result)

		b, _ := ioutil.ReadFile(ct.location)
		expected := string(b)

		if expected != result.String() {
			t.Fatalf("Render went awry! Expected:\n%s \n\nbut got:\n%s", expected, result.String())
		}
	}
}

func TestHttpEndpointRequestShouldYieldCorrectOutput(t *testing.T) {
	initCalFetcherVars()

	req, err := http.NewRequest("GET", "http://bla.com", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	calHandler().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Request should return status 200 but was %d!", w.Code)
	}

	b, _ := ioutil.ReadFile("../../../../testfiles/empty.ical")
	expected := string(b)

	if expected != w.Body.String() {
		t.Fatalf("Request went awry! Expected:\n%s \n\nbut got:\n%s", expected, w.Body.String())
	}
}
