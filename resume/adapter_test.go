package resume

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResumeFromFile_JSON(t *testing.T) {
	tmpDir := t.TempDir()

	validJSON := `{
  "contact": {
    "name": "Jane Smith",
    "email": "jane@example.com"
  }
}`

	tests := []struct {
		name      string
		content   string
		wantErr   bool
		checkFunc func(t *testing.T, data InputData)
	}{
		{
			name:    "valid JSON file",
			content: validJSON,
			wantErr: false,
			checkFunc: func(t *testing.T, data InputData) {
				if data.GetFormat() != "json" {
					t.Errorf("GetFormat() = %q, want json", data.GetFormat())
				}
				r := data.ToResume()
				if r.Contact.Name != "Jane Smith" {
					t.Errorf("Contact.Name = %q, want Jane Smith", r.Contact.Name)
				}
			},
		},
		{
			name:    "invalid JSON syntax",
			content: `{"invalid": json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "resume.json")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			got, err := LoadResumeFromFile(testFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadResumeFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestLoadResumeFromFile_UnknownExtFallsToMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "resume.xml")
	if err := os.WriteFile(testFile, []byte("<resume></resume>"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadResumeFromFile(testFile)
	if err == nil {
		t.Error("LoadResumeFromFile() expected error for non-resume content, got nil")
	}
}

func TestLoadResumeFromFile_FileNotFound(t *testing.T) {
	_, err := LoadResumeFromFile("/nonexistent/resume.json")
	if err == nil {
		t.Error("LoadResumeFromFile() expected error for non-existent file, got nil")
	}
}

func TestResumeAdapter_GetFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"json format", "json", "json"},
		{"md format", "md", "md"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ResumeAdapter{
				Resume:           &Resume{},
				SerializationFmt: tt.format,
			}
			if got := adapter.GetFormat(); got != tt.want {
				t.Errorf("GetFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResumeAdapter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		r       *Resume
		wantErr bool
	}{
		{
			name:    "valid resume",
			r:       &Resume{Contact: Contact{Name: "Test", Email: "test@example.com"}},
			wantErr: false,
		},
		{
			name:    "missing name",
			r:       &Resume{Contact: Contact{Email: "test@example.com"}},
			wantErr: true,
		},
		{
			name:    "missing email",
			r:       &Resume{Contact: Contact{Name: "Test"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ResumeAdapter{Resume: tt.r, SerializationFmt: "json"}
			err := adapter.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
