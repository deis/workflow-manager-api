#!/usr/bin/env bash
#
# Deploy to production
#

export IMAGE_PREFIX=deis
# publish app image to repositories
source publish.sh
# download deis CLI
source get-deis.sh
# deploy to production
./deis login --username=$DEIS_PROD_USERNAME --password=$DEIS_PROD_PASSWORD ${DEIS_PROD_URL}
DEIS_BINARY_NAME=./_scripts/deis DEIS_APP_NAME=${DEIS_PROD_APP_NAME} make -C .. deploy-to-deis
