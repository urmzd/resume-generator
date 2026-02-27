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

func TestMarkdownEscapeText(t *testing.T) {
	formatter := newMarkdownFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"backslash", `a\b`, `a\\b`},
		{"backtick", "a`b", "a\\`b"},
		{"asterisk", `a*b`, `a\*b`},
		{"underscore", `a_b`, `a\_b`},
		{"braces", `a{b}`, `a\{b\}`},
		{"brackets", `a[b]`, `a\[b\]`},
		{"parens", `a(b)`, `a\(b\)`},
		{"hash", `#heading`, `\#heading`},
		{"plus", `a+b`, `a\+b`},
		{"dash", `a-b`, `a\-b`},
		{"dot", `a.b`, `a\.b`},
		{"bang", `a!b`, `a\!b`},
		{"pipe", `a|b`, `a\|b`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.EscapeText(tt.input)
			if got != tt.want {
				t.Errorf("EscapeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMarkdownFormatLink(t *testing.T) {
	formatter := newMarkdownFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"full url", "https://example.com", "[example.com](https://example.com)"},
		{"with www", "https://www.example.com", "[example.com](https://www.example.com)"},
		{"empty", "", ""},
		{"spaces", "  ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.FormatLink(tt.input)
			if got != tt.want {
				t.Errorf("FormatLink(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMarkdownExtractDisplayURL(t *testing.T) {
	formatter := newMarkdownFormatter()

	tests := []struct {
		input string
		want  string
	}{
		{"https://example.com/", "example.com"},
		{"http://www.example.com/", "example.com"},
		{"https://github.com/user/repo", "github.com/user/repo"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := formatter.ExtractDisplayURL(tt.input)
			if got != tt.want {
				t.Errorf("ExtractDisplayURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMarkdownGeneratorGenerate(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewMarkdownGenerator(logger)

	expStart := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	eduStart := time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	r := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+1 (555) 123-4567",
			Links: []resume.Link{
				{URI: "https://github.com/janedoe"},
			},
		},
		Summary: "Experienced software engineer with a passion for building scalable systems.",
		Skills: resume.Skills{
			Categories: []resume.SkillCategory{
				{
					Category: "Languages",
					Items:    []string{"Go", "Rust", "Python"},
				},
			},
		},
		Experience: resume.ExperienceList{
			Positions: []resume.Experience{
				{
					Title:   "Senior Engineer",
					Company: "Acme Corp",
					Highlights: []string{
						"Improved throughput by 50%",
						"Led team of 5 engineers",
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
					Institution: "MIT",
					Degree:      resume.Degree{Name: "B.Sc Computer Science"},
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
					Name: "OpenSource Tool",
					Highlights: []string{
						"Used by 1000+ developers",
					},
					Link: resume.Link{URI: "https://github.com/janedoe/tool"},
				},
			},
		},
	}

	templatePath := filepath.Join("..", "..", "templates", "modern-markdown", "template.md")
	templateContentBytes, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read Markdown template: %v", err)
	}

	got, err := gen.Generate(string(templateContentBytes), r)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	expects := []struct {
		name string
		sub  string
	}{
		{"name heading", "# Jane Doe"},
		{"email link", "[jane@example.com](mailto:jane@example.com)"},
		{"github link", "github.com/janedoe"},
		{"summary text", "Experienced software engineer"},
		{"skills category", "**Languages:**"},
		{"skills items", "Go, Rust, Python"},
		{"experience title", "### Senior Engineer"},
		{"experience company", "**Acme Corp**"},
		{"experience dates", "Jan 2022"},
		{"experience highlight", "Improved throughput by 50%"},
		{"education institution", "### MIT"},
		{"education degree", "B.Sc Computer Science"},
		{"education dates", "Sep 2018"},
		{"project name", "### OpenSource Tool"},
		{"project link", "github.com/janedoe/tool"},
		{"horizontal rule", "---"},
	}

	for _, expectation := range expects {
		if !strings.Contains(got, expectation.sub) {
			t.Errorf("Generate() missing %s: expected %q in output", expectation.name, expectation.sub)
		}
	}
}
