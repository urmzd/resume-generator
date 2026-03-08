package generators

import (
	"testing"

	"github.com/urmzd/resume-generator/pkg/resume"
)

func TestMarkdownFormatLocation(t *testing.T) {
	f := newMarkdownFormatter()

	tests := []struct {
		name string
		loc  *resume.Location
		want string
	}{
		{"nil", nil, ""},
		{"full", &resume.Location{City: "NYC", State: "NY", Country: "USA"}, "NYC, NY, USA"},
		{"country only", &resume.Location{Country: "Canada"}, "Canada"},
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

func TestMarkdownTemplateFuncs(t *testing.T) {
	f := newMarkdownFormatter()
	funcs := f.TemplateFuncs()

	expectedKeys := []string{
		"escape",
		"fmtDateRange", "fmtOptDateRange", "fmtDates", "formatDate", "formatDateShort",
		"fmtLocation",
		"formatList", "join", "skillNames", "filterEmpty",
		"fmtLink", "fmtLinkWithDomain", "extractDisplayURL",
		"formatGPA", "sanitizePhone",
		"title", "upper", "lower",
		"trim", "filterEmpty2", "default",
		"bold", "italic",
		"sortExperienceByOrder", "sortProjectsByOrder", "sortEducationByOrder",
		"add",
	}

	for _, key := range expectedKeys {
		if funcs[key] == nil {
			t.Errorf("TemplateFuncs missing key %q", key)
		}
	}
}

func TestMarkdownBoldItalic(t *testing.T) {
	f := newMarkdownFormatter()
	funcs := f.TemplateFuncs()

	boldFn := funcs["bold"].(func(string) string)
	italicFn := funcs["italic"].(func(string) string)

	tests := []struct {
		name string
		fn   func(string) string
		in   string
		want string
	}{
		{"bold text", boldFn, "hello", "**hello**"},
		{"bold empty", boldFn, "", ""},
		{"italic text", italicFn, "hello", "*hello*"},
		{"italic empty", italicFn, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.in)
			if got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", tt.name, tt.in, got, tt.want)
			}
		})
	}
}

func TestMarkdownFmtLocation(t *testing.T) {
	f := newMarkdownFormatter()
	funcs := f.TemplateFuncs()
	fmtLoc := funcs["fmtLocation"].(func(interface{}) string)

	t.Run("pointer", func(t *testing.T) {
		loc := &resume.Location{City: "NYC"}
		got := fmtLoc(loc)
		if got != "NYC" {
			t.Errorf("fmtLocation(*Location) = %q, want %q", got, "NYC")
		}
	})

	t.Run("value", func(t *testing.T) {
		loc := resume.Location{City: "Berlin", Country: "Germany"}
		got := fmtLoc(loc)
		if got != "Berlin, Germany" {
			t.Errorf("fmtLocation(Location) = %q, want %q", got, "Berlin, Germany")
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		got := fmtLoc("not a location")
		if got != "" {
			t.Errorf("fmtLocation(string) = %q, want empty", got)
		}
	})
}
