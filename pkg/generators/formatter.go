package generators

import (
	"strings"
	"text/template"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
)

// Formatter defines the contract each rendering engine must satisfy to ensure
// consistent sanitisation and formatting behaviour across templates.
type Formatter interface {
	// EscapeText escapes engine-specific control characters.
	EscapeText(value string) string

	// FormatDateRange converts a date range to a human readable string.
	FormatDateRange(definition.DateRange) string

	// FormatOptionalDateRange formats a potentially nil date range pointer.
	FormatOptionalDateRange(*definition.DateRange) string

	// FormatDates handles legacy date representations used in templates.
	FormatDates(value interface{}) string

	// FormatCertificationDate renders certification issue/expiration dates.
	FormatCertificationDate(issue, expiration *time.Time) string

	// FormatLocation renders a user location string.
	FormatLocation(*definition.Location) string

	// FormatList renders a comma separated list after trimming empty values.
	FormatList([]string) string

	// FormatGPA renders GPA with optional denominator.
	FormatGPA(gpa, max string) string

	// SkillNames returns ordered skill names for display.
	SkillNames([]definition.SkillItem) []string

	// Join concatenates strings using a separator.
	Join(sep string, items []string) string

	// FormatLink renders a link using engine-specific markup.
	FormatLink(definition.Link) string

	// Lower converts text to lower-case using engine rules.
	Lower(string) string

	// Upper converts text to upper-case using engine rules.
	Upper(string) string

	// Title converts text to title-case using engine rules.
	Title(string) string

	// SanitizePhone removes unsupported characters from phone numbers.
	SanitizePhone(string) string

	// TemplateFuncs exposes helper functions for templating.
	TemplateFuncs() template.FuncMap
}

func filterStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
