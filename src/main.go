package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	log "gopkg.in/inconshreveable/log15.v2"
	"net/http"
)

const listenAddress = ":7070"

var (
	logger log.Logger
)

func main() {
	logger = log.New()
	logger.Info("Starting raad071cal")

	cron := cron.New()
	cron.Start()

	// Serve the metrics endpoint
	http.Handle("/raad071metrics", prometheus.Handler())
	logger.Info(fmt.Sprintf("Fully initialised and listening on [%s].", listenAddress))

	http.ListenAndServe(listenAddress, nil)
}

func fetchCalender() {

}