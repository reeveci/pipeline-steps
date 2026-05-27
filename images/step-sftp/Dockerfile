FROM alpine

RUN apk add openssh sshpass
COPY --chmod=755 docker-entrypoint.sh /usr/local/bin/

# SSH_LOGIN_USER: SSH user
ENV SSH_LOGIN_USER=
# SSH_LOGIN_PASSWORD: SSH password
ENV SSH_LOGIN_PASSWORD=
# SSH_LOGIN_KEY: SSH private key file
ENV SSH_LOGIN_KEY=
# HOST: Host to connect to
ENV HOST=
# PORT: Custom port to connect to
ENV PORT=
# SSH_OPTIONS: Additional SSH options ("SSH_OPTIONS=-o Option1=value -o Option2=value")
ENV SSH_OPTIONS=
# SCRIPT: Script file to be executed (or - for stdin)
ENV SCRIPT=-

ENTRYPOINT ["docker-entrypoint.sh"]
