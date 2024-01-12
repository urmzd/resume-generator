package generators

import (
	"fmt"
	"strings"

	"github.com/urmzd/generate-resumes/pkg/template"
)

type ProfessionalResumeGenerator struct {
	code []string
}

// Ensure ProfessionalResumeGenerator implements template.ResumeGenerator
var _ template.ResumeGenerator = &ProfessionalResumeGenerator{}

func (gen *ProfessionalResumeGenerator) StartResume(contact *template.Contact) {
	// Starting LaTeX commands for a professional resume
	header := []string{
		`\documentclass[letterpaper,11pt]{article}`,
		`\usepackage[utf8]{inputenc}`,
		`\usepackage{geometry}`,
		`\usepackage{enumitem}`,
		`\usepackage{fancyhdr}`,
		`\usepackage{hyperref}`,
		`\geometry{left=0.75in, top=0.6in, right=0.75in, bottom=0.6in}`,
		`\hypersetup{colorlinks=true, linkcolor=blue, filecolor=magenta, urlcolor=cyan}`,
		`\pagestyle{fancy}`,
		`\fancyhf{}`,
		fmt.Sprintf(`\rhead{%s | \href{mailto:%s}{%s}}`, contact.Name, contact.Email, contact.Email),
		`\begin{document}`,
		fmt.Sprintf(`\begin{center}\textbf{\Huge %s}\\`, contact.Name),
		fmt.Sprintf(`\href{mailto:%s}{%s} | %s\\`, contact.Email, contact.Email, contact.Phone),
		`\end{center}`,
	}
	gen.write(header...)
}

func (gen *ProfessionalResumeGenerator) AddEducation(education *[]template.Education) {
	gen.write(`\section*{Education}`)
	for _, school := range *education {
		details := fmt.Sprintf(`\textbf{%s}, %s, %s - %s`, school.School, school.Degree,
			school.Dates.Start.Format("Jan 2006"), school.Dates.End.Format("Jan 2006"))
		gen.write(details)
	}
}

func (gen *ProfessionalResumeGenerator) AddExperiences(experience *[]template.Experience) {
	gen.write(`\section*{Experience}`)
	for _, xp := range *experience {
		var dateRange string
		if xp.Dates.End != nil {
			dateRange = fmt.Sprintf("%s - %s", xp.Dates.Start.Format("Jan 2006"), xp.Dates.End.Format("Jan 2006"))
		} else {
			dateRange = fmt.Sprintf("%s - Present", xp.Dates.Start.Format("Jan 2006"))
		}

		details := fmt.Sprintf(`\textbf{%s}, %s (%s)`, xp.Company, xp.Title, dateRange)
		gen.write(details)
	}
}

func (gen *ProfessionalResumeGenerator) AddSkills(skills *[]template.Detail) {
	gen.write(`\section*{Skills}`)
	var skillList []string
	for _, skill := range *skills {
		skillList = append(skillList, skill.Value)
	}
	gen.write(strings.Join(skillList, ", "))
}

func (gen *ProfessionalResumeGenerator) EndResume() string {
	gen.write(`\end{document}`)
	return strings.Join(gen.code, "\n")
}

func (gen *ProfessionalResumeGenerator) write(strs ...string) {
	for _, str := range strs {
		gen.code = append(gen.code, strings.TrimSpace(str))
	}
}

func NewProfessionalResumeGenerator() template.ResumeGenerator {
	return &ProfessionalResumeGenerator{}
}

func (gen *ProfessionalResumeGenerator) AddProjects(projects *[]template.Project) {
	gen.write(`\section*{Projects}`)
	for _, project := range *projects {
		projectDetails := fmt.Sprintf(`\textbf{%s}, %s: %s`, project.Name, project.Language, strings.Join(project.Details, "; "))
		if project.Link.Text != "" {
			projectDetails += fmt.Sprintf(` [\href{%s}{Link}]`, project.Link.Ref)
		}
		gen.write(projectDetails)
	}
}
