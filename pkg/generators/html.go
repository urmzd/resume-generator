package generators

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

// HTMLGenerator generates HTML resumes from templates
type HTMLGenerator struct {
	logger *zap.SugaredLogger
	funcs  template.FuncMap
}

type htmlTemplatePayload struct {
	*definition.Resume
	CSS  template.CSS
	View *htmlResumeView
}

type htmlResumeView struct {
	Header         htmlHeaderView
	Summary        string
	Skills         *htmlSkillsView
	Experience     *htmlExperienceView
	Projects       *htmlProjectsView
	Education      *htmlEducationView
	Certifications *htmlCertificationsView
}

type htmlHeaderView struct {
	Name          string
	Title         string
	ContactItems  []htmlInlineItem
	HasContactRow bool
}

type htmlInlineItem struct {
	Text string
	URL  string
}

type htmlSkillsView struct {
	Title      string
	Categories []htmlSkillCategoryView
}

type htmlSkillCategoryView struct {
	Name    string
	Display string
}

type htmlExperienceView struct {
	Title     string
	Positions []htmlExperienceEntry
}

type htmlExperienceEntry struct {
	Title      string
	Company    string
	Location   string
	DateRange  string
	Highlights []string
}

type htmlProjectsView struct {
	Title   string
	Entries []htmlProjectEntry
}

type htmlProjectEntry struct {
	Name         string
	Category     string
	DateRange    string
	Description  []string
	Links        []htmlInlineItem
	Technologies string
}

type htmlEducationView struct {
	Title   string
	Schools []htmlEducationEntry
}

type htmlEducationEntry struct {
	Institution string
	Degree      string
	Field       string
	Location    string
	DateRange   string
	GPA         string
	Honors      string
	Details     []string
}

type htmlCertificationsView struct {
	Title          string
	Certifications []htmlCertificationEntry
}

type htmlCertificationEntry struct {
	Name            string
	Issuer          string
	DateRange       string
	CredentialID    string
	VerificationURL string
}

// NewHTMLGenerator creates a new HTML resume generator
func NewHTMLGenerator(logger *zap.SugaredLogger) *HTMLGenerator {
	generator := &HTMLGenerator{
		logger: logger,
	}
	generator.setupTemplateFunctions()
	return generator
}

// setupTemplateFunctions initializes template helper functions
func (g *HTMLGenerator) setupTemplateFunctions() {
	g.funcs = template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("January 2006")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("Jan 2006")
		},
		"formatDateRange": func(start time.Time, end *time.Time) string {
			startStr := start.Format("Jan 2006")
			if end == nil {
				return startStr + " - Present"
			}
			endStr := end.Format("Jan 2006")
			if startStr == endStr {
				return startStr
			}
			return startStr + " - " + endStr
		},
		"calculateDuration": func(start time.Time, end *time.Time) string {
			var endTime time.Time
			if end == nil {
				endTime = time.Now()
			} else {
				endTime = *end
			}

			diff := endTime.Sub(start)
			years := int(diff.Hours() / 24 / 365)
			months := int((diff.Hours() / 24 / 30)) % 12

			if years > 0 && months > 0 {
				return fmt.Sprintf("%d yr %d mo", years, months)
			} else if years > 0 {
				return fmt.Sprintf("%d yr", years)
			} else if months > 0 {
				return fmt.Sprintf("%d mo", months)
			}
			return "< 1 mo"
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"escape": func(s string) string {
			return template.HTMLEscapeString(s)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"contains":  strings.Contains,
		"sortSkillsByOrder": func(skills []definition.SkillCategory) []definition.SkillCategory {
			sorted := make([]definition.SkillCategory, len(skills))
			copy(sorted, skills)
			sort.Slice(sorted, func(i, j int) bool {
				return sorted[i].Order < sorted[j].Order
			})
			return sorted
		},
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
		"getIconClass": func(linkType string) string {
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
		},
		"formatGPA": func(gpa, maxGPA string) string {
			if gpa == "" {
				return ""
			}
			if maxGPA != "" && maxGPA != "4.0" {
				return fmt.Sprintf("%s/%s", gpa, maxGPA)
			}
			return gpa
		},
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"divide": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"isEven": func(n int) bool {
			return n%2 == 0
		},
		"isOdd": func(n int) bool {
			return n%2 != 0
		},
	}
}

// Generate creates an HTML resume from the resume data and template
func (g *HTMLGenerator) Generate(templateContent string, resume *definition.Resume) (string, error) {
	g.logger.Info("Generating HTML resume")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	payload := g.buildTemplatePayload(resume, "")
	if err := tmpl.Execute(&buf, payload); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume")
	return buf.String(), nil
}

