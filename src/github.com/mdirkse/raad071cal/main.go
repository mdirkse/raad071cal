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
	"fmt"
	"github.com/robfig/cron"
	log "gopkg.in/inconshreveable/log15.v2"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	listenAddress      = ":80"
	raad071CalendarUrl = "http://leiden.notudoc.nl/cgi-bin/calendar.cgi"
)

var (
	logger     log.Logger
	calItems   *[]*CalItem
	CEST       *time.Location
	cronT      *cron.Cron
	cutoffDate time.Time
	itemsRegex *regexp.Regexp
)

func main() {
	initCalFetcherVars()
	logger.Info("Starting raad071cal")

	// Configure periodic polling
	logger.Info(fmt.Sprintf("Polling source calendar [%s] every hour.", raad071CalendarUrl))
	cronT.AddFunc("1 1 * * * *", loadCalendarItems)
	cronT.Start()

	http.Handle("/kalender/alles.ics", loggingHandler(calHandler()))
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("html"))))

	logger.Info(fmt.Sprintf("Fully initialised and listening on [%s].", listenAddress))
	go loadCalendarItems() // do initial load

	http.ListenAndServe(listenAddress, nil)
}

func initCalFetcherVars() {
	logger = log.New()
	CEST, _ = time.LoadLocation("Europe/Amsterdam")
	cronT = cron.New()
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, CEST)
	itemsRegex = regexp.MustCompile(`var vdate='(.+)'.split\(`)

	emptyCal := []*CalItem{}
	calItems = &emptyCal

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

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(fmt.Sprintf("bla: [%s] [%s] %s", r.RemoteAddr, r.Method, r.URL.Path))
		h.ServeHTTP(w, r)
	})
}

func calHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/calendar")
		w.Header().Set("Cache-Control", "max-age=3600")

		if err := renderCalendar(calItems, w); err != nil {
			http.Error(w, "Couldn't render calendar items!", http.StatusInternalServerError)
		}
	})
}

func renderCalendar(items *[]*CalItem, w io.Writer) error {
	_, err := w.Write([]byte("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//mdirkse/raad071cal//NONSGML v1.0//EN\n"))

	if err != nil {
		return fmt.Errorf("Could not write calendar!")
	}

	for _, c := range *items {
		c.RenderItem(w)
		w.Write([]byte("\n"))
	}

	w.Write([]byte("END:VCALENDAR"))

	return nil
}
