package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	dateTimeLayout  = "2006-01-02T15:04"
	dateLayout      = "2006-01-02"
	itemTemplateSrc = `BEGIN:VEVENT
UID:19970610T172345Z-AF23B2@example.com
DTSTAMP:Europe/Amsterdam:{{.CreatedDateTime.Format "20060102T150405"}}
DTSTART:Europe/Amsterdam:{{.StartDateTime.Format "20060102T150405"}}
DTEND:Europe/Amsterdam:{{.EndDateTime.Format "20060102T150405"}}
SUMMARY:{{.Name}}
DESCRIPTION:{{.Description}}
LOCATION:{{.Location}}
END:VEVENT`
)

var (
	itemsRegex   *regexp.Regexp
	itemTemplate *template.Template
)

type CalItem struct {
	CreatedDateTime time.Time
	Description     string
	EndDateTime     time.Time
	Location        string
	Name            string
	Organizer       string
	StartDateTime   time.Time
	Url             string
}

func InitCalItemVars() {
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, CEST)
	itemsRegex = regexp.MustCompile(`var vdate='(.+)'.split\(`)
	itemTemplate, _ = template.New("item").Parse(itemTemplateSrc)
}

func NewItem(i string, runStart time.Time) (*CalItem, error) {
	fields := strings.Split(i, "!")

	var st time.Time
	var err error

	if fields[1] == "" {
		st, err = time.ParseInLocation(dateLayout, fmt.Sprintf("%s", fields[0]), CEST)
	} else {
		// Strip out "uur" if we have to
		timeField := strings.Replace(fields[1], " uur", "", 1)
		st, err = time.ParseInLocation(dateTimeLayout, fmt.Sprintf("%sT%s", fields[0], timeField), CEST)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to parse event date: %+v", err)
	}

	return &CalItem{
		CreatedDateTime: runStart,
		Description:     "bla",
		EndDateTime:     st.Add(2 * time.Hour),
		Location:        fields[6],
		Name:            fields[2],
		Organizer:       fields[7],
		StartDateTime:   st,
		Url:             fields[4],
	}, nil
}

func (i CalItem) RenderItem(w io.Writer) error {
	err := itemTemplate.Execute(w, i)

	if err != nil {
		return fmt.Errorf("Could not render the item [%+v]! (error: [%+v])", i, err)
	}

	return nil
}
