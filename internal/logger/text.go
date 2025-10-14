// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

// Package logger implements text-based logging output.
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Text format constants.
const (
	textTimestampFormat = "15:04:05.000" // HH:MM:SS.mmm format for timestamps
	textLogPrefixFormat = "[%s] %s: %s"  // [timestamp] LEVEL: message
	textFieldSeparator  = " |"           // Separator between message and fields
	textFieldFormat     = " %s=%v"       // key=value format for fields
)

// textLogger implements the Logger interface with human-readable text output.
// Log entries are formatted as: [HH:MM:SS.mmm] LEVEL: message | key1=value1 key2=value2.
type textLogger struct {
	level  LogLevel  // Minimum log level to output
	output io.Writer // Destination for log output
}

// newTextLogger creates a new text format logger with the specified level and output destination.
func newTextLogger(level LogLevel, output io.Writer) Logger {
	return &textLogger{level: level, output: output}
}

// Debug logs a debug-level message with optional structured fields.
// Only outputs if the logger's level is debug or lower.
func (l *textLogger) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.log(LevelDebug.UpperString(), msg, fields...)
	}
}

// Info logs an info-level message with optional structured fields.
// Only outputs if the logger's level is info or lower.
func (l *textLogger) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.log(LevelInfo.UpperString(), msg, fields...)
	}
}

// Warn logs a warning-level message with optional structured fields.
// Only outputs if the logger's level is warn or lower.
func (l *textLogger) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.log(LevelWarn.UpperString(), msg, fields...)
	}
}

// Error logs an error-level message with optional structured fields.
// Always outputs regardless of log level setting.
func (l *textLogger) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.log(LevelError.UpperString(), msg, fields...)
	}
}

// Fatal logs a fatal-level message with optional structured fields and then exits the program.
// Calls osExit(1) after logging. Always outputs regardless of log level setting.
func (l *textLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal.UpperString(), msg, fields...)
	osExit(1)
}

// log formats and writes a log entry in text format.
// Format: [HH:MM:SS.mmm] LEVEL: message | key1=value1 key2=value2
// If any write operation fails, an error message is written to stderr.
func (l *textLogger) log(level, msg string, fields ...Field) {
	timestamp := time.Now().Format(textTimestampFormat)

	if _, err := fmt.Fprintf(l.output, textLogPrefixFormat, timestamp, level, msg); err != nil {
		_, _ = os.Stderr.WriteString(errMsgLogPrefixWriteFailed)
		return
	}

	if len(fields) > 0 {
		if _, err := fmt.Fprint(l.output, textFieldSeparator); err != nil {
			_, _ = os.Stderr.WriteString(errMsgLogSepWriteFailed)
			return
		}
		for _, f := range fields {
			if _, err := fmt.Fprintf(l.output, textFieldFormat, f.Key, f.Value); err != nil {
				_, _ = os.Stderr.WriteString(errMsgLogFieldWriteFailed)
				return
			}
		}
	}

	if _, err := fmt.Fprintln(l.output); err != nil {
		_, _ = os.Stderr.WriteString(errMsgLogNewlineWriteFailed)
	}
}
