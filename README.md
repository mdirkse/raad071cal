[![Build Status](https://travis-ci.org/mdirkse/raad071cal.svg?branch=master)](https://travis-ci.org/mdirkse/raad071cal)

# Raad 071 Cal

For reasons I can't quite fathom, [gemeenteraad.leiden.nl](http://gemeenteraad.leiden.nl) has the entire municipal council meeting calender encoded into a Javascript variable on the homepage.
This simple web service reads that variable, parses the info and serves up an iCal version of the information (so you can add it to Google Calendar for instance).