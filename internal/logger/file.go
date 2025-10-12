package logger

import (
	"io"
	"path/filepath"

	"github.com/LeafLock-Security-Solutions/lazispace/internal/app"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newFileWriter creates an io.Writer that writes to a rotating log file.
//
// The file writer uses lumberjack for automatic log rotation based on:
//   - Size: Rotates when file reaches MaxSizeMB megabytes
//   - Age: Automatically deletes files older than MaxAgeDays days (0 = never delete)
//   - Backups: Keeps at most MaxBackups old log files (0 = keep all)
//   - Compression: Optionally gzip old log files when Compress is true
//
// The log file path is constructed by joining cfg.Log.File.Path and cfg.Log.File.Filename.
func newFileWriter(cfg *app.Config) io.Writer {
	logPath := filepath.Join(cfg.Log.File.Path, cfg.Log.File.Filename)

	Debug("configuring file writer",
		NewField("path", logPath),
		NewField("max_size_mb", cfg.Log.File.MaxSizeMB),
		NewField("max_backups", cfg.Log.File.MaxBackups),
		NewField("max_age_days", cfg.Log.File.MaxAgeDays),
		NewField("compress", cfg.Log.File.Compress),
	)

	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.Log.File.MaxSizeMB,
		MaxBackups: cfg.Log.File.MaxBackups,
		MaxAge:     cfg.Log.File.MaxAgeDays,
		Compress:   cfg.Log.File.Compress,
	}
}
