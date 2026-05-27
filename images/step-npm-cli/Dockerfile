FROM node:lts-alpine

RUN apk add --no-cache bash git openssh-client
RUN mv /usr/local/bin/docker-entrypoint.sh /usr/local/bin/docker-entrypoint-node.sh
COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/

# NPM_LOGIN_REGISTRY: NPM registry to log in to
ENV NPM_LOGIN_REGISTRY=
# NPM_LOGIN_TOKEN: NPM token to log in with (instead of NPM_LOGIN_USER / NPM_LOGIN_PASSWORD)
ENV NPM_LOGIN_TOKEN=
# NPM_LOGIN_USER: NPM user to log in with (instead of NPM_LOGIN_TOKEN)
ENV NPM_LOGIN_USER=
# NPM_LOGIN_PASSWORD: NPM password to log in with (instead of NPM_LOGIN_TOKEN)
ENV NPM_LOGIN_PASSWORD=
# CONTEXT: Context directory (relative to project root)
ENV CONTEXT=.

ENTRYPOINT ["docker-entrypoint.sh"]
