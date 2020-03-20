#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"

source "${BIN}/verbose.sh"

cd "${SRC}/proto"

if [ \! -f 'example.proto' ]
then
  error "Protocol buffer specification files for Example not found in current directory"
fi

log "Generating JS stubs from $(pwd)"

OUT_DIR="${SRC}/present/src/grpc/example"
mkdir -p "${OUT_DIR}"

protoc --js_out="import_style=commonjs:${OUT_DIR}" --grpc-web_out="import_style=commonjs+dts,mode=grpcwebtext:${OUT_DIR}" -I. *.proto

# Add /* eslint-disable */
cd "${SRC}/present/src/grpc/example"
sed -E -i \
  -e '1s:^/\* eslint-disable \*/$:/*@@@ eslint-disable @@@*/:' \
  -e "1i\\
/* eslint-disable */" \
  -e '/^\/\*@@@ eslint-disable @@@\*\//d' \
  *.js
