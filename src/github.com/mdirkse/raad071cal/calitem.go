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
	agendaUrlPrefix = "https://leiden.notudoc.nl"
	itemTemplateSrc = `BEGIN:VEVENT
UID:{{.Uid}}@raad071.mdirkse.nl
DTSTAMP:{{.CreatedDateTime.Format "20060102T150405"}}Z
{{- if .AllDay}}
DTSTART;VALUE=DATE:{{.StartDateTime.Format "20060102"}}
DTEND;VALUE=DATE:{{.EndDateTime.Format "20060102"}}
{{- else}}
DTSTART:{{.StartDateTime.Format "20060102T150405"}}Z
DTEND:{{.EndDateTime.Format "20060102T150405"}}Z
{{- end}}
SUMMARY:{{.Name}}
DESCRIPTION:Organisator: {{.Organizer}}\nStukken: {{.Url}}
LOCATION:{{.Location}}
END:VEVENT`
)

var (
	itemTemplate *template.Template
)

type CalItem struct {
	Uid             string
	AllDay          bool
	CreatedDateTime time.Time
	Url             string
	EndDateTime     time.Time
	Location        string
	Name            string
	Organizer       string
	StartDateTime   time.Time
}

func InitCalItemVars() {
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, CEST)
	itemTemplate, _ = template.New("item").Parse(itemTemplateSrc)
}

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
		st, err = time.ParseInLocation(dateTimeLayout, fmt.Sprintf("%sT%s", fields[0], timeField), CEST)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to parse event date: %+v", err)
	}

	// Change to UTC because that's what the iCal format likes best
	stUTC := st.In(time.UTC)

	return &CalItem{
		Uid:             GenerateId(stUTC, fields[2]),
		AllDay:          allDay,
		CreatedDateTime: runStart.In(time.UTC),
		Url:             RenderUrl(fields[4]),
		EndDateTime:     GetEndTime(allDay, stUTC),
		Location:        RenderLocation(fields[6]),
		Name:            strings.Title(fields[2]),
		Organizer:       strings.Title(fields[7]),
		StartDateTime:   stUTC,
	}, nil
}

func (i CalItem) RenderItem(w io.Writer) error {
	err := itemTemplate.Execute(w, i)

	if err != nil {
		return fmt.Errorf("Could not render the item [%+v]! (error: [%+v])", i, err)
	}

	return nil
}

func GenerateId(startTime time.Time, name string) string {
	timeStamp := startTime.Format("20060102T150405")
	data := []byte(timeStamp + name)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func GetEndTime(allDay bool, t time.Time) time.Time {
	if allDay {
		return t
	}

	return t.Add(2 * time.Hour)
}

func RenderUrl(o string) string {
	// Construct the description
	var description bytes.Buffer
	if o != "" {
		// If it's a relative URL then add the agenda part
		if strings.HasPrefix(o, "/") {
			description.WriteString(agendaUrlPrefix)
		}

		description.WriteString(o)
	}

	return description.String()
}

func RenderLocation(o string) string {
	lo := strings.ToLower(o)

	if lo == "raadzaal" || lo == "commissiekamer" {
		return fmt.Sprintf("%s, Stadhuis, Leiden", strings.Title(o))
	}

	return o
}
