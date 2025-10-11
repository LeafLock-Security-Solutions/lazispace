package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

// Sentinel errors for configuration validation.
var (
	ErrEmptyAppName       = errors.New("app name cannot be empty")
	ErrMissingConfigDir   = errors.New("when storage.useXDG is false, configDir must be explicitly set")
	ErrMissingDataDir     = errors.New("when storage.useXDG is false, dataDir must be explicitly set")
	ErrInvalidLogLevel    = errors.New("invalid log level")
	ErrInvalidLogFormat   = errors.New("invalid log format")
	ErrInvalidMaxSizeMB   = errors.New("log file maxSizeMB must be positive")
	ErrNegativeMaxBackups = errors.New("log file maxBackups cannot be negative")
	ErrNegativeMaxAgeDays = errors.New("log file maxAgeDays cannot be negative")
	ErrEmptyLogFilename   = errors.New("log file filename cannot be empty")
	ErrNoLogOutputEnabled = errors.New("at least one log output (console or file) must be enabled")
)

// --- Application & Installation Settings ---
// These constants define the application's identity and default folder names.
const (
	installationFolderName = ".lazispace" // Using a dot prefix to hide it on UNIX-like systems
	logFolderName          = "logs"
)

// --- Environment Configuration ---
// Environment variable used to select environment-specific configuration.
// If not set, only the base application.yaml is loaded.
// If set to a value (e.g., "prod", "dev", "test"), application.{env}.yaml is merged.
const envVarName = "LSPACE_ENV"

// --- Viper Configuration ---
// These constants configure the behavior of the Viper library, which handles
// reading configuration from files and environment variables.
const (
	viperConfigName  = "application"
	viperConfigType  = "yaml"
	viperEnvPrefix   = "LSPACE"
	viperConfigPath1 = "./configs" // First search path
	viperConfigPath2 = "."         // Second (fallback) search path
)

// --- Filesystem Permissions ---
// Defines standard directory permissions used when creating them.
const dirPermission = 0o755

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

// LoadConfig initializes and loads the complete application configuration.
//
// Configuration Loading Strategy:
//  1. Reads base configuration from application.yaml (if it exists)
//  2. If LSPACE_ENV is set (e.g., "dev", "prod", "test"), merges application.{env}.yaml
//  3. Environment variables override config file values (prefix: LSPACE_)
//
// Configuration Search Paths:
//   - ./configs/application.yaml (primary)
//   - ./application.yaml (fallback)
//
// Environment Variable Overrides:
// Environment variables can override any config value using the LSPACE_ prefix.
// Nested values use underscores. Examples:
//   - LSPACE_LOG_LEVEL overrides log.level
//   - LSPACE_LOG_FORMAT overrides log.format
//   - LSPACE_STORAGE_USEXDG overrides storage.useXDG
//
// Post-Processing:
// After loading configuration, this function:
//   - Validates all configuration values (log levels, file settings, etc.)
//   - Resolves storage paths (XDG or project-relative)
//   - Creates necessary directories (config, data, logs)
//   - Expands ~ to user home directory in paths
//
// Returns:
//   - Fully initialized Config with resolved paths and created directories
//   - Error if config files are malformed, validation fails, or directories cannot be created
func LoadConfig() (*Config, error) {
	// Get environment from LSPACE_ENV (e.g., "dev", "prod", "test")
	env := GetEnvironment()

	// Set up Viper
	v := viper.New()

	// --- Configure Viper to Find the Config File ---
	v.SetConfigName(viperConfigName)
	v.SetConfigType(viperConfigType)
	v.AddConfigPath(viperConfigPath1)
	v.AddConfigPath(viperConfigPath2)

	// Read the base config file (e.g., application.yaml). It's okay if it doesn't exist.
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			// A real error occurred (e.g., malformed YAML), so we should fail.
			return nil, fmt.Errorf("failed to read base config: %w", err)
		}
		// It's just a "file not found" error, which is fine. Defaults will be used.
	}

	// Merge environment-specific config file if environment is set (e.g., application.prod.yaml).
	if env != "" {
		v.SetConfigName(fmt.Sprintf("%s.%s", viperConfigName, env))
		if err := v.MergeInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileNotFoundError) {
				return nil, fmt.Errorf("failed to read environment config '%s': %w", env, err)
			}
		}
	}

	// --- Configure Environment Variable Overrides ---
	// Enable automatic environment variable binding with proper key replacement.
	// This allows LSPACE_LOG_LEVEL to override "log.level", LSPACE_STORAGE_USEXDG to override "storage.useXDG", etc.
	v.SetEnvPrefix(viperEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// --- Unmarshal into Struct ---
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// --- Resolve and Create Storage Paths ---
	if err := resolveStoragePaths(&cfg); err != nil {
		return nil, fmt.Errorf("failed to resolve storage paths: %w", err)
	}

	// Resolve log path
	if err := resolveLogPath(&cfg); err != nil {
		return nil, fmt.Errorf("failed to resolve log path: %w", err)
	}

	return &cfg, nil
}

