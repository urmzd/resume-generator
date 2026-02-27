package resume

import (
	"testing"
	"time"
)

const fullResumeMD = `# Jane Doe

[jane@email.com](mailto:jane@email.com) | +1-555-1234 | Techville, Academia, USA | [linkedin.com/in/janedoe](https://linkedin.com/in/janedoe)

---

## Summary

Experienced engineer with a passion for building scalable systems.

## Core Skills

- **Programming Languages:** Python, Java, C++
- **Tools & Frameworks:** AWS, Docker, React

## Professional Experience

### Software Developer

**Tech Innovations Inc.** | Jul 2021 – Present | Techville, Academia

*Python, Java, Docker*

- Built scalable cloud systems
- Improved data retrieval by 40%

### Systems Engineer

**FutureSoft** | Jan 2023 – Dec 2023 | Innovation City, Futuristan

- Implemented network security upgrades

## Education

### Prestigious University — Ph.D. in Computer Science

Sep 2021 – May 2024 | Techville, Academia | GPA: 3.9 / 4.0

- Dean's List
- **Thesis:** ML for Code Review — [Link](https://example.com/thesis)

### University of Fictional — B.Sc. in Software Engineering

Sep 2017 – Jun 2021 | Imaginary City, Stateville

## Projects

### Finance Tracker — [Link](https://github.com/janedoe/finance-tracker)

Jan 2022 – Jun 2022

*React, Node.js*

- Full-stack personal finance app

### Eco Route Finder

- Calculates eco-friendly routes

## Certifications

- **AWS Solutions Architect** — Amazon Web Services (2023)
- **CKA** — CNCF

## Languages

- **English** — Native
- **French** — Intermediate
`

