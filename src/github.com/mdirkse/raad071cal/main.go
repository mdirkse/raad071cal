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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	listenAddress      = ":80"
	raad071CalendarURL = "http://leiden.notudoc.nl/cgi-bin/calendar.cgi"
	calendarHeader     = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//mdirkse/raad071cal//NONSGML v1.0//EN
URL:http://raad071.mdirkse.nl/kalender/alles.ics
NAME:#raad071 kalender
X-WR-CALNAME:#raad071 kalender
DESCRIPTION:De politieke agenda van de Leidse gemeenteraad
X-WR-CALDESC:De politieke agenda van de Leidse gemeenteraad
X-PUBLISHED-TTL:PT6H
`
)

var (
	calItems   []CalItem
	cestTz     *time.Location
	cronT      *cron.Cron
	cutoffDate time.Time
	itemsRegex *regexp.Regexp
	mutex      sync.RWMutex
)

func main() {
	initCalFetcherVars()
	log.Println("Starting raad071cal")

	// Configure periodic polling
	log.Printf("Polling source calendar [%s] every 6 hours.", raad071CalendarURL)
	cronT.AddFunc("1 1 */6 * * *", loadCalendarItems)
	cronT.Start()

	http.Handle("/kalender/alles.ics", loggingHandler(calHandler()))
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("html"))))

	log.Printf("Fully initialised and listening on [%s].", listenAddress)
	go loadCalendarItems() // do initial load

	http.ListenAndServe(listenAddress, nil)
}

func initCalFetcherVars() {
	cestTz, _ = time.LoadLocation("Europe/Amsterdam")
	cronT = cron.New()
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, cestTz)
	itemsRegex = regexp.MustCompile(`var vdate='(.+)'.split\(`)

	calItems = []CalItem{}

	initCalItemVars()
}

func loadCalendarItems() {
	pageBytes, err := fetchCalenderPage(raad071CalendarURL)
	if err != nil {
		log.Printf("ERROR - Unable to fetch the calendar page: %+v", err)
		return
	}

	newCalItems, err := parseCalendar(pageBytes)
	if err != nil {
		log.Printf("ERROR - Unable to parse the calendar: %+v", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	calItems = newCalItems
}

func fetchCalenderPage(calendarURL string) (*[]byte, error) {
	start := time.Now()

	resp, err := http.Get(calendarURL)
	if err != nil {
		return nil, fmt.Errorf("Could not fetch the calendar from [%s]: %+v", calendarURL, err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read the calender URL contents: %+v", err)
	}

	log.Printf("Fetched calendar in %0.2f seconds.", time.Since(start).Seconds())
	return &bytes, nil
}

func parseCalendar(pageBytes *[]byte) ([]CalItem, error) {
	log.Printf("Parsing the calendar.")
	start := time.Now()

	runStart := time.Now()
	calBytes := itemsRegex.FindSubmatch(*pageBytes)

	if len(calBytes) != 2 {
		return nil, fmt.Errorf("Could not find calendar items in text:\n%s", string(*pageBytes))
	}

	calText := string(calBytes[1])
	rawItems := strings.Split(calText, "@")
	items := make([]CalItem, 0, len(rawItems))

	for _, c := range rawItems {
		i, err := NewItem(c, runStart)

		if err != nil {
			log.Printf("ERROR - Unable to parse item [%s]: %+v", c, err)
			continue
		}

		if i.StartDateTime.After(cutoffDate) {
			items = append(items, i)
		}
	}

	log.Printf("Parsed %d calendar items in %0.3f seconds.", len(items), time.Since(start).Seconds())

	return items, nil
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] [%s] %s", r.RemoteAddr, r.Method, r.URL.Path)
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

func renderCalendar(items []CalItem, w io.Writer) error {
	start := time.Now()

	_, err := io.WriteString(w, calendarHeader)

	if err != nil {
		return fmt.Errorf("Could not write calendar!")
	}

	mutex.RLock()
	for _, c := range items {
		c.RenderItem(w)
		io.WriteString(w, "\n")
	}
	mutex.RUnlock()

	io.WriteString(w, "END:VCALENDAR")

	log.Printf("Rendered iCal calendar in %0.3f seconds.", time.Since(start).Seconds())

	return nil
}
