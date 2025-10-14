// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package logger

import (
	"sync"
	"time"
)

// Field names added to bootstrap logs when replayed.
const (
	bootstrapFieldName = "bootstrap"
	timestampFieldName = "timestamp"
)

// Bootstrap buffer configuration.
const (
	bootstrapBufferCapacity = 20 // Expected number of bootstrap log entries
)

// bufferedLogger stores log entries in memory until they can be replayed.
//
// This logger is used during application bootstrap before the real logger is initialized.
// All log entries are stored in memory with their timestamp, level, message, and fields.
// Once the real logger is ready, entries can be replayed with ReplayTo().
//
// Thread-safety: All operations are protected by a mutex for concurrent access.
type bufferedLogger struct {
	mu      sync.Mutex // Protects entries from concurrent access
	entries []logEntry // Buffered log entries
}

// logEntry represents a single buffered log entry with all its metadata.
type logEntry struct {
	timestamp time.Time // When the log was originally created
	level     string    // Log level as string (debug, info, warn, error, fatal)
	msg       string    // Log message
	fields    []Field   // Structured fields attached to this log entry
}

// NewBuffered creates a logger that buffers all log entries in memory.
//
// The returned logger stores log entries until ReplayTo() is called to write them
// to a real logger. This is useful for bootstrap phase logging where the logging
// configuration isn't yet available.
//
// Pre-allocates space for bootstrap log entries to minimize allocations during startup.
func NewBuffered() BufferedLogger {
	return &bufferedLogger{
		entries: make([]logEntry, 0, bootstrapBufferCapacity),
	}
}

// Debug buffers a debug-level log entry for later replay.
func (b *bufferedLogger) Debug(msg string, fields ...Field) {
	b.append(LevelDebug.String(), msg, fields...)
}

// Info buffers an info-level log entry for later replay.
func (b *bufferedLogger) Info(msg string, fields ...Field) {
	b.append(LevelInfo.String(), msg, fields...)
}

// Warn buffers a warning-level log entry for later replay.
func (b *bufferedLogger) Warn(msg string, fields ...Field) {
	b.append(LevelWarn.String(), msg, fields...)
}

// Error buffers an error-level log entry for later replay.
func (b *bufferedLogger) Error(msg string, fields ...Field) {
	b.append(LevelError.String(), msg, fields...)
}

// Fatal buffers a fatal-level log entry for later replay.
//
// Note: Unlike other logger implementations, this does NOT call os.Exit(1).
// The buffered Fatal entry will be replayed later, and the real logger will
// handle exiting the program at that time.
func (b *bufferedLogger) Fatal(msg string, fields ...Field) {
	b.append(LevelFatal.String(), msg, fields...)
	// Note: We don't actually exit here since we're buffering
	// The real logger will handle the exit when replayed
}

// append stores a log entry in the buffer with the current timestamp.
// Thread-safe: protected by mutex.
func (b *bufferedLogger) append(level, msg string, fields ...Field) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = append(b.entries, logEntry{
		timestamp: time.Now(),
		level:     level,
		msg:       msg,
		fields:    fields,
	})
}

// ReplayTo replays all buffered log entries to the target logger.
//
// Important: If any buffered entry is a Fatal-level log, calling Fatal on the
// target logger will cause the program to exit immediately (via os.Exit(1)),
// preventing any subsequent buffered entries from being replayed.
//
// This is intentional behavior - Fatal logs indicate unrecoverable errors that
// should terminate the program. If you have Fatal logs during bootstrap, they
// will be replayed first and exit before other logs can be written.
func (b *bufferedLogger) ReplayTo(target Logger) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, entry := range b.entries {
		// Add a field indicating this is a replayed bootstrap log
		fieldsWithBootstrap := append(
			[]Field{NewField(bootstrapFieldName, true), NewField(timestampFieldName, entry.timestamp.Format(time.RFC3339))},
			entry.fields...,
		)

		switch entry.level {
		case LevelDebug.String():
			target.Debug(entry.msg, fieldsWithBootstrap...)
		case LevelInfo.String():
			target.Info(entry.msg, fieldsWithBootstrap...)
		case LevelWarn.String():
			target.Warn(entry.msg, fieldsWithBootstrap...)
		case LevelError.String():
			target.Error(entry.msg, fieldsWithBootstrap...)
		case LevelFatal.String():
			// Note: This will call os.Exit(1) and terminate the program immediately.
			// Any remaining buffered logs will NOT be replayed.
			target.Fatal(entry.msg, fieldsWithBootstrap...)
		}
	}

	// Clear buffer after replay
	b.entries = nil
}

// Count returns the number of buffered entries.
func (b *bufferedLogger) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}
