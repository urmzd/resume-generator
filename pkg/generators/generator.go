package generators

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// TemplateType represents the template rendering engine
type TemplateType string

const (
	TemplateTypeLaTeX TemplateType = "latex"
	TemplateTypeHTML  TemplateType = "html"
	TemplateTypeDOCX  TemplateType = "docx"
)

// Template represents a resume template including metadata from config.yml
type Template struct {
	Name        string
	Type        TemplateType
	Path        string // filesystem path (empty when embedded)
	DisplayName string
	Description string
	Version     string
	Author      string
	Tags        []string
	Config      TemplateConfig
	Embedded    bool   // true when loaded from embedded FS
	EmbeddedDir string // e.g. "templates/modern-latex" for embedded reads
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

var embeddedFS embed.FS

// SetEmbeddedFS sets the embedded filesystem used to load templates.
func SetEmbeddedFS(efs embed.FS) {
	embeddedFS = efs
}

// NewGenerator creates a new template-based generator
func NewGenerator(logger *zap.SugaredLogger) *Generator {
	return &Generator{logger: logger}
}

// LoadTemplate loads a template by name.
// It tries the filesystem first (via RESUME_TEMPLATES_DIR or local templates/ dir),
// then falls back to the embedded FS.
func LoadTemplate(templateName string) (*Template, error) {
	// Try filesystem first
	templateDir, err := utils.ResolveAssetPath(filepath.Join("templates", templateName))
	if err == nil && utils.DirExists(templateDir) {
		return loadTemplateFromFS(templateDir, templateName)
	}

	// Fall back to embedded FS
	embeddedDir := "templates/" + templateName
	return loadTemplateFromEmbed(embeddedDir, templateName)
}

// ListTemplates returns all available templates.
// It tries the filesystem first, then falls back to the embedded FS.
func ListTemplates() ([]Template, error) {
	// Try filesystem first
	templatesDir, err := utils.ResolveAssetPath("templates")
	if err == nil && utils.DirExists(templatesDir) {
		return listTemplatesFromFS(templatesDir)
	}

	// Fall back to embedded FS
	return listTemplatesFromEmbed()
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

func listTemplatesFromEmbed() ([]Template, error) {
	entries, err := fs.ReadDir(embeddedFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded templates: %w", err)
	}

	var templates []Template
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		embeddedDir := "templates/" + entry.Name()
		tmpl, err := loadTemplateFromEmbed(embeddedDir, entry.Name())
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
		Embedded:    false,
	}, nil
}

func loadTemplateFromEmbed(embeddedDir, templateName string) (*Template, error) {
	config, err := loadTemplateConfigFromEmbed(embeddedDir, templateName)
	if err != nil {
		return nil, err
	}

	tmplType, err := parseTemplateType(config.Format)
	if err != nil {
		return nil, err
	}

	// Verify template file exists in embedded FS (skip for DOCX)
	if tmplType != TemplateTypeDOCX {
		filename := resolveTemplateFilename(tmplType, config.TemplateFile)
		embeddedPath := embeddedDir + "/" + filename
		if _, err := fs.Stat(embeddedFS, embeddedPath); err != nil {
			return nil, fmt.Errorf("embedded template file not found: %s", embeddedPath)
		}
	}

	return &Template{
		Name:        config.Name,
		Type:        tmplType,
		Path:        "",
		DisplayName: config.DisplayName,
		Description: config.Description,
		Version:     config.Version,
		Author:      config.Author,
		Tags:        config.Tags,
		Config:      config,
		Embedded:    true,
		EmbeddedDir: embeddedDir,
	}, nil
}

// Generate renders a resume using the specified template name.
func (g *Generator) Generate(templateName string, resume *resume.Resume) (string, error) {
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		return "", err
	}
	return g.GenerateWithTemplate(tmpl, resume)
}

// GenerateWithTemplate renders a resume using an already-loaded template.
func (g *Generator) GenerateWithTemplate(tmpl *Template, resume *resume.Resume) (string, error) {
	g.logger.Infof("Generating resume using template: %s (%s)", tmpl.Name, tmpl.Type)

	var content []byte
	var err error

	if tmpl.Embedded {
		filename := resolveTemplateFilename(tmpl.Type, tmpl.Config.TemplateFile)
		content, err = fs.ReadFile(embeddedFS, tmpl.EmbeddedDir+"/"+filename)
	} else {
		content, err = os.ReadFile(tmpl.Path)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

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
func (g *Generator) renderHTML(templateContent string, resume *resume.Resume) (string, error) {
	htmlGen := NewHTMLGenerator(g.logger)
	return htmlGen.Generate(templateContent, resume)
}

// renderLaTeX renders a LaTeX template
func (g *Generator) renderLaTeX(templateContent string, resume *resume.Resume) (string, error) {
	latexGen := NewLaTeXGenerator(g.logger)
	return latexGen.Generate(templateContent, resume)
}

// GenerateDOCX generates a DOCX document from the resume.
func (g *Generator) GenerateDOCX(resume *resume.Resume) ([]byte, error) {
	g.logger.Info("Generating DOCX resume")
	docxGen := NewDOCXGenerator(g.logger)
	return docxGen.Generate(resume)
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
	if strings.Contains(name, "-html") || strings.Contains(name, "-latex") {
		return name
	}

	candidates := []string{
		name + "-html",
		name + "-latex",
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

// ExtractEmbeddedTemplateDir extracts all files from an embedded template directory
// to a temporary directory on disk. Returns the temp directory path.
func ExtractEmbeddedTemplateDir(embeddedDir string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "resume-template-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	entries, err := fs.ReadDir(embeddedFS, embeddedDir)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", fmt.Errorf("failed to read embedded dir %s: %w", embeddedDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := fs.ReadFile(embeddedFS, embeddedDir+"/"+entry.Name())
		if err != nil {
			_ = os.RemoveAll(tmpDir)
			return "", fmt.Errorf("failed to read embedded file %s: %w", entry.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, entry.Name()), data, 0644); err != nil {
			_ = os.RemoveAll(tmpDir)
			return "", fmt.Errorf("failed to write file %s: %w", entry.Name(), err)
		}
	}

	return tmpDir, nil
}

func loadTemplateConfigFromFS(templateDir, templateName string) (TemplateConfig, error) {
	configPath := filepath.Join(templateDir, "config.yml")
	if !utils.FileExists(configPath) {
		return TemplateConfig{}, fmt.Errorf("template %s is missing config.yml (expected at %s)", templateName, configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return TemplateConfig{}, fmt.Errorf("failed to read template config for %s: %w", templateName, err)
	}

	return parseTemplateConfig(data, templateName)
}

func loadTemplateConfigFromEmbed(embeddedDir, templateName string) (TemplateConfig, error) {
	configPath := embeddedDir + "/config.yml"
	data, err := fs.ReadFile(embeddedFS, configPath)
	if err != nil {
		return TemplateConfig{}, fmt.Errorf("embedded template %s is missing config.yml", templateName)
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
	default:
		return ""
	}
}

func resolveTemplateFileFS(templateDir string, tmplType TemplateType, override string) (string, error) {
	if tmplType == TemplateTypeDOCX {
		return "", nil
	}

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
