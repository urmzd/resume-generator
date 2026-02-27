package generators

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// LaTeXGenerator renders LaTeX templates using engine-specific formatting helpers.
type LaTeXGenerator struct {
	logger    *zap.SugaredLogger
	formatter *latexFormatter
}

// NewLaTeXGenerator creates a new LaTeX generator wired with the LaTeX formatter.
func NewLaTeXGenerator(logger *zap.SugaredLogger) *LaTeXGenerator {
	return &LaTeXGenerator{
		logger:    logger,
		formatter: newLaTeXFormatter(),
	}
}

// Generate renders a LaTeX template with resume data using the formatter's helper functions.
func (g *LaTeXGenerator) Generate(templateContent string, r *resume.Resume) (string, error) {
	g.logger.Info("Rendering LaTeX template")

	funcs := g.formatter.TemplateFuncs()

	tmpl, err := template.New("latex").Funcs(funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse LaTeX template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, r); err != nil {
		return "", fmt.Errorf("failed to execute LaTeX template: %w", err)
	}

	g.logger.Info("Successfully rendered LaTeX template")
	return output.String(), nil
}
