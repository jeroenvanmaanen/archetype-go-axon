#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"

source "${BIN}/verbose.sh"

if [ \! -f 'common.proto' ]
then
  error "Protocol buffer specification files for Axon Server not found in current directory"
fi

FRAGMENT="$1" ; shift
if [ -z "${FRAGMENT}" ]
then
  SUFFIX=''
else
  SUFFIX="/${FRAGMENT}"
fi

sed -E -i \
  -e 's/^option/old_option/' \
  -e "3i\\
option go_package = \"src/pkg/grpc/axonserver${SUFFIX}\";" \
  -e '/^old_option go_package =/d' \
  -e 's/^old_option/option/' \
  *.proto

if [ -z "${FRAGMENT}" ]
then
  protoc --go_out=plugins=grpc:/Users/jeroen/src/archetype-nix-go -I. *.proto
else
  protoc --go_out=/Users/jeroen/src/archetype-nix-go -I. "${FRAGMENT}.proto" "common.proto"
fi
