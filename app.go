package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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

// App is the Wails application struct. Its exported methods are
// automatically bound as frontend-callable functions.
type App struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	generator  *generators.Generator
	compiler   *compilers.RodHTMLToPDFCompiler
	resume     *resume.Resume
	resumePath string
	resumeFmt  string
	hasLatex   bool
}

// NewApp creates a new App instance.
func NewApp() *App {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	generators.SetEmbeddedFS(EmbeddedTemplates)

	return &App{
		logger:    sugar,
		generator: generators.NewGenerator(sugar),
		compiler:  compilers.NewRodHTMLToPDFCompiler(sugar),
		hasLatex:  compilers.DetectLaTeXEngine() != "",
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
			{DisplayName: "Resume Files", Pattern: "*.yml;*.yaml;*.json;*.toml"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("file dialog error: %w", err)
	}
	if path == "" {
		return nil, fmt.Errorf("no file selected")
	}

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

	var data []byte
	var err error

	switch a.resumeFmt {
	case "yaml", "yml":
		data, err = yaml.Marshal(a.resume)
	case "json":
		data, err = json.MarshalIndent(a.resume, "", "  ")
	case "toml":
		var buf bytes.Buffer
		err = toml.NewEncoder(&buf).Encode(a.resume)
		data = buf.Bytes()
	default:
		return fmt.Errorf("unsupported format: %s", a.resumeFmt)
	}

	if err != nil {
		return fmt.Errorf("failed to serialize resume: %w", err)
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

// GeneratePDF generates a PDF for the given template and returns base64-encoded bytes.
func (a *App) GeneratePDF(templateName string) (string, error) {
	if a.resume == nil {
		return "", fmt.Errorf("no resume loaded")
	}

	tmpl, err := generators.LoadTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	var pdfBytes []byte

	switch tmpl.Type {
	case generators.TemplateTypeHTML:
		html, err := a.generator.GenerateWithTemplate(tmpl, a.resume)
		if err != nil {
			return "", fmt.Errorf("failed to generate HTML: %w", err)
		}
		pdfBytes, err = a.compiler.CompileToBytes(html)
		if err != nil {
			return "", fmt.Errorf("failed to compile PDF: %w", err)
		}

	case generators.TemplateTypeDOCX:
		// Use HTML fallback for PDF
		htmlTmpl, err := generators.LoadTemplate("modern-html")
		if err != nil {
			return "", fmt.Errorf("no HTML fallback template available: %w", err)
		}
		html, err := a.generator.GenerateWithTemplate(htmlTmpl, a.resume)
		if err != nil {
			return "", fmt.Errorf("failed to generate HTML fallback: %w", err)
		}
		pdfBytes, err = a.compiler.CompileToBytes(html)
		if err != nil {
			return "", fmt.Errorf("failed to compile PDF: %w", err)
		}

	case generators.TemplateTypeLaTeX:
		if a.hasLatex {
			pdfBytes, err = a.compileLaTeXToPDFBytes(tmpl)
			if err != nil {
				return "", err
			}
		} else {
			// Fallback to HTML template
			htmlTmpl, err := generators.LoadTemplate("modern-html")
			if err != nil {
				return "", fmt.Errorf("no HTML fallback for LaTeX: %w", err)
			}
			html, err := a.generator.GenerateWithTemplate(htmlTmpl, a.resume)
			if err != nil {
				return "", fmt.Errorf("failed to generate HTML fallback: %w", err)
			}
			pdfBytes, err = a.compiler.CompileToBytes(html)
			if err != nil {
				return "", fmt.Errorf("failed to compile PDF: %w", err)
			}
		}

	default:
		return "", fmt.Errorf("unsupported template type: %s", tmpl.Type)
	}

	return base64.StdEncoding.EncodeToString(pdfBytes), nil
}

// SavePDF generates a PDF and opens a native Save dialog.
func (a *App) SavePDF(templateName string) error {
	b64, err := a.GeneratePDF(templateName)
	if err != nil {
		return err
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(b64)
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

// SavePDFToPath generates a PDF and writes it to a given path (no native dialog).
// Used by e2e tests and demo automation.
func (a *App) SavePDFToPath(templateName, outputPath string) error {
	b64, err := a.GeneratePDF(templateName)
	if err != nil {
		return err
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return os.WriteFile(outputPath, pdfBytes, 0644)
}

// compileLaTeXToPDFBytes compiles a LaTeX template and returns PDF bytes.
func (a *App) compileLaTeXToPDFBytes(tmpl *generators.Template) ([]byte, error) {
	content, err := a.generator.GenerateWithTemplate(tmpl, a.resume)
	if err != nil {
		return nil, fmt.Errorf("failed to generate LaTeX: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "resume-latex-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Extract embedded template support files if needed
	if tmpl.Embedded && tmpl.EmbeddedDir != "" {
		extractedDir, err := generators.ExtractEmbeddedTemplateDir(tmpl.EmbeddedDir)
		if err != nil {
			return nil, fmt.Errorf("failed to extract template files: %w", err)
		}
		defer func() { _ = os.RemoveAll(extractedDir) }()

		entries, err := os.ReadDir(extractedDir)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					src := filepath.Join(extractedDir, e.Name())
					dst := filepath.Join(tmpDir, e.Name())
					data, readErr := os.ReadFile(src)
					if readErr != nil {
						continue
					}
					_ = os.WriteFile(dst, data, 0644)
				}
			}
		}
	} else if tmpl.Path != "" {
		// Filesystem template — copy support files from template directory
		templateDir := filepath.Dir(tmpl.Path)
		entries, _ := os.ReadDir(templateDir)
		for _, e := range entries {
			if !e.IsDir() {
				src := filepath.Join(templateDir, e.Name())
				dst := filepath.Join(tmpDir, e.Name())
				data, readErr := os.ReadFile(src)
				if readErr != nil {
					continue
				}
				_ = os.WriteFile(dst, data, 0644)
			}
		}
	}

	texPath := filepath.Join(tmpDir, "resume.tex")
	if err := os.WriteFile(texPath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write .tex file: %w", err)
	}

	engine := compilers.DetectLaTeXEngine()
	cmd := exec.Command(engine, "-interaction=nonstopmode", texPath)
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("LaTeX compilation failed: %w\n%s", err, string(out))
	}

	pdfPath := filepath.Join(tmpDir, "resume.pdf")
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read compiled PDF: %w", err)
	}

	return pdfBytes, nil
}
