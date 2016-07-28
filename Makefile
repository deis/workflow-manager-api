SHORT_NAME ?= workflow-manager-api

include versioning.mk
include includes.mk

# Enable vendor/ directory support.
export GO15VENDOREXPERIMENT=1

# dockerized development environment variables
REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.14.0
SWAGGER_IMAGE := quay.io/goswagger/swagger:0.5.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
SWAGGER_CMD := docker run --rm -e GOPATH=/go -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${SWAGGER_IMAGE}

# SemVer with build information is defined in the SemVer 2 spec, but Docker
# doesn't allow +, so we use -.
BINARY_DEST_DIR := rootfs/usr/bin
# Common flags passed into Go's linker.
LDFLAGS := "-s -X main.version=${VERSION}"
# Docker Root FS
BINDIR := ./rootfs

DEIS_REGISTRY ?= ${DEV_REGISTRY}/

DEIS_BINARY_NAME ?= ./deis


all:
	@echo "Use a Makefile to control top-level building of the project."

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

# This illustrates a two-stage Docker build. docker-compile runs inside of
# the Docker environment. Other alternatives are cross-compiling, doing
# the build as a `docker build`.
build:
	${DEV_ENV_PREFIX} -e CGO_ENABLED=0 ${DEV_ENV_IMAGE} go build -a -installsuffix cgo -ldflags ${LDFLAGS} -o ${BINARY_DEST_DIR}/${SHORT_NAME} *.go || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/${SHORT_NAME})
	${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE} upx -9 ${BINARY_DEST_DIR}/${SHORT_NAME} || exit 1

static-html:
	${DEV_ENV_CMD} sh -c 'cp -R static_html ${BINDIR}'

swagger-serverstub:
	${SWAGGER_CMD} generate server -A WorkflowManager -t pkg/swagger -f https://raw.githubusercontent.com/jackfrancis/workflow-manager/oauth/api/swagger-spec/swagger.yml
	mv pkg/swagger/cmd/workflow-manager-server/main.go .

test:
	${DEV_ENV_CMD} sh -c 'go test -tags testonly $$(glide nv)'

test-native:
	go test -tags=testonly $$(glide nv)

docker-build: build static-html
	docker build --rm -t ${IMAGE} rootfs
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

deploy-to-deis:
	${DEIS_BINARY_NAME} pull ${IMAGE} -a ${DEIS_APP_NAME}



.PHONY: all build docker-compile kube-up kube-down deploy
