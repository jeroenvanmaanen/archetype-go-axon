#!/bin/bash

# This script generates a random UUID for use as a secret key

if type uuidgen >/dev/null 2>&1
then
  uuidgen
else
  cat /proc/sys/kernel/random/uuid
fi | tr -d '-' | tr 'A-Z' 'a-z'