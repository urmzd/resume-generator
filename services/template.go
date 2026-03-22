package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/incipit/compilers"
	"github.com/urmzd/incipit/generators"
	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
	"gopkg.in/yaml.v3"
)

// TemplateInfo groups templates by type for listing.
type TemplateInfo struct {
	HTML     []generators.Template
	LaTeX    []generators.Template
	Markdown []generators.Template
}

// ListTemplates returns all available templates grouped by type.
func ListTemplates() (*TemplateInfo, error) {
	tmpls, err := generators.ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("error listing templates: %w", err)
	}

	info := &TemplateInfo{}
	for _, tmpl := range tmpls {
		switch tmpl.Type {
		case generators.TemplateTypeHTML:
			info.HTML = append(info.HTML, tmpl)
		case generators.TemplateTypeMarkdown:
			info.Markdown = append(info.Markdown, tmpl)
		default:
			info.LaTeX = append(info.LaTeX, tmpl)
		}
	}
	return info, nil
}

// TemplateValidationResult holds the result of template file validation.
type TemplateValidationResult struct {
	Format   string
	Size     int
	Warnings []string
	Valid    bool
}

// ValidateTemplate performs basic validation on a template file.
func ValidateTemplate(filePath string) (*TemplateValidationResult, error) {
	templatePath, err := utils.ResolvePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("error resolving template path: %w", err)
	}
	if !utils.FileExists(templatePath) {
		return nil, fmt.Errorf("template file not found: %s", templatePath)
	}

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template: %w", err)
	}

	templateStr := string(content)
	ext := filepath.Ext(templatePath)

	result := &TemplateValidationResult{
		Size:  len(content),
		Valid: true,
	}

	switch ext {
	case ".html":
		result.Format = "html"
		if !strings.Contains(templateStr, "<!DOCTYPE html>") && !strings.Contains(templateStr, "<html") {
			result.Warnings = append(result.Warnings, "Template appears to be a fragment (no DOCTYPE or html tag)")
		}
		if !strings.Contains(templateStr, "{{") {
			result.Warnings = append(result.Warnings, "Template doesn't appear to use Go template syntax")
		}
	case ".tex", ".ltx":
		result.Format = "latex"
		if !strings.Contains(templateStr, "\\documentclass") && !strings.Contains(templateStr, "\\begin{document}") {
			result.Warnings = append(result.Warnings, "Template doesn't appear to be a LaTeX document")
		}
		if !strings.Contains(templateStr, "{{") {
			result.Warnings = append(result.Warnings, "Template doesn't appear to use Go template syntax")
		}
	case ".md":
		result.Format = "markdown"
		if !strings.Contains(templateStr, "{{") {
			result.Warnings = append(result.Warnings, "Template doesn't appear to use Go template syntax")
		}
	default:
		result.Format = ext
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown template type: %s", ext))
	}

	return result, nil
}

// InstallTemplates downloads and registers templates from a GitHub release.
func InstallTemplates(version string, force bool) ([]templates.InstalledTemplate, error) {
	if version == "" {
		latest, err := templates.LatestVersion()
		if err != nil {
			return nil, fmt.Errorf("cannot determine version for template download: %w", err)
		}
		version = latest
	}

	installed, err := templates.Install(templates.InstallOptions{
		Version: version,
		Force:   force,
	})
	if err != nil {
		return nil, err
	}

	cfg, err := templates.LoadConfig()
	if err != nil {
		cfg = &templates.Config{}
	}

	for _, tmpl := range installed {
		_ = cfg.Add(tmpl.Name, tmpl.Version, tmpl.Path)
	}

	if err := templates.SaveConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return installed, nil
}

// UpdateTemplates updates default templates to the latest release.
func UpdateTemplates() ([]templates.InstalledTemplate, error) {
	version, err := templates.LatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to determine latest version: %w", err)
	}

	installed, err := templates.Install(templates.InstallOptions{
		Version: version,
		Force:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to install templates: %w", err)
	}

	cfg, err := templates.LoadConfig()
	if err != nil {
		cfg = &templates.Config{}
	}

	for _, tmpl := range installed {
		_ = cfg.Add(tmpl.Name, tmpl.Version, tmpl.Path)
	}

	if err := templates.SaveConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return installed, nil
}

// AddTemplate registers a local template directory by reference.
func AddTemplate(dirPath string) (name, format string, err error) {
	tmplPath, err := utils.ResolvePath(dirPath)
	if err != nil {
		return "", "", fmt.Errorf("error resolving path: %w", err)
	}
	if !utils.DirExists(tmplPath) {
		return "", "", fmt.Errorf("directory does not exist: %s", tmplPath)
	}

	metadataPath := filepath.Join(tmplPath, "metadata.yml")
	if !utils.FileExists(metadataPath) {
		return "", "", fmt.Errorf("template directory is missing metadata.yml: %s", tmplPath)
	}

	metaData, err := os.ReadFile(metadataPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read metadata.yml: %w", err)
	}

	var meta struct {
		Name         string `yaml:"name"`
		Format       string `yaml:"format"`
		TemplateFile string `yaml:"template_file,omitempty"`
	}
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return "", "", fmt.Errorf("failed to parse metadata.yml: %w", err)
	}
	if meta.Name == "" {
		meta.Name = filepath.Base(tmplPath)
	}

	validFormats := map[string]string{
		"html": "template.html", "latex": "template.tex",
		"markdown": "template.md", "docx": "template.xml",
	}
	if meta.Format == "" {
		return "", "", fmt.Errorf("metadata.yml is missing required 'format' field")
	}
	defaultFile, ok := validFormats[strings.ToLower(meta.Format)]
	if !ok {
		return "", "", fmt.Errorf("unknown format %q in metadata.yml", meta.Format)
	}

	tmplFile := meta.TemplateFile
	if tmplFile == "" {
		tmplFile = defaultFile
	}
	if tmplFile != "" && !utils.FileExists(filepath.Join(tmplPath, tmplFile)) {
		return "", "", fmt.Errorf("template file not found: %s", filepath.Join(tmplPath, tmplFile))
	}

	cfg, err := templates.LoadConfig()
	if err != nil {
		return "", "", fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Add(meta.Name, "", tmplPath); err != nil {
		return "", "", err
	}

	if err := templates.SaveConfig(cfg); err != nil {
		return "", "", fmt.Errorf("failed to save config: %w", err)
	}

	return meta.Name, meta.Format, nil
}

// RemoveTemplate unregisters a template by name or name:version reference.
func RemoveTemplate(ref string) error {
	name, version := templates.ParseTemplateRef(ref)

	cfg, err := templates.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Remove(name, version); err != nil {
		return err
	}

	if err := templates.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// ListLaTeXEngines returns the available LaTeX engines on the system.
func ListLaTeXEngines() []string {
	return compilers.GetAvailableLaTeXEngines()
}
