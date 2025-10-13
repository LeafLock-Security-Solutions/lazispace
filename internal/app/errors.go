package app

import "errors"

// Sentinel errors for configuration validation.
var (
	// App configuration errors.
	ErrEmptyAppName = errors.New("app name cannot be empty")

	// Storage configuration errors.
	ErrMissingConfigDir = errors.New("when storage.useXDG is false, configDir must be explicitly set")
	ErrMissingDataDir   = errors.New("when storage.useXDG is false, dataDir must be explicitly set")

	// Log level and format errors.
	ErrInvalidLogLevel  = errors.New("invalid log level")
	ErrInvalidLogFormat = errors.New("invalid log format")

	// Log file configuration errors.
	ErrInvalidMaxSizeMB   = errors.New("log file maxSizeMB must be positive")
	ErrNegativeMaxBackups = errors.New("log file maxBackups cannot be negative")
	ErrNegativeMaxAgeDays = errors.New("log file maxAgeDays cannot be negative")
	ErrEmptyLogFilename   = errors.New("log file filename cannot be empty")

	// Log output errors.
	ErrNoLogOutputEnabled = errors.New("at least one log output (console or file) must be enabled")
)
