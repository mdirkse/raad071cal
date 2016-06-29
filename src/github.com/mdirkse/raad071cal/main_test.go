package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseCalendarShouldYieldItems(t *testing.T) {
	initCalFetcher()

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
		items := parseCalendar(&b)

		if len(*items) != pt.totalItems {
			t.Fatalf(fmt.Sprintf("Amount returned items does not match! Expected %d but got %d.", pt.totalItems, len(*items)))
		}

	}
}
