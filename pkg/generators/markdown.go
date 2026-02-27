package generators

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// MarkdownGenerator renders Markdown templates using engine-specific formatting helpers.
type MarkdownGenerator struct {
	logger    *zap.SugaredLogger
	formatter *markdownFormatter
}

// NewMarkdownGenerator creates a new Markdown generator wired with the Markdown formatter.
func NewMarkdownGenerator(logger *zap.SugaredLogger) *MarkdownGenerator {
	return &MarkdownGenerator{
		logger:    logger,
		formatter: newMarkdownFormatter(),
	}
}

// Generate renders a Markdown template with resume data using the formatter's helper functions.
func (g *MarkdownGenerator) Generate(templateContent string, r *resume.Resume) (string, error) {
	g.logger.Info("Rendering Markdown template")

	tmpl, err := template.New("markdown").Funcs(g.formatter.TemplateFuncs()).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse Markdown template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, r); err != nil {
		return "", fmt.Errorf("failed to execute Markdown template: %w", err)
	}

	g.logger.Info("Successfully rendered Markdown template")
	return output.String(), nil
}
