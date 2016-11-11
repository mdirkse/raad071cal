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
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, ct.expected, result.String(), "Render went awry!")
	}
}

func TestHttpEndpointRequestShouldYieldCorrectOutput(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://bla.com", nil)

	w := httptest.NewRecorder()
	calHandler().ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Request returned incorrect status!")
	assert.Equal(t, calendarHeader+calendarFooter, w.Body.String(), "Request went awry!")
}
