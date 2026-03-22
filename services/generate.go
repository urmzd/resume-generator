package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/urmzd/incipit/compilers"
	"github.com/urmzd/incipit/generators"
	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
	"go.uber.org/zap"
)

// GenerateOptions configures resume generation.
type GenerateOptions struct {
	InputFile     string
	OutputDir     string
	TemplateNames []string
	LaTeXEngine   string
}

// OutputFormat describes the format of a generated file.
type OutputFormat string

const (
	OutputFormatPDF      OutputFormat = "pdf"
	OutputFormatHTML     OutputFormat = "html"
	OutputFormatLaTeX    OutputFormat = "latex"
	OutputFormatMarkdown OutputFormat = "markdown"
	OutputFormatDOCX     OutputFormat = "docx"
)

// GenerationResult describes a single generated output file.
type GenerationResult struct {
	Template     string
	TemplateType generators.TemplateType
	OutputFormat OutputFormat
	OutputPath   string
	PageCount    int // populated for PDFs, 0 otherwise
}

// Generate runs the full resume generation pipeline and returns the list of
// generated files. Every template produces both its native format and a PDF.
func Generate(opts GenerateOptions) ([]GenerationResult, error) {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	// Load and validate resume
	data, err := LoadResume(opts.InputFile)
	if err != nil {
		return nil, err
	}

	td := generators.NewTemplateData(data.Resume, data.SectionOrder)

	// Apply config defaults
	applyConfigDefaults(&opts)

	// Resolve templates
	normalizedNames := sanitizeTemplateNames(opts.TemplateNames)
	selectedTemplates, err := LoadSelectedTemplates(normalizedNames)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve templates: %w", err)
	}
	if len(selectedTemplates) == 0 {
		return nil, fmt.Errorf("no templates available for generation")
	}

	// Resolve output directory
	runDir, desiredBase, err := prepareOutputDir(opts.OutputDir, opts.InputFile, data.Resume.Contact.Name)
	if err != nil {
		return nil, err
	}

	generator := generators.NewGenerator(sugar)

	// Pre-load the HTML fallback template for DOCX->PDF conversion
	htmlFallbackTmpl, _ := generators.LoadTemplate("modern-html")

	var results []GenerationResult

	for _, tmpl := range selectedTemplates {
		templateResults, err := generateForTemplate(sugar, generator, tmpl, td, runDir, desiredBase, opts.LaTeXEngine, htmlFallbackTmpl)
		if err != nil {
			return nil, err
		}
		results = append(results, templateResults...)
	}

	return results, nil
}

func applyConfigDefaults(opts *GenerateOptions) {
	cfg, err := templates.LoadConfig()
	if err != nil {
		return
	}
	if len(opts.TemplateNames) == 0 && len(cfg.Defaults.Templates) > 0 {
		opts.TemplateNames = cfg.Defaults.Templates
	}
	if opts.LaTeXEngine == "" && cfg.Defaults.LaTeXEngine != "" {
		opts.LaTeXEngine = cfg.Defaults.LaTeXEngine
	}
	if opts.OutputDir == "" && cfg.Defaults.OutputDir != "" {
		opts.OutputDir = cfg.Defaults.OutputDir
	}
}

func prepareOutputDir(outputDir, inputFile, contactName string) (runDir, desiredBase string, err error) {
	rootDirInput := strings.TrimSpace(outputDir)
	if rootDirInput == "" {
		rootDirInput = utils.DefaultOutputDir()
	}

	resolvedDir, err := utils.ResolvePath(rootDirInput)
	if err != nil {
		return "", "", fmt.Errorf("error resolving output directory: %w", err)
	}
	if resolvedDir == "" {
		resolvedDir, err = os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("failed to determine working directory: %w", err)
		}
	}
	if err := utils.EnsureDir(resolvedDir); err != nil {
		return "", "", fmt.Errorf("error creating output directory: %w", err)
	}

	resumeSlug := GenerateFilenameSlug(inputFile)
	runDir = GenerateRunDir(filepath.Join(resolvedDir, resumeSlug), time.Now())
	if err := utils.EnsureDir(runDir); err != nil {
		return "", "", fmt.Errorf("error creating run output directory: %w", err)
	}

	desiredBase = GenerateOutputBaseName(contactName)
	return runDir, desiredBase, nil
}

