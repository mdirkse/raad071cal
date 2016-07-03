#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Build the binary
"${SCRIPTDIR}/wgo.sh" install -x -pkgdir /raad071cal/build/pkg github.com/mdirkse/raad071cal

# Gather other resources for the docker images
mkdir -p "${SCRIPTDIR}/build"
cp "${SCRIPTDIR}/Dockerfile" "${SCRIPTDIR}/build"
cp -R "${SCRIPTDIR}/html" "${SCRIPTDIR}/build"
tar cfz "${SCRIPTDIR}/build/zoneinfo.tar.gz" /usr/share/zoneinfo/Europe/Amsterdam

# Build!
cd "${SCRIPTDIR}/build" || exit
docker build --rm -t "mdirkse/raad071cal" .