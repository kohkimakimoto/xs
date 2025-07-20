.DEFAULT_GOAL := help

SHELL := bash
PATH := $(CURDIR)/.dev/go-tools/bin:$(PATH)
COMMIT_HASH := $(shell git rev-parse HEAD)

VERSION := 0.0.7
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
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'


# --------------------------------------------------------------------------------------
# Development environment
# --------------------------------------------------------------------------------------
.PHONY: setup
setup: ## Setup development environment
	@echo "==> Setting up development environment..."
	@mkdir -p $(CURDIR)/.dev/go-tools
	@export GOPATH=$(CURDIR)/.dev/go-tools && \
		go install github.com/Songmu/goxz/cmd/goxz@latest && \
		go install github.com/tcnksm/ghr@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html/cmd/gocov-html@latest
	@export GOPATH=$(CURDIR)/.dev/go-tools && go clean -modcache && rm -rf $(CURDIR)/.dev/go-tools/pkg

.PHONY: clean
clean: ## Clean up development environment
	@rm -rf .dev


# --------------------------------------------------------------------------------------
# Build
# --------------------------------------------------------------------------------------
.PHONY: build
build: ## Build dev binary
	@mkdir -p .dev/build/dev
	@CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) -o .dev/build/dev/xs ./cmd/xs

.PHONY: build-release
build-release: ## Build release binary
	@mkdir -p .dev/build/release
	@CGO_ENABLED=0 go build -ldflags=$(BUILD_LDFLAGS) -trimpath -o .dev/build/release/xs ./cmd/xs

.PHONY: build-dist
build-dist: ## Build cross-platform binaries for distribution
	@mkdir -p .dev/build/dist
	@CGO_ENABLED=0 goxz -n xs -os=linux,darwin -static -build-ldflags=$(BUILD_LDFLAGS) -trimpath -d=.dev/build/dist ./cmd/xs

.PHONY: build-clean
build-clean: ## Clean up build artifacts
	@rm -rf .dev/build


# --------------------------------------------------------------------------------------
# Testing, Formatting and etc.
# --------------------------------------------------------------------------------------
.PHONY: format
format: ## Format source code
	@go fmt ./...

.PHONY: test
test: ## Run tests
	@go test -race -timeout 30m ./...

.PHONY: test-short
test-short: ## Run short tests
	@go test -short -race -timeout 30m ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose outputting
	@go test -race -timeout 30m -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage report
	@mkdir -p $(CURDIR)/.dev/test
	@go test -race -coverpkg=./... -coverprofile=$(CURDIR)/.dev/test/coverage.out ./...
	@gocov convert $(CURDIR)/.dev/test/coverage.out | gocov-html > $(CURDIR)/.dev/test/coverage.html

.PHONY: test-cover-open
test-cover-open: ## Open coverage report in browser
	@open $(CURDIR)/.dev/test/coverage.html


# --------------------------------------------------------------------------------------
# SSH Server for demo/dev
# --------------------------------------------------------------------------------------
.PHONY: demo-ssh-server-up
demo-ssh-server-up: ## Start dev ssh server
	@cd demo && docker-compose up -d

.PHONY: demo-ssh-server-down
demo-ssh-server-down: ## Stop dev ssh server
	@cd demo && docker-compose down
