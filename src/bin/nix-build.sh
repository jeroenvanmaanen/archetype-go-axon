#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"

function run-with-protoc() {
  if type protoc >/dev/null 2>&1
  then
    (
      cd "${BIN}"
      "$@"
    )
  else
    docker run --rm -v "${PROJECT}:${PROJECT}" -w "${BIN}" jeroenvm/build-protoc
  fi
}

(
  run-with-protoc ./generate-proto-go-package.sh -v
  run-with-protoc ./generate-proto-js-package.sh -v
  cd "${PROJECT}"
  "${BIN}/nix.sh" nix-build -E 'with import <nixpkgs> {}; pkgs.callPackage ./default.nix {}'
  mkdir -p target/bin
  "${BIN}/nix.sh" bash -c "cp -f result/bin/* target/bin/."
  ls -l target/bin
)
