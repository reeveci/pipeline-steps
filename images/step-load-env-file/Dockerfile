FROM golang AS builder

WORKDIR /app
COPY . .

ENV GOFLAGS="-buildvcs=false"
ENV CGO_ENABLED=0
RUN go build -o /usr/local/bin/reeve-step .

FROM alpine

COPY --chmod=755 --from=builder /usr/local/bin/reeve-step /usr/local/bin/

WORKDIR /reeve/src

# FILES: Space separated list of file patterns (see https://pkg.go.dev/github.com/bmatcuk/doublestar/v4#Match) to be included (shell syntax)
ENV FILES=**/*.env
# LOAD_ALL=true|false: Whether to load all variables from the env files
ENV LOAD_ALL=false
# ENV_<name>: Variables to be loaded from the file and their corresponding runtime variable names to be using in Reeve

ENTRYPOINT ["reeve-step"]
