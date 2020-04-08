#!/bin/bash

BIN="$(cd "$(dirname "$0")" ; pwd)"
SRC="$(dirname "${BIN}")"
PROJECT="$(dirname "${SRC}")"
MODULE="${SRC}/pkg/trusted/trusted-generated.go"

source "${BIN}/verbose.sh"
source "${PROJECT}/src/etc/settings-local.sh"

mkdir -p "$(dirname "${MODULE}")"

echo 'package trusted

func Init() {
    TrustedKeys = map[string]string{}
    KeyManagers = map[string]string{}' > "${MODULE}"
(
  cd "${PROJECT}" || exit 1
  N=0
  for F in "${ROOT_PRIVATE_KEY}.pub" "${ADDITIONAL_TRUSTED_KEYS}"
  do
    if [[ -z "${F}" ]]
    then
      continue
    fi
    log ">>> Trusted key: [${F}]"
    KEY="$(cut -d ' ' -f2 "${F}")"
    NAME="$(cut -d ' ' -f3 "${F}")"
    if [[ -z "${KEY}" ]]
    then
      continue
    fi
    if [[ -z "${NAME}" ]]
    then
      N=$((${N} + 1))
      NAME="key-${N}"
    fi
    echo "    TrustedKeys[\"${NAME}\"] = \"${KEY}\""
    echo "    KeyManagers[\"${NAME}\"] = TrustedKeys[\"${NAME}\"]"
  done >> "${MODULE}"
)
echo '}' >> "${MODULE}"

sed -e 's/^/+/' "${MODULE}"
