package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/app"
)

// TestParseLevel tests log level parsing.
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      LogLevel
		wantError bool
	}{
		{"debug level", "debug", LevelDebug, false},
		{"info level", "info", LevelInfo, false},
		{"warn level", "warn", LevelWarn, false},
		{"warning level", "warning", LevelWarn, false},
		{"error level", "error", LevelError, false},
		{"uppercase debug", "DEBUG", LevelDebug, false},
		{"mixed case info", "Info", LevelInfo, false},
		{"invalid level", "invalid", LevelInfo, true},
		{"empty string", "", LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLevel(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("parseLevel(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("parseLevel(%q) unexpected error: %v", tt.input, err)
				}
				if got != tt.want {
					t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.want)
				}
			}
		})
	}
}

// TestTextLogger tests text format logging.
func TestTextLogger(t *testing.T) {
	tests := []struct {
		name        string
		level       LogLevel
		logFunc     func(Logger, string, ...Field)
		msg         string
		fields      []Field
		shouldLog   bool
		wantContain []string
	}{
		{
			name:        "debug log at debug level",
			level:       LevelDebug,
			logFunc:     func(l Logger, msg string, fields ...Field) { l.Debug(msg, fields...) },
			msg:         "debug message",
			shouldLog:   true,
			wantContain: []string{"DEBUG", "debug message"},
		},
		{
			name:        "debug log at info level",
			level:       LevelInfo,
			logFunc:     func(l Logger, msg string, fields ...Field) { l.Debug(msg, fields...) },
			msg:         "debug message",
			shouldLog:   false,
			wantContain: []string{},
		},
		{
			name:        "info log at info level",
			level:       LevelInfo,
			logFunc:     func(l Logger, msg string, fields ...Field) { l.Info(msg, fields...) },
			msg:         "info message",
			shouldLog:   true,
			wantContain: []string{"INFO", "info message"},
		},
		{
			name:        "warn log at warn level",
			level:       LevelWarn,
			logFunc:     func(l Logger, msg string, fields ...Field) { l.Warn(msg, fields...) },
			msg:         "warning message",
			shouldLog:   true,
			wantContain: []string{"WARN", "warning message"},
		},
		{
			name:        "error log at error level",
			level:       LevelError,
			logFunc:     func(l Logger, msg string, fields ...Field) { l.Error(msg, fields...) },
			msg:         "error message",
			shouldLog:   true,
			wantContain: []string{"ERROR", "error message"},
		},
		{
			name:  "log with fields",
			level: LevelInfo,
			logFunc: func(l Logger, msg string, fields ...Field) {
				l.Info(msg, fields...)
			},
			msg:         "message with fields",
			fields:      []Field{NewField("key1", "value1"), NewField("key2", 42)},
			shouldLog:   true,
			wantContain: []string{"INFO", "message with fields", "key1=value1", "key2=42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTextLogger(tt.level, &buf)

			tt.logFunc(logger, tt.msg, tt.fields...)

			output := buf.String()

			if !tt.shouldLog {
				if output != "" {
					t.Errorf("expected no log output, got: %s", output)
				}
				return
			}

			// Verify logged output
			if output == "" {
				t.Errorf("expected log output, got empty string")
			}
			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("log output missing %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestJSONLogger tests JSON format logging.
func TestJSONLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     LogLevel
		logFunc   func(Logger, string, ...Field)
		msg       string
		fields    []Field
		shouldLog bool
		wantLevel string
	}{
		{
			name:      "debug log at debug level",
			level:     LevelDebug,
			logFunc:   func(l Logger, msg string, fields ...Field) { l.Debug(msg, fields...) },
			msg:       "debug message",
			shouldLog: true,
			wantLevel: "debug",
		},
		{
			name:      "debug log at info level",
			level:     LevelInfo,
			logFunc:   func(l Logger, msg string, fields ...Field) { l.Debug(msg, fields...) },
			msg:       "debug message",
			shouldLog: false,
		},
		{
			name:      "info log",
			level:     LevelInfo,
			logFunc:   func(l Logger, msg string, fields ...Field) { l.Info(msg, fields...) },
			msg:       "info message",
			shouldLog: true,
			wantLevel: "info",
		},
		{
			name:      "warn log",
			level:     LevelWarn,
			logFunc:   func(l Logger, msg string, fields ...Field) { l.Warn(msg, fields...) },
			msg:       "warning message",
			shouldLog: true,
			wantLevel: "warn",
		},
		{
			name:      "error log",
			level:     LevelError,
			logFunc:   func(l Logger, msg string, fields ...Field) { l.Error(msg, fields...) },
			msg:       "error message",
			shouldLog: true,
			wantLevel: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newJSONLogger(tt.level, &buf)

			tt.logFunc(logger, tt.msg, tt.fields...)

			output := buf.String()

			if !tt.shouldLog {
				if output != "" {
					t.Errorf("expected no log output, got: %s", output)
				}
				return
			}

			// Verify logged output exists
			if output == "" {
				t.Errorf("expected log output, got empty string")
				return
			}

			// Parse JSON
			var entry map[string]any
			if err := json.Unmarshal([]byte(output), &entry); err != nil {
				t.Fatalf("failed to parse JSON log: %v", err)
			}

			// Verify fields
			if entry["level"] != tt.wantLevel {
				t.Errorf("level = %v, want %v", entry["level"], tt.wantLevel)
			}
			if entry["msg"] != tt.msg {
				t.Errorf("msg = %v, want %v", entry["msg"], tt.msg)
			}
			if _, ok := entry["time"]; !ok {
				t.Error("JSON entry missing 'time' field")
			}
		})
	}
}

// TestJSONLoggerWithFields tests JSON logging with custom fields.
func TestJSONLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(LevelInfo, &buf)

	logger.Info("test message", NewField("user", "alice"), NewField("count", 5))

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["user"] != "alice" {
		t.Errorf("user field = %v, want alice", entry["user"])
	}
	if entry["count"] != float64(5) {
		t.Errorf("count field = %v, want 5", entry["count"])
	}
}

// TestNewLogger tests logger creation from config.
func TestNewLogger(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		cfg       *app.Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "text logger with console output",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "info",
					Format: "text",
					Console: app.ConsoleConfig{
						Enabled: true,
					},
					File: app.FileLogConfig{
						Enabled: false,
					},
				},
			},
			wantError: false,
		},
		{
			name: "json logger with console output",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "debug",
					Format: "json",
					Console: app.ConsoleConfig{
						Enabled: true,
					},
					File: app.FileLogConfig{
						Enabled: false,
					},
				},
			},
			wantError: false,
		},
		{
			name: "logger with file output",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "info",
					Format: "text",
					Console: app.ConsoleConfig{
						Enabled: false,
					},
					File: app.FileLogConfig{
						Enabled:    true,
						Path:       tempDir,
						Filename:   "test.log",
						MaxSizeMB:  10,
						MaxBackups: 3,
						MaxAgeDays: 7,
					},
				},
			},
			wantError: false,
		},
		{
			name: "logger with both outputs",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "warn",
					Format: "json",
					Console: app.ConsoleConfig{
						Enabled: true,
					},
					File: app.FileLogConfig{
						Enabled:    true,
						Path:       tempDir,
						Filename:   "test.log",
						MaxSizeMB:  10,
						MaxBackups: 3,
						MaxAgeDays: 7,
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid log level",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "invalid",
					Format: "text",
					Console: app.ConsoleConfig{
						Enabled: true,
					},
				},
			},
			wantError: true,
			errorMsg:  "unknown log level",
		},
		{
			name: "invalid log format",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "info",
					Format: "invalid",
					Console: app.ConsoleConfig{
						Enabled: true,
					},
				},
			},
			wantError: true,
			errorMsg:  "unknown log format",
		},
		{
			name: "no outputs enabled",
			cfg: &app.Config{
				Log: app.LogConfig{
					Level:  "info",
					Format: "text",
					Console: app.ConsoleConfig{
						Enabled: false,
					},
					File: app.FileLogConfig{
						Enabled: false,
					},
				},
			},
			wantError: true,
			errorMsg:  "no log outputs enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.cfg)

			if tt.wantError {
				if err == nil {
					t.Errorf("New() expected error, got nil")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("error message = %q, want to contain %q", err.Error(), tt.errorMsg)
				}
				return
			}

			// Verify success case
			if err != nil {
				t.Errorf("New() unexpected error: %v", err)
			}
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

// TestFileWriter tests file logging setup.
func TestFileWriter(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &app.Config{
		Log: app.LogConfig{
			File: app.FileLogConfig{
				Enabled:    true,
				Path:       tempDir,
				Filename:   "test.log",
				MaxSizeMB:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
			},
		},
	}

	writer := newFileWriter(cfg)
	if writer == nil {
		t.Fatal("newFileWriter() returned nil writer")
	}

	// Write some data
	testData := []byte("test log entry\n")
	n, err := writer.Write(testData)
	if err != nil {
		t.Errorf("Write() error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Write() wrote %d bytes, want %d", n, len(testData))
	}

	// Verify file exists
	logPath := filepath.Join(tempDir, "test.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("log file not created at %s", logPath)
	}

	// Verify content
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !bytes.Equal(content, testData) {
		t.Errorf("log content = %q, want %q", string(content), string(testData))
	}
}

// TestGlobalLogger tests global logger functions.
func TestGlobalLogger(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "global.log")

	cfg := &app.Config{
		Log: app.LogConfig{
			Level:  "debug",
			Format: "text",
			Console: app.ConsoleConfig{
				Enabled: false,
			},
			File: app.FileLogConfig{
				Enabled:    true,
				Path:       tempDir,
				Filename:   "global.log",
				MaxSizeMB:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
			},
		},
	}

	// Initialize global logger
	if err := Init(cfg); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// Use global logger functions
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	output := string(content)
	expectedMessages := []string{"debug message", "info message", "warn message", "error message"}
	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("log output missing %q", msg)
		}
	}
}

