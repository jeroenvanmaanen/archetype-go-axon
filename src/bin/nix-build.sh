#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"

(
  cd "${PROJECT}"
  "${BIN}/nix.sh" nix-build -E 'with import <nixpkgs> {}; pkgs.callPackage ./default.nix {}'
  mkdir -p target/bin
  "${BIN}/nix.sh" bash -c "cp result/bin/* target/bin/."
  ls -l target/bin
)
