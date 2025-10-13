package app

import "github.com/LeafLock-Security-Solutions/lazispace/internal/constants"

// Application and installation settings.
const (
	// InstallationFolderName is the folder name used for application data.
	// Using a dot prefix to hide it on UNIX-like systems.
	installationFolderName = ".lazispace"

	// LogFolderName is the subfolder name for log files.
	logFolderName = "logs"
)

// Environment configuration.
const (
	// EnvVarName is the environment variable used to select environment-specific configuration.
	// If not set, only the base application.yaml is loaded.
	// If set to a value (e.g., "prod", "dev", "test"), application.{env}.yaml is merged.
	envVarName = "LSPACE_ENV"
)

// Viper configuration constants.
// These constants configure the behavior of the Viper library, which handles
// reading configuration from files and environment variables.
const (
	viperConfigName  = "application" // Base name of config file (without extension)
	viperConfigType  = "yaml"        // Config file type
	viperEnvPrefix   = "LSPACE"      // Prefix for environment variable overrides
	viperConfigPath1 = "./configs"   // Primary search path for config files
	viperConfigPath2 = "."           // Fallback search path for config files
)

// Filesystem permissions.
const (
	// DirPermission defines standard directory permissions used when creating directories.
	dirPermission = 0o755
)

// Valid configuration values for validation.
// These reference the shared constants package to avoid duplication.
var (
	// ValidLogLevels references the shared constants for valid log levels.
	validLogLevels = constants.ValidLogLevels

	// ValidLogFormats references the shared constants for valid log formats.
	validLogFormats = constants.ValidLogFormats
)
