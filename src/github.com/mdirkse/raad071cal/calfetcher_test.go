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
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
	defer func() { httpGet = http.Get }()

	expected := "{}"
	result, _ := fetchCalendarMonthJSON(yearMonth{2016, 8})

	if !strings.Contains(reqURL, "year=2016&month=8") {
		t.Errorf("Incorrect request URL! Was\n[%s]", reqURL)
	}

	if expected != result {
		t.Errorf("Calendar month incorrectly fetched! Expected \n[%s] \nbut got\n[%s]!", expected, result)
	}
}

func TestFailureFetch(t *testing.T) {
	httpGet = func(string) (*http.Response, error) {
		return nil, errors.New("Eep!")
	}
	defer func() { httpGet = http.Get }()

	_, err := fetchCalendarMonthJSON(yearMonth{2016, 8})

	if err == nil {
		t.Fatal("Incorrect calendar month fetch did not result in an error!")
	}
}

func TestGetCalendarItemsFromJSON(t *testing.T) {
	tstJSON, _ := ioutil.ReadFile("../../../../testfiles/tst.json")
	result, err := getCalendarItemsFromJSON(string(tstJSON), GetTestTime())
	expected := []CalItem{GetTestItem1(), GetTestItem2(), GetTestItem3()}

	if err != nil {
		t.Errorf("Unable to get calendar items! Error: [%+v]", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test JSON rendered incorrect items! Expected \n%v \nbut got\n%v!", expected, result)
	}
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

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Test item incorrectly rendered! Expected \n%v \nbut got\n%v!", expected, result)
	}
}

func TestFetchCalendarItemsSuccess(t *testing.T) {
	tstJSON, _ := ioutil.ReadFile("../../../../testfiles/tst.json")
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("callback_function(" + string(tstJSON) + ")")),
		}, nil
	}
	defer func() { httpGet = http.Get }()

	items, err := fetchCalendarItems(GetTestTime())
	if err != nil {
		t.Errorf("Errors were returned. Something went wrong: [%+v]", err)
	}

	// We should have 18 (generated months) * 3 (test items) = 54 items
	expected := 54
	if len(items) != expected {
		t.Errorf("Wrong amount of items returned! Expected %d but got %d!", expected, len(items))
	}
}

func TestFetchCalendarItemsFailureInvalidJson(t *testing.T) {
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString("callback_function({\"bla\": asfasdfasdf})")),
		}, nil
	}
	defer func() { httpGet = http.Get }()

	_, err := fetchCalendarItems(GetTestTime())
	if err == nil {
		t.Fatal("No error was returned even though the input is invalid!")
	}
}
