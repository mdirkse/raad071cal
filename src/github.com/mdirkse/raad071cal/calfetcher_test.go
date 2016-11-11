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
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSuccessfulFetch(t *testing.T) {
	var reqURL string

	httpGet = func(url string) (*http.Response, error) {
		reqURL = url
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("callback_function({})")),
		}, nil
	}
	defer func() {
		httpGet = http.Get
	}()

	expected := "{}"
	result, _ := fetchCalendarMonthJSON(yearMonth{2016, 8})

	assert.Contains(t, reqURL, "year=2016&month=8", "Incorrect request URL!")
	assert.Equal(t, expected, result, "Calendar month incorrectly fetched!")
}

func TestFailureFetch(t *testing.T) {
	httpGet = func(string) (*http.Response, error) {
		return nil, errors.New("Eep!")
	}
	defer func() {
		httpGet = http.Get
	}()

	_, err := fetchCalendarMonthJSON(yearMonth{2016, 8})

	assert.NotNil(t, err, "Incorrect calendar month fetch did not result in an error!")
}

func TestGetCalendarItemsFromJSON(t *testing.T) {
	tstJSON, _ := ioutil.ReadFile("../../../../testfiles/tst.json")
	result, err := getCalendarItemsFromJSON(string(tstJSON), GetTestTime())
	expected := []CalItem{GetTestItem1(), GetTestItem2(), GetTestItem3()}

	assert.Nil(t, err, "Unable to get calendar items!")

	// reset docs to nil so we don't have to fake this in the expected struct as well
	for i := 0; i < len(result); i++ {
		result[i].Documents = nil
	}

	assert.Equal(t, expected, result, "Test JSON rendered incorrect items!")
}

func TestGenerateMonthYearRange(t *testing.T) {
	expected := []yearMonth{
		{2016, 1},
		{2016, 2},
		{2016, 3},
		{2016, 4},
		{2016, 5},
		{2016, 6},
		{2016, 7},
		{2016, 8},
		{2016, 9},
		{2016, 10},
		{2016, 11},
		{2016, 12},
		{2017, 1},
		{2017, 2},
		{2017, 3},
		{2017, 4},
		{2017, 5},
		{2017, 6},
	}

	result := generateMonthYearRange(GetTestTime())

	assert.Equal(t, expected, result, "Test item incorrectly rendered!")
}

func TestFetchCalendarItemsSuccess(t *testing.T) {
	testSet := []struct {
		srcFile     string
		eventsTotal int
	}{
		// We should have 18 (generated months) * 3 (test items) = 54 items
		// {"../../../../testfiles/tst.json", 54},
		// We should have 18 (generated months) * 5 (test items) = 90 items
		{"../../../../testfiles/ghi-1.json", 90},
	}

	// Make sure we reset the httpGet
	defer func() {
		httpGet = http.Get
	}()

	for _, i := range testSet {
		tstJSON, _ := ioutil.ReadFile(i.srcFile)
		httpGet = func(url string) (*http.Response, error) {
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("callback_function(" + string(tstJSON) + ")")),
			}, nil
		}

		items, err := fetchCalendarItems(GetTestTime())
		if err != nil {
			t.Errorf("Errors were returned. Something went wrong: [%+v]", err)
		}

		expected := i.eventsTotal
		if len(items) != expected {
			t.Errorf("Wrong amount of items returned! Expected %d but got %d!", expected, len(items))
		}
	}
}

func TestFetchCalendarItemsFailureInvalidJson(t *testing.T) {
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("callback_function({\"bla\": asfasdfasdf})")),
		}, nil
	}
	defer func() {
		httpGet = http.Get
	}()

	_, err := fetchCalendarItems(GetTestTime())
	if err == nil {
		t.Fatal("No error was returned even though the input is invalid!")
	}
}
