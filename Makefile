GO=$(shell which go)
GOBIN=$(shell $(GO) env GOPATH)/bin
BINPATH=$(CURDIR)/bin

PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 freebsd/386
PLATFORM = $(subst /, ,$@)
OS = $(word 1, $(PLATFORM))
ARCH = $(word 2, $(PLATFORM))

BINNAME=dockerate-ps
CMDPATH=cmd/dockerate-ps/*.go

.PHONY: deps_lint deps lint build build_all clean
.DEFAULT_GOAL := build

deps_lint:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.21.0

deps: deps_lint
	@$(GO) get ./...

lint:
	@golangci-lint run

build:
	@echo building $(BINNAME)...
	@$(GO) build -o $(BINPATH)/$(BINNAME) -ldflags="-w -s" $(CMDPATH)

$(PLATFORMS):
	@echo building $(BINNAME) for $(OS)/$(ARCH)...
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(BINPATH)/$(BINNAME)_$(OS)_$(ARCH) -ldflags="-w -s" $(CMDPATH)

build_all: $(PLATFORMS)

clean:
	@rm -f $(BINPATH)/*
