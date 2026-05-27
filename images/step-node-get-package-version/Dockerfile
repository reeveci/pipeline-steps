FROM alpine

RUN apk add jq
COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/

# FILE: name of the package.json file
ENV FILE=package.json
# RESULT_VAR: Name of the runtime variable for setting the step result to
ENV RESULT_VAR=PACKAGE_VERSION

ENTRYPOINT ["docker-entrypoint.sh"]
