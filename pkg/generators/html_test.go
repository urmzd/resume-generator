package generators

import (
	"strings"
	"testing"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

func TestHTMLGenerator_Generate(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewHTMLGenerator(logger)

	r := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}

	// Template accessing fields directly
	templateContent := `<div id="content">
<h1>{{.Contact.Name}}</h1>
{{if .Contact.Email}}<span>{{.Contact.Email}}</span>{{end}}
</div>`

	got, err := gen.Generate(templateContent, r)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(got, "Jane Doe") {
		t.Error("Generate() missing name")
	}
	if !strings.Contains(got, "jane@example.com") {
		t.Error("Generate() missing email")
	}
}

func TestHTMLGenerator_GenerateStandalone(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewHTMLGenerator(logger)

	r := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}

	templateContent := `<div>{{.Contact.Name}}</div>`
	cssContent := "body { color: #333; }"

	got, err := gen.GenerateStandalone(templateContent, cssContent, r)
	if err != nil {
		t.Fatalf("GenerateStandalone() error = %v", err)
	}

	if !strings.Contains(got, "<!DOCTYPE html>") {
		t.Error("GenerateStandalone() missing <!DOCTYPE html>")
	}
	if !strings.Contains(got, cssContent) {
		t.Errorf("GenerateStandalone() missing CSS content %q", cssContent)
	}
	if !strings.Contains(got, "Jane Doe") {
		t.Error("GenerateStandalone() missing rendered resume content")
	}
}
