package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/pipeline"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

// ParseResult is returned by OpenFile with parsed resume info.
type ParseResult struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Format string `json:"format"`
}

// TemplateInfo describes an available template for the frontend.
type TemplateInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

// GeneratePDFResult is returned by GeneratePDF with the base64-encoded PDF
// and metadata such as the page count.
type GeneratePDFResult struct {
	Data      string `json:"data"`
	PageCount int    `json:"pageCount"`
}

// App is the Wails application struct. Its exported methods are
// automatically bound as frontend-callable functions.
type App struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	generator  *generators.Generator
	pipeline   *pipeline.PDFPipeline
	resume     *resume.Resume
	resumePath string
	resumeFmt  string
}

// NewApp creates a new App instance.
func NewApp() *App {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	generators.SetEmbeddedFS(EmbeddedTemplates)

	gen := generators.NewGenerator(sugar)

	return &App{
		logger:    sugar,
		generator: gen,
		pipeline:  pipeline.NewPDFPipeline(sugar, gen),
	}
}

// startup is the Wails lifecycle hook called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenFile shows a native file dialog, reads and parses the resume file.
func (a *App) OpenFile() (*ParseResult, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Resume File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Resume Files", Pattern: "*.yml;*.yaml;*.json;*.toml;*.md;*.markdown"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("file dialog error: %w", err)
	}
	if path == "" {
		return nil, fmt.Errorf("no file selected")
	}

	return a.loadResume(path)
}

// GetResume returns the full resume struct as JSON.
func (a *App) GetResume() (*resume.Resume, error) {
	if a.resume == nil {
		return nil, fmt.Errorf("no resume loaded")
	}
	return a.resume, nil
}

// UpdateResume validates and replaces the in-memory resume.
// Returns validation errors on failure.
func (a *App) UpdateResume(updated resume.Resume) ([]resume.ValidationError, error) {
	errors := resume.Validate(&updated)
	if len(errors) > 0 {
		return errors, nil
	}
	a.resume = &updated
	return nil, nil
}

// SaveResumeFile serializes the in-memory resume back to the original file.
func (a *App) SaveResumeFile() error {
	if a.resume == nil {
		return fmt.Errorf("no resume loaded")
	}
	if a.resumePath == "" {
		return fmt.Errorf("no file path stored")
	}

	data, canonicalFmt, err := resume.SerializeResume(a.resume, a.resumeFmt)
	if err != nil {
		return fmt.Errorf("failed to serialize resume: %w", err)
	}

	// Markdown falls back to YAML — update path and format
	if (a.resumeFmt == "md" || a.resumeFmt == "markdown") && canonicalFmt == "yaml" {
		a.resumePath = strings.TrimSuffix(a.resumePath, filepath.Ext(a.resumePath)) + ".yml"
		a.resumeFmt = "yaml"
	}

	return os.WriteFile(a.resumePath, data, 0644)
}

// GetTemplates returns the list of available templates.
func (a *App) GetTemplates() ([]TemplateInfo, error) {
	templates, err := generators.ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}

	var result []TemplateInfo
	for _, t := range templates {
		result = append(result, TemplateInfo{
			Name:        t.Name,
			DisplayName: t.DisplayName,
			Format:      string(t.Type),
			Description: t.Description,
		})
	}
	return result, nil
}

