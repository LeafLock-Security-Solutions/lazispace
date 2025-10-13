package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/interfaces"
	"github.com/spf13/viper"
)

// LoadConfig initializes and loads the complete application configuration.
//
// Parameters:
//   - log: Logger for tracking configuration loading progress (required).
//     Typically, pass the bootstrap buffered logger to capture config loading logs.
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
func LoadConfig(log interfaces.Logger) (*Config, error) {
	log.Info("Starting configuration loading")

	// Get environment from LSPACE_ENV (e.g., "dev", "prod", "test")
	env := GetEnvironment()
	if env != "" {
		log.Debug("Environment detected", interfaces.Field{Key: "env", Value: env})
	} else {
		log.Debug("No environment specified, using base configuration only")
	}

	// Set up Viper and configure config file search paths
	v := viper.New()
	v.SetConfigName(viperConfigName)
	v.SetConfigType(viperConfigType)
	v.AddConfigPath(viperConfigPath1)
	v.AddConfigPath(viperConfigPath2)

	log.Debug("Config search paths configured", interfaces.Field{Key: "paths", Value: []string{viperConfigPath1, viperConfigPath2}})

	// Read the base config file (e.g., application.yaml). It's okay if it doesn't exist.
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			// A real error occurred (e.g., malformed YAML), so we should fail
			log.Error("Failed to read base config file", interfaces.Field{Key: "error", Value: err})
			return nil, fmt.Errorf("failed to read base config: %w", err)
		}
		// It's just a "file not found" error, which is fine. Defaults will be used.
		log.Debug("Base config file not found, using defaults")
	} else {
		log.Info("Base config file loaded", interfaces.Field{Key: "file", Value: v.ConfigFileUsed()})
	}

	// Merge environment-specific config file if environment is set (e.g., application.prod.yaml)
	if env != "" {
		v.SetConfigName(fmt.Sprintf("%s.%s", viperConfigName, env))
		if err := v.MergeInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if !errors.As(err, &configFileNotFoundError) {
				log.Error("Failed to read environment config", interfaces.Field{Key: "env", Value: env}, interfaces.Field{Key: "error", Value: err})
				return nil, fmt.Errorf("failed to read environment config '%s': %w", env, err)
			}
			log.Debug("Environment config file not found", interfaces.Field{Key: "env", Value: env})
		} else {
			log.Info("Environment config merged", interfaces.Field{Key: "env", Value: env}, interfaces.Field{Key: "file", Value: v.ConfigFileUsed()})
		}
	}

	// Configure environment variable overrides (LSPACE_ prefix)
	v.SetEnvPrefix(viperEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	log.Debug("Environment variable overrides enabled", interfaces.Field{Key: "prefix", Value: viperEnvPrefix})

	// Unmarshal configuration into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Error("Failed to unmarshal config into struct", interfaces.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	log.Debug("Config unmarshaled successfully")

	if err := validateConfig(&cfg); err != nil {
		log.Error("Config validation failed", interfaces.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	log.Debug("Config validation passed")

	// Resolve and create storage paths
	if err := resolveStoragePaths(&cfg); err != nil {
		log.Error("Failed to resolve storage paths", interfaces.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to resolve storage paths: %w", err)
	}

	log.Info("Storage paths resolved",
		interfaces.Field{Key: "configDir", Value: cfg.Storage.ConfigDir},
		interfaces.Field{Key: "dataDir", Value: cfg.Storage.DataDir})

	// Resolve log path
	if err := resolveLogPath(&cfg); err != nil {
		log.Error("Failed to resolve log path", interfaces.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to resolve log path: %w", err)
	}

	if cfg.Log.File.Enabled {
		log.Info("Log path resolved", interfaces.Field{Key: "path", Value: cfg.Log.File.Path})
	}

	log.Info("Configuration loaded successfully",
		interfaces.Field{Key: "app", Value: cfg.App.Name},
		interfaces.Field{Key: "logLevel", Value: cfg.Log.Level},
		interfaces.Field{Key: "logFormat", Value: cfg.Log.Format})

	return &cfg, nil
}

// GetEnvironment retrieves the current environment name from the LSPACE_ENV variable.
//
// Returns:
//   - Environment name (e.g., "dev", "prod", "test") if LSPACE_ENV is set
//   - Empty string if LSPACE_ENV is not set (only base application.yaml will be loaded)
//
// The environment name is used to load environment-specific configuration files.
// For example, LSPACE_ENV=prod will load application.prod.yaml after application.yaml.
//
// This function is useful for other packages that need to determine
// the current runtime environment.
func GetEnvironment() string {
	return os.Getenv(envVarName)
}
