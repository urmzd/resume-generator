package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

var (
	OutputDir     string
	TemplateNames []string
)

func initRunCmd() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&InputFile, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	runCmd.Flags().StringVarP(&OutputDir, "output-dir", "o", "../outputs", "Root directory where generated resumes will be stored")
	runCmd.Flags().StringVar(&OutputDir, "output-root", "../outputs", "Alias for --output-dir")
	runCmd.Flags().StringSliceVarP(&TemplateNames, "template", "t", nil, "Template name(s). Repeat the flag or use comma-separated values. Defaults to all available templates.")
	runCmd.Flags().StringVarP(&LaTeXEngine, "latex-engine", "e", "", "LaTeX engine to use (xelatex, pdflatex, lualatex, latex). Auto-detects if not specified.")

	runCmd.MarkFlagRequired("input")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate a resume from a data file",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		// Resolve input file path
		inputPath, err := utils.ResolvePath(InputFile)
		if err != nil {
			sugar.Fatalf("Error resolving input path: %s", err)
		}
		if !utils.FileExists(inputPath) {
			sugar.Fatalf("Input file does not exist: %s", inputPath)
		}

		// Load resume data using unified adapter
		inputData, err := resume.LoadResumeFromFile(inputPath)
		if err != nil {
			sugar.Fatalf("Error loading resume data: %s", err)
		}

		// Validate input
		if err := inputData.Validate(); err != nil {
			sugar.Fatalf("Validation error: %s", err)
		}

		// Convert to the runtime resume structure for generation
		resumeData := inputData.ToResume()
		sugar.Infof("Loaded resume for %s (format: %s)", resumeData.Contact.Name, inputData.GetFormat())

		// Generate using unified template system
		generator := generators.NewGenerator(sugar)

		normalizedTemplateNames := sanitizeTemplateNames(TemplateNames)
		selectedTemplates, err := loadSelectedTemplates(normalizedTemplateNames)
		if err != nil {
			sugar.Fatalf("Failed to resolve templates: %v", err)
		}
		if len(selectedTemplates) == 0 {
			sugar.Fatalf("No templates available for generation")
		}
		sugar.Infof("Generating resumes for %d template(s)", len(selectedTemplates))

		// Determine output folder and filenames
		resumeSlug := generateResumeSlug(resumeData)
		currentTime := time.Now()
		dateFolder := currentTime.Format("2006-01-02")

		rootDirInput := strings.TrimSpace(OutputDir)
		resolvedDir, err := utils.ResolvePath(rootDirInput)
		if err != nil {
			sugar.Fatalf("Error resolving output directory: %s", err)
		}
		if resolvedDir == "" {
			if resolvedDir, err = os.Getwd(); err != nil {
				sugar.Fatalf("Failed to determine working directory: %s", err)
			}
		}
		if err := utils.EnsureDir(resolvedDir); err != nil {
			sugar.Fatalf("Error creating output directory: %s", err)
		}

		baseOutputDir := resolvedDir
		desiredPDFBase := defaultResumeBaseName(resumeSlug)
		pdfExt := ".pdf"

		runBaseDir := filepath.Join(baseOutputDir, resumeSlug, dateFolder)
		if err := utils.EnsureDir(runBaseDir); err != nil {
			sugar.Fatalf("Error creating run output directory: %s", err)
		}

		type generationResult struct {
			template string
			tType    generators.TemplateType
			pdfPath  string
			debugDir string
		}

		var results []generationResult

		for _, tmpl := range selectedTemplates {
			templateRunDir, err := resolveTemplateOutputDir(runBaseDir, tmpl)
			if err != nil {
				sugar.Fatalf("Failed to prepare output path for template %s: %v", tmpl.Name, err)
			}

			if err := utils.EnsureDir(templateRunDir); err != nil {
				sugar.Fatalf("Error creating template output directory %s: %v", templateRunDir, err)
			}

			// DOCX has a different flow - it generates bytes directly
			if tmpl.Type == generators.TemplateTypeDOCX {
				docxBytes, err := generator.GenerateDOCX(resumeData)
				if err != nil {
					sugar.Fatalf("Failed to generate DOCX with template %s: %v", tmpl.Name, err)
				}

				docxOutputPath, debugDir, err := ensureUniqueOutputPaths(templateRunDir, desiredPDFBase, ".docx")
				if err != nil {
					sugar.Fatalf("Error determining output filename for template %s: %v", tmpl.Name, err)
				}

				if err := os.WriteFile(docxOutputPath, docxBytes, 0644); err != nil {
					sugar.Fatalf("Failed to write DOCX file: %v", err)
				}

				results = append(results, generationResult{
					template: tmpl.Name,
					tType:    tmpl.Type,
					pdfPath:  docxOutputPath,
					debugDir: debugDir,
				})
				continue
			}

			// Standard template-based generation for HTML and LaTeX
			content, err := generator.GenerateWithTemplate(tmpl, resumeData)
			if err != nil {
				sugar.Fatalf("Failed to generate resume with template %s: %v", tmpl.Name, err)
			}

			pdfOutputPath, debugDir, err := ensureUniqueOutputPaths(templateRunDir, desiredPDFBase, pdfExt)
			if err != nil {
				sugar.Fatalf("Error determining output filename for template %s: %v", tmpl.Name, err)
			}

			if err := utils.EnsureDir(debugDir); err != nil {
				sugar.Fatalf("Error creating debug directory for template %s: %v", tmpl.Name, err)
			}

			templateDir := filepath.Dir(tmpl.Path)

			var compileErr error
			switch tmpl.Type {
			case generators.TemplateTypeLaTeX:
				compileErr = compileLaTeXToPDF(sugar, content, pdfOutputPath, debugDir, templateDir)
			case generators.TemplateTypeHTML:
				compileErr = compileHTMLToPDF(sugar, content, pdfOutputPath, debugDir)
			default:
				sugar.Fatalf("Unknown template type: %s", tmpl.Type)
			}

			if compileErr != nil {
				sugar.Fatalf("Failed to compile template %s: %v", tmpl.Name, compileErr)
			}

			results = append(results, generationResult{
				template: tmpl.Name,
				tType:    tmpl.Type,
				pdfPath:  pdfOutputPath,
				debugDir: debugDir,
			})
		}

		for _, result := range results {
			sugar.Infof("Successfully generated resume (%s) using %s at %s", result.tType, result.template, result.pdfPath)
			sugar.Infof("Render artifacts for %s available in %s", result.template, result.debugDir)
		}
	},
}

