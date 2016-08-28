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
	"strings"
	"testing"
	"time"
	"unicode"
)

const (
	testDateFormat = "02-01-2006"
)

func TestItemsShouldBeEnrichedCorrectly(t *testing.T) {
	testSet := []struct {
		expected CalItem
		input    CalItem
	}{
		{GetTestItem1(), deEnrich(GetTestItem1())},
		{GetTestItem2(), deEnrich(GetTestItem2())},
		{GetTestItem3(), deEnrich(GetTestItem3())},
	}

	for _, i := range testSet {
		result, _ := EnrichItem(i.input, GetTestTime())

		if !reflect.DeepEqual(i.expected, result) {
			t.Errorf("Test item incorrectly parsed! \n Expected [%+v] \n but got  [%+v]!", i.expected, result)
		}
	}
}

func TestParseItemWithInvalidDateShouldYieldAnError(t *testing.T) {
	tstInput := GetTestItem1()
	tstInput.Date = "bladibla"

	_, err := EnrichItem(tstInput, GetTestTime())

	if err == nil {
		t.Fatal("Faulty test item parsed when it shouldn't have been!")
	}
}

func TestVariousSpecialCasesForMeetingDuration(t *testing.T) {
	// Gemeenteraad should always last 'till 21h UTC
	gri := deEnrich(GetTestItem1())
	gri.Description = "Gemeenteraad"
	gri.Time = "16:00"
	gr21, _ := EnrichItem(gri, GetTestTime())
	if gr21.EndDateTime.Hour() != 21 {
		t.Errorf("Gemeenteraad item has wrong end time. Expected 21 but was %d!", gr21.EndDateTime.Hour())
	}

	// College should always last 3h

	ci := deEnrich(GetTestItem1())
	ci.Description = "College Burgemeester en Wethouders"
	ci.Time = "16:00"
	col3h, _ := EnrichItem(ci, GetTestTime())
	col3hDuration := col3h.EndDateTime.Sub(col3h.StartDateTime).Hours()
	if col3hDuration != 3 {
		t.Errorf("College lasted wrong amount of hours. Expected 3 but was %f!", col3hDuration)
	}
}

func TestRenderItemShouldYieldCorrectICalEvent(t *testing.T) {
	testSet := []struct {
		input    CalItem
		expected string
	}{
		{GetTestItem1(), GetRenderedTestItem1()},
		{GetTestItem2(), GetRenderedTestItem2()},
		{GetTestItem3(), GetRenderedTestItem3()},
	}

	for _, i := range testSet {
		var result bytes.Buffer
		err := i.input.RenderItem(&result)

		if err != nil {
			t.Errorf("Unable to render test item! Error: [%+v]", err)
		}

		if i.expected != result.String() {
			t.Errorf("Test item incorrectly rendered! Expected \n[%s] \nbut got\n[%s]!", i.expected, result.String())
		}
	}
}

func deEnrich(i CalItem) CalItem {
	i.UID = ""
	i.AllDay = false
	i.Link = strings.Replace(i.Link, agendaURLPrefix, "", 1)
	i.Location = strings.Split(i.Location, ",")[0]
	i.CreatedDateTime = time.Time{}
	i.StartDateTime = time.Time{}
	i.EndDateTime = time.Time{}

	// Lowercase first letter of description
	d := []rune(i.Description)
	d[0] = unicode.ToLower(d[0])
	i.Description = string(d)

	return i
}

func GetTestItem1() CalItem {
	iTime := GetTestTime().Add(-14 * time.Hour)

	return CalItem{
		UID:             "e058fd25aa867090dd7e25c9455d7156",
		AllDay:          true,
		Documents:       []document{},
		Link:            "",
		Location:        "",
		Description:     "Einde zomerreces",
		Date:            iTime.Format(testDateFormat),
		Time:            "00:00",
		CreatedDateTime: GetTestTime().In(time.UTC),
		StartDateTime:   iTime.In(time.UTC), // Correct to 0 hours for allDay
		EndDateTime:     iTime.In(time.UTC), // Correct to 0 hours for allDay
	}
}

func GetRenderedTestItem1() string {
	return `BEGIN:VEVENT
UID:e058fd25aa867090dd7e25c9455d7156@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART;VALUE=DATE:20160623
DTEND;VALUE=DATE:20160623
SUMMARY:Einde zomerreces
DESCRIPTION:
LOCATION:
END:VEVENT`
}

func GetTestItem2() CalItem {
	return CalItem{
		UID:         "a2dc05212385ac8b98a4ded4e09e952c",
		AllDay:      false,
		Link:        agendaURLPrefix + "/raad071cal.html",
		Location:    "Raadzaal, Stadhuis, Leiden",
		Description: "Instructiebijeenkomst Raad071Cal",
		Documents: []document{
			{
				Title: "iCal spec",
				URL:   "https://www.ietf.org/rfc/rfc2445.txt",
			},
			{
				Title: "History of the calendar",
				URL:   "https://en.wikipedia.org/wiki/Calendar",
			},
		},
		Date:            GetTestTime().Add(4 * time.Hour).Format(testDateFormat),
		Time:            "19:00",
		CreatedDateTime: GetTestTime().In(time.UTC),
		StartDateTime:   GetTestTime().In(time.UTC).Add(3 * time.Hour),
		EndDateTime:     GetTestTime().In(time.UTC).Add(5 * time.Hour),
	}
}

func GetRenderedTestItem2() string {
	return `BEGIN:VEVENT
UID:a2dc05212385ac8b98a4ded4e09e952c@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART:20160623T170000Z
DTEND:20160623T190000Z
SUMMARY:Instructiebijeenkomst Raad071Cal
DESCRIPTION:Notubiz link: https://leiden.notubiz.nl/raad071cal.html\nDocuments:\n- iCal spec https://www.ietf.org/rfc/rfc2445.txt\n- History of the calendar https://en.wikipedia.org/wiki/Calendar\n
LOCATION:Raadzaal, Stadhuis, Leiden
END:VEVENT`
}

func GetTestItem3() CalItem {
	return CalItem{
		UID:             "7599ab178274a0adcbee1b7e80f72bed",
		AllDay:          false,
		Documents:       []document{},
		Link:            agendaURLPrefix + "/vergadering/247980/raadscommissie%20Stedelijke%20Ontwikkeling%2001-09-2016",
		Location:        "Commissiekamer, Stadhuis, Leiden",
		Description:     "Raadscommissie Stedelijke Ontwikkeling",
		Date:            GetTestTime().Add(4 * time.Hour).Format(testDateFormat),
		Time:            "20:00",
		CreatedDateTime: GetTestTime().In(time.UTC),
		StartDateTime:   GetTestTime().In(time.UTC).Add(4 * time.Hour),
		EndDateTime:     GetTestTime().In(time.UTC).Add(7 * time.Hour),
	}
}

func GetRenderedTestItem3() string {
	return `BEGIN:VEVENT
UID:7599ab178274a0adcbee1b7e80f72bed@raad071.mdirkse.nl
DTSTAMP:20160623T140000Z
DTSTART:20160623T180000Z
DTEND:20160623T210000Z
SUMMARY:Raadscommissie Stedelijke Ontwikkeling
DESCRIPTION:Notubiz link: https://leiden.notubiz.nl/vergadering/247980/raadscommissie%20Stedelijke%20Ontwikkeling%2001-09-2016\n
LOCATION:Commissiekamer, Stadhuis, Leiden
END:VEVENT`
}

func GetTestTime() time.Time {
	return time.Date(2016, time.June, 23, 16, 0, 0, 0, cestTz)
}
