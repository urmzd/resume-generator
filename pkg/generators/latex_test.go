package generators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
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

	if got := formatter.FormatDateRange(resume.DateRange{Start: start, End: &end}); got != `Jan 2020 \textendash\ Jun 2021` {
		t.Fatalf("formatDateRange(start, &end, false) = %q, want Jan 2020 \\textendash\\ Jun 2021", got)
	}

	if got := formatter.FormatDateRange(resume.DateRange{Start: start}); got != `Jan 2020 \textendash\ Present` {
		t.Fatalf("formatDateRange(start, nil, true) = %q, want Jan 2020 \\textendash\\ Present", got)
	}

	var zero time.Time
	if got := formatter.FormatDateRange(resume.DateRange{Start: zero}); got != "" {
		t.Fatalf("formatDateRange(zero, nil, false) = %q, want empty string", got)
	}
}

func TestLaTeXGeneratorGenerate(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewLaTeXGenerator(logger)

	expStart := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	eduStart := time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	resume := &resume.Resume{
		Contact: resume.Contact{
			Name:  "John & Co.",
			Email: "john@example.com",
			Phone: "+1 (555) 123-4567",
			Links: []resume.Link{
				{URI: "https://example.com"},
			},
		},
		Skills: resume.Skills{
			Categories: []resume.SkillCategory{
				{
					Category: "Languages",
					Items: []string{
						"Go",
						"Rust",
					},
				},
			},
		},
		Experience: resume.ExperienceList{
			Positions: []resume.Experience{
				{
					Title:   "Engineer",
					Company: "Acme #1",
					Highlights: []string{
						"Improved throughput by 50%",
					},
					Dates: resume.DateRange{
						Start: expStart,
					},
				},
			},
		},
		Education: resume.EducationList{
			Institutions: []resume.Education{
				{
					Institution: "University of {Code}",
					Degree: resume.Degree{
						Name: "B.Sc Computer Science"},
					GPA: &resume.GPA{
						GPA:    "3.9",
						MaxGPA: "4.0",
					},
					Dates: resume.DateRange{
						Start: eduStart,
						End:   &eduEnd,
					},
				},
			},
		},
		Projects: &resume.ProjectList{
			Projects: []resume.Project{
				{
					Name: "Project_One",
					Highlights: []string{
						"Deployed to 100% of regions",
					},
					Link: resume.Link{URI: "https://project.example.com"},
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
		{"website link", `\href{https://example.com}{example.com}`},
		{"experience company escape", `Acme \#1`},
		{"experience dates", `Jan 2022 \textendash\ Present`},
		{"experience highlight escape", `50\%`},
		{"skills display", `Go, Rust`},
		{"education institution escape", `University of \{Code\}`},
		{"education date range", `Sep 2018 \textendash\ Jun 2021`},
		{"project name escape", `Project\_One`},
		{"project link text", `https://project.example.com`},
	}

	for _, expectation := range expects {
		if !strings.Contains(got, expectation.sub) {
			t.Errorf("Generate() missing %s: output = %q", expectation.name, got)
		}
	}
}
