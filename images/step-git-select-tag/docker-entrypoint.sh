#!/bin/sh
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

cd /reeve/src/

RESULT=$(git describe $([[ -n "$MATCH" ]] && printf "%s" "--match "$MATCH"" ||:) --tags --abbrev=0 ${COMMIT:-HEAD})

if [ -n "$SED_COMMAND" ]; then
  RESULT=$(printf "%s" "$RESULT" | sed -r "$SED_COMMAND")
  test -n "$RESULT"
fi

wget -O - -q "$REEVE_API/api/v1/var?key=${RESULT_VAR:-GIT_TAG}&value=$RESULT" >/dev/null
echo Set ${RESULT_VAR:-GIT_TAG}=$RESULT
