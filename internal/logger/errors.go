// Copyright 2025 LeafLock Security Solutions
// SPDX-License-Identifier: Apache-2.0

package logger

import "errors"

// Sentinel errors for logger validation.
var (
	ErrNoLogOutputs     = errors.New("no log outputs enabled")
	ErrUnknownLogFormat = errors.New("unknown log format")
	ErrUnknownLogLevel  = errors.New("unknown log level")
)

// Fallback error messages written to stderr when logging fails.
// These are simple strings (not errors.New) because they're written directly to stderr,
// not returned as error values.
const (
	errMsgLogEntryMarshalFailed = "logger: failed to marshal log entry\n"
	errMsgLogDataWriteFailed    = "logger: failed to write log data\n"
	errMsgLogNewlineWriteFailed = "logger: failed to write newline\n"
	errMsgLogPrefixWriteFailed  = "logger: failed to write text log prefix\n"
	errMsgLogSepWriteFailed     = "logger: failed to write text log separator\n"
	errMsgLogFieldWriteFailed   = "logger: failed to write text log field\n"
)
