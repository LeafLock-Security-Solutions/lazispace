// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package interfaces

// Field represents a key-value pair for structured logging.
// This is duplicated from logger.Field to avoid import cycles.
type Field struct {
	Key   string
	Value any
}

// Logger defines the interface for logging operations.
//
// This interface is defined in a separate package to avoid import cycles
// between packages that need to use logging (like app) and the logger
// package itself.
//
// The logger package implements this interface with various logger types
// (text, JSON, buffered), allowing other packages to use logging without
// directly importing the logger package.
type Logger interface {
	// Debug logs a debug-level message with optional structured fields.
	Debug(msg string, fields ...Field)

	// Info logs an info-level message with optional structured fields.
	Info(msg string, fields ...Field)

	// Warn logs a warning-level message with optional structured fields.
	Warn(msg string, fields ...Field)

	// Error logs an error-level message with optional structured fields.
	Error(msg string, fields ...Field)

	// Fatal logs a fatal-level message with optional structured fields.
	// Implementations should exit the program after logging (typically with os.Exit(1)).
	Fatal(msg string, fields ...Field)
}