func generateForTemplate(
	sugar *zap.SugaredLogger,
	generator *generators.Generator,
	tmpl *generators.Template,
	td *generators.TemplateData,
	runDir, desiredBase, latexEngine string,
	htmlFallbackTmpl *generators.Template,
) ([]GenerationResult, error) {
	switch tmpl.Type {
	case generators.TemplateTypeMarkdown:
		return generateMarkdown(sugar, generator, tmpl, td, runDir, desiredBase)
	case generators.TemplateTypeDOCX:
		return generateDOCX(sugar, generator, tmpl, td, runDir, desiredBase, htmlFallbackTmpl)
	case generators.TemplateTypeHTML:
		return generateHTML(sugar, generator, tmpl, td, runDir, desiredBase)
	case generators.TemplateTypeLaTeX:
		return generateLaTeX(sugar, generator, tmpl, td, runDir, desiredBase, latexEngine)
	default:
		return nil, fmt.Errorf("unknown template type: %s", tmpl.Type)
	}
}

func generateMarkdown(
	sugar *zap.SugaredLogger,
	generator *generators.Generator,
	tmpl *generators.Template,
	td *generators.TemplateData,
	runDir, desiredBase string,
) ([]GenerationResult, error) {
	content, err := generator.GenerateWithTemplate(tmpl, td)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Markdown with template %s: %w", tmpl.Name, err)
	}

	var results []GenerationResult

	// Write native .md
	mdOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".md")
	if err != nil {
		return nil, fmt.Errorf("error determining output filename for template %s: %w", tmpl.Name, err)
	}
	if err := os.WriteFile(mdOutputPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write Markdown file: %w", err)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatMarkdown,
		OutputPath:   mdOutputPath,
	})

	// Convert Markdown -> HTML -> PDF
	htmlContent, err := compilers.MarkdownToHTML(content)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Markdown to HTML for template %s: %w", tmpl.Name, err)
	}

	pdfOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".pdf")
	if err != nil {
		return nil, fmt.Errorf("error determining PDF output filename for template %s: %w", tmpl.Name, err)
	}

	debugDir, err := os.MkdirTemp("", "resume-debug-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp debug directory: %w", err)
	}

	if compileErr := CompileHTMLToPDF(sugar, htmlContent, pdfOutputPath, debugDir); compileErr != nil {
		persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
		_ = os.Rename(debugDir, persistedDebug)
		return nil, fmt.Errorf("failed to compile Markdown template %s to PDF: %w (debug artifacts: %s)", tmpl.Name, compileErr, persistedDebug)
	}
	_ = os.RemoveAll(debugDir)

	pageCount := 0
	if pdfData, readErr := os.ReadFile(pdfOutputPath); readErr == nil {
		pageCount = compilers.CountPDFPages(pdfData)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatPDF,
		OutputPath:   pdfOutputPath,
		PageCount:    pageCount,
	})

	return results, nil
}

func generateHTML(
	sugar *zap.SugaredLogger,
	generator *generators.Generator,
	tmpl *generators.Template,
	td *generators.TemplateData,
	runDir, desiredBase string,
) ([]GenerationResult, error) {
	content, err := generator.GenerateWithTemplate(tmpl, td)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML with template %s: %w", tmpl.Name, err)
	}

	var results []GenerationResult

	// Write native .html
	htmlOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".html")
	if err != nil {
		return nil, fmt.Errorf("error determining output filename for template %s: %w", tmpl.Name, err)
	}
	if err := os.WriteFile(htmlOutputPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write HTML file: %w", err)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatHTML,
		OutputPath:   htmlOutputPath,
	})

	// Compile HTML -> PDF
	pdfOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".pdf")
	if err != nil {
		return nil, fmt.Errorf("error determining PDF output filename for template %s: %w", tmpl.Name, err)
	}

	debugDir, err := os.MkdirTemp("", "resume-debug-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp debug directory: %w", err)
	}

	if compileErr := CompileHTMLToPDF(sugar, content, pdfOutputPath, debugDir); compileErr != nil {
		persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
		_ = os.Rename(debugDir, persistedDebug)
		return nil, fmt.Errorf("failed to compile template %s to PDF: %w (debug artifacts: %s)", tmpl.Name, compileErr, persistedDebug)
	}
	_ = os.RemoveAll(debugDir)

	pageCount := 0
	if pdfData, readErr := os.ReadFile(pdfOutputPath); readErr == nil {
		pageCount = compilers.CountPDFPages(pdfData)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatPDF,
		OutputPath:   pdfOutputPath,
		PageCount:    pageCount,
	})

	return results, nil
}

