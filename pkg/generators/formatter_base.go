package generators

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// baseFormatter contains shared formatting logic used by all output formatters.
// Output-specific formatters embed this and override only what differs (escaping, links).
type baseFormatter struct{}

// FormatDateRange converts a DateRange to a human-readable string like "Jan 2020 – Present".
func (f *baseFormatter) FormatDateRange(dates resume.DateRange) string {
	return f.formatDateRangeInternal(dates.Start, dates.End)
}

// FormatOptionalDateRange handles a potentially nil DateRange pointer.
func (f *baseFormatter) FormatOptionalDateRange(dates *resume.DateRange) string {
	if dates == nil {
		return ""
	}
	return f.FormatDateRange(*dates)
}

// FormatDates handles legacy date representations (string or DateRange).
func (f *baseFormatter) FormatDates(value interface{}) string {
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

// FormatLocation renders a Location as "City, State, Country".
// The escape function is provided by the output-specific formatter.
func (f *baseFormatter) FormatLocation(loc *resume.Location, escape func(string) string) string {
	if loc == nil {
		return ""
	}

	parts := make([]string, 0, 3)
	if city := strings.TrimSpace(loc.City); city != "" {
		parts = append(parts, city)
	}
	if state := strings.TrimSpace(loc.State); state != "" {
		parts = append(parts, state)
	}
	if country := strings.TrimSpace(loc.Country); country != "" {
		// Avoid duplicating country if already in parts (case-insensitive)
		if !f.containsIgnoreCase(parts, country) {
			parts = append(parts, country)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	result := strings.Join(parts, ", ")
	if escape != nil {
		return escape(result)
	}
	return result
}

// FormatGPA renders GPA with optional denominator from string arguments.
func (f *baseFormatter) FormatGPA(gpa, max string) string {
	gpa = strings.TrimSpace(gpa)
	max = strings.TrimSpace(max)
	if gpa == "" {
		return ""
	}
	if max == "" || max == "4.0" {
		return gpa
	}
	return fmt.Sprintf("%s / %s", gpa, max)
}

// FormatGPAStruct renders GPA from a *resume.GPA struct.
func (f *baseFormatter) FormatGPAStruct(g *resume.GPA) string {
	if g == nil {
		return ""
	}
	return f.FormatGPA(g.GPA, g.MaxGPA)
}

// SanitizePhone removes non-numeric characters except +.
func (f *baseFormatter) SanitizePhone(phone string) string {
	var builder strings.Builder
	for _, r := range phone {
		if (r >= '0' && r <= '9') || r == '+' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// FormatList joins non-empty strings with commas.
func (f *baseFormatter) FormatList(values []string) string {
	return strings.Join(filterStrings(values), ", ")
}

// SkillNames returns filtered skill names for display.
func (f *baseFormatter) SkillNames(items []string) []string {
	return filterStrings(items)
}

// Join concatenates strings using a separator.
func (f *baseFormatter) Join(sep string, items []string) string {
	return strings.Join(items, sep)
}

// Lower converts text to lower-case.
func (f *baseFormatter) Lower(value string) string {
	return strings.ToLower(value)
}

// Upper converts text to upper-case.
func (f *baseFormatter) Upper(value string) string {
	return strings.ToUpper(value)
}

// Title converts text to title-case.
func (f *baseFormatter) Title(value string) string {
	return cases.Title(language.English).String(value)
}

// formatDateRangeInternal is the internal implementation for date range formatting.
func (f *baseFormatter) formatDateRangeInternal(start time.Time, end *time.Time) string {
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
	return fmt.Sprintf("%s – %s", startStr, endStr)
}

// formatMonthYear formats a time as "Jan 2006".
func (f *baseFormatter) formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}

// CalculateDuration returns duration as "X yr Y mo" string.
func (f *baseFormatter) CalculateDuration(start time.Time, end *time.Time) string {
	var endTime time.Time
	if end == nil {
		endTime = time.Now()
	} else {
		endTime = *end
	}

	diff := endTime.Sub(start)
	years := int(diff.Hours() / 24 / 365)
	months := int((diff.Hours() / 24 / 30)) % 12

	switch {
	case years > 0 && months > 0:
		return fmt.Sprintf("%d yr %d mo", years, months)
	case years > 0:
		return fmt.Sprintf("%d yr", years)
	case months > 0:
		return fmt.Sprintf("%d mo", months)
	default:
		return "< 1 mo"
	}
}

// SortExperienceByDate returns a copy of experiences sorted by start date descending.
// Entries with no end date (current) sort first.
func (f *baseFormatter) SortExperienceByDate(experiences []resume.Experience) []resume.Experience {
	sorted := make([]resume.Experience, len(experiences))
	copy(sorted, experiences)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Dates.Start.After(sorted[j].Dates.Start)
	})
	return sorted
}

// SortEducationByDate returns a copy of education entries sorted by start date descending.
func (f *baseFormatter) SortEducationByDate(education []resume.Education) []resume.Education {
	sorted := make([]resume.Education, len(education))
	copy(sorted, education)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Dates.Start.After(sorted[j].Dates.Start)
	})
	return sorted
}

// SortProjectsByDate returns a copy of projects sorted by start date descending.
// Projects without dates sort to the end.
func (f *baseFormatter) SortProjectsByDate(projects []resume.Project) []resume.Project {
	sorted := make([]resume.Project, len(projects))
	copy(sorted, projects)
	sort.SliceStable(sorted, func(i, j int) bool {
		di := sorted[i].Dates
		dj := sorted[j].Dates
		if di == nil && dj == nil {
			return false
		}
		if di == nil {
			return false
		}
		if dj == nil {
			return true
		}
		return di.Start.After(dj.Start)
	})
	return sorted
}

// containsIgnoreCase checks if a slice contains a value (case-insensitive).
func (f *baseFormatter) containsIgnoreCase(list []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range list {
		if strings.ToLower(strings.TrimSpace(item)) == value {
			return true
		}
	}
	return false
}
