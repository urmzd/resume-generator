package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

func TestLaTeXEscapeText(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"backslash", `a\b`, `a\textbackslash{}b`},
		{"open brace", `a{b`, `a\{b`},
		{"close brace", `a}b`, `a\}b`},
		{"dollar", `a$b`, `a\$b`},
		{"ampersand", `a&b`, `a\&b`},
		{"percent", `a%b`, `a\%b`},
		{"hash", `a#b`, `a\#b`},
		{"underscore", `a_b`, `a\_b`},
		{"tilde", `a~b`, `a\textasciitilde{}b`},
		{"caret", `a^b`, `a\textasciicircum{}b`},
		{"passthrough", "Hello World", "Hello World"},
		{"empty", "", ""},
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

func TestLaTeXFormatLocation(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name string
		loc  *resume.Location
		want string
	}{
		{"nil", nil, ""},
		{"latex escaped city", &resume.Location{City: "A&B Corp"}, `A\&B Corp`},
		{"standard", &resume.Location{City: "NYC", State: "NY"}, `NYC, NY`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatLocation(tt.loc)
			if got != tt.want {
				t.Errorf("FormatLocation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLaTeXFormatList(t *testing.T) {
	f := newLaTeXFormatter()

	t.Run("items with special chars", func(t *testing.T) {
		got := f.FormatList([]string{"A&B", "C#D"})
		if !strings.Contains(got, `\&`) || !strings.Contains(got, `\#`) {
			t.Errorf("FormatList should escape special chars, got %q", got)
		}
	})

	t.Run("empty items filtered", func(t *testing.T) {
		got := f.FormatList([]string{"Go", "", "Rust"})
		if got != "Go, Rust" {
			t.Errorf("FormatList() = %q, want %q", got, "Go, Rust")
		}
	})

	t.Run("all empty", func(t *testing.T) {
		got := f.FormatList([]string{"", " "})
		if got != "" {
			t.Errorf("FormatList() = %q, want empty", got)
		}
	})
}

func TestLaTeXFormatGPA(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name     string
		gpa, max string
		want     string
	}{
		{"standard", "3.9", "4.0", "3.9"},
		{"custom max", "3.9", "5.0", `3.9 / 5.0`},
		{"empty", "", "4.0", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatGPA(tt.gpa, tt.max)
			if got != tt.want {
				t.Errorf("FormatGPA(%q, %q) = %q, want %q", tt.gpa, tt.max, got, tt.want)
			}
		})
	}
}

func TestLaTeXFormatGPAStruct(t *testing.T) {
	f := newLaTeXFormatter()

	t.Run("nil", func(t *testing.T) {
		if got := f.FormatGPAStruct(nil); got != "" {
			t.Errorf("FormatGPAStruct(nil) = %q, want empty", got)
		}
	})

	t.Run("standard", func(t *testing.T) {
		got := f.FormatGPAStruct(&resume.GPA{GPA: "3.9", MaxGPA: "4.0"})
		if got != "3.9" {
			t.Errorf("FormatGPAStruct() = %q, want %q", got, "3.9")
		}
	})
}

func TestLaTeXFormatterFormatDateRange(t *testing.T) {
	f := newLaTeXFormatter()

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		dr   resume.DateRange
		want string
	}{
		{"uses textendash", resume.DateRange{Start: jan2020, End: &jun2021}, `Jan 2020 \textendash\ Jun 2021`},
		{"start only", resume.DateRange{Start: jan2020, End: nil}, `Jan 2020 \textendash\ Present`},
		{"both zero", resume.DateRange{}, ""},
		{"same month", resume.DateRange{Start: jan2020, End: &jan2020}, "Jan 2020"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatDateRange(tt.dr)
			if got != tt.want {
				t.Errorf("FormatDateRange() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLaTeXFormatLink(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard", "https://example.com", `\href{https://example.com}{https://example.com}`},
		{"empty", "", ""},
		{"url with underscore", "https://example.com/my_page", `\href{https://example.com/my\_page}{https://example.com/my\_page}`},
		{"url with hash", "https://example.com/page#section", `\href{https://example.com/page\#section}{https://example.com/page\#section}`},
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

func TestLaTeXFormatLinkWithDomain(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard", "https://example.com", `\href{https://example.com}{example.com}`},
		{"with www", "https://www.example.com/", `\href{https://www.example.com/}{example.com}`},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatLinkWithDomain(tt.input)
			if got != tt.want {
				t.Errorf("FormatLinkWithDomain(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLaTeXExtractDisplayURL(t *testing.T) {
	f := newLaTeXFormatter()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"https stripped", "https://example.com", "example.com"},
		{"http stripped", "http://example.com", "example.com"},
		{"www stripped", "https://www.example.com/", "example.com"},
		{"empty", "", ""},
		{"no protocol", "example.com", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.ExtractDisplayURL(tt.input)
			if got != tt.want {
				t.Errorf("ExtractDisplayURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLaTeXTemplateFuncs(t *testing.T) {
	f := newLaTeXFormatter()
	funcs := f.TemplateFuncs()

	expectedKeys := []string{
		"escape", "escapeLatexChars",
		"fmtDateRange", "fmtDates", "formatDateRange", "fmtDateLegal",
		"join", "formatList", "skillNames",
		"fmtLink", "fmtLinkWithDomain", "extractDisplayURL",
		"fmtLocation", "formatLocationFull", "formatLocationShort",
		"formatGPA",
		"title", "upper", "lower",
		"trim", "filterEmpty", "default",
		"sortExperienceByOrder", "sortProjectsByOrder", "sortEducationByOrder",
		"add", "employmentType", "now", "linkLabel",
	}

	for _, key := range expectedKeys {
		if funcs[key] == nil {
			t.Errorf("TemplateFuncs missing key %q", key)
		}
	}
}
