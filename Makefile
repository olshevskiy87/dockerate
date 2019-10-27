GO=$(shell which go)

.PHONY: deps lint build
.DEFAULT_GOAL := build

deps:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.21.0
	@$(GO) get ./...

lint:
	@golangci-lint run

build:
	@echo building dockerate-ps...
	@$(GO) build -o bin/dockerate-ps cmd/dockerate-ps/main.go
