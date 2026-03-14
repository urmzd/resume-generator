package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
	defaultOut := utils.DefaultOutputDir()
	runCmd.Flags().StringVarP(&OutputDir, "output-dir", "o", defaultOut, "Root directory where generated resumes will be stored")
	runCmd.Flags().StringVar(&OutputDir, "output-root", defaultOut, "Alias for --output-dir")
	runCmd.Flags().StringSliceVarP(&TemplateNames, "template", "t", nil, "Template name(s). Repeat the flag or use comma-separated values. Defaults to all available templates.")
	runCmd.Flags().StringVarP(&LaTeXEngine, "latex-engine", "e", "", "LaTeX engine to use (xelatex, pdflatex, lualatex, latex). Auto-detects if not specified.")

	_ = runCmd.MarkFlagRequired("input")

	generators.SetEmbeddedFS(EmbeddedTemplatesFS)
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
		resumeSlug := generateFilenameSlug(inputPath)
		currentTime := time.Now()

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

		desiredBase := generateOutputBaseName(resumeData.Contact.Name)
		pdfExt := ".pdf"

		// Create timestamped run directory: <root>/<slug>/<YYYY-MM-DD_HH-MM>/
		runDir := generateRunDir(filepath.Join(resolvedDir, resumeSlug), currentTime)
		if err := utils.EnsureDir(runDir); err != nil {
			sugar.Fatalf("Error creating run output directory: %s", err)
		}

		// Pre-load the HTML fallback template for DOCX->PDF conversion
		htmlFallbackTmpl, htmlFallbackErr := generators.LoadTemplate("modern-html")
		if htmlFallbackErr != nil {
			sugar.Warnf("Could not load HTML fallback template for DOCX PDF generation: %v", htmlFallbackErr)
		}

		type generationResult struct {
			template string
			tType    generators.TemplateType
			outPath  string
		}

		var results []generationResult

		for _, tmpl := range selectedTemplates {
			// Markdown outputs a .md file directly (no PDF compilation)
			if tmpl.Type == generators.TemplateTypeMarkdown {
				content, err := generator.GenerateWithTemplate(tmpl, resumeData)
				if err != nil {
					sugar.Fatalf("Failed to generate Markdown with template %s: %v", tmpl.Name, err)
				}

				mdOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".md")
				if err != nil {
					sugar.Fatalf("Error determining output filename for template %s: %v", tmpl.Name, err)
				}

				if err := os.WriteFile(mdOutputPath, []byte(content), 0644); err != nil {
					sugar.Fatalf("Failed to write Markdown file: %v", err)
				}

				results = append(results, generationResult{
					template: tmpl.Name,
					tType:    tmpl.Type,
					outPath:  mdOutputPath,
				})
				continue
			}

			// DOCX has a different flow - it generates bytes directly
			if tmpl.Type == generators.TemplateTypeDOCX {
				docxBytes, err := generator.GenerateDOCX(resumeData)
				if err != nil {
					sugar.Fatalf("Failed to generate DOCX with template %s: %v", tmpl.Name, err)
				}

				docxOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".docx")
				if err != nil {
					sugar.Fatalf("Error determining output filename for template %s: %v", tmpl.Name, err)
				}

				if err := os.WriteFile(docxOutputPath, docxBytes, 0644); err != nil {
					sugar.Fatalf("Failed to write DOCX file: %v", err)
				}

				// Also generate a PDF via the HTML fallback template
				if htmlFallbackTmpl != nil {
					htmlContent, htmlErr := generator.GenerateWithTemplate(htmlFallbackTmpl, resumeData)
					if htmlErr != nil {
						sugar.Warnf("Failed to generate HTML for DOCX PDF fallback: %v", htmlErr)
					} else {
						pdfOutputPath := strings.TrimSuffix(docxOutputPath, ".docx") + ".pdf"
						debugDir, debugErr := os.MkdirTemp("", "resume-debug-*")
						if debugErr != nil {
							sugar.Warnf("Failed to create temp debug dir for DOCX PDF: %v", debugErr)
						} else {
							if pdfErr := compileHTMLToPDF(sugar, htmlContent, pdfOutputPath, debugDir); pdfErr != nil {
								// Keep debug dir on failure
								persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
								if mvErr := os.Rename(debugDir, persistedDebug); mvErr != nil {
									sugar.Warnf("Failed to persist debug dir: %v (temp dir: %s)", mvErr, debugDir)
								} else {
									sugar.Warnf("Failed to generate PDF for DOCX template %s: %v (debug: %s)", tmpl.Name, pdfErr, persistedDebug)
								}
							} else {
								_ = os.RemoveAll(debugDir)
								sugar.Infof("Generated PDF alongside DOCX: %s", pdfOutputPath)
							}
						}
					}
				}

				results = append(results, generationResult{
					template: tmpl.Name,
					tType:    tmpl.Type,
					outPath:  docxOutputPath,
				})
				continue
			}

			// Standard template-based generation for HTML and LaTeX
			content, err := generator.GenerateWithTemplate(tmpl, resumeData)
			if err != nil {
				sugar.Fatalf("Failed to generate resume with template %s: %v", tmpl.Name, err)
			}

			pdfOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, pdfExt)
			if err != nil {
				sugar.Fatalf("Error determining output filename for template %s: %v", tmpl.Name, err)
			}

			// Use a temp directory for debug artifacts; only persist on failure
			debugDir, err := os.MkdirTemp("", "resume-debug-*")
			if err != nil {
				sugar.Fatalf("Failed to create temp debug directory: %v", err)
			}

			var templateDir string
			if tmpl.Embedded && tmpl.EmbeddedDir != "" {
				extractedDir, extractErr := generators.ExtractEmbeddedTemplateDir(tmpl.EmbeddedDir)
				if extractErr != nil {
					sugar.Fatalf("Failed to extract embedded template files for %s: %v", tmpl.Name, extractErr)
				}
				defer func() { _ = os.RemoveAll(extractedDir) }()
				templateDir = extractedDir
			} else {
				templateDir = filepath.Dir(tmpl.Path)
			}

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
				// Persist debug dir next to output on failure
				persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
				if mvErr := os.Rename(debugDir, persistedDebug); mvErr != nil {
					sugar.Warnf("Failed to persist debug dir: %v (temp dir: %s)", mvErr, debugDir)
				}
				sugar.Fatalf("Failed to compile template %s: %v (debug: %s)", tmpl.Name, compileErr, persistedDebug)
			}

			// Success: clean up debug artifacts
			_ = os.RemoveAll(debugDir)

			results = append(results, generationResult{
				template: tmpl.Name,
				tType:    tmpl.Type,
				outPath:  pdfOutputPath,
			})
		}

		for _, result := range results {
			sugar.Infof("Successfully generated resume (%s) using %s at %s", result.tType, result.template, result.outPath)

			// Warn if the generated PDF exceeds one page
			if strings.HasSuffix(result.outPath, ".pdf") {
				if pdfData, readErr := os.ReadFile(result.outPath); readErr == nil {
					if pages := compilers.CountPDFPages(pdfData); pages > 1 {
						sugar.Warnf("Resume generated with template %s has %d pages (exceeds 1 page)", result.template, pages)
					}
				}
			}
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

	compiler := compilers.NewRodHTMLToPDFCompiler(logger)
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
