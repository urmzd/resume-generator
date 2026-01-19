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
	formatter Formatter
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

	// Parse template with formatter functions
	tmpl, err := template.New("latex").Funcs(g.formatter.TemplateFuncs()).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse LaTeX template: %w", err)
	}

	// Execute template - passing resume directly
	// The template is responsible for using {{escape .Field}} where needed
	var output strings.Builder
	if err := tmpl.Execute(&output, r); err != nil {
		return "", fmt.Errorf("failed to execute LaTeX template: %w", err)
	}

	g.logger.Info("Successfully rendered LaTeX template")
	return output.String(), nil
}
