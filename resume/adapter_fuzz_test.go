package resume

import (
	"os"
	"path/filepath"
	"testing"
)

func FuzzLoadResumeFromFile_JSON(f *testing.F) {
	f.Add([]byte(`{"contact":{"name":"Jane Smith","email":"jane@example.com"}}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpfile := filepath.Join(tmpDir, "resume.json")
		if err := os.WriteFile(tmpfile, data, 0644); err != nil {
			t.Skip()
		}
		_, _ = LoadResumeFromFile(tmpfile)
	})
}

func FuzzLoadResumeFromFile_MD(f *testing.F) {
	f.Add([]byte(`# John Doe
john@example.com

## Experience
### Software Engineer at Acme Corp
- Built things`))
	f.Add([]byte(``))

	f.Fuzz(func(t *testing.T, data []byte) {
		tmpDir := t.TempDir()
		tmpfile := filepath.Join(tmpDir, "resume.md")
		if err := os.WriteFile(tmpfile, data, 0644); err != nil {
			t.Skip()
		}
		_, _ = LoadResumeFromFile(tmpfile)
	})
}
