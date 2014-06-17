#!/bin/bash

set -e

if [ -z $BINTRAY_API_KEY ]; then
    echo "Please set your bintray API key in the BINTRAY_API_KEY env var."
    exit 1
fi

DIR=$(cd $(dirname ${0})/.. && pwd)
cd ${DIR}

VERSION=$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/')

for ARCHIVE in ./pkg/dist/*; do
    ARCHIVE_NAME=$(basename ${ARCHIVE})

    echo Uploading: ${ARCHIVE_NAME}
    curl \
        -T ${ARCHIVE} \
        -utcnksm:${BINTRAY_API_KEY} \
        "https://api.bintray.com/content/tcnksm/dmux/dmux/${VERSION}/${ARCHIVE_NAME}"
done
