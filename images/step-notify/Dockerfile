FROM golang AS builder

WORKDIR /app
COPY . .

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/reeve-step .

FROM alpine

COPY --chmod=755 --from=builder /usr/local/bin/reeve-step /usr/local/bin/

WORKDIR /reeve/src

# URLS: Space separated list of service URLs (see https://containrrr.dev/shoutrrr/services/overview)
ENV URLS=
# MESSAGE: Message to be sent - can contain other params as variables
ENV MESSAGE=
# PARAM_<name>: Additional service params and their corresponding values

ENTRYPOINT ["reeve-step"]
