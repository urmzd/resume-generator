package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/definition"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

var (
	OutputDir string
	Formats   string
)

var (
	TemplateName string
)

func initRunCmd() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&InputFile, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	runCmd.Flags().StringVarP(&OutputFile, "output", "o", "", "Path to the output file (e.g., ./resume.pdf, ~/Documents/resume.pdf)")
	runCmd.Flags().StringVarP(&OutputDir, "output-dir", "", ".", "Path to the output directory (deprecated, use --output instead)")
	runCmd.Flags().StringVarP(&TemplateName, "template", "t", "modern-html", "Template name (e.g., modern-html, base-latex)")
	runCmd.Flags().StringVarP(&ClassesFolder, "classes", "c", "", "Path to LaTeX classes folder (defaults to assets/classes)")
	runCmd.Flags().StringVarP(&LaTeXEngine, "latex-engine", "e", "", "LaTeX engine to use (xelatex, pdflatex, lualatex, latex). Auto-detects if not specified.")
	runCmd.Flags().StringVarP(&Formats, "formats", "f", "pdf", "Output format: pdf (always PDF, template determines HTML vs LaTeX)")

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
		inputData, err := definition.LoadResumeFromFile(inputPath)
		if err != nil {
			sugar.Fatalf("Error loading resume data: %s", err)
		}

		// Validate input
		if err := inputData.Validate(); err != nil {
			sugar.Fatalf("Validation error: %s", err)
		}

		// Convert to enhanced format for generation
		resume := inputData.ToEnhanced()
		sugar.Infof("Loaded resume for %s (format: %s)", resume.Contact.Name, inputData.GetFormat())

		// Generate using unified template system
		generator := generators.NewGenerator(sugar)

		// Generate content using template
		content, err := generator.Generate(TemplateName, resume)
		if err != nil {
			sugar.Fatalf("Failed to generate resume: %v", err)
		}

		// Determine output folder and filenames
		outputFolderName := generateOutputFolderName(resume)

		var baseOutputDir string
		pdfFileName := "resume.pdf"

		if OutputFile != "" {
			resolvedOutput, err := utils.ResolveOutputPath(OutputFile, true)
			if err != nil {
				sugar.Fatalf("Error resolving output path: %s", err)
			}
			if resolvedOutput == "" {
				sugar.Fatalf("Invalid output path provided")
			}

			original := OutputFile
			if utils.DirExists(resolvedOutput) || strings.HasSuffix(original, string(os.PathSeparator)) || filepath.Ext(original) == "" {
				if err := utils.EnsureDir(resolvedOutput); err != nil {
					sugar.Fatalf("Error ensuring output directory: %s", err)
				}
				baseOutputDir = resolvedOutput
			} else {
				baseOutputDir = filepath.Dir(resolvedOutput)
				if err := utils.EnsureDir(baseOutputDir); err != nil {
					sugar.Fatalf("Error creating output directory: %s", err)
				}
				pdfFileName = filepath.Base(resolvedOutput)
			}
		} else {
			resolvedDir, err := utils.ResolvePath(OutputDir)
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
			baseOutputDir = resolvedDir
		}

		if filepath.Ext(pdfFileName) == "" {
			pdfFileName += ".pdf"
		}

		runOutputDir := filepath.Join(baseOutputDir, outputFolderName)
		if err := utils.EnsureDir(runOutputDir); err != nil {
			sugar.Fatalf("Error creating run output directory: %s", err)
		}

		debugDir := filepath.Join(runOutputDir, "debug")
		if err := utils.EnsureDir(debugDir); err != nil {
			sugar.Fatalf("Error creating debug directory: %s", err)
		}

		pdfOutputPath := filepath.Join(runOutputDir, pdfFileName)

		// Get template type to determine if we need compilation
		tmplType, err := generators.GetTemplateType(TemplateName)
		if err != nil {
			sugar.Fatalf("Failed to get template type: %v", err)
		}

		if tmplType == generators.TemplateTypeLaTeX {
			// Compile LaTeX to PDF
			err = compileLaTeXToPDF(sugar, content, pdfOutputPath, debugDir)
			if err != nil {
				sugar.Fatalf("Failed to compile LaTeX to PDF: %v", err)
			}
		} else if tmplType == generators.TemplateTypeHTML {
			// Compile HTML to PDF
			err = compileHTMLToPDF(sugar, content, pdfOutputPath, debugDir)
			if err != nil {
				sugar.Fatalf("Failed to compile HTML to PDF: %v", err)
			}
		} else {
			sugar.Fatalf("Unknown template type: %s", tmplType)
		}

		sugar.Infof("Successfully generated resume PDF at %s", pdfOutputPath)
		sugar.Infof("Render artifacts available in %s", debugDir)
	},
}

// compileHTMLToPDF compiles HTML content to PDF using chromium/wkhtmltopdf
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
func compileLaTeXToPDF(logger *zap.SugaredLogger, latexContent, outputPath, debugDir string) error {
	baseName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	if baseName == "" {
		baseName = "resume"
	}

	// Resolve classes folder path
	classesPath := ClassesFolder
	if classesPath == "" {
		// Default to assets/classes
		classesPath = "assets/classes"
	}
	resolvedClassesPath, err := utils.ResolveAssetPath(classesPath)
	if err != nil {
		return fmt.Errorf("failed to resolve classes path: %w", err)
	}
	if !utils.DirExists(resolvedClassesPath) {
		logger.Warnf("Classes directory not found at %s, LaTeX compilation may fail", resolvedClassesPath)
	}

	// Create compiler based on engine selection
	var compiler definition.Compiler
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

	compiler.LoadClasses(resolvedClassesPath)
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

// generateOutputFolderName creates a folder name in the format:
// {first}_{optional_middle}_{last}_{iso8601}
func generateOutputFolderName(resume *definition.EnhancedResume) string {
	// Sanitize name parts (replace spaces and special chars with underscores)
	sanitize := func(s string) string {
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, " ", "_")
		// Remove any non-alphanumeric characters except underscores
		var result strings.Builder
		for _, r := range s {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
				result.WriteRune(r)
			}
		}
		return result.String()
	}

	// Extract name parts from the full name
	nameParts := strings.Fields(resume.Contact.Name)
	var first, middle, last string

	if len(nameParts) >= 1 {
		first = sanitize(nameParts[0])
	}
	if len(nameParts) >= 3 {
		// Has middle name
		middle = sanitize(nameParts[1])
		last = sanitize(strings.Join(nameParts[2:], "_"))
	} else if len(nameParts) >= 2 {
		// No middle name
		last = sanitize(strings.Join(nameParts[1:], "_"))
	}

	// Build filename
	var parts []string
	if first != "" {
		parts = append(parts, first)
	}
	if middle != "" {
		parts = append(parts, middle)
	}
	if last != "" {
		parts = append(parts, last)
	}

	// Add ISO8601 timestamp
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	parts = append(parts, timestamp)

	if len(parts) == 1 {
		// No name information, fallback to resume timestamp
		return "resume_" + timestamp
	}

	return strings.Join(parts, "_")
}
