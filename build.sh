#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p "${SCRIPTDIR}/build"

"${SCRIPTDIR}/wgo.sh" install -x -pkgdir /raad071cal/build/pkg github.com/mdirkse/raad071cal

cp "${SCRIPTDIR}/Dockerfile" "${SCRIPTDIR}/build"
cd "${SCRIPTDIR}/build" || exit

docker build --rm -t "mdirkse/raad071cal" .