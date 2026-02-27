package generators

import (
	"bytes"
	"strings"

	"github.com/fumiama/go-docx"
	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// DOCXGenerator generates Word documents from resume data.
type DOCXGenerator struct {
	logger    *zap.SugaredLogger
	formatter *docxFormatter
}

// NewDOCXGenerator creates a new DOCX generator.
func NewDOCXGenerator(logger *zap.SugaredLogger) *DOCXGenerator {
	return &DOCXGenerator{
		logger:    logger,
		formatter: newDocxFormatter(),
	}
}

// Generate creates a DOCX document from the resume and returns it as bytes.
func (g *DOCXGenerator) Generate(r *resume.Resume) ([]byte, error) {
	doc := docx.New().WithDefaultTheme()

	g.addHeader(doc, r.Contact)

	if r.Layout != nil && len(r.Layout.Sections) > 0 {
		for _, section := range r.Layout.Sections {
			switch section {
			case "summary":
				g.addSummary(doc, r.Summary)
			case "certifications":
				if r.Certifications != nil {
					g.addCertifications(doc, *r.Certifications)
				}
			case "education":
				g.addEducation(doc, r.Education)
			case "skills":
				g.addSkills(doc, r.Skills)
			case "experience":
				g.addExperience(doc, r.Experience)
			case "projects":
				if r.Projects != nil {
					g.addProjects(doc, *r.Projects)
				}
			case "languages":
				if r.Languages != nil {
					g.addLanguages(doc, *r.Languages)
				}
			}
		}
	} else {
		g.addSummary(doc, r.Summary)
		if r.Certifications != nil {
			g.addCertifications(doc, *r.Certifications)
		}
		g.addEducation(doc, r.Education)
		g.addSkills(doc, r.Skills)
		g.addExperience(doc, r.Experience)
		if r.Projects != nil {
			g.addProjects(doc, *r.Projects)
		}
		if r.Languages != nil {
			g.addLanguages(doc, *r.Languages)
		}
	}

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// addHeader adds the name and contact information.
func (g *DOCXGenerator) addHeader(doc *docx.Docx, contact resume.Contact) {
	// Name - large, bold, centered
	namePara := doc.AddParagraph().Justification("center")
	namePara.AddText(strings.ToUpper(contact.Name)).Bold().Size("36")

	// Contact info line
	var contactParts []string
	if contact.Location != nil {
		loc := g.formatter.FormatLocation(contact.Location)
		if loc != "" {
			contactParts = append(contactParts, loc)
		}
	}
	if contact.Email != "" {
		contactParts = append(contactParts, contact.Email)
	}
	if contact.Phone != "" {
		contactParts = append(contactParts, contact.Phone)
	}
	if contact.Credentials != "" {
		contactParts = append(contactParts, contact.Credentials)
	}
	for _, link := range contact.Links {
		if link.URI != "" {
			contactParts = append(contactParts, link.URI)
		}
	}

	if len(contactParts) > 0 {
		contactPara := doc.AddParagraph().Justification("center")
		contactPara.AddText(strings.Join(contactParts, " | ")).Size("22")
	}

	// Add spacing after header
	doc.AddParagraph()
}

// addSummary adds the professional summary section.
func (g *DOCXGenerator) addSummary(doc *docx.Docx, summary string) {
	if summary == "" {
		return
	}
	g.addSectionHeader(doc, "Professional Summary")
	para := doc.AddParagraph()
	para.AddText(summary).Size("22")
	doc.AddParagraph()
}

// addCertifications adds the certifications section.
func (g *DOCXGenerator) addCertifications(doc *docx.Docx, certs resume.Certifications) {
	if len(certs.Items) == 0 {
		return
	}
	title := certs.Title
	if title == "" {
		title = "Certifications"
	}
	g.addSectionHeader(doc, title)
	for _, cert := range certs.Items {
		bulletPara := doc.AddParagraph()
		line := cert.Name
		if cert.Issuer != "" {
			line += " — " + cert.Issuer
		}
		if cert.Notes != "" {
			line += " (" + cert.Notes + ")"
		}
		bulletPara.AddText("• " + line).Size("22")
	}
	doc.AddParagraph()
}

// addLanguages adds the languages section.
func (g *DOCXGenerator) addLanguages(doc *docx.Docx, languages resume.LanguageList) {
	if len(languages.Languages) == 0 {
		return
	}
	title := languages.Title
	if title == "" {
		title = "Languages"
	}
	g.addSectionHeader(doc, title)
	for _, lang := range languages.Languages {
		bulletPara := doc.AddParagraph()
		line := lang.Name
		if lang.Proficiency != "" {
			line += " — " + lang.Proficiency
		}
		bulletPara.AddText("• " + line).Size("22")
	}
	doc.AddParagraph()
}

// addSectionHeader adds a section title with underline styling.
func (g *DOCXGenerator) addSectionHeader(doc *docx.Docx, title string) {
	para := doc.AddParagraph()
	para.AddText(strings.ToUpper(title)).Bold().Size("24")
	doc.AddParagraph() // Add spacing after section header
}

// addEducation adds the education section.
func (g *DOCXGenerator) addEducation(doc *docx.Docx, education resume.EducationList) {
	if len(education.Institutions) == 0 {
		return
	}

	title := education.Title
	if title == "" {
		title = "Education"
	}
	g.addSectionHeader(doc, title)

	for _, inst := range education.Institutions {
		// Institution and dates on same logical line
		entryPara := doc.AddParagraph()

		// Build institution text
		instText := inst.Institution
		if inst.Degree.Name != "" {
			instText += ", " + inst.Degree.Name
		}

		dates := g.formatter.FormatDateRange(inst.Dates)
		if dates != "" {
			instText += " — " + dates
		}
		entryPara.AddText(instText).Bold().Size("22")

		// GPA
		if inst.GPA != nil {
			detailPara := doc.AddParagraph()
			gpaStr := g.formatter.FormatGPA(inst.GPA.GPA, inst.GPA.MaxGPA)
			detailPara.AddText("GPA: " + gpaStr).Italic().Size("22")
		}

		descriptions := filterStrings(inst.Degree.Descriptions)
		for _, description := range descriptions {
			descPara := doc.AddParagraph()
			descPara.AddText("• " + description).Size("22")
		}

		// Thesis
		if inst.Thesis != nil && inst.Thesis.Title != "" {
			title := strings.TrimSpace(inst.Thesis.Title)
			if title != "" {
				line := "Thesis: " + title
				if url := strings.TrimSpace(inst.Thesis.Link.URI); url != "" {
					line += " (" + url + ")"
				}
				descs := make([]string, 0, len(inst.Thesis.Highlights))
				for _, desc := range inst.Thesis.Highlights {
					if d := strings.TrimSpace(desc); d != "" {
						descs = append(descs, d)
					}
				}
				if len(descs) > 0 {
					line += " — " + strings.Join(descs, "; ")
				}
				thesisPara := doc.AddParagraph()
				thesisPara.AddText(line).Size("22")
			}
		}
	}

	doc.AddParagraph() // spacing
}

// addSkills adds the skills section.
func (g *DOCXGenerator) addSkills(doc *docx.Docx, skills resume.Skills) {
	if len(skills.Categories) == 0 {
		return
	}

	title := skills.Title
	if title == "" {
		title = "Skills"
	}
	g.addSectionHeader(doc, title)

	for _, category := range skills.Categories {
		skillPara := doc.AddParagraph()
		skillPara.AddText("• ").Size("22")
		skillPara.AddText(category.Category + ": ").Bold().Size("22")

		skillNames := g.formatter.SkillNames(category.Items)
		skillPara.AddText(strings.Join(skillNames, ", ")).Size("22")
	}

	doc.AddParagraph() // spacing
}

// addExperience adds the experience section.
func (g *DOCXGenerator) addExperience(doc *docx.Docx, experience resume.ExperienceList) {
	if len(experience.Positions) == 0 {
		return
	}

	title := experience.Title
	if title == "" {
		title = "Experience"
	}
	g.addSectionHeader(doc, title)

	for _, pos := range experience.Positions {
		// Title and dates
		headerPara := doc.AddParagraph()
		dates := g.formatter.FormatDateRange(pos.Dates)
		titleLine := pos.Title
		if dates != "" {
			titleLine += " — " + dates
		}
		headerPara.AddText(titleLine).Bold().Size("22")

		// Company and location
		var companyParts []string
		if pos.Company != "" {
			companyParts = append(companyParts, pos.Company)
		}
		if pos.Location != nil {
			loc := g.formatter.FormatLocation(pos.Location)
			if loc != "" {
				companyParts = append(companyParts, loc)
			}
		}
		if len(companyParts) > 0 {
			companyPara := doc.AddParagraph()
			companyPara.AddText(strings.Join(companyParts, " | ")).Italic().Size("22")
		}

		// Highlights/bullets
		highlights := pos.Highlights
		for _, highlight := range highlights {
			bulletPara := doc.AddParagraph()
			bulletPara.AddText("• " + highlight).Size("22")
		}

		doc.AddParagraph() // spacing between positions
	}
}

// addProjects adds the projects section.
func (g *DOCXGenerator) addProjects(doc *docx.Docx, projects resume.ProjectList) {
	if len(projects.Projects) == 0 {
		return
	}

	title := projects.Title
	if title == "" {
		title = "Projects"
	}
	g.addSectionHeader(doc, title)

	for _, proj := range projects.Projects {
		// Project name
		headerPara := doc.AddParagraph()
		headerPara.AddText(proj.Name).Bold().Size("22")

		// Highlights bullets
		for _, desc := range proj.Highlights {
			bulletPara := doc.AddParagraph()
			bulletPara.AddText("• " + desc).Size("22")
		}

		// Link
		if url := strings.TrimSpace(proj.Link.URI); url != "" {
			linkPara := doc.AddParagraph()
			linkPara.AddText("  → " + url).Size("20")
		}
	}

	doc.AddParagraph() // spacing
}
