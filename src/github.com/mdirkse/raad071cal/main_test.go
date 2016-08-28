// Copyright 2016 Maarten Dirkse
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initCalFetcherVars()
	code := m.Run()
	os.Exit(code)
}

//func TestParseCalendarShouldYieldItems(t *testing.T) {
//	var files = []struct {
//		location   string
//		totalItems int
//	}{
//		{"../../../../testfiles/errorparse.html", 4},
//		{"../../../../testfiles/normalparse.html", 5},
//		{"../../../../testfiles/outdatedparse.html", 3},
//	}
//
//	for _, pt := range files {
//		b, _ := ioutil.ReadFile(pt.location)
//		items, _ := parseCalendar(&b)
//
//		if len(items) != pt.totalItems {
//			t.Errorf("Amount returned items does not match! Expected %d but got %d.", pt.totalItems, len(items))
//		}
//
//	}
//}
//
//func TestParseCalenderWithWrongInputShouldYieldError(t *testing.T) {
//	nonsenseCalSource := []byte("the regex will not trigger on this string")
//	_, err := parseCalendar(&nonsenseCalSource)
//
//	if err == nil {
//		t.Fatal("Nonsense calendar string should have produced an error!")
//	}
//}
//
//func TestFetchCalenderPageWithBrokenUrlShouldYieldError(t *testing.T) {
//	_, err := fetchCalenderPage("http://localhost:60606")
//
//	if err == nil {
//		t.Fatal("Fetching a broken URL should have produced an error!")
//	}
//}

func TestRenderCalendarShouldYieldCorrectOutput(t *testing.T) {
	testCals := []CalItem{GetTestItem1(), GetTestItem2(), GetTestItem3()}

	var testICal bytes.Buffer
	testICal.WriteString(calendarHeader)
	for _, c := range testCals {
		c.RenderItem(&testICal)
		testICal.WriteByte('\n')
	}
	testICal.WriteString(calendarFooter)

	emptyCal := calendarHeader + calendarFooter

	var iCals = []struct {
		expected string
		items    []CalItem
	}{
		{testICal.String(), testCals},
		{emptyCal, []CalItem{}},
	}

	for _, ct := range iCals {
		var result bytes.Buffer
		renderCalendar(ct.items, &result)

		if ct.expected != result.String() {
			t.Errorf("Render went awry! Expected:\n%s \n\nbut got:\n%s", ct.expected, result.String())
		}
	}
}

func TestHttpEndpointRequestShouldYieldCorrectOutput(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://bla.com", nil)

	w := httptest.NewRecorder()
	calHandler().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Request should return status 200 but was %d!", w.Code)
	}

	expected := calendarHeader + calendarFooter

	if expected != w.Body.String() {
		t.Errorf("Request went awry! Expected:\n%s \n\nbut got:\n%s", expected, w.Body.String())
	}
}
