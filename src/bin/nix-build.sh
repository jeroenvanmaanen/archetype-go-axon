#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"

(
  docker run --rm -v "${PROJECT}:${PROJECT}" -w "${BIN}" jeroenvm/build-protoc ./generate-proto-go-package.sh -v
  docker run --rm -v "${PROJECT}:${PROJECT}" -w "${BIN}" jeroenvm/build-protoc ./generate-proto-js-package.sh -v
  cd "${PROJECT}"
  "${BIN}/nix.sh" nix-build -E 'with import <nixpkgs> {}; pkgs.callPackage ./default.nix {}'
  mkdir -p target/bin
  "${BIN}/nix.sh" bash -c "cp -f result/bin/* target/bin/."
  ls -l target/bin
)
