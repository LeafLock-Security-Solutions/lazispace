.PHONY: help setup lint lint-fix fmt

# Variables
GO := go

help: ## Show this help message
	@echo 'LaziSpace Development Commands'
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

setup: ## Setup development environment
	@./scripts/dev-setup.sh

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
