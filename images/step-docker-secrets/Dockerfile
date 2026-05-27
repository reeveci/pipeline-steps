ARG DOCKER_VERSION

FROM golang AS builder

WORKDIR /app
COPY . .

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/reeve-step .

FROM docker:${DOCKER_VERSION}

RUN apk add --no-cache jq
COPY --chmod=755 --from=builder /usr/local/bin/reeve-step /usr/local/bin/

WORKDIR /reeve/src

# VOLUME: name of the volume
ENV VOLUME=
# TARGET_UID=1000
ENV TARGET_UID=1000
# TARGET_GID=1000
ENV TARGET_GID=1000
# FILE_MODE: mode to apply to the secrets
ENV FILE_MODE=0440
# REVISION_VAR: name of a runtime variable for setting the volume's revision to - this value can be applied to the corresponding containers as an environment variable in order to automatically update them when a secret changes
ENV REVISION_VAR=SECRET_REV
# SECRET_<name>: Names and values of the secrets to be written to the volume
# SECRETFILE_<name>: Names and paths to the secrets to be written to the volume

ENTRYPOINT ["reeve-step"]
