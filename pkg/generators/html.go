package generators

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// HTMLGenerator generates HTML resumes from templates
type HTMLGenerator struct {
	logger *zap.SugaredLogger
	funcs  template.FuncMap
	format Formatter
}

// htmlPayload wraps Resume with optional CSS for templates that need it
type htmlPayload struct {
	*resume.Resume
	CSS template.CSS
}

// NewHTMLGenerator creates a new HTML resume generator
func NewHTMLGenerator(logger *zap.SugaredLogger) *HTMLGenerator {
	formatter := newHTMLFormatter()
	return &HTMLGenerator{
		logger: logger,
		funcs:  formatter.TemplateFuncs(),
		format: formatter,
	}
}

// Generate creates an HTML resume from the resume data and template
func (g *HTMLGenerator) Generate(templateContent string, r *resume.Resume) (string, error) {
	g.logger.Info("Generating HTML resume")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template - passing resume directly
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume")
	return buf.String(), nil
}

// GenerateWithCSS creates an HTML resume with embedded CSS
func (g *HTMLGenerator) GenerateWithCSS(templateContent, cssContent string, r *resume.Resume) (string, error) {
	g.logger.Info("Generating HTML resume with embedded CSS")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template with payload containing CSS
	var buf bytes.Buffer
	payload := htmlPayload{
		Resume: r,
		CSS:    template.CSS(cssContent),
	}

	if err := tmpl.Execute(&buf, payload); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume with CSS")
	return buf.String(), nil
}

// GenerateStandalone creates a complete HTML document with all dependencies
func (g *HTMLGenerator) GenerateStandalone(templateContent, cssContent string, r *resume.Resume) (string, error) {
	htmlContent, err := g.GenerateWithCSS(templateContent, cssContent, r)
	if err != nil {
		return "", err
	}

	// Ensure it's a complete HTML document
	if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
		standaloneTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Resume.Contact.Name}} - Resume</title>
    <style>
        {{.CSS}}
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>`

		data := map[string]interface{}{
			"Resume":  r,
			"CSS":     template.CSS(cssContent),
			"Content": template.HTML(htmlContent),
		}

		tmpl, err := template.New("standalone").Funcs(g.funcs).Parse(standaloneTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to parse standalone template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute standalone template: %w", err)
		}

		htmlContent = buf.String()
	}

	return htmlContent, nil
}

// HTMLGeneratorOptions provides configuration for HTML generation
type HTMLGeneratorOptions struct {
	Theme            string
	IncludeCSS       bool
	Standalone       bool
	ResponsiveDesign bool
	PrintOptimized   bool
	FontAwesome      bool
	CustomFonts      []string
	ColorScheme      string
}

// AdvancedHTMLGenerator provides advanced HTML generation capabilities
type AdvancedHTMLGenerator struct {
	logger  *zap.SugaredLogger
	options HTMLGeneratorOptions
	funcs   template.FuncMap
}

// NewAdvancedHTMLGenerator creates a new advanced HTML generator
func NewAdvancedHTMLGenerator(logger *zap.SugaredLogger, options HTMLGeneratorOptions) *AdvancedHTMLGenerator {
	generator := &AdvancedHTMLGenerator{
		logger:  logger,
		options: options,
	}
	generator.setupAdvancedTemplateFunctions()
	return generator
}

// setupAdvancedTemplateFunctions initializes advanced template helper functions
func (g *AdvancedHTMLGenerator) setupAdvancedTemplateFunctions() {
	opts := g.options
	g.funcs = template.FuncMap{
		// Basic formatters available in standard generator
		// We re-implement specific ones here or could reuse the formatter if we had it
		"formatDate": func(t time.Time) string {
			return t.Format("January 2006")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("Jan 2006")
		},
		"formatDateRange": func(start time.Time, end *time.Time) string {
			startStr := start.Format("Jan 2006")
			if end == nil {
				return startStr + " - Present"
			}
			endStr := end.Format("Jan 2006")
			if startStr == endStr {
				return startStr
			}
			return startStr + " - " + endStr
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},

		// Advanced functions
		"getThemeClass": func() string {
			return opts.Theme
		},
		"shouldIncludeCSS": func() bool {
			return opts.IncludeCSS
		},
		"isStandalone": func() bool {
			return opts.Standalone
		},
		"isResponsive": func() bool {
			return opts.ResponsiveDesign
		},
		"isPrintOptimized": func() bool {
			return opts.PrintOptimized
		},
		"shouldIncludeFontAwesome": func() bool {
			return opts.FontAwesome
		},
		"getCustomFonts": func() []string {
			return opts.CustomFonts
		},
		"getColorScheme": func() string {
			return opts.ColorScheme
		},

		// Utility functions
		"generateID": func(prefix string) string {
			return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
		},
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"wordCount": func(s string) int {
			return len(strings.Fields(s))
		},
		"characterCount": func(s string) int {
			return len(s)
		},
	}
}

// Generate creates an advanced HTML resume
func (g *AdvancedHTMLGenerator) Generate(templateContent string, resume *resume.Resume) (string, error) {
	g.logger.Infof("Generating advanced HTML resume with theme: %s", g.options.Theme)

	// Create template data with options
	data := map[string]interface{}{
		"Resume":  resume,
		"Options": g.options,
	}

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse advanced HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute advanced HTML template: %w", err)
	}

	g.logger.Info("Successfully generated advanced HTML resume")
	return buf.String(), nil
}