// GenerateWithCSS creates an HTML resume with embedded CSS
func (g *HTMLGenerator) GenerateWithCSS(templateContent, cssContent string, resume *definition.Resume) (string, error) {
	g.logger.Info("Generating HTML resume with embedded CSS")

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	payload := g.buildTemplatePayload(resume, cssContent)
	if err := tmpl.Execute(&buf, payload); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	g.logger.Info("Successfully generated HTML resume with CSS")
	return buf.String(), nil
}

func (g *HTMLGenerator) buildTemplatePayload(resume *definition.Resume, css string) htmlTemplatePayload {
	payload := htmlTemplatePayload{
		Resume: resume,
		View:   g.buildHTMLView(resume),
	}

	if css != "" {
		payload.CSS = template.CSS(css)
	}

	return payload
}

func (g *HTMLGenerator) buildHTMLView(resume *definition.Resume) *htmlResumeView {
	view := &htmlResumeView{}

	view.Header = buildHeaderView(resume)
	if resume.Contact.Visibility.ShowSummary && strings.TrimSpace(resume.Contact.Summary) != "" {
		view.Summary = strings.TrimSpace(resume.Contact.Summary)
	}

	if skillsView := buildSkillsView(resume); skillsView != nil {
		view.Skills = skillsView
	}

	if experienceView := buildExperienceView(resume); experienceView != nil {
		view.Experience = experienceView
	}

	if projectsView := buildProjectsView(resume); projectsView != nil {
		view.Projects = projectsView
	}

	if educationView := buildEducationView(resume); educationView != nil {
		view.Education = educationView
	}

	if certificationsView := buildCertificationsView(resume); certificationsView != nil {
		view.Certifications = certificationsView
	}

	return view
}

func buildHeaderView(resume *definition.Resume) htmlHeaderView {
	header := htmlHeaderView{
		Name:  strings.TrimSpace(resume.Contact.Name),
		Title: strings.TrimSpace(resume.Contact.Title),
	}

	if resume.Contact.Visibility.ShowEmail && strings.TrimSpace(resume.Contact.Email) != "" {
		email := strings.TrimSpace(resume.Contact.Email)
		header.ContactItems = append(header.ContactItems, htmlInlineItem{
			Text: email,
			URL:  "mailto:" + email,
		})
	}

	if resume.Contact.Visibility.ShowPhone && strings.TrimSpace(resume.Contact.Phone) != "" {
		phone := strings.TrimSpace(resume.Contact.Phone)
		header.ContactItems = append(header.ContactItems, htmlInlineItem{
			Text: phone,
			URL:  "tel:" + sanitizePhone(phone),
		})
	}

	if resume.Contact.Visibility.ShowLocation && resume.Contact.Location != nil {
		if location := formatLocation(resume.Contact.Location); location != "" {
			header.ContactItems = append(header.ContactItems, htmlInlineItem{
				Text: location,
			})
		}
	}

	if strings.TrimSpace(resume.Contact.Website) != "" {
		website := strings.TrimSpace(resume.Contact.Website)
		header.ContactItems = append(header.ContactItems, htmlInlineItem{
			Text: website,
			URL:  website,
		})
	}

	if len(resume.Contact.Links) > 0 {
		links := make([]definition.Link, len(resume.Contact.Links))
		copy(links, resume.Contact.Links)
		sort.SliceStable(links, func(i, j int) bool {
			return links[i].Order < links[j].Order
		})

		for _, link := range links {
			text := strings.TrimSpace(link.Text)
			if text == "" {
				text = strings.TrimSpace(link.URL)
			}
			if text == "" {
				continue
			}
			header.ContactItems = append(header.ContactItems, htmlInlineItem{
				Text: text,
				URL:  strings.TrimSpace(link.URL),
			})
		}
	}

	if len(header.ContactItems) > 0 {
		header.HasContactRow = true
	}

	return header
}

