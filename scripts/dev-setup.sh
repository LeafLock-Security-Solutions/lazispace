#!/bin/bash
set -e

echo "======================================="
echo "LaziSpace Development Environment Setup"
echo "======================================="
echo ""

# Check Go version
echo "→ Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "✗ Go is not installed"
    echo "  Please install Go 1.25.1 or later from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go ${GO_VERSION} found"
echo ""

# Install tools
echo "→ Installing development tools..."
if [ -f ./scripts/install-tools.sh ]; then
    ./scripts/install-tools.sh
else
    echo "✗ scripts/install-tools.sh not found"
    exit 1
fi
echo ""

# Install project dependencies
echo "→ Installing project dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installed"
echo ""

# Create directories
echo "→ Creating necessary directories..."
mkdir -p bin
mkdir -p dev-data/{config,data,logs}
mkdir -p test-data/{config,data,logs}
echo "✓ Directories created"
echo ""

# Setup git hooks (if .git exists)
if [ -d .git ]; then
    echo "→ Configuring git hooks..."
    if [ -d .githooks ]; then
        git config core.hooksPath .githooks
        chmod +x .githooks/* 2>/dev/null || true
        echo "✓ Git hooks configured"
    else
        echo "⚠  No .githooks directory found (will skip)"
    fi
    echo ""
fi

echo "======================================="
echo "✓ Setup complete!"
echo "======================================="
echo ""
echo "Next steps:"
echo "  make dev      - Run in development mode"
echo "  make test     - Run tests"
echo "  make lint     - Check code quality"
echo "  make fmt      - Format code"
echo "  make build    - Build binary"
echo ""
echo "For all commands: make help"
