GO=$(shell which go)
BINPATH=$(CURDIR)/bin

PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 freebsd/386
PLATFORM = $(subst /, ,$@)
OS = $(word 1, $(PLATFORM))
ARCH = $(word 2, $(PLATFORM))

BINNAME=dockerate-ps
CMDPATH=cmd/dockerate-ps/*.go

.PHONY: deps lint build build_all
.DEFAULT_GOAL := build

deps:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.21.0
	@$(GO) get ./...

lint:
	@golangci-lint run

build:
	@echo building $(BINNAME)...
	@$(GO) build -o $(BINPATH)/$(BINNAME) -ldflags="-w -s" $(CMDPATH)

$(PLATFORMS):
	@echo building $(BINNAME) for $(OS)/$(ARCH)...
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(BINPATH)/$(BINNAME)_$(OS)_$(ARCH) $(CMDPATH)

build_all: $(PLATFORMS)
