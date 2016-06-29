package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestParseValidItemShouldYieldCorrectStruct(t *testing.T) {
	initCalFetcherVars()

	tstInput := []string{
		"2016-06-23!16:00!The rent!#!/some-url/here!is-agenda!My house!80's women",
		"2016-06-23!16:00 uur!The rent!#!/some-url/here!is-agenda!My house!80's women",
	}
	expected := GetTestItem()

	for _, i := range tstInput {
		result, _ := NewItem(i, GetTestTime())

		if !reflect.DeepEqual(expected, result) {
			t.Fatalf("Test item incorrectly parsed! \n Expected [%+v] \n but got  [%+v]!", expected, result)
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

func TestRenderItemShouldYieldCorrectICalEvent(t *testing.T) {
	initCalFetcherVars()

	var result bytes.Buffer
	GetTestItem().RenderItem(&result)

	if GetRenderedTestItem() != result.String() {
		t.Fatalf("Test item incorrectly rendered! Expected \n[%s] \nbut got\n[%s]!", GetRenderedTestItem(), result.String())
	}
}

func GetTestItem() *CalItem {
	return &CalItem{
		CreatedDateTime: GetTestTime(),
		Description:     "bla",
		EndDateTime:     GetTestTime().Add(2 * time.Hour),
		Location:        "My house",
		Name:            "The rent",
		Organizer:       "80's women",
		StartDateTime:   GetTestTime(),
		Url:             "/some-url/here",
	}
}

func GetRenderedTestItem() string {
	return `BEGIN:VEVENT
UID:19970610T172345Z-AF23B2@example.com
DTSTAMP:Europe/Amsterdam:20160623T160000
DTSTART:Europe/Amsterdam:20160623T160000
DTEND:Europe/Amsterdam:20160623T180000
SUMMARY:The rent
DESCRIPTION:bla
LOCATION:My house
END:VEVENT`
}

func GetTestTime() time.Time {
	return time.Date(2016, time.June, 23, 16, 0, 0, 0, CEST)
}
