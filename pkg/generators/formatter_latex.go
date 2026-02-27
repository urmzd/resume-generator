package generators

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

// latexFormatter provides LaTeX-specific formatting, embedding shared logic from baseFormatter.
type latexFormatter struct {
	baseFormatter
	// autoEscaped is set when all resume string data has been pre-escaped.
	// When true, EscapeText becomes identity to prevent double-escaping.
	autoEscaped bool
}

func newLaTeXFormatter() *latexFormatter {
	return &latexFormatter{}
}

// latexEscaper replaces LaTeX special characters.
var latexEscaper = strings.NewReplacer(
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

// EscapeText escapes LaTeX special characters.
// When autoEscaped is set, returns the input unchanged (data is already escaped).
func (f *latexFormatter) EscapeText(value string) string {
	if f.autoEscaped {
		return value
	}
	return latexEscaper.Replace(value)
}

// FormatLocation renders a location with LaTeX escaping.
func (f *latexFormatter) FormatLocation(loc *resume.Location) string {
	return f.baseFormatter.FormatLocation(loc, f.EscapeText)
}

// FormatList joins values with LaTeX escaping.
func (f *latexFormatter) FormatList(values []string) string {
	filtered := filterStrings(values)
	for i := range filtered {
		filtered[i] = f.EscapeText(filtered[i])
	}
	return strings.Join(filtered, ", ")
}

// FormatGPA renders GPA with LaTeX escaping.
func (f *latexFormatter) FormatGPA(gpa, max string) string {
	result := f.baseFormatter.FormatGPA(gpa, max)
	return f.EscapeText(result)
}

// FormatGPAStruct renders GPA from a *resume.GPA struct with LaTeX escaping.
func (f *latexFormatter) FormatGPAStruct(g *resume.GPA) string {
	if g == nil {
		return ""
	}
	return f.FormatGPA(g.GPA, g.MaxGPA)
}

// FormatDateRange overrides the base formatter to use LaTeX-specific en-dash.
func (f *latexFormatter) FormatDateRange(dates resume.DateRange) string {
	return f.formatDateRangeLaTeX(dates.Start, dates.End)
}

// formatDateRangeLaTeX formats dates using LaTeX \textendash\ for the en-dash.
func (f *latexFormatter) formatDateRangeLaTeX(start time.Time, end *time.Time) string {
	if start.IsZero() && (end == nil || end.IsZero()) {
		return ""
	}

	startStr := f.formatMonthYear(start)
	var endStr string

	switch {
	case end == nil:
		endStr = "Present"
	case !end.IsZero():
		endStr = f.formatMonthYear(*end)
	default:
		endStr = "Present"
	}

	if startStr == "" {
		return endStr
	}
	if endStr == "" || startStr == endStr {
		return startStr
	}
	return fmt.Sprintf(`%s \textendash\ %s`, startStr, endStr)
}

// formatMonthYear formats a time as "Jan 2006".
func (f *latexFormatter) formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}

// FormatDates overrides the base formatter to use LaTeX-specific en-dash.
func (f *latexFormatter) FormatDates(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case resume.DateRange:
		return f.FormatDateRange(v)
	case *resume.DateRange:
		if v == nil {
			return ""
		}
		return f.FormatDateRange(*v)
	default:
		return ""
	}
}

// FormatLink renders a LaTeX \href command.
func (f *latexFormatter) FormatLink(link string) string {
	url := strings.TrimSpace(link)
	if url == "" {
		return ""
	}
	return fmt.Sprintf(`\href{%s}{%s}`, f.EscapeText(url), f.EscapeText(url))
}

// ExtractDisplayURL removes protocol and www prefix for cleaner display.
func (f *latexFormatter) ExtractDisplayURL(url string) string {
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

// FormatLinkWithDomain renders a link using domain as display text.
func (f *latexFormatter) FormatLinkWithDomain(link string) string {
	url := strings.TrimSpace(link)
	if url == "" {
		return ""
	}
	displayURL := f.ExtractDisplayURL(url)
	return fmt.Sprintf(`\href{%s}{%s}`, f.EscapeText(url), f.EscapeText(displayURL))
}

// TemplateFuncs exposes helper functions for LaTeX templates.
func (f *latexFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Text escaping
		"escape":           f.EscapeText,
		"escapeLatexChars": f.EscapeText,

		// Date formatting
		"fmtDateRange": f.FormatDateRange,
		"fmtDates":     f.FormatDates,
		"formatDateRange": func(start time.Time, end *time.Time) string {
			return f.formatDateRangeInternal(start, end)
		},
		"fmtDateLegal": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("01/02/2006")
		},

		// List formatting
		"join": func(sep string, items []string) string {
			escaped := make([]string, len(items))
			for i, item := range items {
				escaped[i] = f.EscapeText(item)
			}
			return strings.Join(escaped, sep)
		},
		"formatList": f.FormatList,
		"skillNames": f.SkillNames,

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
				return f.FormatLinkWithDomain(v)
			case resume.Link:
				return f.FormatLinkWithDomain(v.URI)
			case *resume.Link:
				if v == nil {
					return ""
				}
				return f.FormatLinkWithDomain(v.URI)
			default:
				return ""
			}
		},
		"extractDisplayURL": f.ExtractDisplayURL,

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
		"formatLocationFull": func(loc *resume.Location) string {
			if loc == nil {
				return ""
			}
			parts := []string{}
			if loc.City != "" {
				parts = append(parts, loc.City)
			}
			// Use Province if set, otherwise State
			region := loc.Province
			if region == "" {
				region = loc.State
			}
			if region != "" {
				parts = append(parts, region)
			}
			if loc.Country != "" {
				parts = append(parts, loc.Country)
			}
			result := strings.Join(parts, ", ")
			if loc.Remote {
				result += " (Remote)"
			}
			return f.EscapeText(result)
		},
		"formatLocationShort": func(loc *resume.Location) string {
			if loc == nil {
				return ""
			}
			if loc.Province != "" {
				return f.EscapeText(loc.Province)
			}
			return f.EscapeText(loc.State)
		},

		// GPA formatting
		"formatGPA": f.FormatGPAStruct,

		// Case transformations
		"title": f.Title,
		"upper": f.Upper,
		"lower": f.Lower,

		// String utilities
		"trim":        strings.TrimSpace,
		"filterEmpty": filterStrings,
		"default": func(defaultVal, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultVal
			}
			return value
		},

		// Math utilities
		"add": func(a, b int) int { return a + b },

		// Employment type helper
		"employmentType": func(et string) string {
			if et == "" {
				return "Full-Time"
			}
			return et
		},

		// Current time helper
		"now": func() time.Time { return time.Now() },

		// Link label helper
		"linkLabel": func(url string) string {
			lower := strings.ToLower(url)
			switch {
			case strings.Contains(lower, "github"):
				return "GitHub"
			case strings.Contains(lower, "linkedin"):
				return "LinkedIn"
			case strings.Contains(lower, "twitter"), strings.Contains(lower, "x.com"):
				return "Twitter"
			default:
				return "Website"
			}
		},
	}
}
