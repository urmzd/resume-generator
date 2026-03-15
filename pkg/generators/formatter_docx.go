package generators

import (
	"bytes"
	"encoding/xml"
	"strings"
	"text/template"

	"github.com/urmzd/resume-generator/pkg/resume"
)

// docxFormatter implements Formatter for DOCX output.
// It provides XML escaping and helper functions for the DOCX Go template.
type docxFormatter struct {
	baseFormatter
}

func newDocxFormatter() *docxFormatter {
	return &docxFormatter{}
}

// xmlEscape escapes text for safe inclusion in XML content.
func xmlEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

// EscapeText XML-escapes special characters for DOCX content.
func (f *docxFormatter) EscapeText(value string) string {
	return xmlEscape(value)
}

// FormatDateRange formats dates using precision-aware formatting with en-dash.
func (f *docxFormatter) FormatDateRange(dr resume.DateRange) string {
	return f.formatDateRangeInternal(dr.Start, dr.End)
}

// FormatLocation renders a user location string.
func (f *docxFormatter) FormatLocation(loc *resume.Location) string {
	return f.baseFormatter.FormatLocation(loc, nil)
}

// FormatGPA renders GPA with optional denominator.
func (f *docxFormatter) FormatGPA(gpa, max string) string {
	return f.baseFormatter.FormatGPA(gpa, max)
}

// FormatGPAStruct renders GPA from a *resume.GPA struct.
func (f *docxFormatter) FormatGPAStruct(g *resume.GPA) string {
	if g == nil {
		return ""
	}
	return f.FormatGPA(g.GPA, g.MaxGPA)
}

// SanitizePhone returns phone as-is for DOCX.
func (f *docxFormatter) SanitizePhone(phone string) string {
	return phone
}

// TemplateFuncs exposes helper functions for the DOCX XML template.
func (f *docxFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Text escaping
		"escape": f.EscapeText,

		// Case transformations
		"upper": f.Upper,
		"lower": f.Lower,
		"title": f.Title,

		// String utilities
		"trim":        strings.TrimSpace,
		"filterEmpty": filterStrings,
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"skillNames": f.SkillNames,
		"default": func(defaultVal, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultVal
			}
			return value
		},

		// Date formatting
		"fmtDateRange": f.FormatDateRange,
		"fmtDates":     f.FormatDates,

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

		// GPA formatting
		"formatGPA": f.FormatGPAStruct,

		// Sort functions
		"sortExperienceByOrder": f.SortExperienceByDate,
		"sortProjectsByOrder":   f.SortProjectsByDate,
		"sortEducationByOrder":  f.SortEducationByDate,

		// dict helper for passing data to sub-templates
		"dict": func(pairs ...interface{}) map[string]interface{} {
			m := make(map[string]interface{}, len(pairs)/2)
			for i := 0; i < len(pairs)-1; i += 2 {
				key, ok := pairs[i].(string)
				if ok {
					m[key] = pairs[i+1]
				}
			}
			return m
		},

		// Contact line builder
		"contactLine": func(c resume.Contact) string {
			var parts []string
			if c.Location != nil {
				loc := f.FormatLocation(c.Location)
				if loc != "" {
					parts = append(parts, loc)
				}
			}
			if c.Email != "" {
				parts = append(parts, c.Email)
			}
			if c.Phone != "" {
				parts = append(parts, c.Phone)
			}
			if c.Credentials != "" {
				parts = append(parts, c.Credentials)
			}
			for _, link := range c.Links {
				if link.URI != "" {
					parts = append(parts, link.URI)
				}
			}
			return strings.Join(parts, " | ")
		},

		// Certification line builder
		"certLine": func(cert resume.Certification) string {
			line := cert.Name
			if cert.Issuer != "" {
				line += " \u2014 " + cert.Issuer
			}
			if cert.Notes != "" {
				line += " (" + cert.Notes + ")"
			}
			return line
		},

		// Education header builder
		"eduHeader": func(inst resume.Education) string {
			text := inst.Institution
			if inst.Degree.Name != "" {
				text += ", " + inst.Degree.Name
			}
			dates := f.FormatDateRange(inst.Dates)
			if dates != "" {
				text += " \u2014 " + dates
			}
			return text
		},

		// Thesis line builder
		"thesisLine": func(thesis *resume.Thesis) string {
			if thesis == nil {
				return ""
			}
			title := strings.TrimSpace(thesis.Title)
			if title == "" {
				return ""
			}
			line := "Thesis: " + title
			if url := strings.TrimSpace(thesis.Link.URI); url != "" {
				line += " (" + url + ")"
			}
			var descs []string
			for _, desc := range thesis.Highlights {
				if d := strings.TrimSpace(desc); d != "" {
					descs = append(descs, d)
				}
			}
			if len(descs) > 0 {
				line += " \u2014 " + strings.Join(descs, "; ")
			}
			return line
		},

		// Experience title line builder
		"expTitleLine": func(pos resume.Experience) string {
			titleLine := pos.Title
			dates := f.FormatDateRange(pos.Dates)
			if dates != "" {
				titleLine += " \u2014 " + dates
			}
			return titleLine
		},

		// Experience company line builder
		"expCompanyLine": func(pos resume.Experience) string {
			var parts []string
			if pos.Company != "" {
				parts = append(parts, pos.Company)
			}
			if pos.Location != nil {
				loc := f.FormatLocation(pos.Location)
				if loc != "" {
					parts = append(parts, loc)
				}
			}
			return strings.Join(parts, " | ")
		},

		// Language line builder
		"langLine": func(lang resume.Language) string {
			line := lang.Name
			if lang.Proficiency != "" {
				line += " \u2014 " + lang.Proficiency
			}
			return line
		},
	}
}

// formatDateShort returns a short date format using the date's precision.
func (f *docxFormatter) formatDateShort(pd resume.PartialDate) string {
	if pd.IsZero() {
		return ""
	}
	if pd.Precision == resume.PrecisionYear {
		return pd.Time.Format("2006")
	}
	return pd.Time.Format("Jan 2006")
}
