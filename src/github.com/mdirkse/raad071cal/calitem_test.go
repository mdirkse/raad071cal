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
	"reflect"
	"testing"
	"time"
)

func TestParseValidItemShouldYieldCorrectStruct(t *testing.T) {
	initCalFetcherVars()

	testSet := []struct {
		expected *CalItem
		input    string
	}{
		{GetTestItem1(), "2016-06-23!!The rent!#!http://what.have.you/done/for/me/lately.pdf!is-agenda!My house!80's women"},
		{GetTestItem2(), "2016-06-23!20:00 uur!gemeenteraad!#!/some-url/here!is-agenda!raadzaal!gemeenteraad"},
	}

	for _, i := range testSet {
		result, _ := NewItem(i.input, GetTestTime())

		if !reflect.DeepEqual(i.expected, result) {
			t.Fatalf("Test item incorrectly parsed! \n Expected [%+v] \n but got  [%+v]!", i.expected, result)
		}
	}
}

func TestParseItemWithInvalidDateShouldYieldAnError(t *testing.T) {
	initCalFetcherVars()

	tstInput := "2016-06-23!16:safdsfd!!!!!!"

	_, err := NewItem(tstInput, GetTestTime())

	if err == nil {
		t.Fatalf("Test item incorrectly parsed when it shouldn't have been!")
	}
}

func TestVariousSpecialCasesForMeetingDuration(t *testing.T) {
	initCalFetcherVars()

	// Gemeenteraad should always last 'till 21h UTC
	gr21, _ := NewItem("2016-06-23!16:00 uur!gemeenteraad!#!/some-url/here!is-agenda!raadzaal!gemeenteraad", GetTestTime())
	if gr21.EndDateTime.Hour() != 21 {
		t.Fatalf("Gemeenteraad item has wrong end time. Expected 21 but was %d!", gr21.EndDateTime.Hour())
	}

	// College should always last 3h
	col3h, _ := NewItem("2016-06-23!16:00 uur!College Burgemeester en Wethouders!#!/some-url/here!is-agenda!raadzaal!gemeenteraad", GetTestTime())
	col3hDuration := col3h.EndDateTime.Sub(col3h.StartDateTime).Hours()
	if col3hDuration != 3 {
		t.Fatalf("College lasted wrong amount of hours. Expected 3 but was %f!", col3hDuration)
	}

	// Commissie should always last 3h
	comm3h, _ := NewItem("2016-06-23!16:00 uur!College Burgemeester en Wethouders!#!/some-url/here!is-agenda!raadzaal!gemeenteraad", GetTestTime())
	comm3hDuration := comm3h.EndDateTime.Sub(comm3h.StartDateTime).Hours()
	if col3hDuration != 3 {
		t.Fatalf("Commission lasted wrong amount of hours. Expected 3 but was %f!", comm3hDuration)
	}
}

func TestRenderItemShouldYieldCorrectICalEvent(t *testing.T) {
	initCalFetcherVars()

	testSet := []struct {
		input    *CalItem
		expected string
	}{
		{GetTestItem1(), GetRenderedTestItem1()},
		{GetTestItem2(), GetRenderedTestItem2()},
	}

	for _, i := range testSet {
		var result bytes.Buffer
		i.input.RenderItem(&result)

		if i.expected != result.String() {
			t.Fatalf("Test item incorrectly rendered! Expected \n[%s] \nbut got\n[%s]!", i.expected, result.String())
		}
	}
}

func GetTestItem1() *CalItem {
	return &CalItem{
		UID:             "14c88381339fffd3963618fe8cf93825",
		AllDay:          true,
		CreatedDateTime: GetTestTime().In(time.UTC),
		URL:             "http://what.have.you/done/for/me/lately.pdf",
		EndDateTime:     GetTestTime().Add(-14 * time.Hour).In(time.UTC), // Correct to 0 hours for allDay
		Location:        "My house",
		Name:            "The Rent",
		Organizer:       "80'S Women",
		StartDateTime:   GetTestTime().Add(-14 * time.Hour).In(time.UTC), // Correct to 0 hours for allDay
	}
}

func GetRenderedTestItem1() string {
	return `BEGIN:VEVENT
UID:14c88381339fffd3963618fe8cf93825@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART;VALUE=DATE:20160623
DTEND;VALUE=DATE:20160623
SUMMARY:The Rent
DESCRIPTION:Organisator: 80'S Women\nStukken: http://what.have.you/done/for/me/lately.pdf
LOCATION:My house
END:VEVENT`
}

func GetTestItem2() *CalItem {
	return &CalItem{
		UID:             "37384d42f0c7fe4ad6103b2b7344bbe7",
		CreatedDateTime: GetTestTime().In(time.UTC),
		URL:             agendaURLPrefix + "/some-url/here",
		EndDateTime:     GetTestTime().In(time.UTC).Add(7 * time.Hour),
		Location:        "Raadzaal, Stadhuis, Leiden",
		Name:            "Gemeenteraad",
		Organizer:       "Gemeenteraad",
		StartDateTime:   GetTestTime().In(time.UTC).Add(4 * time.Hour),
	}
}

func GetRenderedTestItem2() string {
	return `BEGIN:VEVENT
UID:37384d42f0c7fe4ad6103b2b7344bbe7@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART:20160623T180000Z
DTEND:20160623T210000Z
SUMMARY:Gemeenteraad
DESCRIPTION:Organisator: Gemeenteraad\nStukken: https://leiden.notudoc.nl/some-url/here
LOCATION:Raadzaal, Stadhuis, Leiden
END:VEVENT`
}

func GetTestTime() time.Time {
	return time.Date(2016, time.June, 23, 16, 0, 0, 0, cestTz)
}
