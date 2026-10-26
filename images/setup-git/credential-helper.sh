#!/bin/sh
set -e

test "$1" = get

[ -n "${GIT_USER:-$GIT_PASSWORD}" ] && echo username=${GIT_USER:-git}
[ -n "$GIT_PASSWORD" ] && echo password=$GIT_PASSWORD
