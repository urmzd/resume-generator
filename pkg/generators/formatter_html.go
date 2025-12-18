package generators

import (
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
)

type htmlFormatter struct{}

func newHTMLFormatter() *htmlFormatter {
	return &htmlFormatter{}
}

func (f *htmlFormatter) EscapeText(value string) string {
	return template.HTMLEscapeString(value)
}

func (f *htmlFormatter) FormatDateRange(dates definition.DateRange) string {
	return f.formatDateRange(dates.Start, dates.End, dates.Current)
}

func (f *htmlFormatter) FormatOptionalDateRange(dates *definition.DateRange) string {
	if dates == nil {
		return ""
	}
	return f.FormatDateRange(*dates)
}

func (f *htmlFormatter) FormatDates(value interface{}) string {
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

func (f *htmlFormatter) FormatCertificationDate(issue, expiration *time.Time) string {
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
		return fmt.Sprintf("%s – %s", issueStr, expStr)
	}
}

func (f *htmlFormatter) FormatLocation(loc *definition.Location) string {
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
	if strings.TrimSpace(loc.Country) != "" && !f.containsIgnoreCase(parts, loc.Country) {
		parts = append(parts, strings.TrimSpace(loc.Country))
	}

	return strings.Join(parts, ", ")
}

func (f *htmlFormatter) FormatList(values []string) string {
	return strings.Join(filterStrings(values), ", ")
}

func (f *htmlFormatter) FormatGPA(gpa, max string) string {
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

func (f *htmlFormatter) SkillNames(items []definition.SkillItem) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if name := strings.TrimSpace(item.Name); name != "" {
			result = append(result, name)
		}
	}
	return result
}

func (f *htmlFormatter) Join(sep string, items []string) string {
	return strings.Join(items, sep)
}

func (f *htmlFormatter) FormatLink(link definition.Link) string {
	url := strings.TrimSpace(link.URL)
	if url == "" {
		return ""
	}
	text := strings.TrimSpace(link.Text)
	if text == "" {
		text = url
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, template.HTMLEscapeString(url), template.HTMLEscapeString(text))
}

func (f *htmlFormatter) Lower(value string) string {
	return strings.ToLower(value)
}

func (f *htmlFormatter) Upper(value string) string {
	return strings.ToUpper(value)
}

func (f *htmlFormatter) Title(value string) string {
	return strings.Title(value)
}

func (f *htmlFormatter) SanitizePhone(phone string) string {
	var builder strings.Builder
	for _, r := range phone {
		if (r >= '0' && r <= '9') || r == '+' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func (f *htmlFormatter) TemplateFuncs() template.FuncMap {
	funcs := template.FuncMap{
		"escape":            f.EscapeText,
		"safeHTML":          func(value string) template.HTML { return template.HTML(value) },
		"formatDate":        func(t time.Time) string { return t.Format("January 2006") },
		"formatDateShort":   func(t time.Time) string { return t.Format("Jan 2006") },
		"formatDateRange":   f.formatDateRange,
		"calculateDuration": f.calculateDuration,
		"join":              f.Join,
		"lower":             f.Lower,
		"upper":             f.Upper,
		"title":             f.Title,
		"replace":           strings.ReplaceAll,
		"hasPrefix":         strings.HasPrefix,
		"hasSuffix":         strings.HasSuffix,
		"contains":          strings.Contains,
		"sortSkillsByOrder": f.sortSkillCategories,
		"sortExperienceByOrder": func(experiences []definition.Experience) []definition.Experience {
			sorted := make([]definition.Experience, len(experiences))
			copy(sorted, experiences)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortProjectsByOrder": func(projects []definition.Project) []definition.Project {
			sorted := make([]definition.Project, len(projects))
			copy(sorted, projects)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortEducationByOrder": func(education []definition.Education) []definition.Education {
			sorted := make([]definition.Education, len(education))
			copy(sorted, education)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortLinksByOrder": func(links []definition.Link) []definition.Link {
			sorted := make([]definition.Link, len(links))
			copy(sorted, links)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"sortCertificationsByOrder": func(certs []definition.Certification) []definition.Certification {
			sorted := make([]definition.Certification, len(certs))
			copy(sorted, certs)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
		"getIconClass": f.getIconClass,
		"formatGPA":    f.FormatGPA,
		"add":          func(a, b int) int { return a + b },
		"subtract":     func(a, b int) int { return a - b },
		"multiply":     func(a, b int) int { return a * b },
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"isEven":                  func(n int) bool { return n%2 == 0 },
		"isOdd":                   func(n int) bool { return n%2 != 0 },
		"formatCertificationDate": f.FormatCertificationDate,
	}

	return funcs
}

func (f *htmlFormatter) formatDateRange(start time.Time, end *time.Time, current bool) string {
	if start.IsZero() && (end == nil || (end != nil && end.IsZero())) && !current {
		return ""
	}

	startStr := f.formatMonthYear(start)
	var endStr string

	switch {
	case current:
		endStr = "Present"
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

func (f *htmlFormatter) calculateDuration(start time.Time, end *time.Time) string {
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

func (f *htmlFormatter) formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}

func (f *htmlFormatter) containsIgnoreCase(list []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range list {
		if strings.ToLower(strings.TrimSpace(item)) == value {
			return true
		}
	}
	return false
}

func (f *htmlFormatter) sortSkillCategories(skills []definition.SkillCategory) []definition.SkillCategory {
	sorted := make([]definition.SkillCategory, len(skills))
	copy(sorted, skills)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})
	return sorted
}

func (f *htmlFormatter) getIconClass(linkType string) string {
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
}
