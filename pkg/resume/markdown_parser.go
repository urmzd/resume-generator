package resume

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// section represents the current parser state.
type section int

const (
	sectionNone section = iota
	sectionSummary
	sectionSkills
	sectionExperience
	sectionEducation
	sectionProjects
	sectionCertifications
	sectionLanguages
)

// Regex patterns used throughout the parser.
var (
	reH1         = regexp.MustCompile(`^#\s+(.+)$`)
	reH2         = regexp.MustCompile(`^##\s+(.+)$`)
	reH3         = regexp.MustCompile(`^###\s+(.+)$`)
	reBold       = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reItalic     = regexp.MustCompile(`^\*(.+)\*$`)
	reLink       = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
	reMailtoLink = regexp.MustCompile(`\[([^\]]*)\]\(mailto:([^)]+)\)`)
	reBullet     = regexp.MustCompile(`^[-*]\s+(.+)$`)
	reGPA        = regexp.MustCompile(`(?i)GPA:\s*([^\s/]+)\s*/\s*([^\s|]+)`)
	reDateRange  = regexp.MustCompile(`([A-Za-z]+\.?\s+\d{4})\s*[–—-]\s*([A-Za-z]+\.?\s+\d{4}|Present)`)
	reDateSingle = regexp.MustCompile(`^([A-Za-z]+\.?\s+\d{4})$`)
	rePhone      = regexp.MustCompile(`^[+]?[\d()\s.-]+$`)
	reBoldPrefix = regexp.MustCompile(`^\*\*(.+?):?\*\*\s*:?\s*(.*)$`)
	reDashSplit  = regexp.MustCompile(`\s+[—–-]\s+`)
	reThesisLine = regexp.MustCompile(`(?i)^\*\*Thesis:\*\*\s*(.+)$`)
)

// dateFormats lists the time layouts used when parsing month+year strings.
var dateFormats = []string{
	"January 2006",
	"Jan 2006",
	"Jan. 2006",
	"1/2006",
	"01/2006",
}

// parseMarkdown parses a Markdown resume (matching the modern-markdown template)
// into a Resume struct.
func parseMarkdown(data []byte) (*Resume, error) {
	lines := strings.Split(string(data), "\n")
	r := &Resume{}

	var cur section
	var summaryLines []string

	// Experience state
	var curExp *Experience
	// Education state
	var curEdu *Education
	var eduExpectMeta bool
	// Project state
	var curProj *Project

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Skip horizontal rules and empty lines at the top level
		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			continue
		}

		// H1: Contact name (only first H1)
		if m := reH1.FindStringSubmatch(trimmed); m != nil {
			if r.Contact.Name == "" {
				r.Contact.Name = strings.TrimSpace(m[1])
				// Next non-empty line is the contact info line
				for i+1 < len(lines) {
					i++
					next := strings.TrimSpace(lines[i])
					if next == "" || next == "---" {
						if next == "---" {
							break
						}
						continue
					}
					parseContactLine(next, r)
					break
				}
			}
			continue
		}

		// H2: Section header
		if m := reH2.FindStringSubmatch(trimmed); m != nil {
			// Flush any in-progress items
			flushExperience(curExp, r)
			curExp = nil
			flushEducation(curEdu, r)
			curEdu = nil
			eduExpectMeta = false
			flushProject(curProj, r)
			curProj = nil

			if cur == sectionSummary && len(summaryLines) > 0 {
				r.Summary = strings.TrimSpace(strings.Join(summaryLines, "\n"))
				summaryLines = nil
			}

			cur = classifySection(strings.TrimSpace(m[1]), r)
			continue
		}

		// H3: Sub-item header (experience position, education entry, project)
		if m := reH3.FindStringSubmatch(trimmed); m != nil {
			switch cur {
			case sectionExperience:
				flushExperience(curExp, r)
				curExp = &Experience{Title: strings.TrimSpace(m[1])}

			case sectionEducation:
				flushEducation(curEdu, r)
				curEdu = &Education{}
				parseEducationH3(strings.TrimSpace(m[1]), curEdu)
				eduExpectMeta = true

			case sectionProjects:
				flushProject(curProj, r)
				curProj = &Project{}
				parseProjectH3(strings.TrimSpace(m[1]), curProj)
			}
			continue
		}

		// Empty line
		if trimmed == "" {
			if cur == sectionSummary {
				summaryLines = append(summaryLines, "")
			}
			continue
		}

		// Section-specific parsing
		switch cur {
		case sectionSummary:
			summaryLines = append(summaryLines, trimmed)

		case sectionSkills:
			if b := reBullet.FindStringSubmatch(trimmed); b != nil {
				parseSkillLine(b[1], r)
			}

		case sectionExperience:
			if curExp != nil {
				parseExperienceLine(trimmed, curExp)
			}

		case sectionEducation:
			if curEdu != nil {
				if eduExpectMeta && !strings.HasPrefix(trimmed, "-") {
					parseEducationMeta(trimmed, curEdu)
					eduExpectMeta = false
				} else {
					parseEducationLine(trimmed, curEdu)
				}
			}

		case sectionProjects:
			if curProj != nil {
				parseProjectLine(trimmed, curProj)
			}

		case sectionCertifications:
			if b := reBullet.FindStringSubmatch(trimmed); b != nil {
				parseCertificationLine(b[1], r)
			}

		case sectionLanguages:
			if b := reBullet.FindStringSubmatch(trimmed); b != nil {
				parseLanguageLine(b[1], r)
			}
		}
	}

	// Flush remaining items
	if cur == sectionSummary && len(summaryLines) > 0 {
		r.Summary = strings.TrimSpace(strings.Join(summaryLines, "\n"))
	}
	flushExperience(curExp, r)
	flushEducation(curEdu, r)
	flushProject(curProj, r)

	if r.Contact.Name == "" {
		return nil, fmt.Errorf("markdown parse error: no H1 heading found for contact name")
	}

	return r, nil
}

