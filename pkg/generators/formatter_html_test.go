package generators

import (
	"testing"

	"github.com/urmzd/resume-generator/pkg/resume"
)

func TestHTMLEscapeText(t *testing.T) {
	f := newHTMLFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"special chars", `<script>alert("xss")</script>`, "&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;"},
		{"ampersand", "A & B", "A &amp; B"},
		{"passthrough", "Hello World", "Hello World"},
		{"empty", "", ""},
		{"single quote", "it's", "it&#39;s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.EscapeText(tt.input)
			if got != tt.want {
				t.Errorf("EscapeText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHTMLFormatLink(t *testing.T) {
	f := newHTMLFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard url", "https://example.com", `<a href="https://example.com">https://example.com</a>`},
		{"html chars in url", "https://example.com/a&b", `<a href="https://example.com/a&amp;b">https://example.com/a&amp;b</a>`},
		{"empty", "", ""},
		{"whitespace", "  ", ""},
		{"with spaces trimmed", "  https://example.com  ", `<a href="https://example.com">https://example.com</a>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatLink(tt.input)
			if got != tt.want {
				t.Errorf("FormatLink(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLayoutClass(t *testing.T) {
	f := newHTMLFormatter()

	tests := []struct {
		name   string
		layout *resume.Layout
		want   string
	}{
		{"nil layout", nil, "density-standard typo-classic header-centered"},
		{"compact density", &resume.Layout{Density: "compact"}, "density-compact typo-classic header-centered"},
		{"modern typography", &resume.Layout{Typography: "modern"}, "density-standard typo-modern header-centered"},
		{"split header", &resume.Layout{Header: "split"}, "density-standard typo-classic header-split"},
		{"skill columns", &resume.Layout{SkillColumns: 3}, "density-standard typo-classic header-centered skills-columns"},
		{"invalid values default", &resume.Layout{Density: "invalid", Typography: "bad", Header: "nope"}, "density-standard typo-classic header-centered"},
		{"detailed density", &resume.Layout{Density: "detailed"}, "density-detailed typo-classic header-centered"},
		{"elegant typography", &resume.Layout{Typography: "elegant"}, "density-standard typo-elegant header-centered"},
		{"minimal header", &resume.Layout{Header: "minimal"}, "density-standard typo-classic header-minimal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.layoutClass(tt.layout)
			if got != tt.want {
				t.Errorf("layoutClass() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHasSection(t *testing.T) {
	f := newHTMLFormatter()

	fullResume := &resume.Resume{
		Summary: "A summary",
		Certifications: &resume.Certifications{
			Items: []resume.Certification{{Name: "AWS"}},
		},
		Education: resume.EducationList{
			Institutions: []resume.Education{{Institution: "MIT"}},
		},
		Skills: resume.Skills{
			Categories: []resume.SkillCategory{{Category: "Languages", Items: []string{"Go"}}},
		},
		Experience: resume.ExperienceList{
			Positions: []resume.Experience{{Title: "Dev"}},
		},
		Projects: &resume.ProjectList{
			Projects: []resume.Project{{Name: "Proj"}},
		},
		Languages: &resume.LanguageList{
			Languages: []resume.Language{{Name: "English"}},
		},
	}

	emptyResume := &resume.Resume{}

	sections := []string{"summary", "certifications", "education", "skills", "experience", "projects", "languages"}

	for _, section := range sections {
		t.Run(section+" present", func(t *testing.T) {
			if !f.hasSection(section, fullResume) {
				t.Errorf("hasSection(%q) = false for full resume", section)
			}
		})
		t.Run(section+" absent", func(t *testing.T) {
			if f.hasSection(section, emptyResume) {
				t.Errorf("hasSection(%q) = true for empty resume", section)
			}
		})
	}

	t.Run("unknown section", func(t *testing.T) {
		if f.hasSection("nonexistent", fullResume) {
			t.Error("hasSection should return false for unknown section")
		}
	})
}

func TestContainsSection(t *testing.T) {
	f := newHTMLFormatter()

	sections := []string{"summary", "skills", "experience"}

	tests := []struct {
		name string
		sec  string
		want bool
	}{
		{"present", "skills", true},
		{"absent", "projects", false},
		{"empty slice", "summary", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.containsSection(tt.sec, sections)
			if got != tt.want {
				t.Errorf("containsSection(%q) = %v, want %v", tt.sec, got, tt.want)
			}
		})
	}

	t.Run("empty sections slice", func(t *testing.T) {
		if f.containsSection("anything", nil) {
			t.Error("containsSection should return false for nil slice")
		}
	})
}

func TestHTMLTemplateFuncs(t *testing.T) {
	f := newHTMLFormatter()
	funcs := f.TemplateFuncs()

	expectedKeys := []string{
		"escape", "safeHTML",
		"formatDate", "formatDateShort", "formatDateRange", "fmtDateRange", "fmtOptDateRange", "calculateDuration",
		"formatLocation", "fmtLocation",
		"formatList", "join", "skillNames", "filterEmpty",
		"lower", "upper", "title",
		"replace", "hasPrefix", "hasSuffix", "contains", "trim",
		"formatLink", "fmtLink",
		"formatGPA", "sanitizePhone",
		"sortSkillsByOrder", "sortExperienceByOrder", "sortProjectsByOrder", "sortEducationByOrder", "sortLinksByOrder",
		"default",
		"layoutClass", "hasSection", "containsSection",
	}

	for _, key := range expectedKeys {
		if funcs[key] == nil {
			t.Errorf("TemplateFuncs missing key %q", key)
		}
	}
}
