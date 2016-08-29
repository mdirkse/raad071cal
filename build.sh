#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Build the binary
"${SCRIPTDIR}/wgo.sh" install -x -pkgdir /raad071cal/build/pkg github.com/mdirkse/raad071cal

# Gather other resources for the docker images
mkdir -p "${SCRIPTDIR}/build"
cp "${SCRIPTDIR}/Dockerfile" "${SCRIPTDIR}/build"
cp -R "${SCRIPTDIR}/html" "${SCRIPTDIR}/build"

if [ ! -f "${SCRIPTDIR}/build/zoneinfo.tar.gz" ]; then
  tar cfz "${SCRIPTDIR}/build/zoneinfo.tar.gz" /usr/share/zoneinfo/Europe/Amsterdam
fi

if [ ! -f "${SCRIPTDIR}/build/certs.tar.gz" ]; then
  tar cfzh "${SCRIPTDIR}/build/certs.tar.gz" /etc/ssl/certs
fi

# Build!
cd "${SCRIPTDIR}/build" || exit
docker build --rm -t "mdirkse/raad071cal" .