func TestParseMarkdownFullResume(t *testing.T) {
	r, err := parseMarkdown([]byte(fullResumeMD))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Contact
	assertEqual(t, "contact.name", "Jane Doe", r.Contact.Name)
	assertEqual(t, "contact.email", "jane@email.com", r.Contact.Email)
	assertEqual(t, "contact.phone", "+1-555-1234", r.Contact.Phone)

	if r.Contact.Location == nil {
		t.Fatal("expected contact.location to be set")
	}
	assertEqual(t, "contact.location.city", "Techville", r.Contact.Location.City)
	assertEqual(t, "contact.location.state", "Academia", r.Contact.Location.State)
	assertEqual(t, "contact.location.country", "USA", r.Contact.Location.Country)

	if len(r.Contact.Links) != 1 {
		t.Fatalf("expected 1 contact link, got %d", len(r.Contact.Links))
	}
	assertEqual(t, "contact.links[0].uri", "https://linkedin.com/in/janedoe", r.Contact.Links[0].URI)

	// Summary
	assertEqual(t, "summary", "Experienced engineer with a passion for building scalable systems.", r.Summary)

	// Skills
	assertEqual(t, "skills.title", "Core Skills", r.Skills.Title)
	if len(r.Skills.Categories) != 2 {
		t.Fatalf("expected 2 skill categories, got %d", len(r.Skills.Categories))
	}
	assertEqual(t, "skills[0].category", "Programming Languages", r.Skills.Categories[0].Category)
	assertSliceEqual(t, "skills[0].items", []string{"Python", "Java", "C++"}, r.Skills.Categories[0].Items)
	assertEqual(t, "skills[1].category", "Tools & Frameworks", r.Skills.Categories[1].Category)
	assertSliceEqual(t, "skills[1].items", []string{"AWS", "Docker", "React"}, r.Skills.Categories[1].Items)

	// Experience
	assertEqual(t, "experience.title", "Professional Experience", r.Experience.Title)
	if len(r.Experience.Positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(r.Experience.Positions))
	}
	exp0 := r.Experience.Positions[0]
	assertEqual(t, "exp[0].title", "Software Developer", exp0.Title)
	assertEqual(t, "exp[0].company", "Tech Innovations Inc.", exp0.Company)
	assertDateMonth(t, "exp[0].dates.start", 2021, time.July, exp0.Dates.Start)
	if exp0.Dates.End != nil {
		t.Errorf("expected exp[0].dates.end to be nil (Present), got %v", exp0.Dates.End)
	}
	assertSliceEqual(t, "exp[0].technologies", []string{"Python", "Java", "Docker"}, exp0.Technologies)
	if len(exp0.Highlights) != 2 {
		t.Fatalf("expected 2 highlights for exp[0], got %d", len(exp0.Highlights))
	}

	exp1 := r.Experience.Positions[1]
	assertEqual(t, "exp[1].title", "Systems Engineer", exp1.Title)
	assertEqual(t, "exp[1].company", "FutureSoft", exp1.Company)
	assertDateMonth(t, "exp[1].dates.start", 2023, time.January, exp1.Dates.Start)
	if exp1.Dates.End == nil {
		t.Fatal("expected exp[1].dates.end to be set")
	}
	assertDateMonth(t, "exp[1].dates.end", 2023, time.December, *exp1.Dates.End)

	// Education
	if len(r.Education.Institutions) != 2 {
		t.Fatalf("expected 2 education entries, got %d", len(r.Education.Institutions))
	}
	edu0 := r.Education.Institutions[0]
	assertEqual(t, "edu[0].institution", "Prestigious University", edu0.Institution)
	assertEqual(t, "edu[0].degree", "Ph.D. in Computer Science", edu0.Degree.Name)
	assertDateMonth(t, "edu[0].dates.start", 2021, time.September, edu0.Dates.Start)
	if edu0.GPA == nil {
		t.Fatal("expected edu[0].gpa to be set")
	}
	assertEqual(t, "edu[0].gpa.gpa", "3.9", edu0.GPA.GPA)
	assertEqual(t, "edu[0].gpa.max_gpa", "4.0", edu0.GPA.MaxGPA)
	if edu0.Thesis == nil {
		t.Fatal("expected edu[0].thesis to be set")
	}
	assertEqual(t, "edu[0].thesis.title", "ML for Code Review", edu0.Thesis.Title)
	assertEqual(t, "edu[0].thesis.link.uri", "https://example.com/thesis", edu0.Thesis.Link.URI)

	edu1 := r.Education.Institutions[1]
	assertEqual(t, "edu[1].institution", "University of Fictional", edu1.Institution)
	assertEqual(t, "edu[1].degree", "B.Sc. in Software Engineering", edu1.Degree.Name)

	// Projects
	if r.Projects == nil {
		t.Fatal("expected projects to be set")
	}
	if len(r.Projects.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(r.Projects.Projects))
	}
	proj0 := r.Projects.Projects[0]
	assertEqual(t, "proj[0].name", "Finance Tracker", proj0.Name)
	assertEqual(t, "proj[0].link.uri", "https://github.com/janedoe/finance-tracker", proj0.Link.URI)
	assertSliceEqual(t, "proj[0].technologies", []string{"React", "Node.js"}, proj0.Technologies)
	if proj0.Dates == nil {
		t.Fatal("expected proj[0].dates to be set")
	}
	assertDateMonth(t, "proj[0].dates.start", 2022, time.January, proj0.Dates.Start)

	proj1 := r.Projects.Projects[1]
	assertEqual(t, "proj[1].name", "Eco Route Finder", proj1.Name)
	if len(proj1.Highlights) != 1 {
		t.Fatalf("expected 1 highlight for proj[1], got %d", len(proj1.Highlights))
	}

	// Certifications
	if r.Certifications == nil {
		t.Fatal("expected certifications to be set")
	}
	if len(r.Certifications.Items) != 2 {
		t.Fatalf("expected 2 certifications, got %d", len(r.Certifications.Items))
	}
	assertEqual(t, "cert[0].name", "AWS Solutions Architect", r.Certifications.Items[0].Name)
	assertEqual(t, "cert[0].issuer", "Amazon Web Services", r.Certifications.Items[0].Issuer)
	assertEqual(t, "cert[0].notes", "2023", r.Certifications.Items[0].Notes)
	assertEqual(t, "cert[1].name", "CKA", r.Certifications.Items[1].Name)
	assertEqual(t, "cert[1].issuer", "CNCF", r.Certifications.Items[1].Issuer)

	// Languages
	if r.Languages == nil {
		t.Fatal("expected languages to be set")
	}
	if len(r.Languages.Languages) != 2 {
		t.Fatalf("expected 2 languages, got %d", len(r.Languages.Languages))
	}
	assertEqual(t, "lang[0].name", "English", r.Languages.Languages[0].Name)
	assertEqual(t, "lang[0].proficiency", "Native", r.Languages.Languages[0].Proficiency)
	assertEqual(t, "lang[1].name", "French", r.Languages.Languages[1].Name)
	assertEqual(t, "lang[1].proficiency", "Intermediate", r.Languages.Languages[1].Proficiency)
}

