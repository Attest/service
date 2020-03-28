.PHONY: test fmt

goos := linux
VERSION := undefined
gotest := go test -race -timeout 1m
golangci-lint := golangci-lint
FILES := $(shell go list ./... | grep -v /system)

default: test

build-test:
	rm -r build/server || true
	go build -o build/ ./examples/server

test: lint test-unit test-system

test-unit:
	$(gotest) $(FILES)

test-system: build-test
	$(gotest) ./system

lint:
	go vet ./...
	$(golangci-lint) run
	go mod tidy

fmt:
	$(golangci-lint) run --fix --fast
	go mod tidy
	go fmt ./...

