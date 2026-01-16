package generators

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

// LaTeXGenerator renders LaTeX templates using engine-specific formatting helpers.
type LaTeXGenerator struct {
	logger    *zap.SugaredLogger
	formatter Formatter
}

// NewLaTeXGenerator creates a new LaTeX generator wired with the LaTeX formatter.
func NewLaTeXGenerator(logger *zap.SugaredLogger) *LaTeXGenerator {
	return &LaTeXGenerator{
		logger:    logger,
		formatter: newLaTeXFormatter(),
	}
}

// Generate renders a LaTeX template with sanitised resume data.
func (g *LaTeXGenerator) Generate(templateContent string, resume *definition.Resume) (string, error) {
	g.logger.Info("Rendering LaTeX template")

	tmpl, err := template.New("latex").Funcs(g.formatter.TemplateFuncs()).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse LaTeX template: %w", err)
	}

	payload := g.buildTemplateData(resume)

	var output strings.Builder
	if err := tmpl.Execute(&output, payload); err != nil {
		return "", fmt.Errorf("failed to execute LaTeX template: %w", err)
	}

	g.logger.Info("Successfully rendered LaTeX template")
	return output.String(), nil
}

type latexTemplateData struct {
	Contact    latexContact
	Experience latexExperienceList
	Education  latexEducationList
	Skills     latexSkills
	Projects   latexProjectList
}

type latexContact struct {
	Name     string
	Title    string
	Email    string
	Phone    string
	Location *definition.Location
	Links    []definition.Link
	Summary  string
}

type latexExperienceList struct {
	Title     string
	Positions []latexExperience
}

type latexExperience struct {
	Title      string
	Company    string
	Dates      interface{}
	Location   *definition.Location
	Highlights []string
}

type latexEducationList struct {
	Title        string
	Institutions []latexEducation
}

type latexEducation struct {
	Institution string
	Degree      string
	Dates       interface{}
	Location    *definition.Location
	GPA         string
	Honors      []string
	Description []definition.CategoryValuePair
	Thesis      *definition.Thesis
}

type latexSkills struct {
	Title      string
	Categories []latexSkillCategory
}

type latexSkillCategory struct {
	Name  string
	Items []definition.SkillItem
}

type latexProjectList struct {
	Title    string
	Projects []latexProject
}

type latexProject struct {
	Name         string
	Technologies []string
	URL          string
	Highlights   []string
}

func (g *LaTeXGenerator) buildTemplateData(resume *definition.Resume) latexTemplateData {
	return latexTemplateData{
		Contact:    g.sanitiseContact(resume.Contact),
		Experience: g.sanitiseExperience(resume.Experience),
		Education:  g.sanitiseEducation(resume.Education),
		Skills:     g.sanitiseSkills(resume.Skills),
		Projects:   g.sanitiseProjects(resume.Projects),
	}
}

func (g *LaTeXGenerator) sanitiseContact(contact definition.Contact) latexContact {
	sanitised := latexContact{
		Name:     g.formatter.EscapeText(strings.TrimSpace(contact.Name)),
		Title:    g.formatter.EscapeText(strings.TrimSpace(contact.Title)),
		Email:    g.formatter.EscapeText(strings.TrimSpace(contact.Email)),
		Phone:    g.formatter.EscapeText(strings.TrimSpace(contact.Phone)),
		Location: contact.Location,
		Summary:  g.formatter.EscapeText(strings.TrimSpace(contact.Summary)),
	}

	sanitisedLinks := make([]definition.Link, 0, len(contact.Links))
	for _, link := range contact.Links {
		url := strings.TrimSpace(link.URL)
		if url == "" {
			continue
		}
		text := strings.TrimSpace(link.Text)
		sanitisedLinks = append(sanitisedLinks, definition.Link{
			Text:        text,
			URL:         url,
			Type:        link.Type,
			Icon:        link.Icon,
			Description: link.Description,
		})
	}

	sanitised.Links = sanitisedLinks
	return sanitised
}

func (g *LaTeXGenerator) sanitiseExperience(list definition.ExperienceList) latexExperienceList {
	sanitised := latexExperienceList{
		Title: g.formatter.EscapeText(strings.TrimSpace(list.Title)),
	}

	for _, position := range list.Positions {
		entry := latexExperience{
			Title:      g.formatter.EscapeText(strings.TrimSpace(position.Title)),
			Company:    g.formatter.EscapeText(strings.TrimSpace(position.Company)),
			Dates:      position.Dates,
			Location:   position.Location,
			Highlights: g.collectExperienceHighlights(position),
		}
		sanitised.Positions = append(sanitised.Positions, entry)
	}

	return sanitised
}