// classifySection maps an H2 title to a section enum, storing custom titles.
func classifySection(title string, r *Resume) section {
	lower := strings.ToLower(title)

	switch {
	case lower == "summary" || lower == "profile" || lower == "about":
		return sectionSummary

	case strings.Contains(lower, "skill") || lower == "core skills" || lower == "technical skills":
		r.Skills.Title = title
		return sectionSkills

	case strings.Contains(lower, "experience") || lower == "work history" || lower == "employment":
		r.Experience.Title = title
		return sectionExperience

	case strings.Contains(lower, "education"):
		r.Education.Title = title
		return sectionEducation

	case strings.Contains(lower, "project"):
		if r.Projects == nil {
			r.Projects = &ProjectList{}
		}
		r.Projects.Title = title
		return sectionProjects

	case strings.Contains(lower, "certif"):
		if r.Certifications == nil {
			r.Certifications = &Certifications{}
		}
		r.Certifications.Title = title
		return sectionCertifications

	case strings.Contains(lower, "language"):
		if r.Languages == nil {
			r.Languages = &LanguageList{}
		}
		r.Languages.Title = title
		return sectionLanguages
	}

	return sectionNone
}

// parseContactLine parses the pipe-separated contact info line.
func parseContactLine(line string, r *Resume) {
	parts := strings.Split(line, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check for mailto link → email
		if m := reMailtoLink.FindStringSubmatch(part); m != nil {
			r.Contact.Email = strings.TrimSpace(m[2])
			continue
		}

		// Check for regular link → contact link
		if m := reLink.FindStringSubmatch(part); m != nil {
			label := strings.TrimSpace(m[1])
			uri := strings.TrimSpace(m[2])
			r.Contact.Links = append(r.Contact.Links, Link{URI: uri, Label: label})
			continue
		}

		// Check for phone number
		stripped := strings.ReplaceAll(part, " ", "")
		if rePhone.MatchString(stripped) && len(stripped) >= 7 {
			r.Contact.Phone = part
			continue
		}

		// Otherwise treat as location (comma-separated: city, state, country)
		parseLocation(part, r)
	}
}

// parseLocation parses a comma-separated location string into Contact.Location.
func parseLocation(s string, r *Resume) {
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return
	}
	loc := &Location{}
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	switch len(parts) {
	case 1:
		loc.City = parts[0]
	case 2:
		loc.City = parts[0]
		loc.State = parts[1]
	default:
		loc.City = parts[0]
		loc.State = parts[1]
		loc.Country = parts[2]
	}
	r.Contact.Location = loc
}

// parseSkillLine parses a skill bullet: **Category:** item1, item2
func parseSkillLine(content string, r *Resume) {
	m := reBoldPrefix.FindStringSubmatch(content)
	if m == nil {
		return
	}
	category := strings.TrimSpace(m[1])
	items := strings.Split(m[2], ",")
	var cleaned []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			cleaned = append(cleaned, item)
		}
	}
	r.Skills.Categories = append(r.Skills.Categories, SkillCategory{
		Category: category,
		Items:    cleaned,
	})
}

