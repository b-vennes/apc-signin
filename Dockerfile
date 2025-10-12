# syntax=docker/dockerfile:1

FROM oven/bun

WORKDIR /app
COPY package.json .
COPY bun.lock .
RUN bun install
COPY . .
RUN bun run build

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

EXPOSE 3000

# What the container should run when it is started.
ENTRYPOINT [ "bun", "run", "start" ]
