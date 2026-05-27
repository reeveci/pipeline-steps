#!/bin/sh
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

if [ "$MODE" = "v1" ]; then
  # ===== compose v1 =====

  COMMAND="docker-compose $([[ -n "$FILE" ]] && printf "%s" "-f $FILE" ||:) $([[ -n "$NAME" ]] && printf "%s" "-p $NAME" ||:)"

  if [ "$STATE" = "present" ]; then
    echo Deploying compose $NAME with docker compose v1...
    $COMMAND up -d --remove-orphans $([[ "$COMPOSE_BUILD" = "always" ]] && printf "%s" "--build" ||:) $([[ "$COMPOSE_BUILD" = "never" ]] && printf "%s" "--no-build" ||:) $([[ "$COMPOSE_RECREATE" = "always" ]] && printf "%s" "--force-recreate" ||:) $([[ "$COMPOSE_RECREATE" = "never" ]] && printf "%s" "--no-recreate" ||:) --pull $COMPOSE_PULL
  elif [ "$STATE" = "absent" ]; then
    echo Removing compose $NAME with docker compose v1...
    $COMMAND down
  else
    echo Invalid state $STATE
    exit 1
  fi
elif [ "$MODE" = "v2" ]; then
  # ===== compose v2 =====

  COMMAND="docker compose $([[ -n "$FILE" ]] && printf "%s" "-f $FILE" ||:) $([[ -n "$NAME" ]] && printf "%s" "-p $NAME" ||:)"

  if [ "$STATE" = "present" ]; then
    echo Deploying compose $NAME with docker compose v2...
    $COMMAND up -d --remove-orphans $([[ "$COMPOSE_BUILD" = "always" ]] && printf "%s" "--build" ||:) $([[ "$COMPOSE_BUILD" = "never" ]] && printf "%s" "--no-build" ||:) $([[ "$COMPOSE_RECREATE" = "always" ]] && printf "%s" "--force-recreate" ||:) $([[ "$COMPOSE_RECREATE" = "never" ]] && printf "%s" "--no-recreate" ||:) --pull $COMPOSE_PULL
  elif [ "$STATE" = "absent" ]; then
    echo Removing compose $NAME with docker compose v2...
    $COMMAND down
  else
    echo Invalid state $STATE
    exit 1
  fi
elif [ "$MODE" = "swarm" ]; then
  # ===== swarm mode =====

  if [ -z "$NAME" ]; then
    echo Missing name
    exit 1
  fi

  if [ "$STATE" = "present" ]; then
    echo Deploying stack $NAME...
    docker stack deploy --detach=false --prune $([[ -n "$FILE" ]] && printf "%s" "-c $FILE" ||:) --resolve-image $SWARM_PULL $([[ -n "$DOCKER_LOGIN_REGISTRIES$DOCKER_LOGIN_REGISTRY" ]] && printf "%s" "--with-registry-auth") $NAME
  elif [ "$STATE" = "absent" ]; then
    echo Removing stack $NAME...
    docker stack rm $NAME
  else
    echo Invalid state $STATE
    exit 1
  fi
else
  echo Invalid mode $MODE
  exit 1
fi
