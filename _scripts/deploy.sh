#!/usr/bin/env bash
#
# Build and push Docker images to Docker Hub and quay.io.
#

cd "$(dirname "$0")" || exit 1

export IMAGE_PREFIX=deisci
docker login -e="$DOCKER_EMAIL" -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
DEIS_REGISTRY='' make -C .. docker-push
docker login -e="$QUAY_EMAIL" -u="$QUAY_USERNAME" -p="$QUAY_PASSWORD" quay.io
DEIS_REGISTRY=quay.io/ make -C .. docker-build docker-push

# download deis CLI & deploy to deis
curl -sSL http://deis.io/deis-cli/install-v2.sh | bash
./deis login --username=$DEIS_USERNAME --password=$DEIS_PASSWORD ${DEIS_URL}
DEIS_BINARY_NAME=./_scripts/deis DEIS_APP_NAME=${DEIS_APP_NAME} make -C .. deploy-to-deis
