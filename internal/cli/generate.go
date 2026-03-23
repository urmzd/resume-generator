package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/resume"
	"github.com/urmzd/incipit/services"
)

var (
	generateOutputDir     string
	generateTemplateNames []string
	generateLatexEngine   string
	generateDryRun        bool
	generateSchema        bool
)

func initGenerateCmd() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&generateOutputDir, "output-dir", "o", "", "Root directory for generated resumes")
	generateCmd.Flags().StringSliceVarP(&generateTemplateNames, "template", "t", nil, "Template ref(s) as name or name:version")
	generateCmd.Flags().StringVarP(&generateLatexEngine, "latex-engine", "e", "", "LaTeX engine (auto, xelatex, pdflatex, lualatex, latex)")
	generateCmd.Flags().BoolVar(&generateDryRun, "dry-run", false, "Validate and preview resume as JSON without generating files")
	generateCmd.Flags().BoolVar(&generateSchema, "schema", false, "Output JSON schema for the resume input format")
}

var generateCmd = &cobra.Command{
	Use:   "generate <file>",
	Short: "Generate resumes from a JSON resume file",
	Long: `Generate resumes from a structured JSON data file using templates.

Accepts .json resume files directly. For unstructured text (.txt, .md), use
"incipit ai create" to convert to JSON first.

Each template produces both its native format and a PDF.

Use --dry-run to validate and preview the resume without generating files.
Use --schema to output the JSON schema for the resume input format.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if generateSchema {
			runSchema()
			return
		}

		if len(args) == 0 {
			outputError(fmt.Errorf("requires a file argument (use --schema for schema output)"))
		}

		inputPath := args[0]

		// Only accept JSON input; point users to ai create for other formats
		ext := strings.ToLower(filepath.Ext(inputPath))
		if ext != ".json" {
			outputError(fmt.Errorf("unsupported input format %q — generate accepts .json files only\n\nFor unstructured text, use: incipit ai create %s", ext, inputPath))
		}

		if generateDryRun {
			runDryRun(inputPath)
			return
		}

		stderrLog("Generating resumes...")
		results, err := services.Generate(services.GenerateOptions{
			InputFile:     inputPath,
			OutputDir:     generateOutputDir,
			TemplateNames: generateTemplateNames,
			LaTeXEngine:   generateLatexEngine,
		})
		if err != nil {
			outputError(err)
		}

		type resultJSON struct {
			Template string `json:"template"`
			Type     string `json:"type"`
			Format   string `json:"format"`
			Path     string `json:"path"`
			Pages    int    `json:"pages,omitempty"`
		}

		out := make([]resultJSON, len(results))
		for i, r := range results {
			out[i] = resultJSON{
				Template: r.Template,
				Type:     string(r.TemplateType),
				Format:   string(r.OutputFormat),
				Path:     r.OutputPath,
				Pages:    r.PageCount,
			}
		}

		outputJSON(map[string]any{"results": out})
	},
}

func runDryRun(inputPath string) {
	data, err := services.LoadResume(inputPath)
	if err != nil {
		outputError(err)
	}

	validationErrors, valErr := services.ValidateResume(inputPath)
	if valErr != nil {
		outputError(valErr)
	}

	valid := len(validationErrors) == 0
	var valErrs []map[string]string
	for _, e := range validationErrors {
		valErrs = append(valErrs, map[string]string{"field": e.Field, "message": e.Message})
	}

	// Resolve which templates would be used
	templateNames := generateTemplateNames
	selectedTemplates, _ := services.LoadSelectedTemplates(templateNames)
	var tmplNames []string
	for _, t := range selectedTemplates {
		tmplNames = append(tmplNames, t.Name)
	}

	outputJSON(map[string]any{
		"resume":        data.Resume,
		"format":        data.Format,
		"section_order": data.SectionOrder,
		"validation":    map[string]any{"valid": valid, "errors": valErrs},
		"templates":     tmplNames,
	})
}

func runSchema() {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&resume.Resume{})

	schema.Title = "Resume Format (v2.0)"
	schema.Description = `Unified resume format used by the CLI.

Input format: JSON (.json). Use "incipit ai create" to convert unstructured text.

Sections: contact, experience, education, projects, skills, certifications, languages.`

	schemaJSON, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		outputError(fmt.Errorf("failed to marshal schema: %w", err))
	}

	fmt.Println(string(schemaJSON))
}
