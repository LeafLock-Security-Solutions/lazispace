// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package app

// Version information for the application.
// These variables are intended to be set at build time using -ldflags.
//
// Example build command:
//
//	go build -ldflags "-X github.com/charanravela/lazispace/internal/app.Version=1.0.0 \
//	  -X github.com/charanravela/lazispace/internal/app.GitCommit=$(git rev-parse HEAD) \
//	  -X github.com/charanravela/lazispace/internal/app.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
//
// For local development builds without ldflags, these default to "dev"/"unknown".
var (
	// Version is the semantic version of the application (e.g., "1.0.0", "1.2.3-beta").
	// Injected at build time from Git tags or VERSION file.
	// Defaults to "dev" for development builds.
	Version = "dev"

	// GitCommit is the git commit SHA that this binary was built from.
	// Injected at build time.
	// Defaults to "unknown" for development builds.
	GitCommit = "unknown"

	// BuildDate is the timestamp when this binary was built (RFC3339 format).
	// Injected at build time.
	// Defaults to "unknown" for development builds.
	BuildDate = "unknown"
)

// GetVersion returns the current application version.
func GetVersion() string {
	return Version
}

// GetGitCommit returns the git commit SHA this binary was built from.
func GetGitCommit() string {
	return GitCommit
}

// GetBuildDate returns the build timestamp.
func GetBuildDate() string {
	return BuildDate
}

// GetVersionInfo returns a formatted string with all version information.
func GetVersionInfo() string {
	return "Version: " + Version + ", Commit: " + GitCommit + ", Built: " + BuildDate
}
