ARG DOCKER_VERSION

FROM golang AS builder

ARG REEVE_TOOLS_VERSION

ENV CGO_ENABLED=0
RUN go install github.com/reeveci/reeve/reeve-tools@v${REEVE_TOOLS_VERSION}
RUN cp $(go env GOPATH)/bin/reeve-tools /usr/local/bin/

FROM docker:${DOCKER_VERSION}

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/
COPY --chmod=755 --from=builder /usr/local/bin/reeve-tools /usr/local/bin/

# DOCKER_LOGIN_REGISTRIES: Space separated list of docker registries to log in to (user:password@registry)
ENV DOCKER_LOGIN_REGISTRIES=
# MODE=v1|v2|swarm
ENV MODE=swarm
# NAME: compose project name (required for swarm mode)
ENV NAME=
# FILE: compose file name (use - for STDIN)
ENV FILE=
# CONTEXT: compose context directory (relative to project root)
ENV CONTEXT=.
# COMPOSE_BUILD=always|missing|never
ENV COMPOSE_BUILD=missing
# COMPOSE_RECREATE=always|missing|never
ENV COMPOSE_RECREATE=missing
# COMPOSE_PULL=always|missing|never
ENV COMPOSE_PULL=missing
# SWARM_PULL=always|changed|never
ENV SWARM_PULL=always
# STATE=present|absent
ENV STATE=present

ENTRYPOINT ["docker-entrypoint.sh"]
