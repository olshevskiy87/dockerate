GO=$(shell which go)

.PHONY: deps lint build
.DEFAULT_GOAL := build

deps:
	@$(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) get ./...

lint:
	@echo gofmt...
	@if [ ! -z "$(shell gofmt -s -l .)" ]; then echo "gofmt blames these files:" && gofmt -s -l . && exit 1; fi
	@echo go vet...
	@$(GO) vet ./...
	@echo golangci-lint...
	@golangci-lint run --tests=false --enable-all -D prealloc -D gochecknoglobals ./...

build:
	@echo building dockerate-ps...
	@$(GO) build -o bin/dockerate-ps cmd/dockerate-ps/main.go