func generateLaTeX(
	sugar *zap.SugaredLogger,
	generator *generators.Generator,
	tmpl *generators.Template,
	td *generators.TemplateData,
	runDir, desiredBase, latexEngine string,
) ([]GenerationResult, error) {
	content, err := generator.GenerateWithTemplate(tmpl, td)
	if err != nil {
		return nil, fmt.Errorf("failed to generate LaTeX with template %s: %w", tmpl.Name, err)
	}

	var results []GenerationResult

	// Write native .tex
	texOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".tex")
	if err != nil {
		return nil, fmt.Errorf("error determining output filename for template %s: %w", tmpl.Name, err)
	}
	if err := os.WriteFile(texOutputPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write LaTeX file: %w", err)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatLaTeX,
		OutputPath:   texOutputPath,
	})

	// Compile LaTeX -> PDF
	pdfOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".pdf")
	if err != nil {
		return nil, fmt.Errorf("error determining PDF output filename for template %s: %w", tmpl.Name, err)
	}

	debugDir, err := os.MkdirTemp("", "resume-debug-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp debug directory: %w", err)
	}

	templateDir := filepath.Dir(tmpl.Path)
	if compileErr := CompileLaTeXToPDF(sugar, content, pdfOutputPath, debugDir, templateDir, latexEngine); compileErr != nil {
		persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
		_ = os.Rename(debugDir, persistedDebug)
		return nil, fmt.Errorf("failed to compile template %s to PDF: %w (debug artifacts: %s)", tmpl.Name, compileErr, persistedDebug)
	}
	_ = os.RemoveAll(debugDir)

	pageCount := 0
	if pdfData, readErr := os.ReadFile(pdfOutputPath); readErr == nil {
		pageCount = compilers.CountPDFPages(pdfData)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatPDF,
		OutputPath:   pdfOutputPath,
		PageCount:    pageCount,
	})

	return results, nil
}

func generateDOCX(
	sugar *zap.SugaredLogger,
	generator *generators.Generator,
	tmpl *generators.Template,
	td *generators.TemplateData,
	runDir, desiredBase string,
	htmlFallbackTmpl *generators.Template,
) ([]GenerationResult, error) {
	docxBytes, err := generator.GenerateDOCXWithTemplate(tmpl, td)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DOCX with template %s: %w", tmpl.Name, err)
	}

	var results []GenerationResult

	// Write native .docx
	docxOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".docx")
	if err != nil {
		return nil, fmt.Errorf("error determining output filename for template %s: %w", tmpl.Name, err)
	}
	if err := os.WriteFile(docxOutputPath, docxBytes, 0644); err != nil {
		return nil, fmt.Errorf("failed to write DOCX file: %w", err)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatDOCX,
		OutputPath:   docxOutputPath,
	})

	// Generate PDF via HTML fallback
	if htmlFallbackTmpl == nil {
		return results, nil
	}

	htmlContent, htmlErr := generator.GenerateWithTemplate(htmlFallbackTmpl, td)
	if htmlErr != nil {
		return results, nil
	}

	pdfOutputPath, err := EnsureUniqueOutputPath(runDir, desiredBase, tmpl.Name, ".pdf")
	if err != nil {
		return results, nil
	}

	debugDir, debugErr := os.MkdirTemp("", "resume-debug-*")
	if debugErr != nil {
		return results, nil
	}

	if pdfErr := CompileHTMLToPDF(sugar, htmlContent, pdfOutputPath, debugDir); pdfErr != nil {
		persistedDebug := filepath.Join(runDir, desiredBase+"."+tmpl.Name+"_debug")
		_ = os.Rename(debugDir, persistedDebug)
		return results, nil
	}
	_ = os.RemoveAll(debugDir)

	pageCount := 0
	if pdfData, readErr := os.ReadFile(pdfOutputPath); readErr == nil {
		pageCount = compilers.CountPDFPages(pdfData)
	}
	results = append(results, GenerationResult{
		Template:     tmpl.Name,
		TemplateType: tmpl.Type,
		OutputFormat: OutputFormatPDF,
		OutputPath:   pdfOutputPath,
		PageCount:    pageCount,
	})

	return results, nil
}

