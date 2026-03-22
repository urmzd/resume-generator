package generators

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tmplmgr "github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// TemplateType represents the template rendering engine
type TemplateType string

const (
	TemplateTypeLaTeX    TemplateType = "latex"
	TemplateTypeHTML     TemplateType = "html"
	TemplateTypeDOCX     TemplateType = "docx"
	TemplateTypeMarkdown TemplateType = "markdown"
)

// Template represents a resume template including metadata from metadata.yml
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

// TemplateConfig contains metadata about a template loaded from metadata.yml
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

func templatesNotFoundError() error {
	configDir := utils.AppConfigDir()
	if configDir == "" {
		configDir = "<config dir>"
	}
	return fmt.Errorf("templates not found; run: incipit init\n"+
		"  Or set RESUME_TEMPLATES_DIR to a directory containing a templates/ folder.\n"+
		"  Expected config: %s/config.yaml",
		configDir)
}

// LoadTemplate loads a template by name.
// It checks: (1) config manifest entries, (2) RESUME_TEMPLATES_DIR / filesystem resolution.
func LoadTemplate(templateName string) (*Template, error) {
	// Check config manifest first
	if cfg, err := tmplmgr.LoadConfig(); err == nil {
		if entry := cfg.Lookup(templateName); entry != nil && utils.DirExists(entry.Path) {
			return loadTemplateFromFS(entry.Path, templateName)
		}
	}

	// Fall back to filesystem resolution (env var, config dir, cwd, exe dir)
	templateDir, err := utils.ResolveAssetPath(filepath.Join("templates", templateName))
	if err == nil && utils.DirExists(templateDir) {
		return loadTemplateFromFS(templateDir, templateName)
	}
	return nil, templatesNotFoundError()
}

// ListTemplates returns all available templates.
// It merges templates from the config manifest and filesystem resolution.
func ListTemplates() ([]Template, error) {
	seen := make(map[string]bool)
	var result []Template

	// Load from config manifest
	if cfg, err := tmplmgr.LoadConfig(); err == nil {
		for _, entry := range cfg.List() {
			if !utils.DirExists(entry.Path) {
				continue
			}
			tmpl, err := loadTemplateFromFS(entry.Path, entry.Name)
			if err != nil {
				continue
			}
			seen[tmpl.Name] = true
			result = append(result, *tmpl)
		}
	}

	// Also check filesystem resolution (for dev workflows with local templates/)
	templatesDir, err := utils.ResolveAssetPath("templates")
	if err == nil && utils.DirExists(templatesDir) {
		fsTmpls, err := listTemplatesFromFS(templatesDir)
		if err == nil {
			for _, tmpl := range fsTmpls {
				if !seen[tmpl.Name] {
					seen[tmpl.Name] = true
					result = append(result, tmpl)
				}
			}
		}
	}

	if len(result) == 0 {
		return nil, templatesNotFoundError()
	}
	return result, nil
}

func listTemplatesFromFS(templatesDir string) ([]Template, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []Template
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		tmpl, err := loadTemplateFromFS(filepath.Join(templatesDir, entry.Name()), entry.Name())
		if err != nil {
			continue
		}
		templates = append(templates, *tmpl)
	}
	return templates, nil
}

func loadTemplateFromFS(templateDir, templateName string) (*Template, error) {
	config, err := loadTemplateConfigFromFS(templateDir, templateName)
	if err != nil {
		return nil, err
	}

	tmplType, err := parseTemplateType(config.Format)
	if err != nil {
		return nil, err
	}

	templatePath, err := resolveTemplateFileFS(templateDir, tmplType, config.TemplateFile)
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

// Generate renders a resume using the specified template name.
func (g *Generator) Generate(templateName string, td *TemplateData) (string, error) {
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return "", err
	}
	return g.GenerateWithTemplate(tmpl, td)
}

