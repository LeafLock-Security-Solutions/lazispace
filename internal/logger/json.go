package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// JSON field names for log entries.
const (
	jsonFieldTime  = "time"
	jsonFieldLevel = "level"
	jsonFieldMsg   = "msg"
	jsonFieldError = "error"
)

// JSON error messages for fallback scenarios.
const (
	jsonErrMarshalFailed = "failed to marshal log entry"
)

// jsonLogger implements the Logger interface with JSON-formatted output.
// Each log entry is a single-line JSON object with fields: time, level, msg, and any custom fields.
type jsonLogger struct {
	level  LogLevel  // Minimum log level to output
	output io.Writer // Destination for log output
}

// newJSONLogger creates a new JSON format logger with the specified level and output destination.
func newJSONLogger(level LogLevel, output io.Writer) Logger {
	return &jsonLogger{level: level, output: output}
}

// Debug logs a debug-level message as JSON with optional structured fields.
// Only outputs if the logger's level is debug or lower.
func (l *jsonLogger) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.log(LevelDebug.String(), msg, fields...)
	}
}

// Info logs an info-level message as JSON with optional structured fields.
// Only outputs if the logger's level is info or lower.
func (l *jsonLogger) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.log(LevelInfo.String(), msg, fields...)
	}
}

// Warn logs a warning-level message as JSON with optional structured fields.
// Only outputs if the logger's level is warn or lower.
func (l *jsonLogger) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.log(LevelWarn.String(), msg, fields...)
	}
}

// Error logs an error-level message as JSON with optional structured fields.
// Always outputs regardless of log level setting.
func (l *jsonLogger) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.log(LevelError.String(), msg, fields...)
	}
}

// Fatal logs a fatal-level message as JSON with optional structured fields and then exits the program.
// Calls osExit(1) after logging. Always outputs regardless of log level setting.
func (l *jsonLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal.String(), msg, fields...)
	osExit(1)
}

// log formats and writes a log entry as a JSON object.
// Output format: {"time":"2006-01-02T15:04:05Z07:00","level":"info","msg":"message","field1":"value1"}
// If JSON marshaling fails or write operations fail, error messages are written to stderr.
func (l *jsonLogger) log(level, msg string, fields ...Field) {
	entry := map[string]any{
		jsonFieldTime:  time.Now().Format(time.RFC3339),
		jsonFieldLevel: level,
		jsonFieldMsg:   msg,
	}

	for _, f := range fields {
		entry[f.Key] = f.Value
	}

	data, err := json.Marshal(entry)
	if err != nil {
		// If marshaling fails, write a simple error message to stderr
		// We can't do much if this fails, so we just ignore the error
		fallbackMsg := fmt.Sprintf(`{"%s":"%s"}`+"\n", jsonFieldError, jsonErrMarshalFailed)
		if _, writeErr := l.output.Write([]byte(fallbackMsg)); writeErr != nil {
			// Last resort: write to stderr directly
			_, _ = os.Stderr.WriteString(errMsgLogEntryMarshalFailed)
		}
		return
	}

	if _, err := l.output.Write(data); err != nil {
		_, _ = os.Stderr.WriteString(errMsgLogDataWriteFailed)
		return
	}
	if _, err := l.output.Write([]byte("\n")); err != nil {
		_, _ = os.Stderr.WriteString(errMsgLogNewlineWriteFailed)
	}
}
