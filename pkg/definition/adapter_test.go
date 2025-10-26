package definition

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadResumeFromFile_YAML(t *testing.T) {
	tmpDir := t.TempDir()

	validYAML := `meta:
  version: "2.0"
contact:
  name: "John Doe"
  email: "john@example.com"
  phone: "+1-555-1234"`

	tests := []struct {
		name      string
		filename  string
		content   string
		wantErr   bool
		checkFunc func(t *testing.T, data InputData)
	}{
		{
			name:     "valid YAML file",
			filename: "resume.yml",
			content:  validYAML,
			wantErr:  false,
			checkFunc: func(t *testing.T, data InputData) {
				if data.GetFormat() != "yaml" {
					t.Errorf("GetFormat() = %q, want yaml", data.GetFormat())
				}
				resume := data.ToResume()
				if resume.Contact.Name != "John Doe" {
					t.Errorf("Contact.Name = %q, want John Doe", resume.Contact.Name)
				}
				if resume.Contact.Email != "john@example.com" {
					t.Errorf("Contact.Email = %q, want john@example.com", resume.Contact.Email)
				}
			},
		},
		{
			name:     "valid YAML with .yaml extension",
			filename: "resume.yaml",
			content:  validYAML,
			wantErr:  false,
			checkFunc: func(t *testing.T, data InputData) {
				if data.GetFormat() != "yaml" {
					t.Errorf("GetFormat() = %q, want yaml", data.GetFormat())
				}
			},
		},
		{
			name:     "invalid YAML syntax",
			filename: "resume.yml",
			content:  "invalid: yaml: syntax:",
			wantErr:  true,
		},
		{
			name:     "missing required field",
			filename: "resume.yml",
			content: `meta:
  version: "2.0"
contact:
  email: "john@example.com"`,
			wantErr: true,
		},
		{
			name:     "empty file",
			filename: "resume.yml",
			content:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tt.filename)
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

func TestLoadResumeFromFile_JSON(t *testing.T) {
	tmpDir := t.TempDir()

	validJSON := `{
  "meta": {
    "version": "2.0"
  },
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
				resume := data.ToResume()
				if resume.Contact.Name != "Jane Smith" {
					t.Errorf("Contact.Name = %q, want Jane Smith", resume.Contact.Name)
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

func TestLoadResumeFromFile_TOML(t *testing.T) {
	tmpDir := t.TempDir()

	validTOML := `[meta]
version = "2.0"

[contact]
name = "Bob Wilson"
email = "bob@example.com"`

	tests := []struct {
		name      string
		content   string
		wantErr   bool
		checkFunc func(t *testing.T, data InputData)
	}{
		{
			name:    "valid TOML file",
			content: validTOML,
			wantErr: false,
			checkFunc: func(t *testing.T, data InputData) {
				if data.GetFormat() != "toml" {
					t.Errorf("GetFormat() = %q, want toml", data.GetFormat())
				}
				resume := data.ToResume()
				if resume.Contact.Name != "Bob Wilson" {
					t.Errorf("Contact.Name = %q, want Bob Wilson", resume.Contact.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "resume.toml")
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

func TestLoadResumeFromFile_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "resume.xml")
	if err := os.WriteFile(testFile, []byte("<resume></resume>"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := LoadResumeFromFile(testFile)
	if err == nil {
		t.Error("LoadResumeFromFile() expected error for unsupported format, got nil")
	}
}

func TestLoadResumeFromFile_FileNotFound(t *testing.T) {
	_, err := LoadResumeFromFile("/nonexistent/resume.yml")
	if err == nil {
		t.Error("LoadResumeFromFile() expected error for non-existent file, got nil")
	}
}

func TestResumeAdapter_ToResume(t *testing.T) {
	resume := &Resume{
		Meta: ResumeMetadata{
			Version: "2.0",
			Output: OutputPreferences{
				Formats: []string{"pdf"},
			},
		},
		Contact: Contact{
			Name:  "Test User",
			Email: "test@example.com",
		},
	}

	adapter := &ResumeAdapter{
		Resume:           resume,
		SerializationFmt: "yaml",
	}

	if got := adapter.ToResume(); got != resume {
		t.Errorf("ToResume() returned different resume instance")
	}
}

func TestResumeAdapter_GetFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"yaml format", "yaml", "yaml"},
		{"json format", "json", "json"},
		{"toml format", "toml", "toml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ResumeAdapter{
				Resume:           &Resume{},
				SerializationFmt: tt.format,
			}

			got := adapter.GetFormat()
			if got != tt.want {
				t.Errorf("GetFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResumeAdapter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		resume  *Resume
		wantErr bool
	}{
		{
			name: "valid resume",
			resume: &Resume{
				Meta: ResumeMetadata{
					Version: "2.0",
					Output: OutputPreferences{
						Formats: []string{"pdf"},
					},
				},
				Contact: Contact{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			resume: &Resume{
				Meta: ResumeMetadata{
					Version: "2.0",
					Output: OutputPreferences{
						Formats: []string{"pdf"},
					},
				},
				Contact: Contact{
					Email: "test@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "missing email",
			resume: &Resume{
				Meta: ResumeMetadata{
					Version: "2.0",
					Output: OutputPreferences{
						Formats: []string{"pdf"},
					},
				},
				Contact: Contact{
					Name: "Test User",
				},
			},
			wantErr: true,
		},
		{
			name: "missing output formats",
			resume: &Resume{
				Meta: ResumeMetadata{
					Version: "2.0",
				},
				Contact: Contact{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "missing version",
			resume: &Resume{
				Meta: ResumeMetadata{
					Output: OutputPreferences{
						Formats: []string{"pdf"},
					},
				},
				Contact: Contact{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &ResumeAdapter{
				Resume:           tt.resume,
				SerializationFmt: "yaml",
			}

			err := adapter.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
