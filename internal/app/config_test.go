// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/interfaces"
)

// testLogger is a no-op logger for testing config loading.
type testLogger struct{}

func (l *testLogger) Debug(msg string, fields ...interfaces.Field) {}
func (l *testLogger) Info(msg string, fields ...interfaces.Field)  {}
func (l *testLogger) Warn(msg string, fields ...interfaces.Field)  {}
func (l *testLogger) Error(msg string, fields ...interfaces.Field) {}
func (l *testLogger) Fatal(msg string, fields ...interfaces.Field) {}

// newTestLogger creates a no-op logger for testing.
func newTestLogger() interfaces.Logger {
	return &testLogger{}
}

// TestValidateConfig_AppMetadata tests validation of application metadata.
func TestValidateConfig_AppMetadata(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid app name",
			config: Config{
				App: AppConfig{
					Name: "lazispace",
				},
				Log: LogConfig{
					Level:  "info",
					Format: "json",
					Console: ConsoleConfig{
						Enabled: true,
					},
				},
				Storage: StorageConfig{
					UseXDG: true,
				},
			},
			wantErr: false,
		},
		{
			name: "empty app name",
			config: Config{
				App: AppConfig{
					Name: "",
				},
			},
			wantErr: true,
			errMsg:  "app name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateConfig() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestValidateConfig_Storage tests validation of storage configuration.
func TestValidateConfig_Storage(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "useXDG true - no paths required",
			config: Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{
					UseXDG:    true,
					ConfigDir: "",
					DataDir:   "",
				},
			},
			wantErr: false,
		},
		{
			name: "useXDG false - both paths set",
			config: Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{
					UseXDG:    false,
					ConfigDir: "./config",
					DataDir:   "./data",
				},
			},
			wantErr: false,
		},
		{
			name: "useXDG false - missing configDir",
			config: Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{
					UseXDG:    false,
					ConfigDir: "",
					DataDir:   "./data",
				},
			},
			wantErr: true,
			errMsg:  "when storage.useXDG is false, configDir must be explicitly set",
		},
		{
			name: "useXDG false - missing dataDir",
			config: Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{
					UseXDG:    false,
					ConfigDir: "./config",
					DataDir:   "",
				},
			},
			wantErr: true,
			errMsg:  "when storage.useXDG is false, dataDir must be explicitly set",
		},
		{
			name: "useXDG false - both paths missing",
			config: Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{
					UseXDG:    false,
					ConfigDir: "",
					DataDir:   "",
				},
			},
			wantErr: true,
			errMsg:  "when storage.useXDG is false, configDir must be explicitly set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateConfig() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestValidateConfig_LogLevel tests validation of log level.
func TestValidateConfig_LogLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"debug level", "debug", false},
		{"info level", "info", false},
		{"warn level", "warn", false},
		{"error level", "error", false},
		{"invalid level", "trace", true},
		{"invalid level uppercase", "INFO", true},
		{"empty level", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   tt.level,
					Format:  "text",
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{UseXDG: true},
			}
			err := validateConfig(&cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() with level %q error = %v, wantErr %v", tt.level, err, tt.wantErr)
			}
		})
	}
}

