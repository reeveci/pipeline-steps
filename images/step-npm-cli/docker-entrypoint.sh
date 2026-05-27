#!/bin/bash
set -e

if [ -z "$REEVE_API" ]; then
  echo This docker image is a Reeve CI pipeline step and is not intended to be used on its own.
  exit 1
fi

cd /reeve/src/${CONTEXT}

if [ -n "$NPM_LOGIN_REGISTRY" ]; then
  npm config set registry "${NPM_LOGIN_REGISTRY}"

  NPM_LOGIN_REGISTRY_NAME="$(node -e "console.log(/(?<=https?:\/\/)(.[^/]+)/.exec(process.env.NPM_LOGIN_REGISTRY)[0]);" 2>/dev/null)" || (echo "Parsing NPM login registry failed"; exit 1)

  if [ -n "$NPM_LOGIN_TOKEN" ]; then
    npm config set "//$NPM_LOGIN_REGISTRY_NAME/:_authToken" "$NPM_LOGIN_TOKEN"
  else
    if [ -z "$NPM_LOGIN_USER" ]; then
      echo Missing login user or token
      exit 1
    fi
    if [ -z "$NPM_LOGIN_PASSWORD" ]; then
      echo Missing login password
      exit 1
    fi

    echo Login attempt for $NPM_LOGIN_REGISTRY...
    NPM_LOGIN_TOKEN="$(node -e 'const u = process.env.NPM_LOGIN_USER; const p = process.env.NPM_LOGIN_PASSWORD; const response = await fetch(`${process.env.NPM_LOGIN_REGISTRY}/-/user/org.couchdb.user:${u}`, { method: "PUT", headers: { Accept: "application/json", "Content-Type": "application/json", Authorization: `Basic ${Buffer.from(`${u}:${p}`).toString("base64")}` }, body: JSON.stringify({ name: u, password: p }) }); if (!response.ok) { throw new Error(`Login failed (status ${response.status})`); }; console.log((await response.json()).token);')"
    npm config set "//$NPM_LOGIN_REGISTRY_NAME/:_authToken" "$NPM_LOGIN_TOKEN"
    echo Login successful
  fi
fi

unset NPM_LOGIN_REGISTRY NPM_LOGIN_REGISTRY_NAME NPM_LOGIN_TOKEN NPM_LOGIN_USER NPM_LOGIN_PASSWORD

exec docker-entrypoint-node.sh "$@"
