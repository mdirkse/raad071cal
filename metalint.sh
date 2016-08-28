#!/bin/bash
GOPATH="$(pwd)/vendor:$(pwd)"
gometalinter --disable=errcheck \
             --disable=gotype \
             src/github.com/mdirkse/raad071cal
