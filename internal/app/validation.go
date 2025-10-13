package app

import (
	"fmt"
	"slices"
)

// validateConfig validates all configuration values for correctness.
//
// Validation Rules:
//
// App Metadata:
//   - Name cannot be empty
//
// Log Level:
//   - Must be one of: "debug", "info", "warn", "error"
//
// Log Format:
//   - Must be one of: "text", "json"
//
// File Logging (when enabled):
//   - Filename cannot be empty
//   - MaxSizeMB must be positive (> 0)
//   - MaxBackups cannot be negative (>= 0)
//   - MaxAgeDays cannot be negative (>= 0)
//   - Note: MaxBackups = 0 means no backups (only current log file)
//
// Output Configuration:
//   - At least one output (console or file) must be enabled
//
// Storage Configuration:
//   - When UseXDG is false, both ConfigDir and DataDir must be explicitly set
//
// Returns:
//   - nil if all validations pass
//   - Error describing the first validation failure
func validateConfig(cfg *Config) error {
	if err := validateAppConfig(&cfg.App); err != nil {
		return err
	}

	if err := validateLogConfig(&cfg.Log); err != nil {
		return err
	}

	return validateStorageConfig(&cfg.Storage)
}

// validateAppConfig validates application metadata.
func validateAppConfig(cfg *AppConfig) error {
	if cfg.Name == "" {
		return ErrEmptyAppName
	}
	return nil
}

// validateStorageConfig validates storage configuration.
func validateStorageConfig(cfg *StorageConfig) error {
	// When not using XDG, paths must be explicitly set
	if !cfg.UseXDG {
		if cfg.ConfigDir == "" {
			return ErrMissingConfigDir
		}
		if cfg.DataDir == "" {
			return ErrMissingDataDir
		}
	}
	return nil
}

// validateLogConfig validates logging configuration.
func validateLogConfig(cfg *LogConfig) error {
	// Validate log level
	if !slices.Contains(validLogLevels, cfg.Level) {
		return fmt.Errorf("%w: '%s', must be one of: %v", ErrInvalidLogLevel, cfg.Level, validLogLevels)
	}

	// Validate log format
	if !slices.Contains(validLogFormats, cfg.Format) {
		return fmt.Errorf("%w: '%s', must be one of: %v", ErrInvalidLogFormat, cfg.Format, validLogFormats)
	}

	// Validate file log settings if enabled
	if cfg.File.Enabled {
		if err := validateFileLogConfig(&cfg.File); err != nil {
			return err
		}
	}

	// Validate at least one output is enabled
	if !cfg.Console.Enabled && !cfg.File.Enabled {
		return ErrNoLogOutputEnabled
	}

	return nil
}

// validateFileLogConfig validates file logging configuration.
func validateFileLogConfig(cfg *FileLogConfig) error {
	if cfg.MaxSizeMB <= 0 {
		return fmt.Errorf("%w: got %d", ErrInvalidMaxSizeMB, cfg.MaxSizeMB)
	}

	if cfg.MaxBackups < 0 {
		return fmt.Errorf("%w: got %d", ErrNegativeMaxBackups, cfg.MaxBackups)
	}
	// Note: maxBackups = 0 means no backups are kept (only current log file)

	if cfg.MaxAgeDays < 0 {
		return fmt.Errorf("%w: got %d", ErrNegativeMaxAgeDays, cfg.MaxAgeDays)
	}

	if cfg.Filename == "" {
		return ErrEmptyLogFilename
	}

	return nil
}
