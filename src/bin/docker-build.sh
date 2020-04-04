#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
DOCKER="${SRC}/docker"

docker build -t 'jeroenvm/archetype-go-axon' "${DOCKER}"
docker build -t 'jeroenvm/build-protoc' -f "${DOCKER}/Dockerfile-protoc" "${DOCKER}"