// compileHTMLToPDF compiles HTML content to PDF using a Chromium-based browser
func compileHTMLToPDF(logger *zap.SugaredLogger, htmlContent, outputPath, debugDir string) error {
	baseName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	if baseName == "" {
		baseName = "resume"
	}

	debugHTMLPath := filepath.Join(debugDir, baseName+".html")
	if err := os.WriteFile(debugHTMLPath, []byte(htmlContent), 0644); err != nil {
		logger.Warnf("Failed to save HTML debug file: %v", err)
	}

	compiler := compilers.NewHTMLToPDFCompiler(logger)
	return compiler.Compile(htmlContent, outputPath)
}

// compileLaTeXToPDF compiles LaTeX content to PDF using available LaTeX engines
func compileLaTeXToPDF(logger *zap.SugaredLogger, latexContent, outputPath, debugDir, templateDir string) error {
	baseName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	if baseName == "" {
		baseName = "resume"
	}

	resolvedTemplateDir := filepath.Clean(templateDir)
	if resolvedTemplateDir != "" && !utils.DirExists(resolvedTemplateDir) {
		logger.Warnf("Template directory not found at %s, LaTeX compilation may fail", resolvedTemplateDir)
		resolvedTemplateDir = ""
	}

	// Create compiler based on engine selection
	var compiler compilers.Compiler
	if LaTeXEngine != "" {
		// User specified an engine
		logger.Infof("Using specified LaTeX engine: %s", LaTeXEngine)
		compiler = compilers.NewLaTeXCompiler(LaTeXEngine, logger)
	} else {
		// Auto-detect available engine
		autoCompiler, err := compilers.NewAutoLaTeXCompiler(logger)
		if err != nil {
			// List available engines for better error message
			available := compilers.GetAvailableLaTeXEngines()
			if len(available) > 0 {
				return fmt.Errorf("failed to auto-detect LaTeX engine: %w\n\nAvailable engines: %v", err, available)
			}
			return err
		}
		compiler = autoCompiler
	}

	if resolvedTemplateDir != "" {
		compiler.LoadClasses(resolvedTemplateDir)
	}
	compiler.AddOutputFolder(debugDir)

	compiler.Compile(latexContent, baseName)

	// Move compiled PDF to the output location
	generatedPDF := filepath.Join(debugDir, baseName+".pdf")
	if !utils.FileExists(generatedPDF) {
		return fmt.Errorf("expected PDF was not generated at %s", generatedPDF)
	}

	if err := os.Rename(generatedPDF, outputPath); err != nil {
		return fmt.Errorf("failed to move PDF: %w", err)
	}

	return nil
}

