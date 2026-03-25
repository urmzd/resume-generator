package compilers

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// LaTeXEngine implements Engine for LaTeX-to-PDF compilation.
type LaTeXEngine struct {
	command     string
	templateDir string // directory containing .cls/.sty files to copy alongside the .tex
	logger      *zap.SugaredLogger
}

func (e *LaTeXEngine) Name() string { return e.command }

// Compile writes the LaTeX content to a temp directory, runs the LaTeX
// compiler, and moves the resulting PDF to outputPath.
func (e *LaTeXEngine) Compile(content string, outputPath string) error {
	baseName := filenameWithoutExt(outputPath)

	debugDir, err := os.MkdirTemp("", "incipit-latex-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(debugDir) }()

	// Use the old LaTeXCompiler internally for the actual compilation step.
	compiler := NewLaTeXCompiler(e.command, e.logger)
	if e.templateDir != "" {
		compiler.LoadClasses(e.templateDir)
	}
	compiler.AddOutputFolder(debugDir)
	compiler.Compile(content, baseName)

	generatedPDF := filepath.Join(debugDir, baseName+".pdf")
	if _, err := os.Stat(generatedPDF); os.IsNotExist(err) {
		return fmt.Errorf("expected PDF was not generated at %s", generatedPDF)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.Rename(generatedPDF, outputPath); err != nil {
		return fmt.Errorf("failed to move PDF: %w", err)
	}

	return nil
}

func filenameWithoutExt(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	if name == "" {
		return "resume"
	}
	return name
}
