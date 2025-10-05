#!/bin/bash
set -e

echo "Installing development tools..."

# Detect OS
OS="$(uname -s)"
ARCH="$(uname -m)"

# Convert to golangci-lint naming
case "$OS" in
    Linux*)   OS_NAME="linux" ;;
    Darwin*)  OS_NAME="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) OS_NAME="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  ARCH_NAME="amd64" ;;
    aarch64|arm64) ARCH_NAME="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Development Environment Tools
GOLANGCI_VERSION="latest"
GOFUMPT_VERSION="latest"

# Check and install golangci-lint
if command -v golangci-lint &> /dev/null; then
    CURRENT_VERSION=$(golangci-lint version --format short 2>/dev/null || echo "unknown")
    echo "✓ golangci-lint already installed (${CURRENT_VERSION})"
else
    echo "→ Installing golangci-lint (latest stable)..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION}
    echo "✓ golangci-lint installed"
fi

# Check and install gofumpt
if command -v gofumpt &> /dev/null; then
    CURRENT_VERSION=$(gofumpt -version 2>/dev/null || echo "unknown")
    echo "✓ gofumpt already installed (${CURRENT_VERSION})"
else
    echo "→ Installing gofumpt ${GOFUMPT_VERSION}..."
    go install mvdan.cc/gofumpt@${GOFUMPT_VERSION}
    echo "✓ gofumpt installed"
fi

echo ""
echo "✓ Tools installed successfully"
echo ""
echo "Checking installed tool versions:"
echo "→ golangci-lint version:"
golangci-lint version
echo ""
echo "→ gofumpt version:"
gofumpt -version
