[![Build Status](https://travis-ci.org/mdirkse/raad071cal.svg?branch=master)](https://travis-ci.org/mdirkse/raad071cal)
[![Coverage](https://codecov.io/gh/mdirkse/raad071cal/branch/master/graph/badge.svg)](https://codecov.io/gh/mdirkse/raad071cal)

# Raad 071 Cal
For reasons I can't quite fathom, the [gemeenteraad.leiden.nl calendar](http://leiden.notudoc.nl/cgi-bin/calendar.cgi) has the entire municipal council meeting calender encoded into a Javascript variable on the homepage.
This simple web service reads that variable, parses the info and serves up an iCal version of the information (so you can add it to Google Calendar for instance).

### Website
This code runs at http://raad071cal.mdirkse.nl/

### Docker image
The docker image with this service can be found here: https://hub.docker.com/r/mdirkse/raad071cal/
