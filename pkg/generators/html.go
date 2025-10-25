package generators

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

// HTMLGenerator generates HTML resumes from templates
type HTMLGenerator struct {
	logger *zap.SugaredLogger
	funcs  template.FuncMap
}

// NewHTMLGenerator creates a new HTML resume generator
func NewHTMLGenerator(logger *zap.SugaredLogger) *HTMLGenerator {
	generator := &HTMLGenerator{
		logger: logger,
	}
	generator.setupTemplateFunctions()
	return generator
}

// setupTemplateFunctions initializes template helper functions
func (g *HTMLGenerator) setupTemplateFunctions() {
	g.funcs = template.FuncMap{
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
		"calculateDuration": func(start time.Time, end *time.Time) string {
			var endTime time.Time
			if end == nil {
				endTime = time.Now()
			} else {
				endTime = *end
			}

			diff := endTime.Sub(start)
			years := int(diff.Hours() / 24 / 365)
			months := int((diff.Hours() / 24 / 30)) % 12

			if years > 0 && months > 0 {
				return fmt.Sprintf("%d yr %d mo", years, months)
			} else if years > 0 {
				return fmt.Sprintf("%d yr", years)
			} else if months > 0 {
				return fmt.Sprintf("%d mo", months)
			}
			return "< 1 mo"
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"escape": func(s string) string {
			return template.HTMLEscapeString(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"contains": strings.Contains,
		"sortSkillsByOrder": func(skills []definition.EnhancedSkillCategory) []definition.EnhancedSkillCategory {
			sorted := make([]definition.EnhancedSkillCategory, len(skills))
			copy(sorted, skills)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortExperienceByOrder": func(experiences []definition.EnhancedExperience) []definition.EnhancedExperience {
			sorted := make([]definition.EnhancedExperience, len(experiences))
			copy(sorted, experiences)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortProjectsByOrder": func(projects []definition.EnhancedProject) []definition.EnhancedProject {
			sorted := make([]definition.EnhancedProject, len(projects))
			copy(sorted, projects)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortEducationByOrder": func(education []definition.EnhancedEducation) []definition.EnhancedEducation {
			sorted := make([]definition.EnhancedEducation, len(education))
			copy(sorted, education)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortLinksByOrder": func(links []definition.EnhancedLink) []definition.EnhancedLink {
			sorted := make([]definition.EnhancedLink, len(links))
			copy(sorted, links)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortCertificationsByOrder": func(certs []definition.EnhancedCertification) []definition.EnhancedCertification {
			sorted := make([]definition.EnhancedCertification, len(certs))
			copy(sorted, certs)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"getIconClass": func(linkType string) string {
			icons := map[string]string{
				"github":    "fab fa-github",
				"linkedin":  "fab fa-linkedin",
				"twitter":   "fab fa-twitter",
				"website":   "fas fa-globe",
				"portfolio": "fas fa-briefcase",
				"email":     "fas fa-envelope",
				"phone":     "fas fa-phone",
			}
			if icon, exists := icons[strings.ToLower(linkType)]; exists {
				return icon
			}
			return "fas fa-link"
		},
		"formatGPA": func(gpa, maxGPA string) string {
			if gpa == "" {
				return ""
			}
			if maxGPA != "" && maxGPA != "4.0" {
				return fmt.Sprintf("%s/%s", gpa, maxGPA)
			}
			return gpa
		},
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"isEven": func(n int) bool {
			return n%2 == 0
		},
		"isOdd": func(n int) bool {
			return n%2 != 0
		},
	}
}

// Generate creates an HTML resume from the enhanced resume data and template
func (g *HTMLGenerator) Generate(templateContent string, resume *definition.EnhancedResume) (string, error) {
	g.logger.Info("Generating HTML resume")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resume); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume")
	return buf.String(), nil
}

// GenerateWithCSS creates an HTML resume with embedded CSS
func (g *HTMLGenerator) GenerateWithCSS(templateContent, cssContent string, resume *definition.EnhancedResume) (string, error) {
	g.logger.Info("Generating HTML resume with embedded CSS")

	// Create combined template data
	data := map[string]interface{}{
		"Resume": resume,
		"CSS":    template.CSS(cssContent),
	}

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume with CSS")
	return buf.String(), nil
}

// GenerateStandalone creates a complete HTML document with all dependencies
func (g *HTMLGenerator) GenerateStandalone(templateContent, cssContent string, resume *definition.EnhancedResume) (string, error) {
	htmlContent, err := g.GenerateWithCSS(templateContent, cssContent, resume)
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
			"Resume":  resume,
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

// LegacyHTMLGenerator provides compatibility with the old Resume format
type LegacyHTMLGenerator struct {
	logger *zap.SugaredLogger
	funcs  template.FuncMap
}

// NewLegacyHTMLGenerator creates a new legacy HTML resume generator
func NewLegacyHTMLGenerator(logger *zap.SugaredLogger) *LegacyHTMLGenerator {
	generator := &LegacyHTMLGenerator{
		logger: logger,
	}
	generator.setupTemplateFunctions()
	return generator
}

// setupTemplateFunctions initializes template helper functions for legacy format
func (g *LegacyHTMLGenerator) setupTemplateFunctions() {
	g.funcs = template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2006")
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
		"escape": func(s string) string {
			return template.HTMLEscapeString(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
	}
}

// Generate creates an HTML resume from the legacy resume data and template
func (g *LegacyHTMLGenerator) Generate(templateContent string, resume *definition.Resume) (string, error) {
	g.logger.Info("Generating HTML resume from legacy format")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resume); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume from legacy format")
	return buf.String(), nil
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
	g.funcs = template.FuncMap{
		// Include all basic functions
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
			return g.options.Theme
		},
		"shouldIncludeCSS": func() bool {
			return g.options.IncludeCSS
		},
		"isStandalone": func() bool {
			return g.options.Standalone
		},
		"isResponsive": func() bool {
			return g.options.ResponsiveDesign
		},
		"isPrintOptimized": func() bool {
			return g.options.PrintOptimized
		},
		"shouldIncludeFontAwesome": func() bool {
			return g.options.FontAwesome
		},
		"getCustomFonts": func() []string {
			return g.options.CustomFonts
		},
		"getColorScheme": func() string {
			return g.options.ColorScheme
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
func (g *AdvancedHTMLGenerator) Generate(templateContent string, resume *definition.EnhancedResume) (string, error) {
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