# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

ARG BUILD_VERSION=0.1.2
ARG BUILD_ARCH=amd64
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd

RUN case "${BUILD_ARCH}" in \
        aarch64) export GOARCH=arm64 ;; \
        amd64) export GOARCH=amd64 ;; \
        *) echo "unsupported BUILD_ARCH=${BUILD_ARCH}" >&2; exit 1 ;; \
    esac \
    && CGO_ENABLED=0 GOOS=linux go build \
        -trimpath \
        -ldflags="-s -w -X main.version=${BUILD_VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
        -o /out/homesignal-agent ./cmd/agent

FROM alpine:3.22

ARG BUILD_VERSION=0.1.2
ARG BUILD_ARCH=amd64

LABEL \
    io.hass.version="${BUILD_VERSION}" \
    io.hass.type="app" \
    io.hass.arch="${BUILD_ARCH}"

RUN mkdir -p /config /data

COPY --from=builder /out/homesignal-agent /usr/local/bin/homesignal-agent

EXPOSE 8099
ENTRYPOINT ["/usr/local/bin/homesignal-agent"]
