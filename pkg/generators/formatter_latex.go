package generators

import (
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
)

type latexFormatter struct{}

func newLaTeXFormatter() *latexFormatter {
	return &latexFormatter{}
}

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

func (f *latexFormatter) EscapeText(value string) string {
	return latexEscaper.Replace(value)
}

func (f *latexFormatter) FormatDateRange(dates definition.DateRange) string {
	return f.formatDateRange(dates.Start, dates.End, dates.Current)
}

func (f *latexFormatter) FormatOptionalDateRange(dates *definition.DateRange) string {
	if dates == nil {
		return ""
	}
	return f.FormatDateRange(*dates)
}

func (f *latexFormatter) FormatDates(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case definition.DateRange:
		return f.FormatDateRange(v)
	case *definition.DateRange:
		if v == nil {
			return ""
		}
		return f.FormatDateRange(*v)
	default:
		return ""
	}
}

func (f *latexFormatter) FormatCertificationDate(issue, expiration *time.Time) string {
	if (issue == nil || issue.IsZero()) && (expiration == nil || expiration.IsZero()) {
		return ""
	}

	issueStr := ""
	if issue != nil && !issue.IsZero() {
		issueStr = f.formatMonthYear(*issue)
	}

	expStr := ""
	if expiration != nil && !expiration.IsZero() {
		expStr = f.formatMonthYear(*expiration)
	}

	switch {
	case issueStr == "":
		return expStr
	case expStr == "":
		return issueStr
	default:
		return fmt.Sprintf("%s - %s", issueStr, expStr)
	}
}

func (f *latexFormatter) FormatLocation(loc *definition.Location) string {
	if loc == nil {
		return ""
	}

	parts := make([]string, 0, 3)
	if strings.TrimSpace(loc.City) != "" {
		parts = append(parts, strings.TrimSpace(loc.City))
	}
	if strings.TrimSpace(loc.State) != "" {
		parts = append(parts, strings.TrimSpace(loc.State))
	} else if strings.TrimSpace(loc.Region) != "" {
		parts = append(parts, strings.TrimSpace(loc.Region))
	}
	if strings.TrimSpace(loc.Country) != "" {
		parts = append(parts, strings.TrimSpace(loc.Country))
	}

	if len(parts) == 0 {
		return ""
	}

	return f.EscapeText(strings.Join(parts, ", "))
}

func (f *latexFormatter) FormatList(values []string) string {
	filtered := filterStrings(values)
	for i := range filtered {
		filtered[i] = f.EscapeText(filtered[i])
	}
	return strings.Join(filtered, ", ")
}

func (f *latexFormatter) FormatGPA(gpa, max string) string {
	gpa = strings.TrimSpace(gpa)
	max = strings.TrimSpace(max)
	if gpa == "" {
		return ""
	}
	if max == "" || max == "4.0" {
		return f.EscapeText(gpa)
	}
	return fmt.Sprintf("%s / %s", f.EscapeText(gpa), f.EscapeText(max))
}

func (f *latexFormatter) SkillNames(items []definition.SkillItem) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if name := strings.TrimSpace(item.Name); name != "" {
			result = append(result, name)
		}
	}
	return result
}

func (f *latexFormatter) Join(sep string, items []string) string {
	return strings.Join(items, sep)
}

func (f *latexFormatter) FormatLink(link definition.Link) string {
	url := strings.TrimSpace(link.URL)
	if url == "" {
		return ""
	}
	text := strings.TrimSpace(link.Text)
	if text == "" {
		text = url
	}
	return fmt.Sprintf(`\href{%s}{%s}`, f.EscapeText(url), f.EscapeText(text))
}

func (f *latexFormatter) Lower(value string) string {
	return strings.ToLower(value)
}

func (f *latexFormatter) Upper(value string) string {
	return strings.ToUpper(value)
}

func (f *latexFormatter) Title(value string) string {
	return strings.Title(value)
}

func (f *latexFormatter) SanitizePhone(phone string) string {
	var builder strings.Builder
	for _, r := range phone {
		if (r >= '0' && r <= '9') || r == '+' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func (f *latexFormatter) TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"escape":           f.EscapeText,
		"escapeLatexChars": f.EscapeText,
		"fmtDateRange":     f.FormatDateRange,
		"fmtDates":         f.FormatDates,
		"join":             f.Join,
		"skillNames":       f.SkillNames,
		"fmtLink": func(value interface{}) string {
			switch v := value.(type) {
			case definition.Link:
				return f.FormatLink(v)
			default:
				return ""
			}
		},
		"fmtLocation": func(value interface{}) string {
			switch v := value.(type) {
			case *definition.Location:
				return f.FormatLocation(v)
			case definition.Location:
				return f.FormatLocation(&v)
			default:
				return ""
			}
		},
		"title": f.Title,
		"upper": f.Upper,
		"lower": f.Lower,
	}
}

func (f *latexFormatter) formatDateRange(start time.Time, end *time.Time, current bool) string {
	if start.IsZero() && (end == nil || (end != nil && end.IsZero())) && !current {
		return ""
	}

	startStr := f.formatMonthYear(start)
	var endStr string

	switch {
	case current:
		endStr = "Present"
	case end != nil && !end.IsZero():
		endStr = f.formatMonthYear(*end)
	}

	if startStr == "" {
		return endStr
	}
	if endStr == "" || startStr == endStr {
		return startStr
	}
	return fmt.Sprintf("%s - %s", startStr, endStr)
}

func (f *latexFormatter) formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}
