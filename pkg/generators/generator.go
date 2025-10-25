package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/resume-generator/pkg/definition"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

// TemplateType represents the template rendering engine
type TemplateType string

const (
	TemplateTypeLaTeX TemplateType = "latex"
	TemplateTypeHTML  TemplateType = "html"
)

// Template represents a resume template
type Template struct {
	Name string
	Type TemplateType
	Path string
}

// Generator renders resumes to PDF using templates
type Generator struct {
	logger *zap.SugaredLogger
}

// NewGenerator creates a new template-based generator
func NewGenerator(logger *zap.SugaredLogger) *Generator {
	return &Generator{logger: logger}
}

// LoadTemplate loads a template by name from assets/templates/
// Auto-detects whether it's HTML or LaTeX based on file extension
func LoadTemplate(templateName string) (*Template, error) {
	// Resolve template directory path
	templateDir, err := utils.ResolveAssetPath(filepath.Join("assets", "templates", templateName))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template path: %w", err)
	}

	// Check if template directory exists
	if !utils.DirExists(templateDir) {
		return nil, fmt.Errorf("template not found: %s (searched in %s)", templateName, templateDir)
	}

	// Look for template files
	htmlPath := filepath.Join(templateDir, "template.html")
	latexPath := filepath.Join(templateDir, "template.tex")

	// Check which type exists
	if utils.FileExists(htmlPath) {
		return &Template{
			Name: templateName,
			Type: TemplateTypeHTML,
			Path: htmlPath,
		}, nil
	}

	if utils.FileExists(latexPath) {
		return &Template{
			Name: templateName,
			Type: TemplateTypeLaTeX,
			Path: latexPath,
		}, nil
	}

	return nil, fmt.Errorf("template %s has no template.html or template.tex file", templateName)
}

// ListTemplates returns all available templates
func ListTemplates() ([]Template, error) {
	// Resolve templates directory
	templatesDir, err := utils.ResolveAssetPath(filepath.Join("assets", "templates"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve templates directory: %w", err)
	}

	if !utils.DirExists(templatesDir) {
		return nil, fmt.Errorf("templates directory not found: %s", templatesDir)
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []Template
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		tmpl, err := LoadTemplate(entry.Name())
		if err != nil {
			// Skip invalid templates
			continue
		}

		templates = append(templates, *tmpl)
	}

	return templates, nil
}

// Generate renders a resume using the specified template
// Returns the rendered content as a string
func (g *Generator) Generate(templateName string, resume *definition.EnhancedResume) (string, error) {
	// Load template
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return "", err
	}

	g.logger.Infof("Generating resume using template: %s (%s)", tmpl.Name, tmpl.Type)

	// Read template content
	content, err := os.ReadFile(tmpl.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Render based on type
	switch tmpl.Type {
	case TemplateTypeHTML:
		return g.renderHTML(string(content), resume)
	case TemplateTypeLaTeX:
		return g.renderLaTeX(string(content), resume)
	default:
		return "", fmt.Errorf("unknown template type: %s", tmpl.Type)
	}
}

// renderHTML renders an HTML template
func (g *Generator) renderHTML(templateContent string, resume *definition.EnhancedResume) (string, error) {
	htmlGen := NewHTMLGenerator(g.logger)
	return htmlGen.Generate(templateContent, resume)
}

// renderLaTeX renders a LaTeX template
func (g *Generator) renderLaTeX(templateContent string, resume *definition.EnhancedResume) (string, error) {
	latexGen := NewLaTeXGenerator(g.logger)
	return latexGen.Generate(templateContent, resume)
}

// GetTemplateType returns the type of a template
func GetTemplateType(templateName string) (TemplateType, error) {
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return "", err
	}
	return tmpl.Type, nil
}

// FormatTemplateName formats a raw template name
// e.g., "modern" -> "modern-html" (if modern-html exists)
func FormatTemplateName(name string) string {
	// If already properly formatted, return as-is
	if strings.Contains(name, "-html") || strings.Contains(name, "-latex") {
		return name
	}

	// Try common patterns
	candidates := []string{
		name + "-html",
		name + "-latex",
		"modern-html", // fallback
		"base-latex",  // fallback
	}

	for _, candidate := range candidates {
		if _, err := LoadTemplate(candidate); err == nil {
			return candidate
		}
	}

	// Return original if nothing found
	return name
}
