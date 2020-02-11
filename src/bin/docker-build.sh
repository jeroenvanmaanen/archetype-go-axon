#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"

docker build -t 'jeroenvm/archetype-nix-go' "${SRC}/docker"