// TestBufferedLogger tests the buffered logger for bootstrap phase.
func TestBufferedLogger(t *testing.T) {
	buffer := NewBuffered()

	// Log some messages
	buffer.Debug("debug msg", NewField("key", "value"))
	buffer.Info("info msg")
	buffer.Warn("warn msg")
	buffer.Error("error msg")

	// Check count
	if count := buffer.Count(); count != 4 {
		t.Errorf("buffer.Count() = %d, want 4", count)
	}

	// Create a test logger to replay to
	var output bytes.Buffer
	targetLogger := newTextLogger(LevelDebug, &output)

	// Replay
	buffer.ReplayTo(targetLogger)

	// Verify replayed logs
	logs := output.String()
	if !strings.Contains(logs, "debug msg") {
		t.Error("replayed logs missing 'debug msg'")
	}
	if !strings.Contains(logs, "info msg") {
		t.Error("replayed logs missing 'info msg'")
	}
	if !strings.Contains(logs, "warn msg") {
		t.Error("replayed logs missing 'warn msg'")
	}
	if !strings.Contains(logs, "error msg") {
		t.Error("replayed logs missing 'error msg'")
	}

	// Check for bootstrap marker
	if !strings.Contains(logs, "bootstrap=true") {
		t.Error("replayed logs missing bootstrap marker")
	}

	// Verify buffer is cleared after replay
	if count := buffer.Count(); count != 0 {
		t.Errorf("buffer.Count() after replay = %d, want 0", count)
	}
}

