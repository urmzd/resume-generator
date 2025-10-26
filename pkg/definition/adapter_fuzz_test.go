package definition

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzLoadResumeFromFile_YAML tests the YAML parser with random input
// to ensure it doesn't panic or crash on malformed input
func FuzzLoadResumeFromFile_YAML(f *testing.F) {
	// Seed corpus with valid examples
	f.Add([]byte(`meta:
  version: "2.0"
contact:
  name: "John Doe"
  email: "john@example.com"`))

	f.Add([]byte(`meta:
  version: "2.0"
contact:
  name: "Test"
  email: "test@test.com"
experience:
  order: 1
  positions:
    - company: "Tech Corp"
      title: "Engineer"`))

	f.Add([]byte(``)) // Empty file

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpfile := filepath.Join(tmpDir, "resume.yml")

		if err := os.WriteFile(tmpfile, data, 0644); err != nil {
			t.Skip()
		}

		// Test that parser doesn't panic
		// We don't care if it fails, just that it doesn't crash
		_, _ = LoadResumeFromFile(tmpfile)
	})
}

// FuzzLoadResumeFromFile_JSON tests the JSON parser with random input
func FuzzLoadResumeFromFile_JSON(f *testing.F) {
	// Seed corpus
	f.Add([]byte(`{
  "meta": {"version": "2.0"},
  "contact": {
    "name": "Jane Smith",
    "email": "jane@example.com"
  }
}`))

	f.Add([]byte(`{"meta":{"version":"2.0"},"contact":{"name":"Test","email":"test@test.com"}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpfile := filepath.Join(tmpDir, "resume.json")

		if err := os.WriteFile(tmpfile, data, 0644); err != nil {
			t.Skip()
		}

		// Test that parser doesn't panic
		_, _ = LoadResumeFromFile(tmpfile)
	})
}

// FuzzLoadResumeFromFile_TOML tests the TOML parser with random input
func FuzzLoadResumeFromFile_TOML(f *testing.F) {
	// Seed corpus
	f.Add([]byte(`[meta]
version = "2.0"

[contact]
name = "Bob Wilson"
email = "bob@example.com"`))

	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpfile := filepath.Join(tmpDir, "resume.toml")

		if err := os.WriteFile(tmpfile, data, 0644); err != nil {
			t.Skip()
		}

		// Test that parser doesn't panic
		_, _ = LoadResumeFromFile(tmpfile)
	})
}
