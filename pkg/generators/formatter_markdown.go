package generators

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

// markdownFormatter provides Markdown-specific formatting, embedding shared logic from baseFormatter.
type markdownFormatter struct {
	baseFormatter
}

func newMarkdownFormatter() *markdownFormatter {
	return &markdownFormatter{}
}

// markdownEscaper replaces Markdown special characters.
var markdownEscaper = strings.NewReplacer(
	`\`, `\\`,
	"`", "\\`",
	`*`, `\*`,
	`_`, `\_`,
	`{`, `\{`,
	`}`, `\}`,
	`[`, `\[`,
	`]`, `\]`,
	`(`, `\(`,
	`)`, `\)`,
	`#`, `\#`,
	`+`, `\+`,
	`-`, `\-`,
	`.`, `\.`,
	`!`, `\!`,
	`|`, `\|`,
)

// EscapeText escapes Markdown special characters.
func (f *markdownFormatter) EscapeText(value string) string {
	return markdownEscaper.Replace(value)
}

// FormatLocation renders a location without escaping (plain text is fine in markdown context).
func (f *markdownFormatter) FormatLocation(loc *resume.Location) string {
	return f.baseFormatter.FormatLocation(loc, nil)
}

// FormatLink renders a Markdown link.
func (f *markdownFormatter) FormatLink(link string) string {
	url := strings.TrimSpace(link)
	if url == "" {
		return ""
	}
	display := f.ExtractDisplayURL(url)
	return fmt.Sprintf("[%s](%s)", display, url)
}

// ExtractDisplayURL removes protocol and www prefix for cleaner display.
func (f *markdownFormatter) ExtractDisplayURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")
	url = strings.TrimSuffix(url, "/")
	return url
}

// TemplateFuncs exposes helper functions for Markdown templates.
func (f *markdownFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Text escaping
		"escape": f.EscapeText,

		// Date formatting
		"fmtDateRange":    f.FormatDateRange,
		"fmtOptDateRange": f.FormatOptionalDateRange,
		"fmtDates":        f.FormatDates,
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("January 2006")
		},
		"formatDateShort": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("Jan 2006")
		},

		// Location formatting
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
		"formatList": f.FormatList,
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"skillNames":  f.SkillNames,
		"filterEmpty": filterStrings,

		// Link formatting
		"fmtLink": func(value interface{}) string {
			switch v := value.(type) {
			case string:
				return f.FormatLink(v)
			default:
				return ""
			}
		},
		"fmtLinkWithDomain": func(value interface{}) string {
			switch v := value.(type) {
			case string:
				return f.FormatLink(v)
			case resume.Link:
				return f.FormatLink(v.URI)
			case *resume.Link:
				if v == nil {
					return ""
				}
				return f.FormatLink(v.URI)
			default:
				return ""
			}
		},
		"extractDisplayURL": f.ExtractDisplayURL,

		// GPA formatting
		"formatGPA": f.FormatGPAStruct,

		// Phone sanitization
		"sanitizePhone": f.SanitizePhone,

		// Case transformations
		"title": f.Title,
		"upper": f.Upper,
		"lower": f.Lower,

		// String utilities
		"trim":         strings.TrimSpace,
		"filterEmpty2": filterStrings,
		"default": func(defaultVal, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultVal
			}
			return value
		},

		// Markdown helpers
		"bold": func(s string) string {
			if s == "" {
				return ""
			}
			return "**" + s + "**"
		},
		"italic": func(s string) string {
			if s == "" {
				return ""
			}
			return "*" + s + "*"
		},

		// Math utilities
		"add": func(a, b int) int { return a + b },
	}
}