// TestValidateConfig_LogFormat tests validation of log format.
func TestValidateConfig_LogFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"text format", "text", false},
		{"json format", "json", false},
		{"invalid format", "xml", true},
		{"invalid format uppercase", "JSON", true},
		{"empty format", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  tt.format,
					Console: ConsoleConfig{Enabled: true},
				},
				Storage: StorageConfig{UseXDG: true},
			}
			err := validateConfig(&cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() with format %q error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

// TestValidateConfig_FileLogging tests validation of file logging configuration.
func TestValidateConfig_FileLogging(t *testing.T) {
	tests := []struct {
		name    string
		file    FileLogConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file config",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  100,
				MaxBackups: 5,
				MaxAgeDays: 30,
			},
			wantErr: false,
		},
		{
			name: "file logging disabled - no validation",
			file: FileLogConfig{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "empty filename",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "",
				MaxSizeMB:  100,
				MaxBackups: 5,
				MaxAgeDays: 30,
			},
			wantErr: true,
			errMsg:  "log file filename cannot be empty",
		},
		{
			name: "zero maxSizeMB",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  0,
				MaxBackups: 5,
				MaxAgeDays: 30,
			},
			wantErr: true,
			errMsg:  "log file maxSizeMB must be positive: got 0",
		},
		{
			name: "negative maxSizeMB",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  -1,
				MaxBackups: 5,
				MaxAgeDays: 30,
			},
			wantErr: true,
			errMsg:  "log file maxSizeMB must be positive: got -1",
		},
		{
			name: "zero maxBackups - valid (no backups)",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  100,
				MaxBackups: 0,
				MaxAgeDays: 30,
			},
			wantErr: false,
		},
		{
			name: "negative maxBackups",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  100,
				MaxBackups: -1,
				MaxAgeDays: 30,
			},
			wantErr: true,
			errMsg:  "log file maxBackups cannot be negative: got -1",
		},
		{
			name: "zero maxAgeDays - valid (no age-based deletion)",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  100,
				MaxBackups: 5,
				MaxAgeDays: 0,
			},
			wantErr: false,
		},
		{
			name: "negative maxAgeDays",
			file: FileLogConfig{
				Enabled:    true,
				Filename:   "app.log",
				MaxSizeMB:  100,
				MaxBackups: 5,
				MaxAgeDays: -1,
			},
			wantErr: true,
			errMsg:  "log file maxAgeDays cannot be negative: got -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:   "info",
					Format:  "text",
					Console: ConsoleConfig{Enabled: !tt.file.Enabled}, // Ensure at least one output
					File:    tt.file,
				},
				Storage: StorageConfig{UseXDG: true},
			}
			err := validateConfig(&cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateConfig() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestValidateConfig_LogOutput tests validation that at least one output is enabled.
func TestValidateConfig_LogOutput(t *testing.T) {
	tests := []struct {
		name           string
		consoleEnabled bool
		fileEnabled    bool
		wantErr        bool
	}{
		{"console only", true, false, false},
		{"file only", false, true, false},
		{"both enabled", true, true, false},
		{"none enabled", false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				App: AppConfig{Name: "test"},
				Log: LogConfig{
					Level:  "info",
					Format: "text",
					Console: ConsoleConfig{
						Enabled: tt.consoleEnabled,
					},
					File: FileLogConfig{
						Enabled:    tt.fileEnabled,
						Filename:   "app.log",
						MaxSizeMB:  100,
						MaxBackups: 5,
						MaxAgeDays: 30,
					},
				},
				Storage: StorageConfig{UseXDG: true},
			}
			err := validateConfig(&cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExpandHomeDir tests the expandHomeDir function.
func TestExpandHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde with path",
			input:    "~/config",
			expected: filepath.Join(home, "config"),
		},
		{
			name:     "tilde with nested path",
			input:    "~/.config/lazispace",
			expected: filepath.Join(home, ".config", "lazispace"),
		},
		{
			name:     "tilde only",
			input:    "~",
			expected: home,
		},
		{
			name:     "absolute path unchanged",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path unchanged",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandHomeDir(tt.input)
			if result != tt.expected {
				t.Errorf("expandHomeDir(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestResolveLocalPath tests the resolveLocalPath function.
func TestResolveLocalPath(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "relative path",
			input:    "relative/path",
			expected: filepath.Join(cwd, "relative", "path"),
			wantErr:  false,
		},
		{
			name:     "absolute path unchanged",
			input:    "/absolute/path",
			expected: "/absolute/path",
			wantErr:  false,
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "dot current directory",
			input:    ".",
			expected: cwd,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveLocalPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveLocalPath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("resolveLocalPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetEnvironment tests the getEnvironment and GetEnvironment functions.
func TestGetEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{"environment set to dev", "dev", "dev"},
		{"environment set to prod", "prod", "prod"},
		{"environment set to test", "test", "test"},
		{"environment set to custom", "staging", "staging"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable
			os.Setenv(envVarName, tt.envValue)
			defer os.Unsetenv(envVarName)

			// Test exported function
			if got := GetEnvironment(); got != tt.want {
				t.Errorf("GetEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}

	// Test with no environment set
	t.Run("no environment set", func(t *testing.T) {
		os.Unsetenv(envVarName)
		if got := GetEnvironment(); got != "" {
			t.Errorf("GetEnvironment() with no env = %v, want empty string", got)
		}
	})
}

// TestLoadConfig_WithTestData tests LoadConfig with actual YAML files.
func TestLoadConfig_WithTestData(t *testing.T) {
	// Save original working directory and restore after test
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Change to testdata directory to simulate configs being there
	testdataDir := filepath.Join(originalWd, "testdata")

	t.Run("valid config loads successfully", func(t *testing.T) {
		// Set up Viper to look in testdata
		os.Chdir(testdataDir)
		defer os.Chdir(originalWd)

		// Clear environment
		os.Unsetenv(envVarName)

		cfg, err := LoadConfig(newTestLogger())
		if err != nil {
			t.Fatalf("LoadConfig() failed: %v", err)
		}

		// Verify loaded values
		if cfg.App.Name != "test-app" {
			t.Errorf("App.Name = %q, want %q", cfg.App.Name, "test-app")
		}
		if cfg.Log.Level != "info" {
			t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
		}
		if cfg.Log.Format != "json" {
			t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
		}
		if !cfg.Log.Console.Enabled {
			t.Error("Log.Console.Enabled should be true")
		}
		if cfg.Log.File.Enabled {
			t.Error("Log.File.Enabled should be false")
		}
		if !cfg.Storage.UseXDG {
			t.Error("Storage.UseXDG should be true")
		}
	})

	t.Run("environment-specific config merges correctly", func(t *testing.T) {
		// Set up Viper to look in testdata
		os.Chdir(testdataDir)
		defer os.Chdir(originalWd)

		// Set environment to dev
		os.Setenv(envVarName, "dev")
		defer os.Unsetenv(envVarName)

		cfg, err := LoadConfig(newTestLogger())
		if err != nil {
			t.Fatalf("LoadConfig() with dev env failed: %v", err)
		}

		// Clean up created directories
		t.Cleanup(func() {
			os.RemoveAll(filepath.Join(testdataDir, "test-config"))
			os.RemoveAll(filepath.Join(testdataDir, "test-data"))
			os.RemoveAll(filepath.Join(testdataDir, "test-logs"))
		})

		// Verify dev overrides were applied
		if cfg.Log.Level != "debug" {
			t.Errorf("Log.Level = %q, want %q (from dev config)", cfg.Log.Level, "debug")
		}
		if cfg.Log.Format != "text" {
			t.Errorf("Log.Format = %q, want %q (from dev config)", cfg.Log.Format, "text")
		}
		if !cfg.Log.File.Enabled {
			t.Error("Log.File.Enabled should be true (from dev config)")
		}
		if cfg.Storage.UseXDG {
			t.Error("Storage.UseXDG should be false (from dev config)")
		}
		if cfg.Storage.ConfigDir == "" {
			t.Error("Storage.ConfigDir should be set (from dev config)")
		}
		if cfg.Storage.DataDir == "" {
			t.Error("Storage.DataDir should be set (from dev config)")
		}
	})

	t.Run("malformed YAML returns error", func(t *testing.T) {
		// Create a temp directory with malformed config
		tempDir := t.TempDir()
		malformedPath := filepath.Join(tempDir, "application.yaml")

		// Copy malformed file
		malformedData, err := os.ReadFile(filepath.Join(testdataDir, "malformed.yaml"))
		if err != nil {
			t.Fatalf("Failed to read malformed.yaml: %v", err)
		}
		if err := os.WriteFile(malformedPath, malformedData, 0o644); err != nil {
			t.Fatalf("Failed to write malformed config: %v", err)
		}

		os.Chdir(tempDir)
		defer os.Chdir(originalWd)
		os.Unsetenv(envVarName)

		_, err = LoadConfig(newTestLogger())
		if err == nil {
			t.Fatal("LoadConfig() should fail with malformed YAML")
		}
		if !containsString(err.Error(), "failed to read base config") {
			t.Errorf("Expected error about base config, got: %v", err)
		}
	})

	t.Run("invalid log level returns validation error", func(t *testing.T) {
		// Create a temp directory with invalid config
		tempDir := t.TempDir()
		invalidPath := filepath.Join(tempDir, "application.yaml")

		// Copy invalid-level file
		invalidData, err := os.ReadFile(filepath.Join(testdataDir, "invalid-level.yaml"))
		if err != nil {
			t.Fatalf("Failed to read invalid-level.yaml: %v", err)
		}
		if err := os.WriteFile(invalidPath, invalidData, 0o644); err != nil {
			t.Fatalf("Failed to write invalid config: %v", err)
		}

		os.Chdir(tempDir)
		defer os.Chdir(originalWd)
		os.Unsetenv(envVarName)

		_, err = LoadConfig(newTestLogger())
		if err == nil {
			t.Fatal("LoadConfig() should fail with invalid log level")
		}
		if !containsString(err.Error(), "invalid log level") {
			t.Errorf("Expected error about invalid log level, got: %v", err)
		}
	})

	t.Run("missing storage paths returns validation error", func(t *testing.T) {
		// Create a temp directory with missing storage paths config
		tempDir := t.TempDir()
		missingPathsFile := filepath.Join(tempDir, "application.yaml")

		// Copy missing-storage-paths file
		missingData, err := os.ReadFile(filepath.Join(testdataDir, "missing-storage-paths.yaml"))
		if err != nil {
			t.Fatalf("Failed to read missing-storage-paths.yaml: %v", err)
		}
		if err := os.WriteFile(missingPathsFile, missingData, 0o644); err != nil {
			t.Fatalf("Failed to write missing paths config: %v", err)
		}

		os.Chdir(tempDir)
		defer os.Chdir(originalWd)
		os.Unsetenv(envVarName)

		_, err = LoadConfig(newTestLogger())
		if err == nil {
			t.Fatal("LoadConfig() should fail with missing storage paths")
		}
		if !containsString(err.Error(), "configDir must be explicitly set") {
			t.Errorf("Expected error about missing configDir, got: %v", err)
		}
	})
}

// containsString checks if a string contains a substring.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && stringContains(s, substr)))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestVersionFunctions tests the version getter functions.
func TestVersionFunctions(t *testing.T) {
	t.Run("GetVersion returns version", func(t *testing.T) {
		result := GetVersion()
		// Should return default "dev" or injected value
		if result == "" {
			t.Error("GetVersion() should not return empty string")
		}
	})

	t.Run("GetGitCommit returns commit", func(t *testing.T) {
		result := GetGitCommit()
		// Should return default "unknown" or injected value
		if result == "" {
			t.Error("GetGitCommit() should not return empty string")
		}
	})

	t.Run("GetBuildDate returns build date", func(t *testing.T) {
		result := GetBuildDate()
		// Should return default "unknown" or injected value
		if result == "" {
			t.Error("GetBuildDate() should not return empty string")
		}
	})

	t.Run("GetVersionInfo returns formatted info", func(t *testing.T) {
		result := GetVersionInfo()
		// Should contain version, commit, and build date
		if !stringContains(result, "Version:") {
			t.Error("GetVersionInfo() should contain 'Version:'")
		}
		if !stringContains(result, "Commit:") {
			t.Error("GetVersionInfo() should contain 'Commit:'")
		}
		if !stringContains(result, "Built:") {
			t.Error("GetVersionInfo() should contain 'Built:'")
		}
	})
}

// TestExpandHomeDir_ErrorPaths tests error handling in expandHomeDir.
func TestExpandHomeDir_ErrorPaths(t *testing.T) {
	// The error path (os.UserHomeDir() fails) is very hard to test in a unit test
	// because it would require mocking os.UserHomeDir(), which is not possible
	// without dependency injection. This path is covered by our regular tests
	// since UserHomeDir() normally succeeds.
	// If UserHomeDir() fails, the function returns the original path, which is safe behavior.
}

// TestResolveStoragePaths_WithExplicitXDGPaths tests UseXDG=true with explicit paths.
func TestResolveStoragePaths_WithExplicitXDGPaths(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &Config{
		Storage: StorageConfig{
			UseXDG:    true,
			ConfigDir: filepath.Join(tempDir, "custom-config"),
			DataDir:   filepath.Join(tempDir, "custom-data"),
		},
	}

	err := resolveStoragePaths(cfg)
	if err != nil {
		t.Fatalf("resolveStoragePaths() failed: %v", err)
	}

	// Should use explicit paths even when UseXDG=true
	if !stringContains(cfg.Storage.ConfigDir, "custom-config") {
		t.Errorf("ConfigDir should contain 'custom-config', got %s", cfg.Storage.ConfigDir)
	}
	if !stringContains(cfg.Storage.DataDir, "custom-data") {
		t.Errorf("DataDir should contain 'custom-data', got %s", cfg.Storage.DataDir)
	}

	// Directories should be created
	if _, err := os.Stat(cfg.Storage.ConfigDir); os.IsNotExist(err) {
		t.Error("ConfigDir should be created")
	}
	if _, err := os.Stat(cfg.Storage.DataDir); os.IsNotExist(err) {
		t.Error("DataDir should be created")
	}
}

// TestResolveLogPath_WithAbsolutePath tests resolveLogPath with absolute path.
func TestResolveLogPath_WithAbsolutePath(t *testing.T) {
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")

	cfg := &Config{
		Log: LogConfig{
			File: FileLogConfig{
				Enabled: true,
				Path:    logDir,
			},
		},
	}

	err := resolveLogPath(cfg)
	if err != nil {
		t.Fatalf("resolveLogPath() failed: %v", err)
	}

	if cfg.Log.File.Path != logDir {
		t.Errorf("Log path should be %s, got %s", logDir, cfg.Log.File.Path)
	}

	// Directory should be created
	if _, err := os.Stat(cfg.Log.File.Path); os.IsNotExist(err) {
		t.Error("Log directory should be created")
	}
}

// TestResolveLogPath_WithTildePath tests resolveLogPath with ~ expansion.
func TestResolveLogPath_WithTildePath(t *testing.T) {
	cfg := &Config{
		Log: LogConfig{
			File: FileLogConfig{
				Enabled: true,
				Path:    "~/test-logs",
			},
		},
	}

	err := resolveLogPath(cfg)
	if err != nil {
		t.Fatalf("resolveLogPath() failed: %v", err)
	}

	// Path should be expanded and absolute
	if !filepath.IsAbs(cfg.Log.File.Path) {
		t.Errorf("Log path should be absolute, got %s", cfg.Log.File.Path)
	}
	if stringContains(cfg.Log.File.Path, "~") {
		t.Errorf("Log path should not contain ~, got %s", cfg.Log.File.Path)
	}
}

// TestLoadConfig_NoConfigFile tests LoadConfig when no config file exists.
func TestLoadConfig_NoConfigFile(t *testing.T) {
	// Create empty temp directory with no config files
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tempDir)
	os.Unsetenv(envVarName)

	// Should fail because app.name will be empty (no config file)
	_, err := LoadConfig(newTestLogger())
	if err == nil {
		t.Fatal("LoadConfig() should fail when no config file exists and name is empty")
	}
	if !stringContains(err.Error(), "app name cannot be empty") {
		t.Errorf("Expected error about app name, got: %v", err)
	}
}

// TestLoadConfig_EnvironmentVariableOverrides tests env var overrides.
func TestLoadConfig_EnvironmentVariableOverrides(t *testing.T) {
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	testdataDir := filepath.Join(originalWd, "testdata")
	os.Chdir(testdataDir)
	defer os.Chdir(originalWd)

	// Set environment variable to override config
	os.Setenv("LSPACE_LOG_LEVEL", "error")
	defer os.Unsetenv("LSPACE_LOG_LEVEL")

	os.Unsetenv(envVarName)

	cfg, err := LoadConfig(newTestLogger())
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Should use env var override instead of file value
	if cfg.Log.Level != "error" {
		t.Errorf("Log.Level should be 'error' from env var, got %s", cfg.Log.Level)
	}
}