// parseExperienceLine handles lines within an experience entry.
func parseExperienceLine(line string, exp *Experience) {
	trimmed := strings.TrimSpace(line)

	// Bold company line: **Company** | DateRange | Location
	if strings.HasPrefix(trimmed, "**") {
		parts := strings.Split(trimmed, "|")
		for i, p := range parts {
			p = strings.TrimSpace(p)

			if i == 0 {
				// Extract company from bold
				if m := reBold.FindStringSubmatch(p); m != nil {
					exp.Company = strings.TrimSpace(m[1])
				}
				continue
			}

			// Try date range
			if dr := reDateRange.FindStringSubmatch(p); dr != nil {
				start := parseDate(dr[1])
				if !start.IsZero() {
					exp.Dates.Start = start
				}
				if !strings.EqualFold(dr[2], "Present") {
					end := parseDate(dr[2])
					if !end.IsZero() {
						exp.Dates.End = &end
					}
				}
				continue
			}

			// Try single date (start only)
			if sd := reDateSingle.FindStringSubmatch(strings.TrimSpace(p)); sd != nil {
				start := parseDate(sd[1])
				if !start.IsZero() {
					exp.Dates.Start = start
				}
				continue
			}

			// Otherwise location
			parseExpLocation(p, exp)
		}
		return
	}

	// Italic technologies line: *tech1, tech2*
	if m := reItalic.FindStringSubmatch(trimmed); m != nil {
		techs := strings.Split(m[1], ",")
		for _, t := range techs {
			t = strings.TrimSpace(t)
			if t != "" {
				exp.Technologies = append(exp.Technologies, t)
			}
		}
		return
	}

	// Bullet: highlight
	if b := reBullet.FindStringSubmatch(trimmed); b != nil {
		exp.Highlights = append(exp.Highlights, strings.TrimSpace(b[1]))
		return
	}
}

// parseExpLocation parses a comma-separated location for an experience entry.
func parseExpLocation(s string, exp *Experience) {
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return
	}
	loc := &Location{}
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	switch len(parts) {
	case 1:
		loc.City = parts[0]
	case 2:
		loc.City = parts[0]
		loc.State = parts[1]
	default:
		loc.City = parts[0]
		loc.State = parts[1]
		loc.Country = parts[2]
	}
	exp.Location = loc
}

// parseEducationH3 handles: ### Institution — Degree
func parseEducationH3(content string, edu *Education) {
	parts := reDashSplit.Split(content, 2)
	edu.Institution = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		edu.Degree.Name = strings.TrimSpace(parts[1])
	}
}

// parseEducationMeta handles the metadata line after the H3:
// DateRange | Location | GPA: X / Y
func parseEducationMeta(line string, edu *Education) {
	parts := strings.Split(line, "|")
	for _, p := range parts {
		p = strings.TrimSpace(p)

		// GPA
		if gm := reGPA.FindStringSubmatch(p); gm != nil {
			edu.GPA = &GPA{
				GPA:    strings.TrimSpace(gm[1]),
				MaxGPA: strings.TrimSpace(gm[2]),
			}
			continue
		}

		// Date range
		if dr := reDateRange.FindStringSubmatch(p); dr != nil {
			start := parseDate(dr[1])
			if !start.IsZero() {
				edu.Dates.Start = start
			}
			if !strings.EqualFold(dr[2], "Present") {
				end := parseDate(dr[2])
				if !end.IsZero() {
					edu.Dates.End = &end
				}
			}
			continue
		}

		// Single date
		if sd := reDateSingle.FindStringSubmatch(strings.TrimSpace(p)); sd != nil {
			start := parseDate(sd[1])
			if !start.IsZero() {
				edu.Dates.Start = start
			}
			continue
		}

		// Location (comma-separated)
		if strings.Contains(p, ",") || (p != "" && !strings.ContainsAny(p, "0123456789")) {
			locParts := strings.Split(p, ",")
			loc := &Location{}
			for i, lp := range locParts {
				locParts[i] = strings.TrimSpace(lp)
			}
			switch len(locParts) {
			case 1:
				loc.City = locParts[0]
			case 2:
				loc.City = locParts[0]
				loc.State = locParts[1]
			default:
				loc.City = locParts[0]
				loc.State = locParts[1]
				loc.Country = locParts[2]
			}
			edu.Location = loc
		}
	}
}

// parseEducationLine handles bullet lines and thesis within education.
func parseEducationLine(line string, edu *Education) {
	trimmed := strings.TrimSpace(line)

	b := reBullet.FindStringSubmatch(trimmed)
	if b == nil {
		return
	}
	content := strings.TrimSpace(b[1])

	// Check for thesis line: **Thesis:** Title — [Label](URL)
	if tm := reThesisLine.FindStringSubmatch(content); tm != nil {
		thesis := &Thesis{}
		rest := strings.TrimSpace(tm[1])

		// Split on dash to get title and link
		dashParts := reDashSplit.Split(rest, 2)
		thesis.Title = strings.TrimSpace(dashParts[0])

		if len(dashParts) > 1 {
			if lm := reLink.FindStringSubmatch(dashParts[1]); lm != nil {
				thesis.Link = Link{
					Label: strings.TrimSpace(lm[1]),
					URI:   strings.TrimSpace(lm[2]),
				}
			}
		}

		edu.Thesis = thesis
		return
	}

	// Regular bullet → degree description
	edu.Degree.Descriptions = append(edu.Degree.Descriptions, content)
}

