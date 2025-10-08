package examples_test

import (
	"os"
	"path/filepath"
	"testing"
)

// This file demonstrates testing patterns and conventions used in LaziSpace.
// It serves as a template for writing tests across the codebase.

// ============================================================================
// TABLE-DRIVEN TESTS
// ============================================================================

// ExampleValidator demonstrates a simple validator function for testing
type ExampleValidator struct {
	MinLength int
	MaxLength int
}

func (v *ExampleValidator) Validate(input string) error {
	if len(input) < v.MinLength {
		return &ValidationError{Field: "input", Message: "too short"}
	}
	if len(input) > v.MaxLength {
		return &ValidationError{Field: "input", Message: "too long"}
	}
	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// TestValidate demonstrates table-driven tests - the standard Go pattern
func TestValidate(t *testing.T) {
	validator := &ExampleValidator{
		MinLength: 3,
		MaxLength: 10,
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid input",
			input:   "hello",
			wantErr: false,
		},
		{
			name:    "too short",
			input:   "hi",
			wantErr: true,
			errMsg:  "input: too short",
		},
		{
			name:    "too long",
			input:   "this is way too long",
			wantErr: true,
			errMsg:  "input: too long",
		},
		{
			name:    "minimum length",
			input:   "abc",
			wantErr: false,
		},
		{
			name:    "maximum length",
			input:   "abcdefghij",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		// Use t.Run to create subtests - each runs independently
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.input)

			// Check if error occurred as expected
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// ============================================================================
// SUBTESTS FOR ORGANIZATION
// ============================================================================

func TestComplexFunctionality(t *testing.T) {
	// Group related tests using subtests
	t.Run("input validation", func(t *testing.T) {
		t.Run("empty input", func(t *testing.T) {
			validator := &ExampleValidator{MinLength: 1, MaxLength: 10}
			err := validator.Validate("")
			if err == nil {
				t.Error("expected error for empty input")
			}
		})

		t.Run("whitespace only", func(t *testing.T) {
			validator := &ExampleValidator{MinLength: 3, MaxLength: 10}
			err := validator.Validate("   ")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		validator := &ExampleValidator{MinLength: 0, MaxLength: 0}

		t.Run("zero length allowed", func(t *testing.T) {
			err := validator.Validate("")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	})
}

// ============================================================================
// TEST HELPERS
// ============================================================================

// setupTestValidator creates a validator for testing
// Notice the t.Helper() call - this marks the function as a helper
// so test failures point to the calling line, not this helper
func setupTestValidator(t *testing.T, min, max int) *ExampleValidator {
	t.Helper()

	return &ExampleValidator{
		MinLength: min,
		MaxLength: max,
	}
}

func TestWithHelper(t *testing.T) {
	validator := setupTestValidator(t, 5, 15)

	err := validator.Validate("test")
	if err == nil {
		t.Error("expected error for short input")
	}
}

// ============================================================================
// TEST FIXTURES
// ============================================================================

// createTempTestFile creates a temporary file for testing
func createTempTestFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir() // Automatically cleaned up after test
	filePath := filepath.Join(tmpDir, "test-file.txt")

	err := os.WriteFile(filePath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	return filePath
}

func TestWithFixture(t *testing.T) {
	content := "test content"
	filePath := createTempTestFile(t, content)

	// File exists and contains expected content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected content %q, got %q", content, string(data))
	}

	// No need to clean up - t.TempDir() handles it
}

// ============================================================================
// PARALLEL TESTS
// ============================================================================

// Tests that don't share state can run in parallel for faster execution
func TestParallelExample(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "hello", "HELLO"},
		{"uppercase", "WORLD", "WORLD"},
		{"mixed", "HeLLo", "HELLO"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Mark this subtest as parallel

			// Test logic here
			// got := strings.ToUpper(tt.input)
			// if got != tt.want { ... }
		})
	}
}

// ============================================================================
// BENCHMARK EXAMPLE
// ============================================================================

// BenchmarkValidate demonstrates how to write benchmarks
func BenchmarkValidate(b *testing.B) {
	validator := &ExampleValidator{MinLength: 3, MaxLength: 10}
	input := "hello"

	b.ResetTimer() // Reset timer after setup

	for i := 0; i < b.N; i++ {
		_ = validator.Validate(input)
	}
}

// ============================================================================
// EXAMPLE TEST (SHOWS IN GODOC)
// ============================================================================

// ExampleExampleValidator demonstrates how to use the validator
// The output comment is checked by go test
func ExampleExampleValidator() {
	validator := &ExampleValidator{MinLength: 3, MaxLength: 10}

	err := validator.Validate("hello")
	if err != nil {
		panic(err)
	}

	// Output:
}