func buildSkillsView(resume *definition.Resume) *htmlSkillsView {
	if len(resume.Skills.Categories) == 0 {
		return nil
	}

	view := &htmlSkillsView{
		Title: defaultString(resume.Skills.Title, "Skills"),
	}

	categories := make([]definition.SkillCategory, len(resume.Skills.Categories))
	copy(categories, resume.Skills.Categories)
	sort.SliceStable(categories, func(i, j int) bool {
		return categories[i].Order < categories[j].Order
	})

	for _, category := range categories {
		items := make([]definition.SkillItem, len(category.Items))
		copy(items, category.Items)
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].Order < items[j].Order
		})

		names := make([]string, 0, len(items))
		for _, item := range items {
			name := strings.TrimSpace(item.Name)
			if name != "" {
				names = append(names, name)
			}
		}

		if len(names) == 0 {
			continue
		}

		view.Categories = append(view.Categories, htmlSkillCategoryView{
			Name:    strings.TrimSpace(category.Name),
			Display: strings.Join(names, ", "),
		})
	}

	if len(view.Categories) == 0 {
		return nil
	}

	return view
}

func buildExperienceView(resume *definition.Resume) *htmlExperienceView {
	if len(resume.Experience.Positions) == 0 {
		return nil
	}

	view := &htmlExperienceView{
		Title: defaultString(resume.Experience.Title, "Experience"),
	}

	positions := make([]definition.Experience, len(resume.Experience.Positions))
	copy(positions, resume.Experience.Positions)
	sort.SliceStable(positions, func(i, j int) bool {
		return positions[i].Order < positions[j].Order
	})

	for _, pos := range positions {
		entry := htmlExperienceEntry{
			Title:      strings.TrimSpace(pos.Title),
			Company:    strings.TrimSpace(pos.Company),
			Location:   formatLocation(pos.Location),
			DateRange:  formatDateRange(pos.Dates.Start, pos.Dates.End, pos.Dates.Current),
			Highlights: filterStrings(pos.Description),
		}

		if len(entry.Highlights) == 0 {
			entry.Highlights = filterStrings(pos.Highlights)
		}

		if len(entry.Highlights) == 0 && len(pos.Achievements) > 0 {
			for _, achievement := range pos.Achievements {
				line := strings.TrimSpace(achievement.Description)
				if line == "" {
					line = strings.TrimSpace(achievement.Title)
				}
				if line != "" {
					entry.Highlights = append(entry.Highlights, line)
				}
			}
		}

		view.Positions = append(view.Positions, entry)
	}

	if len(view.Positions) == 0 {
		return nil
	}

	return view
}

func buildProjectsView(resume *definition.Resume) *htmlProjectsView {
	if len(resume.Projects.Projects) == 0 {
		return nil
	}

	view := &htmlProjectsView{
		Title: defaultString(resume.Projects.Title, "Projects"),
	}

	projects := make([]definition.Project, len(resume.Projects.Projects))
	copy(projects, resume.Projects.Projects)
	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Order < projects[j].Order
	})

	for _, project := range projects {
		entry := htmlProjectEntry{
			Name:         strings.TrimSpace(project.Name),
			Category:     strings.TrimSpace(project.Category),
			DateRange:    formatOptionalDateRange(project.Dates),
			Description:  filterStrings(project.Description),
			Technologies: formatList(project.Technologies),
		}

		if len(entry.Description) == 0 && len(project.Achievements) > 0 {
			for _, achievement := range project.Achievements {
				line := strings.TrimSpace(achievement.Description)
				if line == "" {
					line = strings.TrimSpace(achievement.Title)
				}
				if line != "" {
					entry.Description = append(entry.Description, line)
				}
			}
		}

		if len(project.Links) > 0 {
			links := make([]definition.Link, len(project.Links))
			copy(links, project.Links)
			sort.SliceStable(links, func(i, j int) bool {
				return links[i].Order < links[j].Order
			})
			for _, link := range links {
				text := strings.TrimSpace(link.Text)
				if text == "" {
					text = strings.TrimSpace(link.URL)
				}
				if text == "" {
					continue
				}
				entry.Links = append(entry.Links, htmlInlineItem{
					Text: text,
					URL:  strings.TrimSpace(link.URL),
				})
			}
		}

		view.Entries = append(view.Entries, entry)
	}

	if len(view.Entries) == 0 {
		return nil
	}

	return view
}

