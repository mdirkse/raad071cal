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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const generatedMonths = 18

var httpGet = http.Get

type calendarMonth struct {
	Meetings []CalItem `json:"meetings"`
}

type yearMonth struct {
	year  int
	month int
}

type fetchResult struct {
	items []CalItem
	errs  []error
}

func fetchCalendarItems(fetchStart time.Time) ([]CalItem, error) {
	// Concurrency stuff
	var wg sync.WaitGroup
	cc := make(chan []CalItem)
	ec := make(chan error)
	result := make(chan fetchResult)
	cancel := make(chan bool)
	defer close(cc)
	defer close(ec)
	defer close(cancel)
	defer close(result)

	// Get all the months that we want to fetch
	yms := generateMonthYearRange(fetchStart)

	go readChannels(cc, ec, result, cancel)

	for _, ym := range yms {
		wg.Add(1)
		go func(ym yearMonth) {
			defer wg.Done()

			json, err := fetchCalendarMonthJSON(ym)
			if err != nil {
				ec <- err
				return
			}

			items, err := getCalendarItemsFromJSON(json, fetchStart)
			if err != nil {
				ec <- err
				return
			}

			cc <- items

		}(ym)
	}

	// Wait for all fetches to finish
	wg.Wait()

	// Cancel the gathering goroutine
	cancel <- true

	// Get the result
	fr := <-result

	if len(fr.errs) == 0 {
		return fr.items, nil
	}

	var errString bytes.Buffer
	errString.WriteString("Could not update items because the following error(s) occurred:")
	for _, err := range fr.errs {
		errString.WriteString("\n" + err.Error())
	}

	return nil, errors.New(errString.String())

}

func generateMonthYearRange(locus time.Time) []yearMonth {
	var generated []yearMonth

	// Always start the range 5 months before the locus
	for i := -5; i < generatedMonths-5; i++ {
		d := locus.AddDate(0, i, 0)
		generated = append(generated, yearMonth{year: d.Year(), month: int(d.Month())})
	}

	return generated
}

func readChannels(cc <-chan []CalItem, ec <-chan error, result chan fetchResult, cancel <-chan bool) {
	var cancelled = false
	var items []CalItem
	var errs []error

	for !cancelled {
		select {
		case fetchedItem := <-cc:
			items = append(items, fetchedItem...)
		case fetchErr := <-ec:
			errs = append(errs, fetchErr)
		case <-cancel:
			cancelled = true
		}

	}

	result <- fetchResult{
		items: items,
		errs:  errs,
	}
}

func fetchCalendarMonthJSON(ym yearMonth) (string, error) {
	calendarURL := fmt.Sprintf("https://leiden.notubiz.nl/api/calendar/callback_function?year=%d&month=%d&callback=raad071cal", ym.year, ym.month)
	resp, err := httpGet(calendarURL)
	if err != nil {
		return "", fmt.Errorf("Could not fetch the calendar from [%s]: %+v", calendarURL, err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read the calender URL contents: %+v", err)
	}

	text := strings.TrimSuffix(strings.TrimPrefix(string(bytes[:]), "callback_function("), ")")

	return text, nil
}

func getCalendarItemsFromJSON(cpJSON string, fetchStart time.Time) ([]CalItem, error) {
	var cp calendarMonth

	err := json.Unmarshal([]byte(cpJSON), &cp)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse JSON calendar items! Error: %+v", err)
	}

	items := make([]CalItem, 0, len(cp.Meetings))

	for _, i := range cp.Meetings {
		if strings.ToLower(i.Description) == "fractievergadering" || i.Canceled {
			continue
		}

		ei, err := EnrichItem(i, fetchStart)
		if err != nil {
			log.Printf("ERROR - Unable to enrich item [%+v]: %+v", ei, err)
			continue
		}

		items = append(items, ei)
	}

	return items, nil
}
