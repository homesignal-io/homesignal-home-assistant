# syntax=docker/dockerfile:1

ARG BUILDPLATFORM=linux/amd64

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=0.1.0
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
    -o /out/homesignal-agent ./cmd/agent

FROM alpine:3.22

RUN addgroup -S homesignal \
    && adduser -S -D -H -u 1000 -G homesignal homesignal \
    && mkdir -p /config /data \
    && chown -R homesignal:homesignal /config /data

COPY --from=builder /out/homesignal-agent /usr/local/bin/homesignal-agent

USER homesignal:homesignal
EXPOSE 8099
ENTRYPOINT ["/usr/local/bin/homesignal-agent"]
