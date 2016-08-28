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
	"unicode"
)

const (
	dateTimeLayout  = "02-01-2006 15:04"
	agendaURLPrefix = "https://leiden.notubiz.nl"
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
SUMMARY:{{.Description}}
DESCRIPTION:{{- if .Link}}Notubiz link: {{.Link}}\n{{- end}}{{- if .Documents}}Documents:\n{{range .Documents}}- {{.Title}} {{.URL}}\n{{end}}{{- end}}
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
	Canceled        bool       `json:"canceled"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	Link            string     `json:"link"`
	Documents       []document `json:"documents"`
	Date            string     `json:"date"`
	CreatedDateTime time.Time
	Time            string `json:"time"`
	StartDateTime   time.Time
	EndDateTime     time.Time
}

type document struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func initCalItemVars() {
	itemTemplate = template.Must(template.New("item").Parse(itemTemplateSrc))
}

// EnrichItem creates a new calendar item from a string input
func EnrichItem(i CalItem, runStart time.Time) (CalItem, error) {
	i.CreatedDateTime = runStart.In(time.UTC)

	// Figure out what timezone to parse the date/time with
	// If we've got a full day event (ie start time == 00:00) then
	// parse it with UTC because if we convert it from CEST (ie subtract
	// 2 hours) we will end up on the day before
	var tz *time.Location

	if i.Time == "00:00" {
		i.AllDay = true
		tz = time.UTC
	} else {
		tz = cestTz
	}

	st, err := time.ParseInLocation(dateTimeLayout, fmt.Sprintf("%s %s", i.Date, i.Time), tz)

	if err != nil {
		return CalItem{}, fmt.Errorf("Unable to parse event date: %+v", err)
	}

	if st.Location() == time.UTC {
		i.StartDateTime = st
	} else {
		i.StartDateTime = st.In(time.UTC)
	}

	i.EndDateTime = getEndTime(i)

	i.Description = upperCaseFirstLetter(i.Description)
	i.UID = generateID(i)
	i.Link = renderLink(i)
	i.Location = renderLocation(i.Location)

	return i, nil
}
func upperCaseFirstLetter(i string) string {
	s := []rune(i)
	s[0] = unicode.ToUpper(s[0])
	return string(s)
}

// RenderItem renders a calendar item in iCalendar format
func (i CalItem) RenderItem(w io.Writer) error {
	err := itemTemplate.Execute(w, i)

	if err != nil {
		return fmt.Errorf("Could not render the item [%+v]! (error: [%+v])", i, err)
	}

	return nil
}

func generateID(i CalItem) string {
	timeStamp := i.StartDateTime.Format(dateTimeLayout)
	data := []byte(timeStamp + i.Description)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func getEndTime(i CalItem) time.Time {
	if i.AllDay {
		return i.StartDateTime
	}

	// Switch on the first word of the name
	switch strings.ToLower(strings.Split(i.Description, " ")[0]) {
	case "gemeenteraad":
		return time.Date(i.StartDateTime.Year(), i.StartDateTime.Month(), i.StartDateTime.Day(), 21, 0, 0, 0, time.UTC)
	case "raadscommissie", "college":
		return i.StartDateTime.Add(3 * time.Hour)
	default:
		return i.StartDateTime.Add(2 * time.Hour)
	}
}

func renderLink(i CalItem) string {
	// Construct the description
	var description bytes.Buffer
	if i.Link != "" {
		// If it's a relative URL then add the agenda part
		if strings.HasPrefix(i.Link, "/") {
			description.WriteString(agendaURLPrefix)
		}

		description.WriteString(strings.Replace(i.Link, " ", "%20", -1))
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
