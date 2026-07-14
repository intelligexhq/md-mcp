SHELL := /bin/bash
.PHONY: setup format lint test build run clean ci

GOFUMPT := go run mvdan.cc/gofumpt@latest
GOIMPORTS := go run golang.org/x/tools/cmd/goimports@latest
GOLANGCI := go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest

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
	go build -ldflags="-s -w" -o bin/mcp-md ./cmd/mcp-md

run: build
	./bin/app

clean:
	rm -rf bin/ coverage.*

ci: format lint test
	@echo "✅ All checks passed!"