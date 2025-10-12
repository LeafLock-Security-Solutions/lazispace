package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/app"
)

// Logger defines the interface for logging operations.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value any
}

// NewField creates a new log field with the given key and value.
func NewField(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// New creates a logger based on config.
func New(cfg *app.Config) (Logger, error) {
	Debug("initializing logger", NewField("level", cfg.Log.Level), NewField("format", cfg.Log.Format))

	level, err := parseLevel(cfg.Log.Level)
	if err != nil {
		Error("failed to parse log level", NewField("level", cfg.Log.Level), NewField("error", err))
		return nil, err
	}

	// Determine outputs
	var writers []io.Writer
	var outputTypes []string

	// Console output
	if cfg.Log.Console.Enabled {
		writers = append(writers, os.Stdout)
		outputTypes = append(outputTypes, "console")
	}

	// File output
	if cfg.Log.File.Enabled {
		fileWriter := newFileWriter(cfg)
		writers = append(writers, fileWriter)
		outputTypes = append(outputTypes, "file")
	}

	if len(writers) == 0 {
		Error("no log outputs enabled in configuration")
		return nil, ErrNoLogOutputs
	}

	Debug("logger outputs configured", NewField("outputs", outputTypes))

	output := io.MultiWriter(writers...)

	switch cfg.Log.Format {
	case FormatText:
		Info("logger initialized successfully", NewField("format", "text"), NewField("level", level.String()))
		return newTextLogger(level, output), nil
	case FormatJSON:
		Info("logger initialized successfully", NewField("format", "json"), NewField("level", level.String()))
		return newJSONLogger(level, output), nil
	default:
		Error("unknown log format", NewField("format", cfg.Log.Format))
		return nil, fmt.Errorf("%w: %s", ErrUnknownLogFormat, cfg.Log.Format)
	}
}

func parseLevel(levelStr string) (LogLevel, error) {
	switch strings.ToLower(levelStr) {
	case levelDebugStr:
		return LevelDebug, nil
	case levelInfoStr:
		return LevelInfo, nil
	case levelWarnStr, levelWarningStr:
		return LevelWarn, nil
	case levelErrorStr:
		return LevelError, nil
	default:
		return LevelInfo, fmt.Errorf("%w: %s", ErrUnknownLogLevel, levelStr)
	}
}

var (
	global   Logger
	globalMu sync.RWMutex
)

// Init initializes the global logger with the given configuration.
func Init(cfg *app.Config) error {
	Debug("initializing global logger")

	logger, err := New(cfg)
	if err != nil {
		Error("failed to initialize global logger", NewField("error", err))
		return err
	}

	globalMu.Lock()
	global = logger
	globalMu.Unlock()

	Info("global logger initialized")
	return nil
}

// Debug logs a debug message using the global logger.
func Debug(msg string, fields ...Field) {
	globalMu.RLock()
	logger := global
	globalMu.RUnlock()

	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Info logs an info message using the global logger.
func Info(msg string, fields ...Field) {
	globalMu.RLock()
	logger := global
	globalMu.RUnlock()

	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Warn logs a warning message using the global logger.
func Warn(msg string, fields ...Field) {
	globalMu.RLock()
	logger := global
	globalMu.RUnlock()

	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Error logs an error message using the global logger.
func Error(msg string, fields ...Field) {
	globalMu.RLock()
	logger := global
	globalMu.RUnlock()

	if logger != nil {
		logger.Error(msg, fields...)
	}
}

// Fatal logs a fatal message using the global logger and exits.
func Fatal(msg string, fields ...Field) {
	globalMu.RLock()
	logger := global
	globalMu.RUnlock()

	if logger != nil {
		logger.Fatal(msg, fields...)
	}
	os.Exit(1)
}

// BufferedLogger is a logger that buffers log entries in memory until replayed.
type BufferedLogger interface {
	Logger
	ReplayTo(target Logger)
	Count() int
}

// InitBootstrap initializes a buffered logger for bootstrap phase.
// Returns the buffered logger so it can be replayed later.
func InitBootstrap() BufferedLogger {
	buffer := NewBuffered()

	globalMu.Lock()
	global = buffer
	globalMu.Unlock()

	Debug("bootstrap logger initialized")
	return buffer
}

// UpgradeFromBootstrap initializes the real logger and replays buffered logs.
func UpgradeFromBootstrap(cfg *app.Config, buffer BufferedLogger) error {
	bufferedCount := 0
	if buffer != nil {
		bufferedCount = buffer.Count()
	}
	Debug("upgrading from bootstrap logger", NewField("buffered_logs", bufferedCount))

	// Create the real logger
	realLogger, err := New(cfg)
	if err != nil {
		Error("failed to create real logger during upgrade", NewField("error", err))
		return err
	}

	// Replay buffered logs to real logger
	if buffer != nil {
		buffer.ReplayTo(realLogger)
	}

	// Switch to real logger
	globalMu.Lock()
	global = realLogger
	globalMu.Unlock()

	Info("upgraded from bootstrap logger", NewField("replayed_logs", bufferedCount))
	return nil
}
