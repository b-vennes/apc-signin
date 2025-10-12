# syntax=docker/dockerfile:1

FROM oven/bun

WORKDIR /app
COPY package.json .
COPY bun.lock .
RUN bun install
COPY . .
RUN bun run build

RUN adduser appuser
USER appuser

EXPOSE 3000

# What the container should run when it is started.
ENTRYPOINT [ "bun", "run", "start" ]
