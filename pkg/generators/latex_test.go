package generators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

func TestEscapeLaTeX(t *testing.T) {
	formatter := newLaTeXFormatter()
	input := `50% {value} #test_~^&`
	got := formatter.EscapeText(input)

	assertions := []struct {
		name string
		sub  string
	}{
		{"percent", `50\%`},
		{"brace open", `\{value`},
		{"brace close", `value\}`},
		{"hash", `\#test`},
		{"underscore", `\_`},
		{"tilde", `\textasciitilde{}`},
		{"caret", `\textasciicircum{}`},
		{"ampersand", `\&`},
	}

	for _, assertion := range assertions {
		if !strings.Contains(got, assertion.sub) {
			t.Errorf("escapeLaTeX missing %s escape: %q", assertion.name, got)
		}
	}
}

func TestLaTeXFormatDateRange(t *testing.T) {
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	formatter := newLaTeXFormatter()

	if got := formatter.FormatDateRange(definition.DateRange{Start: start, End: &end}); got != "Jan 2020 - Jun 2021" {
		t.Fatalf("formatDateRange(start, &end, false) = %q, want Jan 2020 - Jun 2021", got)
	}

	if got := formatter.FormatDateRange(definition.DateRange{Start: start, Current: true}); got != "Jan 2020 - Present" {
		t.Fatalf("formatDateRange(start, nil, true) = %q, want Jan 2020 - Present", got)
	}

	var zero time.Time
	if got := formatter.FormatDateRange(definition.DateRange{Start: zero}); got != "" {
		t.Fatalf("formatDateRange(zero, nil, false) = %q, want empty string", got)
	}
}

func TestLaTeXGeneratorGenerate(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewLaTeXGenerator(logger)

	expStart := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	eduStart := time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	resume := &definition.Resume{
		Contact: definition.Contact{
			Name:  "John & Co.",
			Email: "john@example.com",
			Phone: "+1 (555) 123-4567",
			Links: []definition.Link{
				{Order: 1, URL: "https://example.com", Text: "Website"},
			},
		},
		Skills: definition.Skills{
			Categories: []definition.SkillCategory{
				{
					Order: 1,
					Name:  "Languages",
					Items: []definition.SkillItem{
						{Order: 1, Name: "Go"},
						{Order: 2, Name: "Rust"},
					},
				},
			},
		},
		Experience: definition.ExperienceList{
			Positions: []definition.Experience{
				{
					Order:   1,
					Title:   "Engineer",
					Company: "Acme #1",
					Description: []string{
						"Improved throughput by 50%",
					},
					Dates: definition.DateRange{
						Start:   expStart,
						Current: true,
					},
				},
			},
		},
		Education: definition.EducationList{
			Institutions: []definition.Education{
				{
					Order:       1,
					Institution: "University of {Code}",
					Degree:      "B.Sc Computer Science",
					GPA:         "3.9",
					MaxGPA:      "4.0",
					Dates: definition.DateRange{
						Start: eduStart,
						End:   &eduEnd,
					},
				},
			},
		},
		Projects: definition.ProjectList{
			Projects: []definition.Project{
				{
					Order:        1,
					Name:         "Project_One",
					Technologies: []string{"Go", "Terraform"},
					Description: []string{
						"Deployed to 100% of regions",
					},
					Links: []definition.Link{
						{Order: 1, URL: "https://project.example.com", Text: "Repo"},
					},
				},
			},
		},
	}

	templatePath := filepath.Join("..", "..", "templates", "modern-latex", "template.tex")
	templateContentBytes, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read LaTeX template: %v", err)
	}

	templateContent := string(templateContentBytes)

	got, err := gen.Generate(templateContent, resume)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	expects := []struct {
		name string
		sub  string
	}{
		{"escaped name", `John \& Co.`},
		{"email macro", `\email{john@example.com}`},
		{"phone macro", `\phone{+1 (555) 123-4567}`},
		{"website link", `\href{https://example.com}{Website}`},
		{"experience company escape", `Acme \#1`},
		{"experience dates", `Jan 2022 - Present`},
		{"experience highlight escape", `50\%`},
		{"skills display", `Go, Rust`},
		{"education institution escape", `University of \{Code\}`},
		{"education date range", `Sep 2018 - Jun 2021`},
		{"project name escape", `Project\_One`},
		{"project link text", `https://project.example.com`},
	}

	for _, expectation := range expects {
		if !strings.Contains(got, expectation.sub) {
			t.Errorf("Generate() missing %s: output = %q", expectation.name, got)
		}
	}
}