// CompileHTMLToPDF compiles HTML content to PDF using a Chromium-based browser.
func CompileHTMLToPDF(logger *zap.SugaredLogger, htmlContent, outputPath, debugDir string) error {
	baseName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	if baseName == "" {
		baseName = "resume"
	}

	debugHTMLPath := filepath.Join(debugDir, baseName+".html")
	_ = os.WriteFile(debugHTMLPath, []byte(htmlContent), 0644)

	compiler := compilers.NewRodHTMLToPDFCompiler(logger)
	return compiler.Compile(htmlContent, outputPath)
}

// CompileLaTeXToPDF compiles LaTeX content to PDF using available LaTeX engines.
func CompileLaTeXToPDF(logger *zap.SugaredLogger, latexContent, outputPath, debugDir, templateDir, latexEngine string) error {
	baseName := strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
	if baseName == "" {
		baseName = "resume"
	}

	resolvedTemplateDir := filepath.Clean(templateDir)
	if resolvedTemplateDir != "" && !utils.DirExists(resolvedTemplateDir) {
		resolvedTemplateDir = ""
	}

	var compiler compilers.Compiler
	if latexEngine != "" {
		compiler = compilers.NewLaTeXCompiler(latexEngine, logger)
	} else {
		autoCompiler, err := compilers.NewAutoLaTeXCompiler(logger)
		if err != nil {
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

	generatedPDF := filepath.Join(debugDir, baseName+".pdf")
	if !utils.FileExists(generatedPDF) {
		return fmt.Errorf("expected PDF was not generated at %s", generatedPDF)
	}

	if err := os.Rename(generatedPDF, outputPath); err != nil {
		return fmt.Errorf("failed to move PDF: %w", err)
	}

	return nil
}

// LoadSelectedTemplates loads the specified templates or all available templates if none specified.
func LoadSelectedTemplates(templateNames []string) ([]*generators.Template, error) {
	if len(templateNames) == 0 {
		allTemplates, err := generators.ListTemplates()
		if err != nil {
			return nil, fmt.Errorf("failed to list templates: %w", err)
		}

		var result []*generators.Template
		for i := range allTemplates {
			result = append(result, &allTemplates[i])
		}

		sort.Slice(result, func(i, j int) bool {
			return result[i].Name < result[j].Name
		})

		return result, nil
	}

	var tmpls []*generators.Template
	for _, name := range templateNames {
		tmpl, err := generators.LoadTemplate(name)
		if err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", name, err)
		}
		tmpls = append(tmpls, tmpl)
	}

	return tmpls, nil
}

// GenerateFilenameSlug creates a filesystem-safe slug from an input filename.
func GenerateFilenameSlug(inputPath string) string {
	base := filepath.Base(inputPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	slug := sanitizeNameComponent(name)
	if slug == "" {
		return "resume"
	}
	return slug
}

// GenerateOutputBaseName creates a proper-cased base name from a contact name.
func GenerateOutputBaseName(contactName string) string {
	parts := strings.Fields(contactName)
	if len(parts) == 0 {
		return "Resume"
	}
	var nameParts []string
	for _, p := range parts {
		nameParts = append(nameParts, toProperCase(p))
	}
	return strings.Join(nameParts, "_")
}

// GenerateRunDir creates a timestamped run directory path.
func GenerateRunDir(baseDir string, t time.Time) string {
	return filepath.Join(baseDir, t.Format("2006-01-02_15-04"))
}

// EnsureUniqueOutputPath returns a unique file path within the run directory.
func EnsureUniqueOutputPath(runDir, desiredBase, templateName, extension string) (string, error) {
	base := strings.TrimSpace(desiredBase)
	if base == "" {
		base = "Resume"
	}

	tmplSlug := strings.TrimSpace(templateName)
	if tmplSlug == "" {
		tmplSlug = "template"
	}

	ext := extension
	if ext == "" {
		ext = ".pdf"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	candidateBase := base + "." + tmplSlug
	candidate := filepath.Join(runDir, candidateBase+ext)
	if !utils.FileExists(candidate) {
		return candidate, nil
	}

	for attempt := 2; attempt <= 9999; attempt++ {
		suffix := fmt.Sprintf("_%d", attempt)
		candidate = filepath.Join(runDir, candidateBase+suffix+ext)
		if !utils.FileExists(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("failed to find unique output filename in %s", runDir)
}

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

	sort.Strings(result)
	return result
}

func toProperCase(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
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
