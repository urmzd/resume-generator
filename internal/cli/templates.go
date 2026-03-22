package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/compilers"
	"github.com/urmzd/incipit/generators"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
	"gopkg.in/yaml.v3"
)

var (
	templatesInstallVersion string
	templatesInstallForce   bool
)

func initTemplatesCmd() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesValidateCmd)
	templatesCmd.AddCommand(templatesInstallCmd)
	templatesCmd.AddCommand(templatesAddCmd)
	templatesCmd.AddCommand(templatesRemoveCmd)
	templatesCmd.AddCommand(latexEnginesCmd)
	rootCmd.AddCommand(templatesCmd)

	templatesInstallCmd.Flags().StringVar(&templatesInstallVersion, "version", "", "Template version to install (default: binary version)")
	templatesInstallCmd.Flags().BoolVar(&templatesInstallForce, "force", false, "Overwrite existing templates")
}

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates list")

		// Use the new template system
		templates, err := generators.ListTemplates()
		if err != nil {
			ui.Errorf("Error listing templates: %v", err)
			os.Exit(1)
		}

		if len(templates) == 0 {
			ui.Warn("No templates found in templates/")
			return
		}

		// Group by type
		htmlTemplates := []generators.Template{}
		latexTemplates := []generators.Template{}
		markdownTemplates := []generators.Template{}

		for _, tmpl := range templates {
			switch tmpl.Type {
			case generators.TemplateTypeHTML:
				htmlTemplates = append(htmlTemplates, tmpl)
			case generators.TemplateTypeMarkdown:
				markdownTemplates = append(markdownTemplates, tmpl)
			default:
				latexTemplates = append(latexTemplates, tmpl)
			}
		}

		// Display HTML templates
		if len(htmlTemplates) > 0 {
			ui.Blank()
			ui.Section("HTML Templates")
			for _, tmpl := range htmlTemplates {
				name := tmpl.DisplayName
				if name == "" {
					name = tmpl.Name
				}
				ui.PhaseOk(fmt.Sprintf("%s (%s)", name, tmpl.Name), "")
				if tmpl.Description != "" {
					ui.Detail(tmpl.Description)
				}
			}
		}

		// Display LaTeX templates
		if len(latexTemplates) > 0 {
			ui.Blank()
			ui.Section("LaTeX Templates (PDF)")
			for _, tmpl := range latexTemplates {
				name := tmpl.DisplayName
				if name == "" {
					name = tmpl.Name
				}
				ui.PhaseOk(fmt.Sprintf("%s (%s)", name, tmpl.Name), "")
				if tmpl.Description != "" {
					ui.Detail(tmpl.Description)
				}
			}
		}

		// Display Markdown templates
		if len(markdownTemplates) > 0 {
			ui.Blank()
			ui.Section("Markdown Templates")
			for _, tmpl := range markdownTemplates {
				name := tmpl.DisplayName
				if name == "" {
					name = tmpl.Name
				}
				ui.PhaseOk(fmt.Sprintf("%s (%s)", name, tmpl.Name), "")
				if tmpl.Description != "" {
					ui.Detail(tmpl.Description)
				}
			}
		}

		ui.Blank()
		ui.Info("Usage:")
		ui.Detail("incipit run -i resume.yml -t modern-html")
		ui.Detail("incipit run -i resume.yml -t modern-latex")
		ui.Detail("incipit run -i resume.yml -t modern-markdown")
	},
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a template file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates validate")

		// Resolve template path
		templatePath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving template path: %v", err)
			os.Exit(1)
		}

		ui.Infof("Validating %s", templatePath)

		// Check if file exists
		if !utils.FileExists(templatePath) {
			ui.Errorf("Template file not found: %s", templatePath)
			os.Exit(1)
		}

		// Read template content
		content, err := os.ReadFile(templatePath)
		if err != nil {
			ui.Errorf("Error reading template: %v", err)
			os.Exit(1)
		}

		// Basic validation checks
		templateStr := string(content)
		ext := filepath.Ext(templatePath)

		switch ext {
		case ".html":
			if !strings.Contains(templateStr, "<!DOCTYPE html>") && !strings.Contains(templateStr, "<html") {
				ui.Warn("Template appears to be a fragment (no DOCTYPE or html tag)")
			}
			if !strings.Contains(templateStr, "{{") {
				ui.Warn("Template doesn't appear to use Go template syntax")
			}
			ui.PhaseOk("HTML template appears valid", "")

		case ".tex", ".ltx":
			if !strings.Contains(templateStr, "\\documentclass") && !strings.Contains(templateStr, "\\begin{document}") {
				ui.Warn("Template doesn't appear to be a LaTeX document")
			}
			if !strings.Contains(templateStr, "{{") {
				ui.Warn("Template doesn't appear to use Go template syntax")
			}
			ui.PhaseOk("LaTeX template appears valid", "")

		case ".md":
			if !strings.Contains(templateStr, "{{") {
				ui.Warn("Template doesn't appear to use Go template syntax")
			}
			ui.PhaseOk("Markdown template appears valid", "")

		default:
			ui.Warnf("Unknown template type: %s", ext)
		}

		ui.Infof("Template size: %d bytes", len(content))
		ui.PhaseOk("Validation complete", "")
	},
}

var templatesInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Download and install templates to the config directory",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates install")

		version := templatesInstallVersion
		if version == "" {
			version = Version
		}
		if version == "" || version == "dev" {
			ui.Error("Cannot determine version for template download.")
			ui.Detail("Specify a version with --version, e.g.: incipit templates install --version 1.0.0")
			os.Exit(1)
		}

		ui.Infof("Downloading templates (%s)...", version)

		installedPath, err := templates.Install(templates.InstallOptions{
			Version: version,
			Force:   templatesInstallForce,
		})
		if err != nil {
			ui.Errorf("%v", err)
			os.Exit(1)
		}

		ui.PhaseOk("Templates installed", installedPath)
	},
}

var templatesAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Register a template directory by reference",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates add")

		tmplPath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving path: %v", err)
			os.Exit(1)
		}

		if !utils.DirExists(tmplPath) {
			ui.Errorf("Directory does not exist: %s", tmplPath)
			os.Exit(1)
		}

		// Read metadata.yml to get template name
		metadataPath := filepath.Join(tmplPath, "metadata.yml")
		if !utils.FileExists(metadataPath) {
			ui.Errorf("Template directory is missing metadata.yml: %s", tmplPath)
			os.Exit(1)
		}

		metaData, err := os.ReadFile(metadataPath)
		if err != nil {
			ui.Errorf("Failed to read metadata.yml: %v", err)
			os.Exit(1)
		}

		// Parse full metadata for validation
		var meta struct {
			Name         string `yaml:"name"`
			Format       string `yaml:"format"`
			TemplateFile string `yaml:"template_file,omitempty"`
		}
		if err := yaml.Unmarshal(metaData, &meta); err != nil {
			ui.Errorf("Failed to parse metadata.yml: %v", err)
			os.Exit(1)
		}
		if meta.Name == "" {
			meta.Name = filepath.Base(tmplPath)
		}

		// Validate format
		validFormats := map[string]string{
			"html": "template.html", "latex": "template.tex",
			"markdown": "template.md", "docx": "template.xml",
		}
		if meta.Format == "" {
			ui.Error("metadata.yml is missing required 'format' field")
			os.Exit(1)
		}
		defaultFile, ok := validFormats[strings.ToLower(meta.Format)]
		if !ok {
			ui.Warnf("Unknown format %q in metadata.yml", meta.Format)
		}

		// Verify template file exists
		tmplFile := meta.TemplateFile
		if tmplFile == "" {
			tmplFile = defaultFile
		}
		if tmplFile != "" && !utils.FileExists(filepath.Join(tmplPath, tmplFile)) {
			ui.Errorf("Template file not found: %s", filepath.Join(tmplPath, tmplFile))
			os.Exit(1)
		}

		cfg, err := templates.LoadConfig()
		if err != nil {
			ui.Errorf("Failed to load config: %v", err)
			os.Exit(1)
		}

		if err := cfg.Add(meta.Name, tmplPath); err != nil {
			ui.Errorf("%v", err)
			os.Exit(1)
		}

		if err := templates.SaveConfig(cfg); err != nil {
			ui.Errorf("Failed to save config: %v", err)
			os.Exit(1)
		}

		ui.PhaseOk(fmt.Sprintf("Registered template %q (%s)", meta.Name, meta.Format), tmplPath)
	},
}

var templatesRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Unregister a template (does not delete files)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates remove")

		name := args[0]

		cfg, err := templates.LoadConfig()
		if err != nil {
			ui.Errorf("Failed to load config: %v", err)
			os.Exit(1)
		}

		if err := cfg.Remove(name); err != nil {
			ui.Errorf("%v", err)
			os.Exit(1)
		}

		if err := templates.SaveConfig(cfg); err != nil {
			ui.Errorf("Failed to save config: %v", err)
			os.Exit(1)
		}

		ui.PhaseOk(fmt.Sprintf("Unregistered template %q", name), "")
	},
}

var latexEnginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "List available LaTeX engines on the system",
	Long:  `List all LaTeX compilation engines available on your system (xelatex, pdflatex, lualatex, latex)`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit templates engines")
		ui.Info("Checking for LaTeX engines...")

		available := compilers.GetAvailableLaTeXEngines()

		if len(available) == 0 {
			ui.Blank()
			ui.Error("No LaTeX engines found on your system.")
			ui.Blank()
			ui.Info("To install LaTeX, use one of the following:")
			ui.Detail("TeX Live:   https://www.tug.org/texlive/")
			ui.Detail("MiKTeX:     https://miktex.org/")
			ui.Detail("MacTeX:     https://www.tug.org/mactex/ (macOS)")
			return
		}

		ui.Blank()
		ui.PhaseOk(fmt.Sprintf("Found %d LaTeX engine(s)", len(available)), "")
		for i, engine := range available {
			if i == 0 {
				ui.Detail(fmt.Sprintf("%s (default - will be used if no engine is specified)", engine))
			} else {
				ui.Detail(engine)
			}
		}

		ui.Blank()
		ui.Info("Usage:")
		ui.Detail("# Use default engine (auto-detected)")
		ui.Detail("incipit run -i resume.yml -t modern-latex")
		ui.Blank()
		ui.Detail("# Specify a particular engine")
		ui.Detail(fmt.Sprintf("incipit run -i resume.yml -t modern-latex --latex-engine %s", available[0]))
	},
}
