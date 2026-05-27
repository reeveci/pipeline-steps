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
# CONTEXT: Context directory (relative to project root)
ENV CONTEXT=.

ENTRYPOINT ["docker-entrypoint.sh"]
