package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// PDFPipeline unifies PDF generation from HTML and LaTeX templates.
type PDFPipeline struct {
	logger    *zap.SugaredLogger
	generator *generators.Generator
	htmlToPDF *compilers.RodHTMLToPDFCompiler
	hasLatex  bool
}

// NewPDFPipeline creates a pipeline that can compile templates to PDF bytes.
func NewPDFPipeline(logger *zap.SugaredLogger, generator *generators.Generator) *PDFPipeline {
	return &PDFPipeline{
		logger:    logger,
		generator: generator,
		htmlToPDF: compilers.NewRodHTMLToPDFCompiler(logger),
		hasLatex:  compilers.DetectLaTeXEngine() != "",
	}
}

// HasLaTeX reports whether a LaTeX engine is available.
func (p *PDFPipeline) HasLaTeX() bool {
	return p.hasLatex
}

// CompileToPDFBytes generates a PDF from a template and resume.
// For DOCX templates, it falls back to the HTML template.
func (p *PDFPipeline) CompileToPDFBytes(tmpl *generators.Template, r *resume.Resume) ([]byte, error) {
	switch tmpl.Type {
	case generators.TemplateTypeHTML:
		return p.compileHTMLTemplateToPDF(tmpl, r)

	case generators.TemplateTypeDOCX:
		return p.compileHTMLFallbackToPDF(r)

	case generators.TemplateTypeLaTeX:
		if p.hasLatex {
			return p.CompileLaTeXToPDFBytes(tmpl, r)
		}
		return p.compileHTMLFallbackToPDF(r)

	default:
		return nil, fmt.Errorf("unsupported template type for PDF: %s", tmpl.Type)
	}
}

// CompileHTMLToPDFBytes compiles raw HTML content to PDF.
func (p *PDFPipeline) CompileHTMLToPDFBytes(html string) ([]byte, error) {
	return p.htmlToPDF.CompileToBytes(html)
}

// CompileLaTeXToPDFBytes compiles a LaTeX template to PDF.
func (p *PDFPipeline) CompileLaTeXToPDFBytes(tmpl *generators.Template, r *resume.Resume) ([]byte, error) {
	content, err := p.generator.GenerateWithTemplate(tmpl, r)
	if err != nil {
		return nil, fmt.Errorf("failed to generate LaTeX: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "resume-latex-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Copy template support files into the temp directory
	if err := p.copyTemplateFiles(tmpl, tmpDir); err != nil {
		return nil, err
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

func (p *PDFPipeline) compileHTMLTemplateToPDF(tmpl *generators.Template, r *resume.Resume) ([]byte, error) {
	html, err := p.generator.GenerateWithTemplate(tmpl, r)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML: %w", err)
	}
	pdfBytes, err := p.htmlToPDF.CompileToBytes(html)
	if err != nil {
		return nil, fmt.Errorf("failed to compile PDF: %w", err)
	}
	return pdfBytes, nil
}

func (p *PDFPipeline) compileHTMLFallbackToPDF(r *resume.Resume) ([]byte, error) {
	htmlTmpl, err := generators.LoadTemplate("modern-html")
	if err != nil {
		return nil, fmt.Errorf("no HTML fallback template available: %w", err)
	}
	return p.compileHTMLTemplateToPDF(htmlTmpl, r)
}

func (p *PDFPipeline) copyTemplateFiles(tmpl *generators.Template, destDir string) error {
	var sourceDir string

	if tmpl.Embedded && tmpl.EmbeddedDir != "" {
		extractedDir, err := generators.ExtractEmbeddedTemplateDir(tmpl.EmbeddedDir)
		if err != nil {
			return fmt.Errorf("failed to extract template files: %w", err)
		}
		defer func() { _ = os.RemoveAll(extractedDir) }()
		sourceDir = extractedDir
	} else if tmpl.Path != "" {
		sourceDir = filepath.Dir(tmpl.Path)
	} else {
		return nil
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil // non-fatal: template may not have support files
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		src := filepath.Join(sourceDir, e.Name())
		dst := filepath.Join(destDir, e.Name())
		data, readErr := os.ReadFile(src)
		if readErr != nil {
			continue
		}
		_ = os.WriteFile(dst, data, 0644)
	}

	return nil
}
