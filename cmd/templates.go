package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/utils"
)

func initTemplatesCmd() {
	templatesCmd.AddCommand(templatesListCmd)
	templatesCmd.AddCommand(templatesValidateCmd)
	templatesCmd.AddCommand(latexEnginesCmd)
	rootCmd.AddCommand(templatesCmd)
}

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Manage templates",
}

var templatesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available Resume Templates:")
		fmt.Println()

		// Use the new template system
		templates, err := generators.ListTemplates()
		if err != nil {
			fmt.Printf("Error listing templates: %v\n", err)
			return
		}

		if len(templates) == 0 {
			fmt.Println("No templates found in templates/")
			return
		}

		// Group by type
		htmlTemplates := []generators.Template{}
		latexTemplates := []generators.Template{}

		for _, tmpl := range templates {
			if tmpl.Type == generators.TemplateTypeHTML {
				htmlTemplates = append(htmlTemplates, tmpl)
			} else {
				latexTemplates = append(latexTemplates, tmpl)
			}
		}

		// Display HTML templates
		if len(htmlTemplates) > 0 {
			fmt.Println("HTML Templates:")
			for _, tmpl := range htmlTemplates {
				name := tmpl.DisplayName
				if name == "" {
					name = tmpl.Name
				}
				fmt.Printf("  📄 %s (%s)\n", name, tmpl.Name)
				if tmpl.Description != "" {
					fmt.Printf("      %s\n", tmpl.Description)
				}
			}
			fmt.Println()
		}

		// Display LaTeX templates
		if len(latexTemplates) > 0 {
			fmt.Println("LaTeX Templates (PDF):")
			for _, tmpl := range latexTemplates {
				name := tmpl.DisplayName
				if name == "" {
					name = tmpl.Name
				}
				fmt.Printf("  📝 %s (%s)\n", name, tmpl.Name)
				if tmpl.Description != "" {
					fmt.Printf("      %s\n", tmpl.Description)
				}
			}
			fmt.Println()
		}

		fmt.Println("Usage:")
		fmt.Println("  resume-generator run -i resume.yml -t modern-html")
		fmt.Println("  resume-generator run -i resume.yml -t modern-latex")
	},
}

var templatesValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a template file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Resolve template path
		templatePath, err := utils.ResolvePath(args[0])
		if err != nil {
			fmt.Printf("Error resolving template path: %v\n", err)
			return
		}

		fmt.Printf("Validating template: %s\n", templatePath)

		// Check if file exists
		if !utils.FileExists(templatePath) {
			fmt.Printf("Error: Template file not found: %s\n", templatePath)
			return
		}

		// Read template content
		content, err := os.ReadFile(templatePath)
		if err != nil {
			fmt.Printf("Error reading template: %v\n", err)
			return
		}

		// Basic validation checks
		templateStr := string(content)
		ext := filepath.Ext(templatePath)

		switch ext {
		case ".html":
			// Check for basic HTML structure
			if !strings.Contains(templateStr, "<!DOCTYPE html>") && !strings.Contains(templateStr, "<html") {
				fmt.Println("Warning: Template appears to be a fragment (no DOCTYPE or html tag)")
			}
			if !strings.Contains(templateStr, "{{") {
				fmt.Println("Warning: Template doesn't appear to use Go template syntax")
			}
			fmt.Println("✓ HTML template appears valid")

		case ".tex", ".ltx":
			// Assume LaTeX template
			if !strings.Contains(templateStr, "\\documentclass") && !strings.Contains(templateStr, "\\begin{document}") {
				fmt.Println("Warning: Template doesn't appear to be a LaTeX document")
			}
			if !strings.Contains(templateStr, "{{") {
				fmt.Println("Warning: Template doesn't appear to use Go template syntax")
			}
			fmt.Println("✓ LaTeX template appears valid")

		default:
			fmt.Printf("Warning: Unknown template type: %s\n", ext)
		}

		fmt.Printf("\nTemplate size: %d bytes\n", len(content))
		fmt.Println("Validation complete!")
	},
}

var latexEnginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "List available LaTeX engines on the system",
	Long:  `List all LaTeX compilation engines available on your system (xelatex, pdflatex, lualatex, latex)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for LaTeX engines...")
		fmt.Println()

		available := compilers.GetAvailableLaTeXEngines()

		if len(available) == 0 {
			fmt.Println("❌ No LaTeX engines found on your system.")
			fmt.Println()
			fmt.Println("To install LaTeX, please use one of the following:")
			fmt.Println("  - TeX Live:   https://www.tug.org/texlive/")
			fmt.Println("  - MiKTeX:     https://miktex.org/")
			fmt.Println("  - MacTeX:     https://www.tug.org/mactex/ (macOS)")
			fmt.Println()
			fmt.Println("Or use Docker which includes LaTeX:")
			fmt.Println("  docker run --rm -v $(pwd):/work texlive/texlive")
			return
		}

		fmt.Printf("✓ Found %d LaTeX engine(s):\n\n", len(available))
		for i, engine := range available {
			prefix := "  "
			if i == 0 {
				prefix = "✓ "
				fmt.Printf("%s%s (default - will be used if no engine is specified)\n", prefix, engine)
			} else {
				fmt.Printf("%s%s\n", prefix, engine)
			}
		}

		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  # Use default engine (auto-detected)")
		fmt.Println("  resume-generator run -i resume.yml -t modern-latex")
		fmt.Println()
		fmt.Println("  # Specify a particular engine")
		fmt.Printf("  resume-generator run -i resume.yml -t modern-latex --latex-engine %s\n", available[0])
	},
}
