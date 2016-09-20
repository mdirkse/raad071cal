[![Build Status](https://travis-ci.org/mdirkse/raad071cal.svg?branch=master)](https://travis-ci.org/mdirkse/raad071cal)
[![Coverage](https://codecov.io/gh/mdirkse/raad071cal/branch/master/graph/badge.svg)](https://codecov.io/gh/mdirkse/raad071cal)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdirkse/raad071cal)](https://goreportcard.com/report/github.com/mdirkse/raad071cal)

# Raad 071 Cal
The Notubiz Politiek Portaal system for Leiden does not expose public municipal council events in iCal format. This small and simple web service remedies that omission by
scraping the [calendar](https://leiden.notubiz.nl/) and providing an URL that lists the past 6, and future 12, months of events in iCal format.

### Website
This code runs at http://raad071.mdirkse.nl/

### Docker image
The docker image with this service can be found here: https://hub.docker.com/r/mdirkse/raad071cal/
