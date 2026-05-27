FROM golang AS builder

ARG REEVE_CLI_VERSION

ENV CGO_ENABLED=0
RUN go install github.com/reeveci/reeve-cli/reeve@v${REEVE_CLI_VERSION}
RUN cp $(go env GOPATH)/bin/reeve /usr/local/bin/

FROM alpine

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/
COPY --chmod=755 --from=builder /usr/local/bin/reeve /usr/local/bin/

# URL: Reeve server URL
ENV URL=
# SECRET: Reeve CLI secret
ENV SECRET=

ENTRYPOINT ["docker-entrypoint.sh"]
