GO=$(shell which go)

.PHONY: deps lint build
.DEFAULT_GOAL := build

deps:
	@$(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) get ./...

lint:
	@golangci-lint run --tests=false --enable-all -D wsl ./...

build:
	@echo building dockerate-ps...
	@$(GO) build -o bin/dockerate-ps cmd/dockerate-ps/main.go
