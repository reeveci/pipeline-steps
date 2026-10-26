#!/bin/sh
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

mkdir -p /reeve/src
cd /reeve/src

git clone --filter=tree:0 "$GIT_REPOSITORY" .
git checkout ${GIT_COMMIT:-HEAD}
git remote | while read remote; do
  git remote rm $remote
done

exec "$@"