// TestBootstrapWorkflow tests the complete bootstrap workflow.
func TestBootstrapWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "bootstrap.log")

	// 1. Initialize bootstrap logger
	buffer := InitBootstrap()

	// 2. Log during bootstrap
	Info("bootstrap started")
	Debug("loading config", NewField("path", "/config"))
	Info("bootstrap complete")

	// Verify logs are buffered (including the debug log from InitBootstrap itself)
	if count := buffer.Count(); count != 4 {
		t.Errorf("buffer count = %d, want 4", count)
	}

	// 3. Upgrade to real logger
	cfg := &app.Config{
		Log: app.LogConfig{
			Level:  "debug",
			Format: "text",
			Console: app.ConsoleConfig{
				Enabled: false,
			},
			File: app.FileLogConfig{
				Enabled:    true,
				Path:       tempDir,
				Filename:   "bootstrap.log",
				MaxSizeMB:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
			},
		},
	}

	if err := UpgradeFromBootstrap(cfg, buffer); err != nil {
		t.Fatalf("UpgradeFromBootstrap() error: %v", err)
	}

	// 4. Log after upgrade
	Info("using real logger now")

	// 5. Verify all logs are in file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logs := string(content)
	expectedMessages := []string{
		"bootstrap started",
		"loading config",
		"bootstrap complete",
		"using real logger now",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(logs, msg) {
			t.Errorf("log file missing %q", msg)
		}
	}

	// Verify bootstrap marker exists for replayed logs
	if !strings.Contains(logs, "bootstrap=true") {
		t.Error("log file missing bootstrap marker")
	}
}

