#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p "${SCRIPTDIR}/build"

"${SCRIPTDIR}/wgo.sh" build -x

cp "${SCRIPTDIR}/src/raad071cal" \
   "${SCRIPTDIR}/Dockerfile" \
   "${SCRIPTDIR}/build"

cd "${SCRIPTDIR}/build" || exit

docker build --rm -t "mdirkse/raad071cal" .