// getEnvironment retrieves the current environment name from the LSPACE_ENV variable.
//
// Returns:
//   - Environment name (e.g., "dev", "prod", "test") if LSPACE_ENV is set
//   - Empty string if LSPACE_ENV is not set (only base application.yaml will be loaded)
//
// The environment name is used to load environment-specific configuration files.
// For example, LSPACE_ENV=prod will load application.prod.yaml after application.yaml.
func getEnvironment() string {
	return os.Getenv(envVarName)
}

// resolveStoragePaths resolves and creates storage directories for config and data.
//
// Path Resolution Strategy:
//
// When UseXDG is true (production/user installations):
//   - ConfigDir defaults to XDG ConfigHome (~/.config/lazispace on Linux/macOS)
//   - DataDir defaults to XDG DataHome (~/.local/share/lazispace on Linux/macOS)
//   - Explicitly configured paths override XDG defaults
//
// When UseXDG is false (development/testing):
//   - Uses project-relative paths from configuration
//   - Converts relative paths to absolute paths
//
// Post-Processing:
//   - Expands ~ to user's home directory in all paths
//   - Creates directories if they don't exist (with 0755 permissions)
//
// Returns error if:
//   - Relative paths cannot be resolved to absolute paths
//   - Directories cannot be created
func resolveStoragePaths(cfg *Config) error {
	// Resolve paths based on UseXDG setting
	if err := resolveXDGOrLocalPaths(cfg); err != nil {
		return err
	}

	// Expand tilde `~` to the user's home directory if present.
	cfg.Storage.ConfigDir = expandHomeDir(cfg.Storage.ConfigDir)
	cfg.Storage.DataDir = expandHomeDir(cfg.Storage.DataDir)

	// Create directories if they don't exist
	if err := os.MkdirAll(cfg.Storage.ConfigDir, dirPermission); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.MkdirAll(cfg.Storage.DataDir, dirPermission); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	return nil
}

// resolveXDGOrLocalPaths resolves storage paths based on UseXDG setting.
func resolveXDGOrLocalPaths(cfg *Config) error {
	if cfg.Storage.UseXDG {
		return resolveXDGPaths(cfg)
	}
	return resolveLocalStoragePaths(cfg)
}

// resolveXDGPaths applies XDG Base Directory defaults if paths are empty.
func resolveXDGPaths(cfg *Config) error {
	if cfg.Storage.ConfigDir == "" {
		cfg.Storage.ConfigDir = filepath.Join(xdg.ConfigHome, installationFolderName)
	}
	if cfg.Storage.DataDir == "" {
		cfg.Storage.DataDir = filepath.Join(xdg.DataHome, installationFolderName)
	}
	return nil
}

// resolveLocalStoragePaths converts project-relative paths to absolute paths.
func resolveLocalStoragePaths(cfg *Config) error {
	var err error
	cfg.Storage.ConfigDir, err = resolveLocalPath(cfg.Storage.ConfigDir)
	if err != nil {
		return err
	}
	cfg.Storage.DataDir, err = resolveLocalPath(cfg.Storage.DataDir)
	if err != nil {
		return err
	}
	return nil
}

// resolveLogPath resolves and creates the log directory.
//
// Path Resolution:
//   - If Path is empty, defaults to XDG StateHome (~/.local/state/lazispace/logs)
//   - Expands ~ to user's home directory if present
//   - Converts relative paths to absolute paths
//   - Creates directory with 0755 permissions
//
// Returns:
//   - nil if file logging is disabled
//   - nil if directory is successfully created
//   - Error if path cannot be resolved or directory cannot be created
func resolveLogPath(cfg *Config) error {
	if !cfg.Log.File.Enabled {
		return nil
	}

	// If path is empty, use XDG StateHome
	if cfg.Log.File.Path == "" {
		// XDG State is for logs, cache, runtime data
		cfg.Log.File.Path = filepath.Join(xdg.StateHome, installationFolderName, logFolderName)
	}

	// Expand ~ if present
	cfg.Log.File.Path = expandHomeDir(cfg.Log.File.Path)

	// Make absolute if relative
	if !filepath.IsAbs(cfg.Log.File.Path) {
		absPath, err := filepath.Abs(cfg.Log.File.Path)
		if err != nil {
			return fmt.Errorf("failed to resolve log path: %w", err)
		}
		cfg.Log.File.Path = absPath
	}

	// Create directory
	if err := os.MkdirAll(cfg.Log.File.Path, dirPermission); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	return nil
}

