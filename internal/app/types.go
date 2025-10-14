// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package app

// Config holds the complete application configuration loaded from YAML files.
// Configuration is loaded from ./configs/ directory with support for environment-specific overrides.
type Config struct {
	// App contains application metadata like name and version.
	App AppConfig `mapstructure:"app"`

	// Log configures logging behavior including level, format, and outputs.
	Log LogConfig `mapstructure:"log"`

	// Storage configures where LaziSpace stores configuration and data files.
	Storage StorageConfig `mapstructure:"storage"`
}

// AppConfig holds application metadata.
type AppConfig struct {
	// Name is the application name (e.g., "lazispace").
	Name string `mapstructure:"name"`
}

// LogConfig configures the logging system including level, format, and output destinations.
type LogConfig struct {
	// Level sets the minimum log level: "debug", "info", "warn", or "error".
	Level string `mapstructure:"level"`

	// Format sets the log output format: "text" or "json".
	Format string `mapstructure:"format"`

	// Console configures console (stdout/stderr) logging.
	Console ConsoleConfig `mapstructure:"console"`

	// File configures file-based logging with rotation.
	File FileLogConfig `mapstructure:"file"`
}

// ConsoleConfig configures console (stdout/stderr) logging output.
type ConsoleConfig struct {
	// Enabled controls whether logs are written to console.
	Enabled bool `mapstructure:"enabled"`
}

// FileLogConfig configures file-based logging with automatic rotation.
type FileLogConfig struct {
	// Enabled controls whether logs are written to a file.
	Enabled bool `mapstructure:"enabled"`

	// Path is the directory where log files are stored.
	// If empty, defaults to XDG StateHome (e.g., ~/.local/state/lazispace/logs).
	Path string `mapstructure:"path"`

	// Filename is the name of the log file (e.g., "lazispace.log").
	Filename string `mapstructure:"filename"`

	// MaxSizeMB is the maximum size of a log file in megabytes before rotation.
	// Must be positive.
	MaxSizeMB int `mapstructure:"maxSizeMB"`

	// MaxBackups is the maximum number of old log files to retain.
	// 0 means no backups are kept (only current log file exists).
	MaxBackups int `mapstructure:"maxBackups"`

	// MaxAgeDays is the maximum number of days to retain old log files.
	// 0 means no age-based deletion.
	MaxAgeDays int `mapstructure:"maxAgeDays"`

	// Compress controls whether rotated log files are compressed with gzip.
	Compress bool `mapstructure:"compress"`
}

// StorageConfig configures where LaziSpace stores configuration and data files.
type StorageConfig struct {
	// UseXDG controls whether to use XDG Base Directory specification.
	// If true, uses XDG standard paths (e.g., ~/.config/lazispace, ~/.local/share/lazispace).
	// If false, uses project-relative or explicitly configured paths.
	UseXDG bool `mapstructure:"useXDG"`

	// ConfigDir is where workspace configurations are stored.
	// If empty and UseXDG is true, defaults to XDG ConfigHome (~/.config/lazispace).
	ConfigDir string `mapstructure:"configDir"`

	// DataDir is where application data is stored.
	// If empty and UseXDG is true, defaults to XDG DataHome (~/.local/share/lazispace).
	DataDir string `mapstructure:"dataDir"`
}
