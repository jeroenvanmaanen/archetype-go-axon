#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"
source "${PROJECT}/src/etc/settings-local.sh"

if [[ ".${ROOT_PRIVATE_KEY#/}" != ".${ROOT_PRIVATE_KEY}" ]]
then
  info "Absolute path: ${ROOT_PRIVATE_KEY}. Not generating anything."
  exit 0
fi

if type nix-env >/dev/null 2>&1
then
  :
else
  exec docker run --rm -v "${PROJECT}:${PROJECT}" -w "${BIN}" jeroenvm/build-protoc "$0" "$@"
  # Unreachable
  exit "$?"
fi

echo "HOME=[${HOME}]"
type nix-env

SECURE="${PROJECT}/data/secure"
mkdir -p "${SECURE}"
cd "${SECURE}"
pwd

if [[ -f "./id_rsa" ]]
then
  info "Private key already exists. Not generating a new one."
else
  nix-env -f '<.>' -iA openssl
  openssl genpkey -algorithm RSA -out id_rsa -pkeyopt rsa_keygen_bits:2048
fi

nix-env -f '<.>' -iA openssh
ssh-keygen -f id_rsa -y > id_rsa.pub

sed -i -e 's/$/ demokey/' id_rsa.pub
