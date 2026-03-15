package generators

import (
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

func TestDocxEscapeText(t *testing.T) {
	f := newDocxFormatter()

	t.Run("passthrough", func(t *testing.T) {
		got := f.EscapeText("Hello <World> & 'Friends'")
		if got != "Hello <World> & 'Friends'" {
			t.Errorf("EscapeText should be passthrough, got %q", got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		if got := f.EscapeText(""); got != "" {
			t.Errorf("EscapeText('') = %q, want empty", got)
		}
	})
}

func TestDocxFormatDateRange(t *testing.T) {
	f := newDocxFormatter()

	jan2020 := pd(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))
	jun2021 := pd(time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC))
	dec2022 := pd(time.Date(2022, time.December, 15, 0, 0, 0, 0, time.UTC))
	year2020 := resume.NewYearDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name string
		dr   resume.DateRange
		want string
	}{
		{"month dates", resume.DateRange{Start: jan2020, End: func() *resume.PartialDate { d := jun2021; return &d }()}, "Jan 2020 Jun 2021"},
		{"regular dates", resume.DateRange{Start: jun2021, End: func() *resume.PartialDate { d := dec2022; return &d }()}, "Jun 2021 Dec 2022"},
		{"nil end", resume.DateRange{Start: jun2021}, "Jun 2021"},
		{"both zero", resume.DateRange{}, ""},
		{"year-only start", resume.DateRange{Start: year2020}, "2020"},
		{"year-only with end", resume.DateRange{Start: year2020, End: func() *resume.PartialDate { d := jun2021; return &d }()}, "2020 Jun 2021"},
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

func TestDocxFormatDateShort(t *testing.T) {
	f := newDocxFormatter()

	tests := []struct {
		name string
		pd   resume.PartialDate
		want string
	}{
		{"zero", resume.PartialDate{}, ""},
		{"year-only", resume.NewYearDate(time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)), "2020"},
		{"month date", resume.NewMonthDate(time.Date(2021, time.June, 15, 0, 0, 0, 0, time.UTC)), "Jun 2021"},
		{"full date", resume.NewPartialDate(time.Date(2022, time.December, 31, 0, 0, 0, 0, time.UTC), resume.PrecisionFull), "Dec 2022"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.formatDateShort(tt.pd)
			if got != tt.want {
				t.Errorf("formatDateShort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDocxFormatGPA(t *testing.T) {
	f := newDocxFormatter()

	tests := []struct {
		name     string
		gpa, max string
		want     string
	}{
		{"default max", "3.9", "4.0", "3.9"},
		{"custom max uses slash", "3.9", "5.0", "3.9/5.0"},
		{"empty gpa", "", "4.0", ""},
		{"whitespace gpa", "  ", "4.0", ""},
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

func TestDocxSanitizePhone(t *testing.T) {
	f := newDocxFormatter()

	t.Run("returns input unchanged", func(t *testing.T) {
		input := "+1 (555) 123-4567"
		if got := f.SanitizePhone(input); got != input {
			t.Errorf("SanitizePhone() = %q, want %q (unchanged)", got, input)
		}
	})

	t.Run("empty", func(t *testing.T) {
		if got := f.SanitizePhone(""); got != "" {
			t.Errorf("SanitizePhone('') = %q, want empty", got)
		}
	})
}

func TestDocxTemplateFuncs(t *testing.T) {
	f := newDocxFormatter()
	funcs := f.TemplateFuncs()

	if len(funcs) != 0 {
		t.Errorf("TemplateFuncs() should return empty FuncMap, got %d entries", len(funcs))
	}
}
