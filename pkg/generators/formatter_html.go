package generators

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

// htmlFormatter provides HTML-specific formatting, embedding shared logic from baseFormatter.
type htmlFormatter struct {
	baseFormatter
}

func newHTMLFormatter() *htmlFormatter {
	return &htmlFormatter{}
}

// EscapeText escapes HTML special characters.
func (f *htmlFormatter) EscapeText(value string) string {
	return template.HTMLEscapeString(value)
}

// FormatLocation renders a location with HTML escaping.
func (f *htmlFormatter) FormatLocation(loc *resume.Location) string {
	return f.baseFormatter.FormatLocation(loc, nil) // No escaping needed for plain text display
}

// FormatLink renders an HTML anchor tag.
func (f *htmlFormatter) FormatLink(link string) string {
	url := strings.TrimSpace(link)
	if url == "" {
		return ""
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, template.HTMLEscapeString(url), template.HTMLEscapeString(url))
}

// TemplateFuncs exposes helper functions for HTML templates.
func (f *htmlFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Text escaping
		"escape":   f.EscapeText,
		"safeHTML": func(value string) template.HTML { return template.HTML(value) },

		// Date formatting
		"formatDate":        func(t time.Time) string { return t.Format("January 2006") },
		"formatDateShort":   func(t time.Time) string { return t.Format("Jan 2006") },
		"formatDateRange":   f.formatDateRange,
		"fmtDateRange":      f.FormatDateRange,
		"calculateDuration": f.CalculateDuration,

		// Location formatting
		"formatLocation": func(loc *resume.Location) string { return f.FormatLocation(loc) },
		"fmtLocation": func(value interface{}) string {
			switch v := value.(type) {
			case *resume.Location:
				return f.FormatLocation(v)
			case resume.Location:
				return f.FormatLocation(&v)
			default:
				return ""
			}
		},

		// List formatting
		"formatList":  f.FormatList,
		"join":        f.Join,
		"skillNames":  f.SkillNames,
		"filterEmpty": filterStrings,

		// Case transformations
		"lower": f.Lower,
		"upper": f.Upper,
		"title": f.Title,

		// String utilities
		"replace":   strings.ReplaceAll,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"contains":  strings.Contains,
		"trim":      strings.TrimSpace,

		// Link formatting
		"formatLink": f.FormatLink,
		"fmtLink": func(value interface{}) string {
			switch v := value.(type) {
			case string:
				return f.FormatLink(v)
			default:
				return ""
			}
		},

		// GPA formatting
		"formatGPA": f.FormatGPAStruct,

		// Phone sanitization
		"sanitizePhone": f.SanitizePhone,

		// Sort functions (preserved for template compatibility - return input unchanged)
		"sortSkillsByOrder":     func(categories []resume.SkillCategory) []resume.SkillCategory { return categories },
		"sortExperienceByOrder": func(experiences []resume.Experience) []resume.Experience { return experiences },
		"sortProjectsByOrder":   func(projects []resume.Project) []resume.Project { return projects },
		"sortEducationByOrder":  func(education []resume.Education) []resume.Education { return education },
		"sortLinksByOrder":      func(links []string) []string { return links },

		// Default value helper
		"default": func(defaultVal, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultVal
			}
			return value
		},
	}
}

// formatDateRange is a template-friendly version accepting individual args.
func (f *htmlFormatter) formatDateRange(start time.Time, end *time.Time) string {
	return f.baseFormatter.formatDateRangeInternal(start, end)
}
