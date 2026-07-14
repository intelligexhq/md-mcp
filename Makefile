SHELL := /bin/bash
.PHONY: setup format lint test build run clean ci

GOFUMPT := go run mvdan.cc/gofumpt@latest
GOIMPORTS := go run golang.org/x/tools/cmd/goimports@latest
GOLANGCI := go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

setup:
	go mod tidy

format:
	$(GOIMPORTS) -w .
	$(GOFUMPT) -w .

lint:
	$(GOLANGCI) run ./...

test:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build -ldflags="$(LDFLAGS)" -o bin/md-mcp ./cmd/md-mcp

run: build
	./bin/md-mcp

clean:
	rm -rf bin/ coverage.*

ci: format lint test
	@echo "✅ All checks passed!"