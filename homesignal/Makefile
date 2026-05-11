VERSION ?= 0.1.3
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
IMAGE ?= homesignal:$(VERSION)
GO_IMAGE ?= golang:1.25-alpine

.PHONY: test build docker-build run-agent

test:
	@if command -v go >/dev/null 2>&1; then \
		go test ./...; \
	else \
		docker run --rm -v "$$(pwd)":/src -w /src $(GO_IMAGE) go test ./...; \
	fi

build:
	@if command -v go >/dev/null 2>&1; then \
		go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)" -o bin/homesignal-agent ./cmd/agent; \
	else \
		docker run --rm -v "$$(pwd)":/src -w /src $(GO_IMAGE) go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)" -o bin/homesignal-agent ./cmd/agent; \
	fi

docker-build:
	docker build \
		--build-arg BUILD_VERSION=$(VERSION) \
		--build-arg BUILD_ARCH=amd64 \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		-t $(IMAGE) .

run-agent:
	mkdir -p .local/config .local/data
	@if command -v go >/dev/null 2>&1; then \
		CONFIG_DIR=.local/config DATA_DIR=.local/data LISTEN_ADDR=:8099 go run ./cmd/agent; \
	else \
		docker build \
			--build-arg BUILD_VERSION=$(VERSION) \
			--build-arg BUILD_ARCH=amd64 \
			--build-arg COMMIT=$(COMMIT) \
			--build-arg BUILD_TIME=$(BUILD_TIME) \
			-t $(IMAGE) .; \
		docker run --rm -p 127.0.0.1:8099:8099 \
			-v "$$(pwd)/.local/config":/config \
			-v "$$(pwd)/.local/data":/data \
			$(IMAGE); \
	fi
