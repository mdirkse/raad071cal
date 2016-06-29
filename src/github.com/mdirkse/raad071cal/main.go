package main

import (
	"fmt"
	"github.com/robfig/cron"
	log "gopkg.in/inconshreveable/log15.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	listenAddress = ":7070"
)

var (
	logger    log.Logger
	cronT     *cron.Cron
	iCalItems string
)

func main() {
	initCalFetcher()

	logger.Info("Starting raad071cal")

	cronT.Start()

	go loadCalendarItems()

	//var iCal bytes.Buffer
	//for _, i := range *alItems {

	//}

	//println("BEGIN:VCALENDAR")
	//println("VERSION:2.0")
	//print("PRODID:-//hacksw/handcal//NONSGML v1.0//EN")
	//println(iCal.String())
	//println("END:VCALENDAR")

	//http.Handle("/raad071metrics", prometheus.Handler())

	//logger.Info(fmt.Sprintf("Fully initialised and listening on [%s].", listenAddress))
	//http.ListenAndServe(listenAddress, nil)
}

func initCalFetcher() {
	logger = log.New()
	CEST, _ = time.LoadLocation("Europe/Amsterdam")
	cronT = cron.New()
	cutoffDate = time.Date(2015, 1, 1, 0, 0, 0, 0, CEST)
	itemsRegex = regexp.MustCompile(`var vdate='(.+)'.split\(`)

	var err error
	itemTemplate, err = template.New("item").Parse(itemTemplateSrc)

	if err != nil {
		panic("Could not parse item template!")
	}
}

func loadCalendarItems() {
	//pageBytes := fetchCalenderPage()
	//calItems := parseCalendar(pageBytes)
}

func fetchCalenderPage() *[]byte {
	resp, _ := http.Get("http://leiden.notudoc.nl/cgi-bin/calendar.cgi")
	bytes, _ := ioutil.ReadAll(resp.Body)

	return &bytes
}

func parseCalendar(pageBytes *[]byte) *[]*CalItem {
	logger.Info("Parsing the calendar.")

	calBytes := itemsRegex.FindSubmatch(*pageBytes)

	calText := string(calBytes[1])

	rawItems := strings.Split(calText, "@")

	items := make([]*CalItem, 0, len(rawItems))

	runStart := time.Now()

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

	return &items
}