func (g *LaTeXGenerator) sanitiseEducation(list definition.EducationList) latexEducationList {
	sanitised := latexEducationList{
		Title: g.formatter.EscapeText(strings.TrimSpace(list.Title)),
	}

	for _, institution := range list.Institutions {
		entry := latexEducation{
			Institution: g.formatter.EscapeText(strings.TrimSpace(institution.Institution)),
			Degree:      g.formatter.EscapeText(strings.TrimSpace(institution.Degree)),
			Dates:       institution.Dates,
			Location:    institution.Location,
			GPA:         g.formatter.FormatGPA(institution.GPA, institution.MaxGPA),
			Honors:      institution.Honors,
			Description: institution.Description,
			Thesis:      institution.Thesis,
		}
		if institution.Thesis != nil && institution.Thesis.Url != "" {
			entry.Thesis.Url = g.formatter.EscapeText(strings.TrimSpace(institution.Thesis.Url))
		}
		sanitised.Institutions = append(sanitised.Institutions, entry)
	}

	return sanitised
}

func (g *LaTeXGenerator) sanitiseSkills(skills definition.Skills) latexSkills {
	sanitised := latexSkills{
		Title: g.formatter.EscapeText(strings.TrimSpace(skills.Title)),
	}

	for _, category := range skills.Categories {
		sanitisedItems := make([]definition.SkillItem, 0, len(category.Items))
		for _, item := range category.Items {
			name := strings.TrimSpace(item.Name)
			if name == "" {
				continue
			}
			sanitisedItems = append(sanitisedItems, definition.SkillItem{
				Name:              name,
				Level:             strings.TrimSpace(item.Level),
				Years:             item.Years,
				YearsOfExperience: item.YearsOfExperience,
				Certification:     strings.TrimSpace(item.Certification),
				Keywords:          filterStrings(item.Keywords),
			})
		}

		if len(sanitisedItems) == 0 {
			continue
		}

		sanitised.Categories = append(sanitised.Categories, latexSkillCategory{
			Name:  strings.TrimSpace(category.Name),
			Items: sanitisedItems,
		})
	}

	return sanitised
}

func (g *LaTeXGenerator) sanitiseProjects(projects definition.ProjectList) latexProjectList {
	sanitised := latexProjectList{
		Title: g.formatter.EscapeText(strings.TrimSpace(projects.Title)),
	}

	for _, project := range projects.Projects {
		entry := latexProject{
			Name:         g.formatter.EscapeText(strings.TrimSpace(project.Name)),
			Technologies: g.sanitiseTechnologies(project.Technologies),
			Highlights:   g.collectProjectHighlights(project),
		}

		if link := selectPrimaryLink(project.Links); link != nil {
			entry.URL = g.formatter.EscapeText(strings.TrimSpace(link.URL))
		}

		sanitised.Projects = append(sanitised.Projects, entry)
	}

	return sanitised
}

func (g *LaTeXGenerator) sanitiseTechnologies(values []string) []string {
	filtered := filterStrings(values)
	if len(filtered) == 0 {
		return nil
	}
	result := make([]string, len(filtered))
	for i, value := range filtered {
		result[i] = g.formatter.EscapeText(value)
	}
	return result
}

func (g *LaTeXGenerator) collectExperienceHighlights(exp definition.Experience) []string {
	highlights := filterStrings(exp.Description)
	if len(highlights) == 0 {
		highlights = filterStrings(exp.Highlights)
	}
	if len(highlights) == 0 {
		for _, achievement := range exp.Achievements {
			line := strings.TrimSpace(achievement.Description)
			if line == "" {
				line = strings.TrimSpace(achievement.Title)
			}
			if line != "" {
				highlights = append(highlights, line)
			}
		}
	}
	return highlights
}

func (g *LaTeXGenerator) collectProjectHighlights(project definition.Project) []string {
	highlights := filterStrings(project.Description)
	if len(highlights) == 0 {
		for _, achievement := range project.Achievements {
			line := strings.TrimSpace(achievement.Description)
			if line == "" {
				line = strings.TrimSpace(achievement.Title)
			}
			if line != "" {
				highlights = append(highlights, line)
			}
		}
	}
	return highlights
}

func selectPrimaryLink(links []definition.Link) *definition.Link {
	for _, link := range links {
		url := strings.TrimSpace(link.URL)
		if url == "" {
			continue
		}
		text := strings.TrimSpace(link.Text)
		if text == "" {
			text = url
		}
		return &definition.Link{
			Text: text,
			URL:  url,
			Type: link.Type,
			Icon: link.Icon,
		}
	}
	return nil
}