func TestParseMarkdownMinimal(t *testing.T) {
	md := `# John Smith

[john@example.com](mailto:john@example.com)

---
`
	r, err := parseMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "name", "John Smith", r.Contact.Name)
	assertEqual(t, "email", "john@example.com", r.Contact.Email)
}

func TestParseMarkdownNoName(t *testing.T) {
	md := `Some random text without headings`
	_, err := parseMarkdown([]byte(md))
	if err == nil {
		t.Fatal("expected error for missing H1 name")
	}
}

func TestParseMarkdownSectionOrderIndependence(t *testing.T) {
	md := `# Test Person

[test@test.com](mailto:test@test.com)

---

## Languages

- **Spanish** — Fluent

## Summary

A brief summary.

## Core Skills

- **Languages:** Go, Rust
`
	r, err := parseMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "summary", "A brief summary.", r.Summary)
	if r.Languages == nil || len(r.Languages.Languages) != 1 {
		t.Fatal("expected 1 language")
	}
	assertEqual(t, "lang.name", "Spanish", r.Languages.Languages[0].Name)
	if len(r.Skills.Categories) != 1 {
		t.Fatal("expected 1 skill category")
	}
}

func TestParseMarkdownDateFormats(t *testing.T) {
	tests := []struct {
		input string
		year  int
		month time.Month
	}{
		{"Jan 2024", 2024, time.January},
		{"January 2024", 2024, time.January},
		{"Dec 2023", 2023, time.December},
		{"September 2021", 2021, time.September},
	}
	for _, tc := range tests {
		d := parseDate(tc.input)
		if d.IsZero() {
			t.Errorf("parseDate(%q) returned zero", tc.input)
			continue
		}
		if d.Year() != tc.year || d.Month() != tc.month {
			t.Errorf("parseDate(%q) = %v, want %d-%s", tc.input, d, tc.year, tc.month)
		}
	}
}

func TestParseMarkdownEmptySections(t *testing.T) {
	md := `# Empty Resume

[empty@test.com](mailto:empty@test.com)

---

## Skills

## Experience

## Education
`
	r, err := parseMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "name", "Empty Resume", r.Contact.Name)
	if len(r.Skills.Categories) != 0 {
		t.Errorf("expected 0 skill categories, got %d", len(r.Skills.Categories))
	}
	if len(r.Experience.Positions) != 0 {
		t.Errorf("expected 0 positions, got %d", len(r.Experience.Positions))
	}
	if len(r.Education.Institutions) != 0 {
		t.Errorf("expected 0 education entries, got %d", len(r.Education.Institutions))
	}
}

func TestParseMarkdownExtraWhitespace(t *testing.T) {
	md := `#   Jane Doe

  [jane@test.com](mailto:jane@test.com)  |  +1-555-0000

---

##   Summary

  Some summary text with leading spaces.
`
	r, err := parseMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqual(t, "name", "Jane Doe", r.Contact.Name)
	assertEqual(t, "email", "jane@test.com", r.Contact.Email)
	assertEqual(t, "phone", "+1-555-0000", r.Contact.Phone)
	assertEqual(t, "summary", "Some summary text with leading spaces.", r.Summary)
}

func TestParseMarkdownDashVariants(t *testing.T) {
	// Test em-dash, en-dash, and regular hyphen in date ranges
	md := `# Test

[t@t.com](mailto:t@t.com)

---

## Experience

### Developer

**Acme Corp** | Jan 2020 — Dec 2021

- Did stuff

### Analyst

**BigCo** | Feb 2019 – Jan 2020

- Analyzed things
`
	r, err := parseMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Experience.Positions) != 2 {
		t.Fatalf("expected 2 positions, got %d", len(r.Experience.Positions))
	}
	assertDateMonth(t, "pos[0].start", 2020, time.January, r.Experience.Positions[0].Dates.Start)
	assertDateMonth(t, "pos[1].start", 2019, time.February, r.Experience.Positions[1].Dates.Start)
}

// --- Helpers ---

func assertEqual(t *testing.T, field, expected, actual string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %q, got %q", field, expected, actual)
	}
}

func assertSliceEqual(t *testing.T, field string, expected, actual []string) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("%s: expected %d items %v, got %d items %v", field, len(expected), expected, len(actual), actual)
		return
	}
	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("%s[%d]: expected %q, got %q", field, i, expected[i], actual[i])
		}
	}
}

func assertDateMonth(t *testing.T, field string, year int, month time.Month, d time.Time) {
	t.Helper()
	if d.Year() != year || d.Month() != month {
		t.Errorf("%s: expected %d-%s, got %v", field, year, month, d)
	}
}