// parseProjectH3 handles: ### Name — [Label](URL)
func parseProjectH3(content string, proj *Project) {
	parts := reDashSplit.Split(content, 2)
	proj.Name = strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		if lm := reLink.FindStringSubmatch(parts[1]); lm != nil {
			proj.Link = Link{
				Label: strings.TrimSpace(lm[1]),
				URI:   strings.TrimSpace(lm[2]),
			}
		}
	}
}

// parseProjectLine handles lines within a project entry.
func parseProjectLine(line string, proj *Project) {
	trimmed := strings.TrimSpace(line)

	// Date line
	if dr := reDateRange.FindStringSubmatch(trimmed); dr != nil {
		start := parseDate(dr[1])
		if !start.IsZero() {
			dates := &DateRange{Start: start}
			if !strings.EqualFold(dr[2], "Present") {
				end := parseDate(dr[2])
				if !end.IsZero() {
					dates.End = &end
				}
			}
			proj.Dates = dates
		}
		return
	}

	// Italic technologies
	if m := reItalic.FindStringSubmatch(trimmed); m != nil {
		techs := strings.Split(m[1], ",")
		for _, t := range techs {
			t = strings.TrimSpace(t)
			if t != "" {
				proj.Technologies = append(proj.Technologies, t)
			}
		}
		return
	}

	// Bullet highlight
	if b := reBullet.FindStringSubmatch(trimmed); b != nil {
		proj.Highlights = append(proj.Highlights, strings.TrimSpace(b[1]))
		return
	}
}

// parseCertificationLine parses: **Name** — Issuer (Notes)
func parseCertificationLine(content string, r *Resume) {
	if r.Certifications == nil {
		r.Certifications = &Certifications{}
	}

	cert := Certification{}

	// Extract bold name
	bm := reBold.FindStringSubmatch(content)
	if bm == nil {
		// Plain text fallback
		cert.Name = strings.TrimSpace(content)
		r.Certifications.Items = append(r.Certifications.Items, cert)
		return
	}
	cert.Name = strings.TrimSpace(bm[1])

	// Get rest after bold+dash
	rest := content[len(bm[0]):]
	rest = strings.TrimSpace(rest)
	rest = strings.TrimLeft(rest, "—–-")
	rest = strings.TrimSpace(rest)

	if rest == "" {
		r.Certifications.Items = append(r.Certifications.Items, cert)
		return
	}

	// Check for parenthesized notes at end
	if idx := strings.LastIndex(rest, "("); idx >= 0 && strings.HasSuffix(rest, ")") {
		cert.Notes = strings.TrimSpace(rest[idx+1 : len(rest)-1])
		rest = strings.TrimSpace(rest[:idx])
	}

	cert.Issuer = strings.TrimSpace(rest)
	r.Certifications.Items = append(r.Certifications.Items, cert)
}

// parseLanguageLine parses: **Name** — Proficiency
func parseLanguageLine(content string, r *Resume) {
	if r.Languages == nil {
		r.Languages = &LanguageList{}
	}

	lang := Language{}
	bm := reBold.FindStringSubmatch(content)
	if bm == nil {
		lang.Name = strings.TrimSpace(content)
		r.Languages.Languages = append(r.Languages.Languages, lang)
		return
	}
	lang.Name = strings.TrimSpace(bm[1])

	rest := content[len(bm[0]):]
	rest = strings.TrimSpace(rest)
	rest = strings.TrimLeft(rest, "—–-")
	rest = strings.TrimSpace(rest)

	lang.Proficiency = rest
	r.Languages.Languages = append(r.Languages.Languages, lang)
}

// flushExperience appends the current experience entry to the resume.
func flushExperience(exp *Experience, r *Resume) {
	if exp == nil {
		return
	}
	r.Experience.Positions = append(r.Experience.Positions, *exp)
}

// flushEducation appends the current education entry to the resume.
func flushEducation(edu *Education, r *Resume) {
	if edu == nil {
		return
	}
	r.Education.Institutions = append(r.Education.Institutions, *edu)
}

// flushProject appends the current project to the resume.
func flushProject(proj *Project, r *Resume) {
	if proj == nil {
		return
	}
	if r.Projects == nil {
		r.Projects = &ProjectList{}
	}
	r.Projects.Projects = append(r.Projects.Projects, *proj)
}

// parseDate tries to parse a month+year string using known formats.
func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	// Normalize period abbreviations: "Jan." → "Jan"
	s = strings.ReplaceAll(s, ".", "")
	for _, layout := range dateFormats {
		clean := strings.ReplaceAll(layout, ".", "")
		if t, err := time.Parse(clean, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
