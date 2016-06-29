package main

import (
	"fmt"
	"github.com/robfig/cron"
	log "gopkg.in/inconshreveable/log15.v2"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	listenAddress      = ":7070"
	raad071CalendarUrl = "http://leiden.notudoc.nl/cgi-bin/calendar.cgi"
)

var (
	logger     log.Logger
	calItems   *[]*CalItem
	CEST       *time.Location
	cronT      *cron.Cron
	cutoffDate time.Time
)

func main() {
	initCalFetcherVars()
	logger.Info("Starting raad071cal")

	// Configure periodic polling
	cronT.Start()

	go loadCalendarItems() // do initial load

	//http.Handle("/raad071metrics", prometheus.Handler())

	//logger.Info(fmt.Sprintf("Fully initialised and listening on [%s].", listenAddress))
	//http.ListenAndServe(listenAddress, nil)
}

func initCalFetcherVars() {
	logger = log.New()
	CEST, _ = time.LoadLocation("Europe/Amsterdam")
	cronT = cron.New()
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, CEST)

	InitCalItemVars()
}

func loadCalendarItems() {
	pageBytes, err := fetchCalenderPage(raad071CalendarUrl)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to fetch the calendar page: %+v", err))
		return
	}

	newCalItems, err := parseCalendar(pageBytes)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to parse the calendar: %+v", err))
		return
	}

	calItems = newCalItems // IS THIS THREAD-SAFE?
}

func fetchCalenderPage(calendarUrl string) (*[]byte, error) {
	resp, err := http.Get(calendarUrl)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch the calendar from [%s]: %+v", calendarUrl, err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read the calender URL contents: %+v", err)
	}

	return &bytes, nil
}

func parseCalendar(pageBytes *[]byte) (*[]*CalItem, error) {
	logger.Info("Parsing the calendar.")

	runStart := time.Now()
	calBytes := itemsRegex.FindSubmatch(*pageBytes)

	if len(calBytes) != 2 {
		return nil, fmt.Errorf("Could not find calendar items in text:\n%s", string(*pageBytes))
	}

	calText := string(calBytes[1])
	rawItems := strings.Split(calText, "@")
	items := make([]*CalItem, 0, len(rawItems))

	for _, c := range rawItems {
		i, err := NewItem(c, runStart)

		if err != nil {
			logger.Error(fmt.Sprintf("Unable to parse item [%s]: %+v", c, err))
			continue
		}

		if i.StartDateTime.After(cutoffDate) {
			items = append(items, i)
		}
	}

	logger.Info(fmt.Sprintf("Parsed %d calendar items.", len(items)))

	return &items, nil
}

func renderCalendar(items *[]*CalItem, w io.Writer) {
	w.Write([]byte("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//mdirkse/raad071cal//NONSGML v1.0//EN\n"))

	for _, c := range *items {
		c.RenderItem(w)
		w.Write([]byte("\n"))
	}

	w.Write([]byte("END:VCALENDAR"))
}
