FROM golang AS builder

WORKDIR /app
COPY . .

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/reeve-step .

FROM alpine

RUN apk add --no-cache envsubst
COPY --chmod=755 --from=builder /usr/local/bin/reeve-step /usr/local/bin/

WORKDIR /reeve/src

# FILES: Space separated list of file patterns (see https://pkg.go.dev/github.com/bmatcuk/doublestar/v4#Match) to be included (shell syntax)
ENV FILES=
# SUBSTITUTE_ALL=true|false: Whether to substitute all variables in the files, even if no value was provided for them
ENV SUBSTITUTE_ALL=false
# ENV_<name>: Variables to be substituted and their corresponding values

ENTRYPOINT ["reeve-step"]
