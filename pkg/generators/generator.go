package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/resume-generator/pkg/definition"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// TemplateType represents the template rendering engine
type TemplateType string

const (
	TemplateTypeLaTeX TemplateType = "latex"
	TemplateTypeHTML  TemplateType = "html"
)

// Template represents a resume template including metadata from config.yml
type Template struct {
	Name        string
	Type        TemplateType
	Path        string
	DisplayName string
	Description string
	Version     string
	Author      string
	Tags        []string
	Config      TemplateConfig
}

// TemplateConfig contains metadata about a template loaded from config.yml
type TemplateConfig struct {
	Name         string   `yaml:"name"`
	DisplayName  string   `yaml:"display_name"`
	Description  string   `yaml:"description"`
	Format       string   `yaml:"format"`
	Version      string   `yaml:"version,omitempty"`
	Author       string   `yaml:"author,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	TemplateFile string   `yaml:"template_file,omitempty"`
}

// Generator renders resumes to PDF using templates
type Generator struct {
	logger *zap.SugaredLogger
}

// NewGenerator creates a new template-based generator
func NewGenerator(logger *zap.SugaredLogger) *Generator {
	return &Generator{logger: logger}
}

// LoadTemplate loads a template by name from templates/
// Uses config.yml metadata to determine template format and additional details
func LoadTemplate(templateName string) (*Template, error) {
	// Resolve template directory path
	templateDir, err := utils.ResolveAssetPath(filepath.Join("templates", templateName))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve template path: %w", err)
	}

	// Check if template directory exists
	if !utils.DirExists(templateDir) {
		return nil, fmt.Errorf("template not found: %s (searched in %s)", templateName, templateDir)
	}

	config, err := loadTemplateConfig(templateDir, templateName)
	if err != nil {
		return nil, err
	}

	tmplType, err := parseTemplateType(config.Format)
	if err != nil {
		return nil, err
	}

	templatePath, err := resolveTemplateFile(templateDir, tmplType, config.TemplateFile)
	if err != nil {
		return nil, err
	}

	return &Template{
		Name:        config.Name,
		Type:        tmplType,
		Path:        templatePath,
		DisplayName: config.DisplayName,
		Description: config.Description,
		Version:     config.Version,
		Author:      config.Author,
		Tags:        config.Tags,
		Config:      config,
	}, nil
}

// ListTemplates returns all available templates
func ListTemplates() ([]Template, error) {
	// Resolve templates directory
	templatesDir, err := utils.ResolveAssetPath("templates")
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

// Generate renders a resume using the specified template name.
// Returns the rendered content as a string.
func (g *Generator) Generate(templateName string, resume *definition.EnhancedResume) (string, error) {
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return "", err
	}
	return g.GenerateWithTemplate(tmpl, resume)
}

// GenerateWithTemplate renders a resume using an already-loaded template.
// Returns the rendered content as a string without re-loading template metadata.
func (g *Generator) GenerateWithTemplate(tmpl *Template, resume *definition.EnhancedResume) (string, error) {
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

func loadTemplateConfig(templateDir, templateName string) (TemplateConfig, error) {
	configPath := filepath.Join(templateDir, "config.yml")
	if !utils.FileExists(configPath) {
		return TemplateConfig{}, fmt.Errorf("template %s is missing config.yml (expected at %s)", templateName, configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return TemplateConfig{}, fmt.Errorf("failed to read template config for %s: %w", templateName, err)
	}

	var cfg TemplateConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return TemplateConfig{}, fmt.Errorf("failed to parse template config for %s: %w", templateName, err)
	}

	cfg.Name = strings.TrimSpace(cfg.Name)
	if cfg.Name == "" {
		cfg.Name = templateName
	}

	if cfg.DisplayName == "" {
		cfg.DisplayName = cfg.Name
	}

	cfg.Format = strings.ToLower(strings.TrimSpace(cfg.Format))
	if cfg.Format == "" {
		return TemplateConfig{}, fmt.Errorf("template %s config missing format", cfg.Name)
	}

	return cfg, nil
}

func parseTemplateType(format string) (TemplateType, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case string(TemplateTypeHTML):
		return TemplateTypeHTML, nil
	case string(TemplateTypeLaTeX):
		return TemplateTypeLaTeX, nil
	default:
		return "", fmt.Errorf("unsupported template format: %s", format)
	}
}

func resolveTemplateFile(templateDir string, tmplType TemplateType, override string) (string, error) {
	filename := strings.TrimSpace(override)
	if filename == "" {
		switch tmplType {
		case TemplateTypeHTML:
			filename = "template.html"
		case TemplateTypeLaTeX:
			filename = "template.tex"
		default:
			return "", fmt.Errorf("unknown template type: %s", tmplType)
		}
	}

	templatePath := filepath.Join(templateDir, filename)
	if !utils.FileExists(templatePath) {
		return "", fmt.Errorf("template file not found at %s", templatePath)
	}

	return templatePath, nil
}
