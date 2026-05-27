#!/bin/bash
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

cd /reeve/src/${CONTEXT}

if [ -n "$DOCKER_LOGIN_REGISTRIES" ]; then
  reeve-tools login-docker $DOCKER_LOGIN_REGISTRIES
elif [ -n "$DOCKER_LOGIN_REGISTRY" ]; then
  echo WARNING: The DOCKER_LOGIN_REGISTRY, DOCKER_LOGIN_USER and DOCKER_LOGIN_PASSWORD params are deprecated and will stop working in a future version! Use DOCKER_LOGIN_REGISTRIES instead.

  if [ -z "$DOCKER_LOGIN_USER" ]; then
    echo Missing login user
    exit 1
  fi
  if [ -z "$DOCKER_LOGIN_PASSWORD" ]; then
    echo Missing login password
    exit 1
  fi

  echo Login attempt for $DOCKER_LOGIN_REGISTRY...
  printf "%s\n" "$DOCKER_LOGIN_PASSWORD" | docker login -u "$DOCKER_LOGIN_USER" --password-stdin $DOCKER_LOGIN_REGISTRY
fi

exec "$@"
