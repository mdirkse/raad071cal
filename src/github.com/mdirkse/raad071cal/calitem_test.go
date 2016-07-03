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
		{GetTestItem2(), "2016-06-23!16:00 uur!gemeenteraad!#!/some-url/here!is-agenda!raadzaal!gemeenteraad"},
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
		Uid:             "14c88381339fffd3963618fe8cf93825",
		AllDay:          true,
		CreatedDateTime: GetTestTime().In(time.UTC),
		Url:             "http://what.have.you/done/for/me/lately.pdf",
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
		Uid:             "c3bf3a489d812fd932c843a5b791d958",
		CreatedDateTime: GetTestTime().In(time.UTC),
		Url:             agendaUrlPrefix + "/some-url/here",
		EndDateTime:     GetTestTime().In(time.UTC).Add(2 * time.Hour),
		Location:        "Raadzaal, Stadhuis, Leiden",
		Name:            "Gemeenteraad",
		Organizer:       "Gemeenteraad",
		StartDateTime:   GetTestTime().In(time.UTC),
	}
}

func GetRenderedTestItem2() string {
	return `BEGIN:VEVENT
UID:c3bf3a489d812fd932c843a5b791d958@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART:20160623T140000Z
DTEND:20160623T160000Z
SUMMARY:Gemeenteraad
DESCRIPTION:Organisator: Gemeenteraad\nStukken: https://leiden.notudoc.nl/some-url/here
LOCATION:Raadzaal, Stadhuis, Leiden
END:VEVENT`
}

func GetTestTime() time.Time {
	return time.Date(2016, time.June, 23, 16, 0, 0, 0, CEST)
}
