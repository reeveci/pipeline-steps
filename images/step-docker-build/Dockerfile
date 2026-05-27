ARG DOCKER_VERSION

FROM golang AS builder

ARG REEVE_TOOLS_VERSION

ENV CGO_ENABLED=0
RUN go install github.com/reeveci/reeve/reeve-tools@v${REEVE_TOOLS_VERSION}
RUN cp $(go env GOPATH)/bin/reeve-tools /usr/local/bin/

FROM docker:${DOCKER_VERSION}

RUN apk add --no-cache bash
COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/
COPY --chmod=755 --from=builder /usr/local/bin/reeve-tools /usr/local/bin/

# DOCKER_LOGIN_REGISTRIES: Space separated list of docker registries to log in to (user:password@registry)
ENV DOCKER_LOGIN_REGISTRIES=
# NAME: image name (required)
ENV NAME=
# TAG: image tag
ENV TAG=
# TAG_ALIASES: Space separated list of additional image tags to use - can contain other params as variables
ENV TAG_ALIASES=latest
# FILE: Dockerfile name (use - for STDIN)
ENV FILE=
# CONTEXT: Context directory (relative to project root)
ENV CONTEXT=.
# BUILD_ARGS: Space separated list of ARG=VALUE pairs - can contain other params as variables
ENV BUILD_ARGS=
# NETWORK: Networking mode for RUN instructions during build
ENV NETWORK=default
# USE_CACHE=true|false
ENV USE_CACHE=true
# PLATFORM: Platform to be used if server is multi-platform capable
ENV PLATFORM=
# PULL=always|missing
ENV PULL=missing
# PUSH=true|false
ENV PUSH=true
# TEST=true|fail|false
ENV TEST=true
# TEST_PULL=true|false
ENV TEST_PULL=false
# RESULT_VAR: Name of a runtime variable for setting the step result (failure|exists|success) to
ENV RESULT_VAR=
# BUILD_DRIVER: Build driver to be used
ENV BUILD_DRIVER=docker

ENTRYPOINT ["docker-entrypoint.sh"]
