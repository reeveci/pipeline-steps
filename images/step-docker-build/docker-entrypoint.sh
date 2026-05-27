#!/bin/bash
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

cd /reeve/src/${CONTEXT}

if [ -n "$RESULT_VAR" ]; then
  wget -O - -q "$REEVE_API/api/v1/var?key=$RESULT_VAR&value=failure" >/dev/null
fi

if [ -z "$NAME" ]; then
  echo Missing name
  exit 1
fi

if [ -n "$DOCKER_LOGIN_REGISTRIES" ]; then
  reeve-tools login-docker $DOCKER_LOGIN_REGISTRIES
elif [ -n "$DOCKER_LOGIN_USER" ]; then
  echo WARNING: The DOCKER_LOGIN_REGISTRY, DOCKER_LOGIN_USER and DOCKER_LOGIN_PASSWORD params are deprecated and will stop working in a future version! Use DOCKER_LOGIN_REGISTRIES instead.

  if [ -z "$DOCKER_LOGIN_PASSWORD" ]; then
    echo Missing login password
    exit 1
  fi

  echo Login attempt for ${DOCKER_LOGIN_REGISTRY:-$NAME}...
  printf "%s\n" "$DOCKER_LOGIN_PASSWORD" | docker login -u "$DOCKER_LOGIN_USER" --password-stdin ${DOCKER_LOGIN_REGISTRY:-$NAME}
fi

FULL_NAME="$NAME:${TAG:-latest}"

if [ "$TEST" = "true" ] || [ "$TEST" = "fail" ]; then
  if ! [ "$TEST_PULL" = "true" ]; then
    echo Testing manifest for image $FULL_NAME...
    if DOCKER_CLI_EXPERIMENTAL=enabled docker manifest inspect $FULL_NAME >/dev/null; then
      echo Image already exists - done
      if [ -n "$RESULT_VAR" ]; then
        wget -O - -q "$REEVE_API/api/v1/var?key=$RESULT_VAR&value=exists" >/dev/null
      fi
      ! [ "$TEST" = "fail" ]; exit
    fi
  else
    if [ "$PUSH" = "true" ]; then
      echo Removing local image $FULL_NAME to prepare testing...
      docker rmi $FULL_NAME >/dev/null 2>&1 ||:
    fi
    echo Pulling image $FULL_NAME for testing...
    if docker pull $FULL_NAME >/dev/null 2>&1; then
      echo Image already exists - done
      if [ -n "$RESULT_VAR" ]; then
        wget -O - -q "$REEVE_API/api/v1/var?key=$RESULT_VAR&value=exists" >/dev/null
      fi
      ! [ "$TEST" = "fail" ]; exit
    fi
  fi
  echo Image does not exist - continuing...
fi

if [ "${BUILD_DRIVER}" != "docker" ]; then
  docker buildx create --name reeve --driver "${BUILD_DRIVER}" --bootstrap --use
fi

declare -a args
IFS=$'\n'
for arg in $(printf "%s" "$BUILD_ARGS" | xargs -n1); do
  args+=("--build-arg" "$(eval printf \"%s\" \"$arg\")")
done
unset IFS

declare -a tags
declare -A tagNames
tags+=("-t" "$FULL_NAME")
tagNames["$FULL_NAME"]=1
if [ -n "$TAG_ALIASES" ]; then
  IFS=$'\n'
  for tag in $(printf "%s" "$TAG_ALIASES" | xargs -n1); do
    nextTag="$NAME:$(eval printf \"%s\" \"$tag\")"
    if [ -z "${tagNames[$nextTag]}" ]; then
      tags+=("-t" "$nextTag")
      tagNames["$nextTag"]=1
    fi
  done
  unset IFS
fi

COMMAND="docker buildx build $([[ -n "$NETWORK" ]] && printf "%s" "--network $NETWORK" ||:) $([[ "$USE_CACHE" = "false" ]] && printf "%s" "--no-cache" ||:) $([[ -n "$PLATFORM" ]] && printf "%s" "--platform $PLATFORM" ||:) $([[ "$PULL" = "always" ]] && printf "%s" "--pull" ||:) $([[ "$PUSH" = "true" ]] && printf "%s" "--push" ||:)"

echo Building image $FULL_NAME...
if [ "$FILE" = "-" ]; then
  $COMMAND "${tags[@]}" "${args[@]}" -
else
  $COMMAND "${tags[@]}" "${args[@]}" $([[ -n "$FILE" ]] && printf "%s" "-f $FILE" ||:) .
fi

if [ -n "$RESULT_VAR" ]; then
  wget -O - -q "$REEVE_API/api/v1/var?key=$RESULT_VAR&value=success" >/dev/null
fi
