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
	"crypto/md5"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"
)

const (
	dateTimeLayout  = "2006-01-02T15:04"
	dateLayout      = "2006-01-02"
	agendaURLPrefix = "https://leiden.notudoc.nl"
	itemTemplateSrc = `BEGIN:VEVENT
UID:{{.UID}}@raad071.mdirkse.nl
DTSTAMP:{{.CreatedDateTime.Format "20060102T150405"}}Z
{{- if .AllDay}}
DTSTART;VALUE=DATE:{{.StartDateTime.Format "20060102"}}
DTEND;VALUE=DATE:{{.EndDateTime.Format "20060102"}}
{{- else}}
DTSTART:{{.StartDateTime.Format "20060102T150405"}}Z
DTEND:{{.EndDateTime.Format "20060102T150405"}}Z
{{- end}}
SUMMARY:{{.Name}}
DESCRIPTION:Organisator: {{.Organizer}}\nStukken: {{.URL}}
LOCATION:{{.Location}}
END:VEVENT`
)

var (
	itemTemplate *template.Template
)

// CalItem represents a calendar item that can be rendered to iCal.
type CalItem struct {
	UID             string
	AllDay          bool
	CreatedDateTime time.Time
	URL             string
	EndDateTime     time.Time
	Location        string
	Name            string
	Organizer       string
	StartDateTime   time.Time
}

func initCalItemVars() {
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, cestTz)
	itemTemplate, _ = template.New("item").Parse(itemTemplateSrc)
}

// NewItem creates a new calendar item from a string input
func NewItem(i string, runStart time.Time) (*CalItem, error) {
	fields := strings.Split(i, "!")

	var st time.Time
	var err error
	var allDay bool

	if fields[1] == "" {
		// No time set, we have a full-day event
		// Parse it with UTC because if we convert it from CEST (ie subtract 2 hours) we will end up on the day before
		st, err = time.ParseInLocation(dateLayout, fields[0], time.UTC)
		allDay = true
	} else {
		// Strip out "uur" if we have to
		timeField := strings.Replace(fields[1], " uur", "", 1)
		st, err = time.ParseInLocation(dateTimeLayout, fmt.Sprintf("%sT%s", fields[0], timeField), cestTz)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to parse event date: %+v", err)
	}

	// Change to UTC because that's what the iCal format likes best
	stUTC := st.In(time.UTC)

	return &CalItem{
		UID:             generateID(stUTC, fields[2]),
		AllDay:          allDay,
		CreatedDateTime: runStart.In(time.UTC),
		URL:             renderURL(fields[4]),
		EndDateTime:     getEndTime(allDay, stUTC),
		Location:        renderLocation(fields[6]),
		Name:            strings.Title(fields[2]),
		Organizer:       strings.Title(fields[7]),
		StartDateTime:   stUTC,
	}, nil
}

// RenderItem renders a calendar item in iCalendar format
func (i CalItem) RenderItem(w io.Writer) error {
	err := itemTemplate.Execute(w, i)

	if err != nil {
		return fmt.Errorf("Could not render the item [%+v]! (error: [%+v])", i, err)
	}

	return nil
}

func generateID(startTime time.Time, name string) string {
	timeStamp := startTime.Format("20060102T150405")
	data := []byte(timeStamp + name)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func getEndTime(allDay bool, t time.Time) time.Time {
	if allDay {
		return t
	}

	return t.Add(3 * time.Hour)
}

func renderURL(o string) string {
	// Construct the description
	var description bytes.Buffer
	if o != "" {
		// If it's a relative URL then add the agenda part
		if strings.HasPrefix(o, "/") {
			description.WriteString(agendaURLPrefix)
		}

		description.WriteString(o)
	}

	return description.String()
}

func renderLocation(o string) string {
	lo := strings.ToLower(o)

	if lo == "raadzaal" || lo == "commissiekamer" {
		return fmt.Sprintf("%s, Stadhuis, Leiden", strings.Title(o))
	}

	return o
}
