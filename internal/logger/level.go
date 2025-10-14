// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"strings"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/constants"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

// Log level constants represent the severity hierarchy.
const (
	LevelDebug LogLevel = iota // Lowest: detailed debugging information
	LevelInfo                  // Informational messages
	LevelWarn                  // Warning messages
	LevelError                 // Error messages
	LevelFatal                 // Highest: fatal errors that terminate the program
)

// Log level string representations.
const (
	levelDebugStr   = constants.LogLevelDebug
	levelInfoStr    = constants.LogLevelInfo
	levelWarnStr    = constants.LogLevelWarn
	levelWarningStr = "warning" // Alias for warn
	levelErrorStr   = constants.LogLevelError
	levelFatalStr   = "fatal"
	levelUnknownStr = "unknown"
)

// String returns the canonical string representation of the log level (lowercase).
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return levelDebugStr
	case LevelInfo:
		return levelInfoStr
	case LevelWarn:
		return levelWarnStr
	case LevelError:
		return levelErrorStr
	case LevelFatal:
		return levelFatalStr
	default:
		return levelUnknownStr
	}
}

// UpperString returns the uppercase string representation of the log level.
// This is used by text formatters that display levels in uppercase (e.g., "DEBUG", "INFO").
func (l LogLevel) UpperString() string {
	return strings.ToUpper(l.String())
}