// GeneratePDF generates a PDF for the given template and returns base64-encoded bytes
// along with metadata such as page count.
func (a *App) GeneratePDF(templateName string) (*GeneratePDFResult, error) {
	if a.resume == nil {
		return nil, fmt.Errorf("no resume loaded")
	}

	tmpl, err := generators.LoadTemplate(templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	pdfBytes, err := a.pipeline.CompileToPDFBytes(tmpl, a.resume)
	if err != nil {
		return nil, err
	}

	pageCount := compilers.CountPDFPages(pdfBytes)
	if pageCount > 1 {
		a.logger.Warnf("Resume exceeds 1 page (%d pages) with template %s", pageCount, templateName)
	}

	return &GeneratePDFResult{
		Data:      base64.StdEncoding.EncodeToString(pdfBytes),
		PageCount: pageCount,
	}, nil
}

// SavePDF generates a PDF and opens a native Save dialog.
func (a *App) SavePDF(templateName string) error {
	result, err := a.GeneratePDF(templateName)
	if err != nil {
		return err
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	defaultDir := utils.DefaultOutputDir()
	_ = utils.EnsureDir(defaultDir)

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: defaultDir,
		Title:            "Save PDF",
		DefaultFilename:  "resume.pdf",
		Filters: []runtime.FileFilter{
			{DisplayName: "PDF Files", Pattern: "*.pdf"},
		},
	})
	if err != nil {
		return fmt.Errorf("save dialog error: %w", err)
	}
	if path == "" {
		return nil // User cancelled
	}

	return os.WriteFile(path, pdfBytes, 0644)
}

// SaveNative saves the native format (.docx/.html/.tex) via native Save dialog.
func (a *App) SaveNative(templateName string) error {
	if a.resume == nil {
		return fmt.Errorf("no resume loaded")
	}

	tmpl, err := generators.LoadTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	var data []byte
	var ext string
	var filterName string

	switch tmpl.Type {
	case generators.TemplateTypeHTML:
		content, err := a.generator.GenerateWithTemplate(tmpl, a.resume)
		if err != nil {
			return fmt.Errorf("failed to generate HTML: %w", err)
		}
		data = []byte(content)
		ext = ".html"
		filterName = "HTML Files"

	case generators.TemplateTypeDOCX:
		docxBytes, err := a.generator.GenerateDOCX(a.resume)
		if err != nil {
			return fmt.Errorf("failed to generate DOCX: %w", err)
		}
		data = docxBytes
		ext = ".docx"
		filterName = "Word Documents"

	case generators.TemplateTypeLaTeX:
		content, err := a.generator.GenerateWithTemplate(tmpl, a.resume)
		if err != nil {
			return fmt.Errorf("failed to generate LaTeX: %w", err)
		}
		data = []byte(content)
		ext = ".tex"
		filterName = "LaTeX Files"

	default:
		return fmt.Errorf("unsupported template type: %s", tmpl.Type)
	}

	defaultDir := utils.DefaultOutputDir()
	_ = utils.EnsureDir(defaultDir)

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: defaultDir,
		Title:            "Save " + strings.ToUpper(strings.TrimPrefix(ext, ".")),
		DefaultFilename:  "resume" + ext,
		Filters: []runtime.FileFilter{
			{DisplayName: filterName, Pattern: "*" + ext},
		},
	})
	if err != nil {
		return fmt.Errorf("save dialog error: %w", err)
	}
	if path == "" {
		return nil // User cancelled
	}

	return os.WriteFile(path, data, 0644)
}

// LoadFileFromPath loads a resume from a given path (no native dialog).
// Used by e2e tests and demo automation.
func (a *App) LoadFileFromPath(path string) (*ParseResult, error) {
	return a.loadResume(path)
}

// SavePDFToPath generates a PDF and writes it to a given path (no native dialog).
// Used by e2e tests and demo automation.
func (a *App) SavePDFToPath(templateName, outputPath string) error {
	result, err := a.GeneratePDF(templateName)
	if err != nil {
		return err
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return os.WriteFile(outputPath, pdfBytes, 0644)
}

// loadResume is shared logic for OpenFile and LoadFileFromPath.
func (a *App) loadResume(path string) (*ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	format := strings.TrimPrefix(filepath.Ext(path), ".")
	inputData, err := resume.LoadResumeFromBytes(data, format)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resume: %w", err)
	}

	if err := inputData.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	a.resume = inputData.ToResume()
	a.resumePath = path
	a.resumeFmt = inputData.GetFormat()

	return &ParseResult{
		Name:   a.resume.Contact.Name,
		Email:  a.resume.Contact.Email,
		Format: inputData.GetFormat(),
	}, nil
}
