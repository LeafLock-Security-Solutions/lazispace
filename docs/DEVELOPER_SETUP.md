# Developer Setup Guide

This guide will help you set up your development environment for LaziSpace.

## Prerequisites

- Go 1.25 or later
- Git
- Make

### Installing Go

**macOS:**
```bash
brew install go
```

**Linux:**
```bash
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**Windows:**  
Download and install from https://go.dev/dl/

### Verify Installation

```bash
go version  # Should show Go 1.25 or later
```

## Quick Setup

Clone the repository and run the setup script:

```bash
git clone https://github.com/LeafLock-Security-Solutions/lazispace.git
cd lazispace
make setup
```

The setup script will:
- Check Go installation
- Install development tools (golangci-lint, gofumpt)
- Download project dependencies
- Create necessary directories
- Configure git hooks (if applicable)

## Verify Setup

```bash
make help        # Show available commands
make lint        # Run linter
make fmt         # Format code
```

## Development Workflow

### Available Commands

```bash
make help        # Show all available commands
make setup       # Run the development setup script
make lint        # Run golangci-lint
make lint-fix    # Run golangci-lint with auto-fix
make fmt         # Format code with gofmt and gofumpt
```

### Before Committing

```bash
make fmt         # Format code
make lint        # Check for issues
```

## Troubleshooting

### "command not found: golangci-lint"

Ensure GOPATH/bin is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Add to ~/.bashrc or ~/.zshrc to persist.

### "no go files to analyze"

Make sure you're in the project root directory where go.mod exists.

### Import errors

Run:
```bash
go mod tidy
go mod download
```

## Getting Help

- Check the Makefile: `make help`
- Read the documentation in docs/
- Check existing issues on GitHub
- Review commit guidelines: `docs/COMMIT_GUIDELINES.md`

## Next Steps

After setup is complete:
1. Review the commit guidelines: `docs/COMMIT_GUIDELINES.md`
2. Read the workspace config guide: `docs/WORKSPACE_CONFIG_GUIDE.md`
3. Check the file locations: `docs/FILE_LOCATIONS.md`
4. Start contributing by following `CONTRIBUTING.md`
