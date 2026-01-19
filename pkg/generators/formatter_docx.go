package generators

import (
	"strings"
	"text/template"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

// docxFormatter implements Formatter for DOCX output.
// DOCX generation is programmatic, so most methods return plain text
// without special escaping since go-docx handles text encoding.
type docxFormatter struct {
	baseFormatter
}

func newDocxFormatter() *docxFormatter {
	return &docxFormatter{}
}

// EscapeText returns text as-is since go-docx handles encoding.
func (f *docxFormatter) EscapeText(value string) string {
	return value
}

// FormatDateRange formats dates, detecting year-only dates.
func (f *docxFormatter) FormatDateRange(dr resume.DateRange) string {
	if dr.End != nil {
		return f.formatDateShort(dr.Start) + " " + f.formatDateShort(*dr.End)
	}
	return f.formatDateShort(dr.Start)
}

// FormatLocation renders a user location string.
func (f *docxFormatter) FormatLocation(loc *resume.Location) string {
	return f.baseFormatter.FormatLocation(loc, nil)
}

// FormatGPA renders GPA with optional denominator (DOCX uses "/" style).
func (f *docxFormatter) FormatGPA(gpa, max string) string {
	gpa = strings.TrimSpace(gpa)
	max = strings.TrimSpace(max)
	if gpa == "" {
		return ""
	}
	if max != "" && max != "4.0" {
		return gpa + "/" + max
	}
	return gpa
}

// SanitizePhone returns phone as-is for DOCX.
func (f *docxFormatter) SanitizePhone(phone string) string {
	return phone
}

// TemplateFuncs returns an empty FuncMap since DOCX doesn't use Go templates.
func (f *docxFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{}
}

// formatDateShort returns a short date format (Jan 2006 or just 2006 for year-only dates).
func (f *docxFormatter) formatDateShort(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// If it's January 1st, likely a year-only date
	if t.Month() == time.January && t.Day() == 1 {
		return t.Format("2006")
	}
	return t.Format("Jan 2006")
}