// GenerateWithTemplate renders a resume using an already-loaded template.
func (g *Generator) GenerateWithTemplate(tmpl *Template, td *TemplateData) (string, error) {
	g.logger.Infof("Generating resume using template: %s (%s)", tmpl.Name, tmpl.Type)

	content, err := os.ReadFile(tmpl.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	switch tmpl.Type {
	case TemplateTypeHTML:
		return g.renderHTML(string(content), td)
	case TemplateTypeLaTeX:
		return g.renderLaTeX(string(content), td)
	case TemplateTypeMarkdown:
		return g.renderMarkdown(string(content), td)
	case TemplateTypeDOCX:
		return "", fmt.Errorf("DOCX templates produce binary output; use GenerateDOCXWithTemplate instead")
	default:
		return "", fmt.Errorf("unknown template type: %s", tmpl.Type)
	}
}

// renderHTML renders an HTML template
func (g *Generator) renderHTML(templateContent string, td *TemplateData) (string, error) {
	htmlGen := NewHTMLGenerator(g.logger)
	return htmlGen.Generate(templateContent, td)
}

// renderLaTeX renders a LaTeX template
func (g *Generator) renderLaTeX(templateContent string, td *TemplateData) (string, error) {
	latexGen := NewLaTeXGenerator(g.logger)
	return latexGen.Generate(templateContent, td)
}

// renderMarkdown renders a Markdown template
func (g *Generator) renderMarkdown(templateContent string, td *TemplateData) (string, error) {
	mdGen := NewMarkdownGenerator(g.logger)
	return mdGen.Generate(templateContent, td)
}

// GenerateDOCXWithTemplate generates a DOCX document using an already-loaded template.
func (g *Generator) GenerateDOCXWithTemplate(tmpl *Template, td *TemplateData) ([]byte, error) {
	g.logger.Infof("Generating DOCX resume using template: %s", tmpl.Name)

	templateContent, err := os.ReadFile(tmpl.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX template: %w", err)
	}

	scaffoldDir := "docx_scaffold"
	var scaffoldFS fs.FS
	scaffoldFS = os.DirFS(filepath.Dir(tmpl.Path))

	docxGen := NewDOCXGenerator(g.logger)
	return docxGen.Generate(string(templateContent), td, scaffoldFS, scaffoldDir)
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
func FormatTemplateName(name string) string {
	if strings.Contains(name, "-html") || strings.Contains(name, "-latex") || strings.Contains(name, "-markdown") {
		return name
	}

	candidates := []string{
		name + "-html",
		name + "-latex",
		name + "-markdown",
		"modern-html",
		"modern-latex",
	}

	for _, candidate := range candidates {
		if _, err := LoadTemplate(candidate); err == nil {
			return candidate
		}
	}

	return name
}

func loadTemplateConfigFromFS(templateDir, templateName string) (TemplateConfig, error) {
	configPath := filepath.Join(templateDir, "metadata.yml")
	if !utils.FileExists(configPath) {
		return TemplateConfig{}, fmt.Errorf("template %s is missing metadata.yml (expected at %s)", templateName, configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return TemplateConfig{}, fmt.Errorf("failed to read template config for %s: %w", templateName, err)
	}

	return parseTemplateConfig(data, templateName)
}

func parseTemplateConfig(data []byte, templateName string) (TemplateConfig, error) {
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
	case string(TemplateTypeDOCX):
		return TemplateTypeDOCX, nil
	case string(TemplateTypeMarkdown):
		return TemplateTypeMarkdown, nil
	default:
		return "", fmt.Errorf("unsupported template format: %s", format)
	}
}

func resolveTemplateFilename(tmplType TemplateType, override string) string {
	filename := strings.TrimSpace(override)
	if filename != "" {
		return filename
	}
	switch tmplType {
	case TemplateTypeHTML:
		return "template.html"
	case TemplateTypeLaTeX:
		return "template.tex"
	case TemplateTypeMarkdown:
		return "template.md"
	case TemplateTypeDOCX:
		return "template.xml"
	default:
		return ""
	}
}

func resolveTemplateFileFS(templateDir string, tmplType TemplateType, override string) (string, error) {
	filename := resolveTemplateFilename(tmplType, override)
	if filename == "" {
		return "", fmt.Errorf("unknown template type: %s", tmplType)
	}

	templatePath := filepath.Join(templateDir, filename)
	if !utils.FileExists(templatePath) {
		return "", fmt.Errorf("template file not found at %s", templatePath)
	}

	return templatePath, nil
}
