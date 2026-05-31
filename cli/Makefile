BINARY  := devdeck
PKG     := github.com/jazho76/devdeck/cli/internal/version

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(PKG).Version=$(VERSION) -X $(PKG).Commit=$(COMMIT) -X $(PKG).Date=$(DATE)

.PHONY: build vet test tidy install clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

install:
	go install -ldflags "$(LDFLAGS)" .

clean:
	rm -rf bin