func buildEducationView(resume *definition.Resume) *htmlEducationView {
	if len(resume.Education.Institutions) == 0 {
		return nil
	}

	view := &htmlEducationView{
		Title: defaultString(resume.Education.Title, "Education"),
	}

	institutions := make([]definition.Education, len(resume.Education.Institutions))
	copy(institutions, resume.Education.Institutions)
	sort.SliceStable(institutions, func(i, j int) bool {
		return institutions[i].Order < institutions[j].Order
	})

	for _, institution := range institutions {
		entry := htmlEducationEntry{
			Institution: strings.TrimSpace(institution.Institution),
			Degree:      strings.TrimSpace(institution.Degree),
			Field:       strings.TrimSpace(institution.Field),
			Location:    formatLocation(institution.Location),
			DateRange:   formatDateRange(institution.Dates.Start, institution.Dates.End, institution.Dates.Current),
			GPA:         formatGPAValue(institution.GPA, institution.MaxGPA),
			Honors:      formatList(institution.Honors),
		}

		if len(institution.Description) > 0 {
			for _, pair := range institution.Description {
				label := strings.TrimSpace(pair.Category)
				value := strings.TrimSpace(pair.Value)
				switch {
				case label != "" && value != "":
					entry.Details = append(entry.Details, fmt.Sprintf("%s: %s", label, value))
				case value != "":
					entry.Details = append(entry.Details, value)
				case label != "":
					entry.Details = append(entry.Details, label)
				}
			}
		}

		if len(institution.Coursework) > 0 {
			courses := make([]string, 0, len(institution.Coursework))
			for _, course := range institution.Coursework {
				name := strings.TrimSpace(course.Name)
				if name == "" && course.Code != "" {
					name = strings.TrimSpace(course.Code)
				}
				if name != "" {
					courses = append(courses, name)
				}
			}
			if len(courses) > 0 {
				entry.Details = append(entry.Details, "Coursework: "+strings.Join(courses, ", "))
			}
		}

		if institution.Thesis != nil {
			title := strings.TrimSpace(institution.Thesis.Title)
			if title != "" {
				entry.Details = append(entry.Details, "Thesis: "+title)
			}
		}

		view.Schools = append(view.Schools, entry)
	}

	if len(view.Schools) == 0 {
		return nil
	}

	return view
}

