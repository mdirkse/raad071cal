#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

gopath="/vendorgopath"

docker run --rm \
           -u "$(id -u):$(id -g)" \
           -e "CGO_ENABLED=0" \
           -e "GOOS=linux" \
           -e "GOPATH=${gopath}" \
           -v "${SCRIPTDIR}/vendor:/vendorgopath" \
           -v "${SCRIPTDIR}/src:/raad071cal" \
           -w "/raad071cal" \
           --net=none \
           golang:1.6.0-alpine \
           go "$@"
