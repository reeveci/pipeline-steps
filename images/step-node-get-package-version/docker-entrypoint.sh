#!/bin/sh
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

cd /reeve/src/

RESULT=$(jq -r '.version | strings' "${FILE:-package.json}")

if [ -z "$RESULT" ]; then
  echo 'The selected file does not contain a "version" field.'
  exit 1
fi

wget -O - -q "$REEVE_API/api/v1/var?key=${RESULT_VAR:-PACKAGE_VERSION}&value=$RESULT" >/dev/null
echo Set ${RESULT_VAR:-PACKAGE_VERSION}=$RESULT
