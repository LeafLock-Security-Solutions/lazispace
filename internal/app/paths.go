// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

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

	// Expand tilde `~` to the user's home directory if present
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
	// Check if the path starts with '~'
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			// If we can't get home dir, return path unmodified
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
