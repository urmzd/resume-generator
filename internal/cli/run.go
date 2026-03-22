package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/compilers"
	"github.com/urmzd/incipit/generators"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/resume"
	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
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
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Generate a resume from a data file",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		ui.Header("incipit run")

		// Resolve input file path
		inputPath, err := utils.ResolvePath(InputFile)
		if err != nil {
			ui.Errorf("Error resolving input path: %s", err)
			os.Exit(1)
		}
		if !utils.FileExists(inputPath) {
			ui.Errorf("Input file does not exist: %s", inputPath)
			os.Exit(1)
		}

		// Load resume data using unified adapter
		inputData, err := resume.LoadResumeFromFile(inputPath)
		if err != nil {
			ui.Errorf("Error loading resume data: %s", err)
			os.Exit(1)
		}

		// Validate input
		if err := inputData.Validate(); err != nil {
			ui.Errorf("Validation error: %s", err)
			os.Exit(1)
		}

		// Convert to the runtime resume structure for generation
		resumeData := inputData.ToResume()
		sectionOrder := inputData.GetSectionOrder()
		td := generators.NewTemplateData(resumeData, sectionOrder)
		ui.PhaseOk("Loaded resume", fmt.Sprintf("%s (format: %s)", resumeData.Contact.Name, inputData.GetFormat()))

		// Apply config defaults for unset flags
		if cfg, cfgErr := templates.LoadConfig(); cfgErr == nil {
			if !cmd.Flags().Changed("template") && len(cfg.Defaults.Templates) > 0 {
				TemplateNames = cfg.Defaults.Templates
			}
			if !cmd.Flags().Changed("latex-engine") && cfg.Defaults.LaTeXEngine != "" {
				LaTeXEngine = cfg.Defaults.LaTeXEngine
			}
			if !cmd.Flags().Changed("output-dir") && !cmd.Flags().Changed("output-root") && cfg.Defaults.OutputDir != "" {
				OutputDir = cfg.Defaults.OutputDir
			}
		}

		// Generate using unified template system
		generator := generators.NewGenerator(sugar)

		normalizedTemplateNames := sanitizeTemplateNames(TemplateNames)
		selectedTemplates, err := loadSelectedTemplates(normalizedTemplateNames)
		if err != nil {
			ui.Errorf("Failed to resolve templates: %v", err)
			os.Exit(1)
		}
		if len(selectedTemplates) == 0 {
			ui.Error("No templates available for generation")
			os.Exit(1)
		}
		ui.Infof("Generating resumes for %d template(s)", len(selectedTemplates))

		// Determine output folder and filenames
		resumeSlug := generateFilenameSlug(inputPath)
		currentTime := time.Now()

		rootDirInput := strings.TrimSpace(OutputDir)
		resolvedDir, err := utils.ResolvePath(rootDirInput)
		if err != nil {
			ui.Errorf("Error resolving output directory: %s", err)
			os.Exit(1)
		}
		if resolvedDir == "" {
			if resolvedDir, err = os.Getwd(); err != nil {
				ui.Errorf("Failed to determine working directory: %s", err)
				os.Exit(1)
			}
		}
		if err := utils.EnsureDir(resolvedDir); err != nil {
			ui.Errorf("Error creating output directory: %s", err)
			os.Exit(1)
		}

		desiredBase := generateOutputBaseName(resumeData.Contact.Name)
		pdfExt := ".pdf"

		// Create timestamped run directory: <root>/<slug>/<YYYY-MM-DD_HH-MM>/
		runDir := generateRunDir(filepath.Join(resolvedDir, resumeSlug), currentTime)
		if err := utils.EnsureDir(runDir); err != nil {
			ui.Errorf("Error creating run output directory: %s", err)
			os.Exit(1)
		}

		// Pre-load the HTML fallback template for DOCX->PDF conversion
		htmlFallbackTmpl, htmlFallbackErr := generators.LoadTemplate("modern-html")
		if htmlFallbackErr != nil {
			ui.Warn("Could not load HTML fallback template for DOCX PDF generation")
			sugar.Debugf("HTML fallback error: %v", htmlFallbackErr)
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
				content, err := generator.GenerateWithTemplate(tmpl, td)
				if err != nil {
					ui.Errorf("Failed to generate Markdown with template %s: %v", tmpl.Name, err)
					os.Exit(1)
				}

				mdOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".md")
				if err != nil {
					ui.Errorf("Error determining output filename for template %s: %v", tmpl.Name, err)
					os.Exit(1)
				}

				if err := os.WriteFile(mdOutputPath, []byte(content), 0644); err != nil {
					ui.Errorf("Failed to write Markdown file: %v", err)
					os.Exit(1)
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
				docxBytes, err := generator.GenerateDOCXWithTemplate(tmpl, td)
				if err != nil {
					ui.Errorf("Failed to generate DOCX with template %s: %v", tmpl.Name, err)
					os.Exit(1)
				}

				docxOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".docx")
				if err != nil {
					ui.Errorf("Error determining output filename for template %s: %v", tmpl.Name, err)
					os.Exit(1)
				}

				if err := os.WriteFile(docxOutputPath, docxBytes, 0644); err != nil {
					ui.Errorf("Failed to write DOCX file: %v", err)
					os.Exit(1)
				}

				// Also generate a PDF via the HTML fallback template
				if htmlFallbackTmpl != nil {
					htmlContent, htmlErr := generator.GenerateWithTemplate(htmlFallbackTmpl, td)
					if htmlErr != nil {
						ui.Warnf("Failed to generate HTML for DOCX PDF fallback: %v", htmlErr)
					} else {
						pdfOutputPath := strings.TrimSuffix(docxOutputPath, ".docx") + ".pdf"
						debugDir, debugErr := os.MkdirTemp("", "resume-debug-*")
						if debugErr != nil {
							ui.Warnf("Failed to create temp debug dir for DOCX PDF: %v", debugErr)
						} else {
							if pdfErr := compileHTMLToPDF(sugar, htmlContent, pdfOutputPath, debugDir); pdfErr != nil {
								persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
								if mvErr := os.Rename(debugDir, persistedDebug); mvErr != nil {
									ui.Warnf("Failed to persist debug dir: %v (temp dir: %s)", mvErr, debugDir)
								} else {
									ui.Warnf("Failed to generate PDF for DOCX template %s: %v", tmpl.Name, pdfErr)
								}
							} else {
								_ = os.RemoveAll(debugDir)
								ui.PhaseOk("Generated PDF alongside DOCX", pdfOutputPath)
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
			content, err := generator.GenerateWithTemplate(tmpl, td)
			if err != nil {
				ui.Errorf("Failed to generate resume with template %s: %v", tmpl.Name, err)
				os.Exit(1)
			}

			pdfOutputPath, err := ensureUniqueOutputPath(runDir, desiredBase, tmpl.Name, pdfExt)
			if err != nil {
				ui.Errorf("Error determining output filename for template %s: %v", tmpl.Name, err)
				os.Exit(1)
			}

			// Use a temp directory for debug artifacts; only persist on failure
			debugDir, err := os.MkdirTemp("", "resume-debug-*")
			if err != nil {
				ui.Errorf("Failed to create temp debug directory: %v", err)
				os.Exit(1)
			}

			templateDir := filepath.Dir(tmpl.Path)

			var compileErr error
			switch tmpl.Type {
			case generators.TemplateTypeLaTeX:
				compileErr = compileLaTeXToPDF(sugar, content, pdfOutputPath, debugDir, templateDir)
			case generators.TemplateTypeHTML:
				compileErr = compileHTMLToPDF(sugar, content, pdfOutputPath, debugDir)
			default:
				ui.Errorf("Unknown template type: %s", tmpl.Type)
				os.Exit(1)
			}

			if compileErr != nil {
				// Persist debug dir next to output on failure
				persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
				if mvErr := os.Rename(debugDir, persistedDebug); mvErr != nil {
					ui.Warnf("Failed to persist debug dir: %v (temp dir: %s)", mvErr, debugDir)
				}
				ui.Errorf("Failed to compile template %s: %v", tmpl.Name, compileErr)
				ui.Infof("Debug artifacts: %s", persistedDebug)
				os.Exit(1)
			}

			// Success: clean up debug artifacts
			_ = os.RemoveAll(debugDir)

			results = append(results, generationResult{
				template: tmpl.Name,
				tType:    tmpl.Type,
				outPath:  pdfOutputPath,
			})
		}

		ui.Blank()
		for _, result := range results {
			ui.PhaseOk(fmt.Sprintf("Generated (%s) using %s", result.tType, result.template), result.outPath)

			// Warn if the generated PDF exceeds one page
			if strings.HasSuffix(result.outPath, ".pdf") {
				if pdfData, readErr := os.ReadFile(result.outPath); readErr == nil {
					if pages := compilers.CountPDFPages(pdfData); pages > 1 {
						ui.Warnf("Resume with template %s has %d pages (exceeds 1 page)", result.template, pages)
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
		ui.Warnf("Template directory not found at %s, LaTeX compilation may fail", resolvedTemplateDir)
		resolvedTemplateDir = ""
	}

	// Create compiler based on engine selection
	var compiler compilers.Compiler
	if LaTeXEngine != "" {
		// User specified an engine
		ui.Infof("Using specified LaTeX engine: %s", LaTeXEngine)
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
