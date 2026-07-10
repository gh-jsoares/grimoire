VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build install clean test lint

build:
	go build -ldflags "$(LDFLAGS)" -o grimoire ./cmd/grimoire

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/grimoire

clean:
	rm -f grimoire
	rm -rf dist/

test:
	go test ./...

lint:
	golangci-lint run ./...
