FROM golang:1.15-alpine3.12

LABEL maintainer="https://github.com/nekochans"

WORKDIR /go/app

COPY . .

RUN set -eux && \
  apk update && \
  apk add --no-cache git curl make

ENV CGO_ENABLED 0
