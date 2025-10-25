package generators

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

// LaTeXGenerator renders LaTeX templates
type LaTeXGenerator struct {
	logger *zap.SugaredLogger
}

// NewLaTeXGenerator creates a new LaTeX generator
func NewLaTeXGenerator(logger *zap.SugaredLogger) *LaTeXGenerator {
	return &LaTeXGenerator{logger: logger}
}

// Generate renders a LaTeX template with resume data
func (g *LaTeXGenerator) Generate(templateContent string, resume *definition.EnhancedResume) (string, error) {
	g.logger.Info("Rendering LaTeX template")

	// Create template with helper functions
	tmpl, err := template.New("latex").Funcs(g.templateFuncs()).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse LaTeX template: %w", err)
	}

	// Execute template
	var output strings.Builder
	if err := tmpl.Execute(&output, resume); err != nil {
		return "", fmt.Errorf("failed to execute LaTeX template: %w", err)
	}

	g.logger.Info("Successfully rendered LaTeX template")
	return output.String(), nil
}

// templateFuncs returns helper functions for LaTeX templates
func (g *LaTeXGenerator) templateFuncs() template.FuncMap {
	// Escape special LaTeX characters function
	escapeFunc := func(str string) string {
		replacer := strings.NewReplacer(
			`\`, `\textbackslash{}`,
			`{`, `\{`,
			`}`, `\}`,
			`$`, `\$`,
			`&`, `\&`,
			`%`, `\%`,
			`#`, `\#`,
			`_`, `\_`,
			`~`, `\textasciitilde{}`,
			`^`, `\textasciicircum{}`,
		)
		return replacer.Replace(str)
	}

	return template.FuncMap{
		// Escape special LaTeX characters
		"escape":           escapeFunc,
		"escapeLatexChars": escapeFunc, // Alias

		// Format date ranges
		"fmtDateRange": func(dates definition.EnhancedDateRange) string {
			start := dates.Start.Format("Jan 2006")
			if dates.End != nil {
				return fmt.Sprintf("%s - %s", start, dates.End.Format("Jan 2006"))
			}
			if dates.Current {
				return fmt.Sprintf("%s - Present", start)
			}
			return start
		},

		// Format location
		"fmtLocation": func(loc *definition.Location) string {
			if loc == nil {
				return "Remote"
			}
			parts := []string{}
			if loc.City != "" {
				parts = append(parts, loc.City)
			}
			if loc.State != "" {
				parts = append(parts, loc.State)
			}
			if len(parts) > 0 {
				return strings.Join(parts, ", ")
			}
			return "Remote"
		},

		// Join strings with separator
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},

		// Extract skill names from EnhancedSkillItem slice
		"skillNames": func(items []definition.EnhancedSkillItem) []string {
			names := make([]string, len(items))
			for i, item := range items {
				names[i] = item.Name
			}
			return names
		},

		// Format link for LaTeX (handles both Link and EnhancedLink types)
		"fmtLink": func(link interface{}) string {
			var url, text string

			switch v := link.(type) {
			case definition.Link:
				url = v.Ref
				text = v.Text
			case definition.EnhancedLink:
				url = v.URL
				text = v.Text
			default:
				return ""
			}

			if url == "" {
				return ""
			}
			// Format as \href{url}{text}
			if text == "" {
				text = url
			}
			return fmt.Sprintf("\\href{%s}{%s}", url, text)
		},

		// Format dates (simple string representation for legacy format)
		"fmtDates": func(dates interface{}) string {
			// For legacy format, dates might be strings
			switch v := dates.(type) {
			case string:
				return v
			case definition.EnhancedDateRange:
				start := v.Start.Format("Jan 2006")
				if v.End != nil {
					return fmt.Sprintf("%s - %s", start, v.End.Format("Jan 2006"))
				}
				if v.Current {
					return fmt.Sprintf("%s - Present", start)
				}
				return start
			default:
				return ""
			}
		},

		// Convert to title case
		"title": strings.Title,

		// Convert to uppercase
		"upper": strings.ToUpper,

		// Convert to lowercase
		"lower": strings.ToLower,
	}
}
