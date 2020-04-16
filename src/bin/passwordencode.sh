#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"

function run-with-protoc() {
  if type protoc >/dev/null 2>&1
  then
    (
      cd "${PROJECT}"
      "$@"
    )
  else
    docker run --rm -v "${PROJECT}:${PROJECT}" -w "${PROJECT}" jeroenvm/build-protoc "$@"
  fi
}

run-with-protoc target/bin/passwordencode "$@"
