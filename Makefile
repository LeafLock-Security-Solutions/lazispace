.PHONY: help setup lint lint-fix fmt hooks test test-verbose test-coverage test-clean

# Variables
GO := go

help: ## Show this help message
	@echo 'LaziSpace Development Commands'
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

setup: ## Setup development environment
	@./scripts/dev-setup.sh

hooks: ## Configure git hooks
	@echo "Configuring git hooks..."
	@git config core.hooksPath .githooks
	@echo "Git hooks configured"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	@golangci-lint run --fix

fmt: ## Format code
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@gofumpt -l -w .
	@echo "Code formatted"

test: ## Run tests
	@echo "Running tests..."
	@$(GO) test -v -race ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@$(GO) test -v -race -count=1 ./...

test-coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	@$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-clean: ## Remove test artifacts
	@echo "Cleaning test artifacts..."
	@rm -f coverage.out coverage.html
	@echo "Test artifacts cleaned"
