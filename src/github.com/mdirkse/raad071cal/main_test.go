package main

import (
	"bytes"
	"io/ioutil"
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
		{"../../../../testfiles/test.ical", []*CalItem{GetTestItem(), GetTestItem()}},
		{"../../../../testfiles/empty.ical", []*CalItem{}},
	}

	iCals[0].items[1].Name = "Muhahahaha"

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
