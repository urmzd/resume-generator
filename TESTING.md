# Testing Strategy for Resume Generator

This document outlines the testing approach for the resume generator project, following Go and Google best practices.

## Testing Philosophy

We follow Google's Go testing best practices:

1. **Test failures occur in Test functions** - Not in helper functions
2. **Prefer standard library** - Avoid third-party assertion libraries unless absolutely needed
3. **Table-driven tests** - For testing multiple scenarios efficiently
4. **Use subtests** - For better organization and parallel execution
5. **Fuzz testing** - For discovering edge cases in parsers and generators
6. **Real transports** - When testing integrations, use real connections to test doubles

## Test Organization

```
pkg/
├── utils/
│   ├── paths.go
│   └── paths_test.go           # Unit tests for path utilities
├── definition/
│   ├── resume.go
│   ├── resume_test.go          # Unit tests for Resume types
│   ├── adapter.go
│   ├── adapter_test.go         # Unit tests for file loading
│   └── adapter_fuzz_test.go    # Fuzz tests for parsing YAML/JSON/TOML
├── generators/
│   ├── generator.go
│   ├── generator_test.go       # Unit tests for template loading
│   ├── html.go
│   ├── html_test.go            # Unit tests for HTML generation
│   ├── latex.go
│   └── latex_test.go           # Unit tests for LaTeX generation
└── compilers/
    ├── latex.go
    ├── latex_test.go           # Unit tests for LaTeX compilation
    ├── html.go
    └── html_test.go            # Unit tests for HTML to PDF

tests/
└── integration/
    └── e2e_test.go             # End-to-end integration tests
```

## Types of Tests

### 1. Unit Tests

Test individual functions and methods in isolation.

**Pattern: Table-Driven with Subtests**

```go
func TestResolvePath(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "absolute path",
            input:   "/tmp/test",
            want:    "/tmp/test",
            wantErr: false,
        },
        {
            name:    "home directory expansion",
            input:   "~/test",
            want:    "", // Will be computed dynamically
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ResolvePath(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ResolvePath() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if tt.want != "" && got != tt.want {
                t.Errorf("ResolvePath() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Key Principles:**
- Use `t.Errorf` for failures (continues testing other cases)
- Use `t.Fatalf` only for setup failures
- Each test case is a complete scenario
- Subtests enable parallel execution with `t.Parallel()`

### 2. Fuzz Tests

Discover edge cases through randomized input generation.

**Pattern: Property-Based Testing**

```go
func FuzzLoadResumeFromFile(f *testing.F) {
    // Seed corpus with valid examples
    f.Add([]byte(`meta:
  version: "2.0"
contact:
  name: "John Doe"
  email: "john@example.com"`))

    f.Fuzz(func(t *testing.T, data []byte) {
        // Create temp file
        tmpfile, err := os.CreateTemp("", "resume-*.yml")
        if err != nil {
            t.Skip()
        }
        defer os.Remove(tmpfile.Name())

        if _, err := tmpfile.Write(data); err != nil {
            t.Skip()
        }
        tmpfile.Close()

        // Test that parser doesn't panic
        _, err = LoadResumeFromFile(tmpfile.Name())
        // We don't care if it fails, just that it doesn't crash
        // This tests for panics and unexpected behaviors
    })
}
```

**When to Use Fuzz Tests:**
- Parsers (YAML, JSON, TOML input)
- Template rendering functions
- Path resolution logic
- Any function processing user input

**Running Fuzz Tests:**
```bash
# Run with fuzzing for 30 seconds
go test -fuzz=FuzzLoadResumeFromFile -fuzztime 30s

# Run specific failing case
go test -run=FuzzLoadResumeFromFile/abc123

