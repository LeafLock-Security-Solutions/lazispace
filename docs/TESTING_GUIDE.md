# Testing Guide

This guide explains testing patterns, conventions, and best practices for LaziSpace.

## Table of Contents

- [Testing Philosophy](#testing-philosophy)
- [Running Tests](#running-tests)
- [Test Organization](#test-organization)
- [Writing Tests](#writing-tests)
- [Test Fixtures](#test-fixtures)
- [Best Practices](#best-practices)

## Testing Philosophy

LaziSpace follows standard Go testing conventions:

- Tests live alongside the code they test (`*_test.go` files)
- Use table-driven tests for multiple test cases
- Use subtests with `t.Run()` for better organization
- Tests should be fast, isolated, and deterministic
- Mock external dependencies

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage

# Run tests for specific package
go test ./internal/config/...

# Run specific test
go test -run TestValidateWorkspace ./internal/config/
```

### Coverage

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
go tool cover -html=coverage.out
```

## Test Organization

### File Structure

Tests live next to the code they test:

```
internal/config/
├── loader.go
├── loader_test.go
├── validator.go
├── validator_test.go
└── testdata/
    ├── valid-workspace.yml
    └── invalid-workspace.yml
```

### Test Naming

- Test files: `*_test.go`
- Test functions: `TestFunctionName(t *testing.T)`
- Benchmark functions: `BenchmarkFunctionName(b *testing.B)`
- Example functions: `ExampleFunctionName()`

### Package Names: Black-Box vs White-Box Testing

Go supports two types of testing approaches based on package naming:

#### Black-Box Testing (package_test)

**What it is**: Testing from outside the package, as a user would. You can only access exported (public) functions, types, and methods.

**When to use**: 
- Testing the public API of your package
- Ensuring the package works as documented
- Preventing tests from relying on internal implementation details

**Example**:
```go
package config_test  // Note the _test suffix

import (
    "testing"
    "github.com/LeafLock-Security-Solutions/lazispace/internal/config"
)

func TestLoadWorkspace(t *testing.T) {
    // Can only use config.LoadWorkspace() and other exported functions
    // Cannot access unexported variables or functions
    ws, err := config.LoadWorkspace("workspace.yml")
    // ...
}
```

**Advantages**:
- Forces you to test through the public API
- Tests remain valid even if internal implementation changes
- Mimics how actual users will use your code
- Catches issues with unclear or insufficient public APIs

#### White-Box Testing (package)

**What it is**: Testing from inside the package. You have access to both exported and unexported (private) functions and variables.

**When to use**:
- Testing internal helper functions that don't need to be public
- Testing complex internal logic in isolation
- When you need to set up internal state for testing
- Testing edge cases that are hard to trigger through the public API

**Example**:
```go
package config  // Same package as the code being tested

import "testing"

func TestParseYAMLInternal(t *testing.T) {
    // Can access unexported functions like parseYAMLInternal()
    result := parseYAMLInternal(data)
    // ...
}
```

**Advantages**:
- Can test internal functions directly
- Can access and manipulate internal state
- Useful for testing implementation details

#### Which Should You Use?

**Default to black-box testing** (`package_test`) because:
- It ensures your public API is sufficient
- Tests are more maintainable (won't break from internal refactoring)
- Better simulates real usage

**Use white-box testing** (`package`) only when:
- You need to test unexported helper functions
- The internal logic is complex enough to warrant isolated testing
- Black-box testing would be overly complicated

#### Practical Example

```go
// internal/config/loader.go
package config

// Exported - users can call this
func LoadWorkspace(path string) (*Workspace, error) {
    data, err := readFile(path)
    if err != nil {
        return nil, err
    }
    return parseYAML(data)
}

// Unexported - internal helper
func parseYAML(data []byte) (*Workspace, error) {
    // Complex YAML parsing logic
}

func readFile(path string) ([]byte, error) {
    // File reading logic
}
```

**Black-box test** (recommended):
```go
package config_test

func TestLoadWorkspace(t *testing.T) {
    // Test the entire LoadWorkspace flow through public API
    ws, err := config.LoadWorkspace("testdata/valid.yml")
    if err != nil {
        t.Fatalf("LoadWorkspace failed: %v", err)
    }
    // Assert on workspace contents
}
```

**White-box test** (if parseYAML is very complex):
```go
package config

func TestParseYAML(t *testing.T) {
    // Test parseYAML directly with various edge cases
    ws, err := parseYAML([]byte(`name: test`))
    // ...
}
```

**Recommendation**: Write most tests as black-box tests. Only use white-box testing for complex internal functions that deserve their own focused tests.

## Writing Tests

### Table-Driven Tests

Table-driven tests are the standard Go pattern for testing multiple cases:

```go
func TestValidateWorkspace(t *testing.T) {
    tests := []struct {
        name    string
        input   *types.Workspace
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid workspace",
            input: &types.Workspace{
                Name:    "test-workspace",
                Version: "1.0",
                Layout:  validLayout(),
            },
            wantErr: false,
        },
        {
            name: "missing name",
            input: &types.Workspace{
                Version: "1.0",
                Layout:  validLayout(),
            },
            wantErr: true,
            errMsg:  "workspace name is required",
        },
        {
            name: "invalid name format",
            input: &types.Workspace{
                Name:    "invalid name!",
                Version: "1.0",
                Layout:  validLayout(),
            },
            wantErr: true,
            errMsg:  "workspace name contains invalid characters",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateWorkspace(tt.input)
            
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
```

### Subtests

Use subtests to group related test cases:

```go
func TestWorkspaceValidation(t *testing.T) {
    t.Run("name validation", func(t *testing.T) {
        // Test name-specific validation
    })

    t.Run("version validation", func(t *testing.T) {
        // Test version-specific validation
    })

    t.Run("layout validation", func(t *testing.T) {
        // Test layout-specific validation
    })
}
```

### Helper Functions

Create test helpers for common setup/teardown:

```go
// Helper function (notice the t.Helper() call)
func setupTestWorkspace(t *testing.T) *types.Workspace {
    t.Helper()
    
    return &types.Workspace{
        Name:    "test-workspace",
        Version: "1.0",
        Layout:  &types.Layout{
            Tabs: []types.Tab{
                {Title: "main", Active: true},
            },
        },
    }
}

func TestSomething(t *testing.T) {
    ws := setupTestWorkspace(t)
    // Use ws in test
}
```

### Test Fixtures

Use `testdata` directories for test files:

```go
func TestLoadConfig(t *testing.T) {
    data, err := os.ReadFile("testdata/valid-workspace.yml")
    if err != nil {
        t.Fatalf("failed to read test fixture: %v", err)
    }
    
    ws, err := LoadWorkspace(data)
    if err != nil {
        t.Fatalf("failed to load workspace: %v", err)
    }
    
    // Assert expectations
}
```

### Golden Files

For comparing complex output, use golden files:

```go
func TestGenerateConfig(t *testing.T) {
    result := GenerateConfig(input)
    
    goldenPath := "testdata/expected-config.yml"
    
    // Update golden file with -update flag
    if *update {
        os.WriteFile(goldenPath, []byte(result), 0644)
    }
    
    expected, err := os.ReadFile(goldenPath)
    if err != nil {
        t.Fatalf("failed to read golden file: %v", err)
    }
    
    if result != string(expected) {
        t.Errorf("output mismatch\ngot:\n%s\nwant:\n%s", result, expected)
    }
}
```

## Test Fixtures

### Directory Structure

```
internal/config/
└── testdata/
    ├── valid-workspace.yml
    ├── invalid-workspace.yml
    ├── minimal.json
    └── complex.yml
```

### Fixture Files

Keep fixtures minimal and focused:

```yaml
# testdata/valid-workspace.yml
name: test-workspace
version: "1.0"
layout:
  tabs:
    - title: main
      active: true
      split:
        direction: horizontal
        panes:
          - commands:
              - type: run
                value: echo "test"
```

### Loading Fixtures

```go
func loadFixture(t *testing.T, filename string) []byte {
    t.Helper()
    
    data, err := os.ReadFile(filepath.Join("testdata", filename))
    if err != nil {
        t.Fatalf("failed to load fixture %s: %v", filename, err)
    }
    
    return data
}
```

## Best Practices

### Do's

- **Use table-driven tests** for multiple similar test cases
- **Use subtests** with `t.Run()` for organization
- **Call `t.Helper()`** in test helper functions
- **Use `t.Fatalf()`** for setup failures
- **Use `t.Errorf()`** for test failures
- **Test edge cases** (empty, nil, zero values)
- **Test error paths** as thoroughly as success paths
- **Keep tests focused** - one concept per test
- **Use descriptive test names** that explain what's being tested
- **Make tests deterministic** - no random data, no time.Now()

### Don'ts

- **Don't skip cleanup** - always clean up resources
- **Don't use sleeps** - use channels or polling with timeout
- **Don't test implementation details** - test behavior
- **Don't share state** between tests
- **Don't use init()** for test setup
- **Don't write flaky tests** - ensure reliability
- **Don't ignore errors** in test code

### Test Coverage Goals

- **Unit tests**: Aim for 80%+ coverage of business logic
- **Integration tests**: Cover critical paths and workflows
- **Edge cases**: Test boundary conditions and error handling

### When to Mock

Mock when:
- Testing external dependencies (APIs, databases, filesystems)
- Testing error conditions that are hard to trigger
- Testing code that depends on time or randomness

Don't mock when:
- Testing pure functions
- Testing simple structs or data structures
- Mocking would make the test more complex than the code

### Parallel Tests

Mark tests as parallel when they don't share state:

```go
func TestIndependentFunction(t *testing.T) {
    t.Parallel()
    
    // Test code
}
```

### Test Timeouts

Use context with timeout for tests that might hang:

```go
func TestLongRunning(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Test code using ctx
}
```

## Example Test Template

See `test/examples/example_test.go` for a complete example showing:
- Table-driven tests
- Subtests
- Test helpers
- Fixture usage
- Error handling

## Continuous Integration

Tests run automatically on every push and pull request via GitHub Actions.
See `.github/workflows/ci.yml` for the CI configuration.

## Further Reading

- [Go Testing Package](https://pkg.go.dev/testing)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Advanced Testing with Go](https://www.youtube.com/watch?v=8hQG7QlcLBk)
- [Testable Examples](https://go.dev/blog/examples)