// TestLogLevelString tests the String() method for LogLevel.
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{"debug level string", LevelDebug, "debug"},
		{"info level string", LevelInfo, "info"},
		{"warn level string", LevelWarn, "warn"},
		{"error level string", LevelError, "error"},
		{"fatal level string", LevelFatal, "fatal"},
		{"unknown level string", LogLevel(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("LogLevel(%d).String() = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

// TestLogLevelUpperString tests the UpperString() method for LogLevel.
func TestLogLevelUpperString(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{"debug level uppercase", LevelDebug, "DEBUG"},
		{"info level uppercase", LevelInfo, "INFO"},
		{"warn level uppercase", LevelWarn, "WARN"},
		{"error level uppercase", LevelError, "ERROR"},
		{"fatal level uppercase", LevelFatal, "FATAL"},
		{"unknown level uppercase", LogLevel(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.UpperString()
			if got != tt.want {
				t.Errorf("LogLevel(%d).UpperString() = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}

// TestFormatConstants tests that format constants are correct.
func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{"text format constant", FormatText, "text"},
		{"json format constant", FormatJSON, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.want)
			}
		})
	}
}

// TestTextFormatConstants tests text format constants.
func TestTextFormatConstants(t *testing.T) {
	// Test that text timestamp format is valid
	var buf bytes.Buffer
	logger := newTextLogger(LevelInfo, &buf)
	logger.Info("test")

	output := buf.String()

	// Should contain timestamp in HH:MM:SS format
	if !strings.Contains(output, ":") {
		t.Error("text log output missing timestamp separator ':'")
	}

	// Should contain the field separator when fields are present
	buf.Reset()
	logger.Info("test", NewField("key", "value"))
	output = buf.String()

	if !strings.Contains(output, "|") {
		t.Error("text log output with fields missing separator '|'")
	}

	if !strings.Contains(output, "key=value") {
		t.Error("text log output missing key=value pair")
	}
}

// TestJSONFieldConstants tests JSON field name constants.
func TestJSONFieldConstants(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(LevelInfo, &buf)

	logger.Info("test message", NewField("custom", "field"))

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Verify required JSON fields exist using the constants
	requiredFields := []string{"time", "level", "msg"}
	for _, field := range requiredFields {
		if _, ok := entry[field]; !ok {
			t.Errorf("JSON log entry missing required field %q", field)
		}
	}

	// Verify custom field exists
	if entry["custom"] != "field" {
		t.Errorf("custom field = %v, want 'field'", entry["custom"])
	}
}

// TestBootstrapFieldConstants tests bootstrap field name constants.
func TestBootstrapFieldConstants(t *testing.T) {
	buffer := NewBuffered()
	buffer.Info("test message")

	var output bytes.Buffer
	targetLogger := newJSONLogger(LevelDebug, &output)

	buffer.ReplayTo(targetLogger)

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	// Verify bootstrap fields exist
	if entry["bootstrap"] != true {
		t.Errorf("bootstrap field = %v, want true", entry["bootstrap"])
	}

	if _, ok := entry["timestamp"]; !ok {
		t.Error("replayed log missing 'timestamp' field")
	}
}

// TestSentinelErrors tests that sentinel errors are correctly defined.
func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantContains string
	}{
		{"no log outputs error", ErrNoLogOutputs, "no log outputs enabled"},
		{"unknown format error", ErrUnknownLogFormat, "unknown log format"},
		{"unknown level error", ErrUnknownLogLevel, "unknown log level"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil, want error", tt.name)
				return
			}
			if !strings.Contains(tt.err.Error(), tt.wantContains) {
				t.Errorf("%s = %q, want to contain %q", tt.name, tt.err.Error(), tt.wantContains)
			}
		})
	}
}

