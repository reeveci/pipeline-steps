FROM alpine/git

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/

# MATCH: glob to be matched against the tag
ENV MATCH=
# COMMIT: commit-ish to be used as base point
ENV COMMIT=HEAD
# SED_COMMAND: optional SED command for editing the selected tag name
ENV SED_COMMAND=
# RESULT_VAR: Name of a runtime variable for setting the step result to
ENV RESULT_VAR=GIT_TAG

ENTRYPOINT ["docker-entrypoint.sh"]