# Run without fuzzing (just seed corpus)
go test
```

### 3. Integration Tests

Test multiple components working together.

**Pattern: Real Resources with Cleanup**

```go
func TestGenerateResume_EndToEnd(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }

    // Setup
    tmpDir := t.TempDir() // Automatic cleanup
    inputFile := filepath.Join(tmpDir, "resume.yml")

    // Create test input
    resumeData := `meta:
  version: "2.0"
contact:
  name: "Test User"
  email: "test@example.com"`

    if err := os.WriteFile(inputFile, []byte(resumeData), 0644); err != nil {
        t.Fatalf("Failed to create test input: %v", err)
    }

    // Test actual generation flow
    inputData, err := definition.LoadResumeFromFile(inputFile)
    if err != nil {
        t.Fatalf("LoadResumeFromFile() error = %v", err)
    }

    resume := inputData.ToResume()
    generator := generators.NewGenerator(logger)

    // Test multiple templates
    templates := []string{"modern-html", "modern-latex"}
    for _, tmplName := range templates {
        t.Run(tmplName, func(t *testing.T) {
            content, err := generator.Generate(tmplName, resume)
            if err != nil {
                t.Errorf("Generate(%s) error = %v", tmplName, err)
            }
            if len(content) == 0 {
                t.Errorf("Generate(%s) returned empty content", tmplName)
            }
        })
    }
}
```

**Running Integration Tests:**
```bash
# Run all tests including integration
go test ./...

# Skip integration tests (fast)
go test -short ./...

# Run only integration tests
go test -run Integration ./tests/integration
```

### 4. Benchmark Tests

Measure performance of critical paths.

```go
func BenchmarkGenerateHTML(b *testing.B) {
    resume := &definition.Resume{
        Contact: definition.Contact{
            Name:  "Test User",
            Email: "test@example.com",
        },
    }

    generator := generators.NewHTMLGenerator(logger)
    tmpl, _ := os.ReadFile("templates/modern-html/template.html")
    templateContent := string(tmpl)

    b.ResetTimer() // Start timing after setup
    for i := 0; i < b.N; i++ {
        _, err := generator.Generate(templateContent, resume)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Test Coverage Goals

| Package | Target Coverage | Priority |
|---------|----------------|----------|
| pkg/utils | 90%+ | High |
| pkg/definition | 85%+ | High |
| pkg/generators | 80%+ | High |
| pkg/compilers | 70%+ | Medium |
| cmd/ | 60%+ | Low |

**Checking Coverage:**
```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Summary by package
go test -cover ./...
```

## Testing Commands

```bash
# Run all tests
just test

# Run tests with coverage
just test-coverage

# Run fuzz tests
just test-fuzz

# Run benchmarks
just test-bench

# Run only fast tests
just test-fast
```

## Continuous Integration

Tests run automatically on:
- Pull requests
- Main branch commits
- Release tags

**CI Pipeline:**
1. Unit tests (fast)
2. Fuzz tests (5 min limit)
3. Integration tests
4. Coverage report generation
5. Benchmark comparison (for PRs)

## Writing New Tests - Checklist

- [ ] Does the test have a clear name describing what it tests?
- [ ] Is it table-driven if testing multiple scenarios?
- [ ] Does it use subtests for organization?
- [ ] Does it clean up resources (use `t.TempDir()`, `defer`)?
- [ ] Does it use `t.Errorf` instead of `t.Fatalf` for assertion failures?
- [ ] Does the fuzz test check invariants, not exact outputs?
- [ ] Are integration tests marked with `testing.Short()` check?
- [ ] Is the test deterministic and reproducible?

## Anti-Patterns to Avoid

❌ **Don't use assertion libraries** (testify, etc.) - Use standard library
❌ **Don't create assertion helpers** - Return errors and let tests decide
❌ **Don't use `t.Fatal` for assertion failures** - Use `t.Error`/`t.Errorf`
❌ **Don't share mutable state** between tests
❌ **Don't rely on test execution order**
❌ **Don't mock when you can use real implementations**
❌ **Don't test implementation details** - Test behavior

## Resources

- [Official Go Testing Docs](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests Wiki](https://go.dev/wiki/TableDrivenTests)
- [Fuzzing Tutorial](https://go.dev/doc/tutorial/fuzz)
- [Google Go Style Guide](https://google.github.io/styleguide/go/best-practices.html)
- [Dave Cheney: Prefer table driven tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
