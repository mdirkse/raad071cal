#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

projectnamedir="/raad071cal"

docker run --rm \
           -u "$(id -u):$(id -g)" \
           -e "CGO_ENABLED=0" \
           -e "GOOS=linux" \
           -e "GOBIN=${projectnamedir}/build" \
           -e "GOPATH=${projectnamedir}:${projectnamedir}/vendor" \
           -e "PKGDIR=${projectnamedir}/build/pkg" \
           -v "${SCRIPTDIR}:${projectnamedir}" \
           -w "${projectnamedir}" \
           --net=none \
           --log-driver=none \
           golang:1.13.6-alpine3.11 \
           go "$@"
