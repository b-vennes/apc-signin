# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o "./apc-signin-api"

RUN adduser appuser
RUN chown -R appuser /app

USER appuser

EXPOSE 3339

ENTRYPOINT ["./apc-signin-api"]
