FROM alpine/git

COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/
COPY --chmod=755 credential-helper.sh /usr/local/bin/

RUN git config --global credential.helper /usr/local/bin/credential-helper.sh
RUN git config --global advice.detachedHead false

# GIT_USER: optional repository credentials
ENV GIT_USER=
# GIT_PASSWORD: optional repository credentials
ENV GIT_PASSWORD=
# GIT_REPOSITORY: repository URL
ENV GIT_REPOSITORY=
# GIT_COMMIT: commit to checkout
ENV GIT_COMMIT=

ENTRYPOINT ["docker-entrypoint.sh"]
