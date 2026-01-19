package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

func getTestLogger(t *testing.T) *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return logger.Sugar()
}

func TestLoadTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test template structure
	templateDir := filepath.Join(tmpDir, "templates", "test-template")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Create config.yml
	configContent := `name: test-template
display_name: Test Template
description: A test template
format: html
version: "1.0"
`
	configPath := filepath.Join(templateDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config.yml: %v", err)
	}

	// Create template.html
	templatePath := filepath.Join(templateDir, "template.html")
	if err := os.WriteFile(templatePath, []byte("<html>Test</html>"), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	// Set environment variable to point to our test root directory
	t.Setenv("RESUME_TEMPLATES_DIR", tmpDir)

	tests := []struct {
		name         string
		templateName string
		wantErr      bool
		checkFunc    func(t *testing.T, tmpl *Template)
	}{
		{
			name:         "valid template",
			templateName: "test-template",
			wantErr:      false,
			checkFunc: func(t *testing.T, tmpl *Template) {
				if tmpl.Name != "test-template" {
					t.Errorf("Template.Name = %q, want test-template", tmpl.Name)
				}
				if tmpl.Type != TemplateTypeHTML {
					t.Errorf("Template.Type = %q, want %q", tmpl.Type, TemplateTypeHTML)
				}
				if tmpl.DisplayName != "Test Template" {
					t.Errorf("Template.DisplayName = %q, want Test Template", tmpl.DisplayName)
				}
			},
		},
		{
			name:         "non-existent template",
			templateName: "nonexistent",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadTemplate(tt.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	templatesDir := filepath.Join(tmpDir, "templates")

	// Create multiple test templates
	templates := []struct {
		name   string
		format string
	}{
		{"template1", "html"},
		{"template2", "latex"},
	}

	for _, tmpl := range templates {
		templateDir := filepath.Join(templatesDir, tmpl.name)
		if err := os.MkdirAll(templateDir, 0755); err != nil {
			t.Fatalf("Failed to create template directory: %v", err)
		}

		configContent := `name: ` + tmpl.name + `
display_name: ` + tmpl.name + `
description: Test
format: ` + tmpl.format + `
`
		configPath := filepath.Join(templateDir, "config.yml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to create config.yml: %v", err)
		}

		ext := ".html"
		if tmpl.format == "latex" {
			ext = ".tex"
		}
		templatePath := filepath.Join(templateDir, "template"+ext)
		if err := os.WriteFile(templatePath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create template file: %v", err)
		}
	}

	t.Setenv("RESUME_TEMPLATES_DIR", tmpDir)

	got, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	if len(got) != 2 {
		t.Errorf("ListTemplates() returned %d templates, want 2", len(got))
	}
}

func TestGenerator_GenerateWithTemplate(t *testing.T) {
	logger := getTestLogger(t)
	generator := NewGenerator(logger)

	tmpDir := t.TempDir()

	// Create a simple HTML template
	templateDir := filepath.Join(tmpDir, "templates", "simple-html")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	configContent := `name: simple-html
display_name: Simple HTML
description: Simple test template
format: html
`
	configPath := filepath.Join(templateDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config.yml: %v", err)
	}

	templateContent := `<html>
<head><title>{{.Contact.Name}}</title></head>
<body>
<h1>{{.Contact.Name}}</h1>
<p>{{.Contact.Email}}</p>
</body>
</html>`

	templatePath := filepath.Join(templateDir, "template.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}

	t.Setenv("RESUME_TEMPLATES_DIR", tmpDir)

	// Load the template
	tmpl, err := LoadTemplate("simple-html")
	if err != nil {
		t.Fatalf("LoadTemplate() error = %v", err)
	}

	// Create test resume
	resume := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Test User",
			Email: "test@example.com",
		},
	}

	// Generate content
	got, err := generator.GenerateWithTemplate(tmpl, resume)
	if err != nil {
		t.Fatalf("GenerateWithTemplate() error = %v", err)
	}

	if len(got) == 0 {
		t.Error("GenerateWithTemplate() returned empty content")
	}

	// Check that template was rendered with data
	if !contains(got, "Test User") {
		t.Error("Generated content does not contain contact name")
	}
	if !contains(got, "test@example.com") {
		t.Error("Generated content does not contain email")
	}
}

func TestFormatTemplateName(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("RESUME_TEMPLATES_DIR", tmpDir)

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDir)
	})

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "already formatted html",
			input: "modern-html",
			want:  "modern-html",
		},
		{
			name:  "already formatted latex",
			input: "modern-latex",
			want:  "modern-latex",
		},
		{
			name:  "bare name",
			input: "custom",
			want:  "custom", // Returns original if nothing found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTemplateName(tt.input)
			if got != tt.want {
				t.Errorf("FormatTemplateName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
