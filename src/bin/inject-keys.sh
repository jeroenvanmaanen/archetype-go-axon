#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"

source "${BIN}/verbose.sh"
source "${SRC}/etc/settings-local.sh"

if [[ -z "${SIGN_PRIVATE_KEY}" ]]
then
  SIGN_PRIVATE_KEY="${ROOT_PRIVATE_KEY}"
fi

(
  cd "${PROJECT}" || exit 1

  ROOT_PUBLIC_KEY="$(cat "${ROOT_PRIVATE_KEY}.pub")"
  ROOT_KEY_NAME="$(cat "${ROOT_PRIVATE_KEY}.pub" | cut -d ' ' -f 3)"
  SIGN_KEY_NAME="$(cat "${SIGN_PRIVATE_KEY}.pub" | cut -d ' ' -f 3)"

  source "${BIN}/verbose.sh"
  source "${SRC}/etc/settings-local.sh"
  (
    echo ">>> Manager: ${ROOT_KEY_NAME}"
    cat "${ROOT_PRIVATE_KEY}"
    if [[ ".${SIGN_PRIVATE_KEY}" != ".${ROOT_PRIVATE_KEY}" ]]
    then
      echo '>>> Trusted:'
      cat "${SIGN_PRIVATE_KEY}.pub"
    fi
    echo ">>> Identity Provider: ${SIGN_KEY_NAME}"
    cat "${SIGN_PRIVATE_KEY}"
    echo '>>> End'
  ) | docker exec -i example_example-command-api_1 target/bin/keymanager
)