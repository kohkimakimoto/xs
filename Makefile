.DEFAULT_GOAL := help

SHELL := bash
PATH := $(CURDIR)/.dev/gopath/bin:$(PATH)
VERSION := 0.0.1
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/kohkimakimoto/xs/internal.CommitHash=$(COMMIT_HASH) -X github.com/kohkimakimoto/xs/internal.Version=$(VERSION)"

# Load .env file if it exists.
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

.PHONY: help
help: ## Show help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'


# --------------------------------------------------------------------------------------
# Development environment
# --------------------------------------------------------------------------------------
.PHONY: setup
setup: ## Setup development environment
	@echo "==> Setting Go tools up..."
	@mkdir -p .dev/gopath
	@export GOPATH=$(CURDIR)/.dev/gopath && \
		go install honnef.co/go/tools/cmd/staticcheck@latest && \
		go install github.com/Songmu/goxz/cmd/goxz@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@export GOPATH=$(CURDIR)/.dev/gopath && go clean -modcache

.PHONY: clean
clean: ## Clean up development environment
	@rm -rf .dev

.PHONY: clean/build
clean/build: ## Clean up build files
	@rm -rf .dev/build



# --------------------------------------------------------------------------------------
# SSH Server for demo/dev
# --------------------------------------------------------------------------------------
.PHONY: demo/ssh-server/up
demo/ssh-server/up: ## Start dev ssh server
	@cd demo && docker-compose up -d

.PHONY: demo/ssh-server/down
demo/ssh-server/down: ## Stop dev ssh server
	@cd demo && docker-compose down

# --------------------------------------------------------------------------------------
# Testing, Formatting and etc.
# --------------------------------------------------------------------------------------
.PHONY: format
format: ## Format go code
	@go fmt ./...

.PHONY: lint
lint: ## Lint source code
	@go vet ./... ; staticcheck ./...

.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m ./...

.PHONY: test/verbos
test/verbose: ## Run tests with verbose outputting
	@go test -race -timeout 30m -v ./...

.PHONY: test/cover
test/cover: ## Run tests with coverage
	@echo "==> Run tests with coverage report..."
	@mkdir -p $(CURDIR)/.dev/coverage
	@go test -coverpkg=./... -coverprofile=$(CURDIR)/.dev/coverage/coverage.out ./...
	@gocov convert $(CURDIR)/.dev/coverage/coverage.out | gocov-html > $(CURDIR)/.dev/coverage/coverage.html
	@echo "==> Open $(CURDIR)/.dev/coverage/coverage.html to see the coverage report."

.PHONY: open/coverage
open/coverage: ## Open coverage report
	@open $(CURDIR)/.dev/coverage/coverage.html


# --------------------------------------------------------------------------------------
# Build
# --------------------------------------------------------------------------------------
.PHONY: build
build: ## Build gold executable
	@mkdir -p .dev/build/dev
	@go build -ldflags=$(BUILD_LDFLAGS) -o .dev/build/dev/xs ./cmd/xs

.PHONY: build/release
build/release: ## build release binary
	@mkdir -p .dev/build/release
	@goxz -n xs -pv=v$(VERSION) -os=linux,darwin -static -build-ldflags=$(BUILD_LDFLAGS) -d=.dev/build/release ./cmd/xs


# --------------------------------------------------------------------------------------
# etc.
# --------------------------------------------------------------------------------------
.PHONY: go-mod-tidy
go-mod-tidy: ## Run go mod tidy
	@go mod tidy

