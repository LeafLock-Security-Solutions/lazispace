// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

// Package constants provides shared constant values used across the application.
// This prevents duplication and ensures consistency across packages.
package constants

// Log level string values.
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Log format string values.
const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

// ValidLogLevels contains all acceptable log level values.
// Do not modify this slice at runtime.
var ValidLogLevels = []string{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

// ValidLogFormats contains all acceptable log format values.
// Do not modify this slice at runtime.
var ValidLogFormats = []string{LogFormatText, LogFormatJSON}