func generateResumeSlug(resume *resume.Resume) string {
	nameParts := strings.Fields(resume.Contact.Name)

	var components []string
	if len(nameParts) >= 1 {
		if first := sanitizeNameComponent(nameParts[0]); first != "" {
			components = append(components, first)
		}
	}
	if len(nameParts) >= 3 {
		if middle := sanitizeNameComponent(nameParts[1]); middle != "" {
			components = append(components, middle)
		}
		remaining := sanitizeNameComponent(strings.Join(nameParts[2:], "_"))
		if remaining != "" {
			components = append(components, remaining)
		}
	} else if len(nameParts) >= 2 {
		remaining := sanitizeNameComponent(strings.Join(nameParts[1:], "_"))
		if remaining != "" {
			components = append(components, remaining)
		}
	}

	if len(components) == 0 {
		return "resume"
	}

	return strings.Join(components, "_")
}

func defaultResumeBaseName(resumeSlug string) string {
	slug := strings.TrimSpace(resumeSlug)
	if slug == "" || slug == "resume" {
		return "resume"
	}
	return slug + "_resume"
}

func sanitizeNameComponent(value string) string {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, " ", "_")
	value = strings.ReplaceAll(value, "-", "_")

	var builder strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// sanitizeTemplateNames cleans and normalizes template names
func sanitizeTemplateNames(names []string) []string {
	var result []string
	seen := make(map[string]bool)

	for _, name := range names {
		cleaned := strings.TrimSpace(name)
		if cleaned != "" && !seen[cleaned] {
			result = append(result, cleaned)
			seen[cleaned] = true
		}
	}

	// Sort for consistent ordering
	sort.Strings(result)
	return result
}

// loadSelectedTemplates loads the specified templates or all available templates if none specified
func loadSelectedTemplates(templateNames []string) ([]*generators.Template, error) {
	if len(templateNames) == 0 {
		// Load all available templates
		allTemplates, err := generators.ListTemplates()
		if err != nil {
			return nil, fmt.Errorf("failed to list templates: %w", err)
		}

		// Convert to pointers
		var result []*generators.Template
		for i := range allTemplates {
			result = append(result, &allTemplates[i])
		}

		// Sort by name for consistent ordering
		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		return result, nil
	}

	// Load specified templates
	var templates []*generators.Template
	for _, name := range templateNames {
		tmpl, err := generators.LoadTemplate(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", name, err)
		}
		templates = append(templates, tmpl)
	}

	return templates, nil
}

// resolveTemplateOutputDir computes the output directory for a template
// It creates a subdirectory based on the template name to keep outputs organized
func resolveTemplateOutputDir(runBaseDir string, tmpl *generators.Template) (string, error) {
	// Use the template name as the subdirectory
	templateSubdir := sanitizeNameComponent(tmpl.Name)
	if templateSubdir == "" {
		templateSubdir = "template"
	}

	return filepath.Join(runBaseDir, templateSubdir), nil
}

func ensureUniqueOutputPaths(runDir, desiredBase, extension string) (string, string, error) {
	base := strings.TrimSpace(desiredBase)
	if base == "" {
		base = "resume"
	}

	ext := extension
	if ext == "" {
		ext = ".pdf"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	for attempt := 1; attempt <= 9999; attempt++ {
		suffix := ""
		if attempt > 1 {
			suffix = fmt.Sprintf("_%d", attempt)
		}

		candidateBase := base + suffix
		pdfPath := filepath.Join(runDir, candidateBase+ext)
		debugDir := filepath.Join(runDir, candidateBase+"_debug")

		if !utils.FileExists(pdfPath) && !utils.DirExists(debugDir) {
			return pdfPath, debugDir, nil
		}
	}

	return "", "", fmt.Errorf("failed to find unique output filename in %s", runDir)
}
