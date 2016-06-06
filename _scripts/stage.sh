#!/usr/bin/env bash
#
# Stage a build
#

export IMAGE_PREFIX=deisci
# publish app image to repositories
source publish.sh
# download deis CLI
source get-deis.sh
# deploy to production
./deis login --username=$DEIS_STAGING_USERNAME --password=$DEIS_STAGING_PASSWORD ${DEIS_STAGING_URL}
DEIS_BINARY_NAME=./_scripts/deis DEIS_APP_NAME=${DEIS_STAGING_APP_NAME} make -C .. deploy-to-deis