func buildCertificationsView(resume *definition.Resume) *htmlCertificationsView {
	if len(resume.Certifications.Certifications) == 0 {
		return nil
	}

	view := &htmlCertificationsView{
		Title: defaultString(resume.Certifications.Title, "Certifications"),
	}

	certs := make([]definition.Certification, len(resume.Certifications.Certifications))
	copy(certs, resume.Certifications.Certifications)
	sort.SliceStable(certs, func(i, j int) bool {
		return certs[i].Order < certs[j].Order
	})

	for _, cert := range certs {
		entry := htmlCertificationEntry{
			Name:            strings.TrimSpace(cert.Name),
			Issuer:          strings.TrimSpace(cert.Issuer),
			DateRange:       formatCertificationDate(cert.IssueDate, cert.ExpirationDate),
			CredentialID:    strings.TrimSpace(cert.CredentialID),
			VerificationURL: strings.TrimSpace(cert.VerificationURL),
		}

		view.Certifications = append(view.Certifications, entry)
	}

	if len(view.Certifications) == 0 {
		return nil
	}

	return view
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func sanitizePhone(phone string) string {
	var b strings.Builder
	for _, r := range phone {
		if (r >= '0' && r <= '9') || r == '+' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func formatLocation(loc *definition.Location) string {
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
	if strings.TrimSpace(loc.Country) != "" && !containsIgnoreCase(parts, loc.Country) {
		parts = append(parts, strings.TrimSpace(loc.Country))
	}

	parts = filterStrings(parts)
	return strings.Join(parts, ", ")
}

func containsIgnoreCase(list []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range list {
		if strings.ToLower(strings.TrimSpace(item)) == value {
			return true
		}
	}
	return false
}

func formatDateRange(start time.Time, end *time.Time, current bool) string {
	if start.IsZero() && (end == nil || (end != nil && end.IsZero())) {
		return ""
	}

	startStr := formatMonthYear(start)
	var endStr string

	switch {
	case current:
		endStr = "Present"
	case end == nil || end.IsZero():
		endStr = "Present"
	default:
		endStr = formatMonthYear(*end)
	}

	if startStr == "" {
		return endStr
	}
	if endStr == "" || startStr == endStr {
		return startStr
	}
	return fmt.Sprintf("%s – %s", startStr, endStr)
}

func formatOptionalDateRange(dr *definition.DateRange) string {
	if dr == nil {
		return ""
	}
	return formatDateRange(dr.Start, dr.End, dr.Current)
}

func formatCertificationDate(issue, expiration *time.Time) string {
	if (issue == nil || issue.IsZero()) && (expiration == nil || expiration.IsZero()) {
		return ""
	}

	issueStr := ""
	if issue != nil && !issue.IsZero() {
		issueStr = formatMonthYear(*issue)
	}

	expStr := ""
	if expiration != nil && !expiration.IsZero() {
		expStr = formatMonthYear(*expiration)
	}

	if issueStr == "" {
		return expStr
	}
	if expStr == "" {
		return issueStr
	}
	return fmt.Sprintf("%s – %s", issueStr, expStr)
}

func formatMonthYear(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("Jan 2006")
}

func filterStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, val := range values {
		if trimmed := strings.TrimSpace(val); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func formatList(values []string) string {
	filtered := filterStrings(values)
	return strings.Join(filtered, ", ")
}

func formatGPAValue(gpa, max string) string {
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

// GenerateStandalone creates a complete HTML document with all dependencies
func (g *HTMLGenerator) GenerateStandalone(templateContent, cssContent string, resume *definition.Resume) (string, error) {
	htmlContent, err := g.GenerateWithCSS(templateContent, cssContent, resume)
	if err != nil {
		return "", err
	}

	// Ensure it's a complete HTML document
	if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
		standaloneTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Resume.Contact.Name}} - Resume</title>
    <style>
        {{.CSS}}
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>`

		data := map[string]interface{}{
			"Resume":  resume,
			"CSS":     template.CSS(cssContent),
			"Content": template.HTML(htmlContent),
		}

		tmpl, err := template.New("standalone").Funcs(g.funcs).Parse(standaloneTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to parse standalone template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute standalone template: %w", err)
		}

		htmlContent = buf.String()
	}

	return htmlContent, nil
}

// HTMLGeneratorOptions provides configuration for HTML generation
type HTMLGeneratorOptions struct {
	Theme            string
	IncludeCSS       bool
	Standalone       bool
	ResponsiveDesign bool
	PrintOptimized   bool
	FontAwesome      bool
	CustomFonts      []string
	ColorScheme      string
}

// AdvancedHTMLGenerator provides advanced HTML generation capabilities
type AdvancedHTMLGenerator struct {
	logger  *zap.SugaredLogger
	options HTMLGeneratorOptions
	funcs   template.FuncMap
}

// NewAdvancedHTMLGenerator creates a new advanced HTML generator
func NewAdvancedHTMLGenerator(logger *zap.SugaredLogger, options HTMLGeneratorOptions) *AdvancedHTMLGenerator {
	generator := &AdvancedHTMLGenerator{
		logger:  logger,
		options: options,
	}
	generator.setupAdvancedTemplateFunctions()
	return generator
}

// setupAdvancedTemplateFunctions initializes advanced template helper functions
func (g *AdvancedHTMLGenerator) setupAdvancedTemplateFunctions() {
	g.funcs = template.FuncMap{
		// Include all basic functions
		"formatDate": func(t time.Time) string {
			return t.Format("January 2006")
		},
		"formatDateShort": func(t time.Time) string {
			return t.Format("Jan 2006")
		},
		"formatDateRange": func(start time.Time, end *time.Time) string {
			startStr := start.Format("Jan 2006")
			if end == nil {
				return startStr + " - Present"
			}
			endStr := end.Format("Jan 2006")
			if startStr == endStr {
				return startStr
			}
			return startStr + " - " + endStr
		},
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},

		// Advanced functions
		"getThemeClass": func() string {
			return g.options.Theme
		},
		"shouldIncludeCSS": func() bool {
			return g.options.IncludeCSS
		},
		"isStandalone": func() bool {
			return g.options.Standalone
		},
		"isResponsive": func() bool {
			return g.options.ResponsiveDesign
		},
		"isPrintOptimized": func() bool {
			return g.options.PrintOptimized
		},
		"shouldIncludeFontAwesome": func() bool {
			return g.options.FontAwesome
		},
		"getCustomFonts": func() []string {
			return g.options.CustomFonts
		},
		"getColorScheme": func() string {
			return g.options.ColorScheme
		},

		// Utility functions
		"generateID": func(prefix string) string {
			return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
		},
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"wordCount": func(s string) int {
			return len(strings.Fields(s))
		},
		"characterCount": func(s string) int {
			return len(s)
		},
	}
}

// Generate creates an advanced HTML resume
func (g *AdvancedHTMLGenerator) Generate(templateContent string, resume *definition.Resume) (string, error) {
	g.logger.Infof("Generating advanced HTML resume with theme: %s", g.options.Theme)

	// Create template data with options
	data := map[string]interface{}{
		"Resume":  resume,
		"Options": g.options,
	}

	// Parse the template
	tmpl, err := template.New("resume").Funcs(g.funcs).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse advanced HTML template: %w", err)
	}

	// Execute the template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute advanced HTML template: %w", err)
	}

	g.logger.Info("Successfully generated advanced HTML resume")
	return buf.String(), nil
}
