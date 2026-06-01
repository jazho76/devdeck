BINARY  := devdeck
PKG     := github.com/jazho76/devdeck/cli/internal/version

VERSION := $(shell git describe --tags --match 'v*' --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(PKG).Version=$(VERSION) -X $(PKG).Commit=$(COMMIT) -X $(PKG).Date=$(DATE)

RUN_IMAGE     := devdeck-sandbox
RUN_CONTAINER := devdeck-sandbox

RUN_IMAGE_READY     := devdeck-sandbox-ready
RUN_CONTAINER_READY := devdeck-sandbox-ready

.PHONY: build dist vet test tidy install clean release image run image-ready run-ready

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

dist:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY)-linux-amd64 .

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

release:
	@echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$' || { echo "Pass a semver: make release VERSION=v0.1.0"; exit 1; }
	git tag $(VERSION)
	git push origin $(VERSION)

image:
	docker build -t $(RUN_IMAGE) -f sandbox/Dockerfile sandbox

run: dist image
	setsid alacritty -e docker run --rm -it --name $(RUN_CONTAINER) \
		-v "$(CURDIR)/bin/$(BINARY)-linux-amd64:/home/dev/.local/bin/devdeck:ro" \
		-v "$(CURDIR)/sandbox/scripts:/opt/sandbox/bin:ro" \
		$(RUN_IMAGE) >/dev/null 2>&1 &
	@echo "Sandbox opening in a new Alacritty window; removed automatically when you exit the shell."

image-ready:
	docker build -t $(RUN_IMAGE_READY) -f sandbox/Dockerfile.ready sandbox

run-ready: dist image-ready
	setsid alacritty -e docker run --rm -it --name $(RUN_CONTAINER_READY) \
		-v "$(CURDIR)/bin/$(BINARY)-linux-amd64:/home/dev/.local/bin/devdeck:ro" \
		$(RUN_IMAGE_READY) >/dev/null 2>&1 &
	@echo "Ready sandbox opening in a new Alacritty window; removed automatically when you exit the shell."