// resolveLocalPath converts a project-relative path to an absolute path.
//
// This function is used when UseXDG is false (development/testing scenarios).
//
// Behavior:
//   - If path is empty, returns empty string
//   - If path is already absolute, returns it unchanged
//   - If path is relative, converts it to absolute based on current working directory
//
// Returns:
//   - Absolute path
//   - Error if absolute path cannot be determined
func resolveLocalPath(path string) (string, error) {
	if path != "" && !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path for '%s': %w", path, err)
		}
		return absPath, nil
	}
	return path, nil
}

// expandHomeDir expands the ~ character to the user's home directory.
//
// Examples:
//   - "~/config" becomes "/Users/username/config" (macOS)
//   - "~/config" becomes "/home/username/config" (Linux)
//   - "~/.config" becomes "/Users/username/.config" (macOS)
//   - "/absolute/path" remains "/absolute/path" (no tilde)
//
// Behavior:
//   - Empty path returns empty string
//   - Path without leading ~ is returned unchanged
//   - If home directory cannot be determined, returns original path
//
// Returns:
//   - Path with ~ expanded to home directory
//   - Original path if expansion not needed or fails
func expandHomeDir(path string) string {
	if path == "" {
		return ""
	}
	// Check if the path starts with '~'.
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			// If we can't get home dir, return path unmodified.
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// validateConfig validates all configuration values for correctness.
//
// Validation Rules:
//
// App Metadata:
//   - Name cannot be empty
//
// Storage Configuration:
//   - When UseXDG is false, both ConfigDir and DataDir must be explicitly set
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
// Returns:
//   - nil if all validations pass
//   - Error describing the first validation failure
func validateConfig(cfg *Config) error {
	// Validate app metadata
	if cfg.App.Name == "" {
		return ErrEmptyAppName
	}

	// Validate storage configuration
	if !cfg.Storage.UseXDG {
		if cfg.Storage.ConfigDir == "" {
			return ErrMissingConfigDir
		}
		if cfg.Storage.DataDir == "" {
			return ErrMissingDataDir
		}
	}

	// Validate log level
	validLevels := []string{"debug", "info", "warn", "error"}
	if !slices.Contains(validLevels, cfg.Log.Level) {
		return fmt.Errorf("%w: '%s', must be one of: %v", ErrInvalidLogLevel, cfg.Log.Level, validLevels)
	}

	// Validate log format
	validFormats := []string{"text", "json"}
	if !slices.Contains(validFormats, cfg.Log.Format) {
		return fmt.Errorf("%w: '%s', must be one of: %v", ErrInvalidLogFormat, cfg.Log.Format, validFormats)
	}

	// Validate file log settings if enabled
	if cfg.Log.File.Enabled {
		if cfg.Log.File.MaxSizeMB <= 0 {
			return fmt.Errorf("%w: got %d", ErrInvalidMaxSizeMB, cfg.Log.File.MaxSizeMB)
		}
		if cfg.Log.File.MaxBackups < 0 {
			return fmt.Errorf("%w: got %d", ErrNegativeMaxBackups, cfg.Log.File.MaxBackups)
		}
		// Note: maxBackups = 0 means no backups are kept (only current log file)

		if cfg.Log.File.MaxAgeDays < 0 {
			return fmt.Errorf("%w: got %d", ErrNegativeMaxAgeDays, cfg.Log.File.MaxAgeDays)
		}
		if cfg.Log.File.Filename == "" {
			return ErrEmptyLogFilename
		}
	}

	// Validate at least one output is enabled
	if !cfg.Log.Console.Enabled && !cfg.Log.File.Enabled {
		return ErrNoLogOutputEnabled
	}

	return nil
}

// GetEnvironment is the exported version of getEnvironment.
//
// Returns the current environment name from the LSPACE_ENV variable.
//
// Returns:
//   - Environment name (e.g., "dev", "prod", "test") if LSPACE_ENV is set
//   - Empty string if LSPACE_ENV is not set
//
// This function is useful for other packages that need to determine
// the current runtime environment.
func GetEnvironment() string {
	return getEnvironment()
}

// GetLogFilePath returns the complete path to the active log file.
//
// Returns:
//   - Full path to log file (e.g., "/Users/username/.local/state/lazispace/logs/lazispace.log")
//   - Empty string if file logging is disabled
//
// The returned path combines the log directory (Path) and filename (Filename)
// from the file logging configuration.
func (c *Config) GetLogFilePath() string {
	if !c.Log.File.Enabled {
		return ""
	}
	return filepath.Join(c.Log.File.Path, c.Log.File.Filename)
}