// TestTextLoggerFatal tests the Fatal() method for text logger.
func TestTextLoggerFatal(t *testing.T) {
	// Save original exit function
	originalExit := osExit
	defer func() { osExit = originalExit }()

	// Mock exit function
	var exitCode int
	exitCalled := false
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
	}

	var buf bytes.Buffer
	logger := newTextLogger(LevelInfo, &buf)

	logger.Fatal("fatal error", NewField("error", "critical"))

	// Verify log was written
	output := buf.String()
	if !strings.Contains(output, "fatal error") {
		t.Error("Fatal log not written")
	}
	if !strings.Contains(output, "FATAL") {
		t.Error("Fatal log missing FATAL level")
	}
	if !strings.Contains(output, "error=critical") {
		t.Error("Fatal log missing field")
	}

	// Verify exit was called with code 1
	if !exitCalled {
		t.Error("osExit was not called")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// TestJSONLoggerFatal tests the Fatal() method for JSON logger.
func TestJSONLoggerFatal(t *testing.T) {
	// Save original exit function
	originalExit := osExit
	defer func() { osExit = originalExit }()

	// Mock exit function
	var exitCode int
	exitCalled := false
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
	}

	var buf bytes.Buffer
	logger := newJSONLogger(LevelInfo, &buf)

	logger.Fatal("fatal error", NewField("error", "critical"))

	// Verify log was written
	output := buf.String()
	if !strings.Contains(output, "fatal error") {
		t.Error("Fatal log not written")
	}

	// Parse JSON to verify structure
	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON log: %v", err)
	}

	if entry["level"] != "fatal" {
		t.Errorf("level = %v, want fatal", entry["level"])
	}
	if entry["msg"] != "fatal error" {
		t.Errorf("msg = %v, want 'fatal error'", entry["msg"])
	}
	if entry["error"] != "critical" {
		t.Errorf("error field = %v, want 'critical'", entry["error"])
	}

	// Verify exit was called with code 1
	if !exitCalled {
		t.Error("osExit was not called")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

// TestGlobalFatal tests the global Fatal() function.
func TestGlobalFatal(t *testing.T) {
	// Save original exit function
	originalExit := osExit
	defer func() { osExit = originalExit }()

	// Mock exit function
	var exitCode int
	exitCalled := false
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
	}

	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "fatal.log")

	cfg := &app.Config{
		Log: app.LogConfig{
			Level:  "info",
			Format: "text",
			Console: app.ConsoleConfig{
				Enabled: false,
			},
			File: app.FileLogConfig{
				Enabled:    true,
				Path:       tempDir,
				Filename:   "fatal.log",
				MaxSizeMB:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
			},
		},
	}

	// Initialize global logger
	if err := Init(cfg); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	// Call global Fatal
	Fatal("fatal message", NewField("code", 500))

	// Verify exit was called with code 1
	if !exitCalled {
		t.Error("osExit was not called")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	// Verify log was written to file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	output := string(content)
	if !strings.Contains(output, "fatal message") {
		t.Error("log file missing 'fatal message'")
	}
	if !strings.Contains(output, "code=500") {
		t.Error("log file missing field")
	}
}

// TestBufferedLoggerFatal tests the Fatal() method for buffered logger.
func TestBufferedLoggerFatal(t *testing.T) {
	// Save original exit function
	originalExit := osExit
	defer func() { osExit = originalExit }()

	// Mock exit function
	var exitCode int
	exitCalled := false
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
	}

	buffer := NewBuffered()

	// Log a fatal message to buffer
	buffer.Fatal("fatal during bootstrap", NewField("stage", "init"))

	// Verify exit was NOT called (buffered logger doesn't exit)
	if exitCalled {
		t.Error("buffered logger should not call osExit")
	}

	// Verify message was buffered
	if count := buffer.Count(); count != 1 {
		t.Errorf("buffer.Count() = %d, want 1", count)
	}

	// Now replay to a real logger and verify it exits
	var output bytes.Buffer
	targetLogger := newTextLogger(LevelDebug, &output)

	buffer.ReplayTo(targetLogger)

	// Now exit should have been called during replay
	if !exitCalled {
		t.Error("osExit should be called during replay of Fatal log")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}

	// Verify the fatal message was replayed
	logs := output.String()
	if !strings.Contains(logs, "fatal during bootstrap") {
		t.Error("replayed logs missing 'fatal during bootstrap'")
	}
	if !strings.Contains(logs, "stage=init") {
		t.Error("replayed logs missing field")
	}
}
