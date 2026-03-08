package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
)

func TestFormatDateRange(t *testing.T) {
	f := &baseFormatter{}

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	jan2020dup := time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		dr   resume.DateRange
		want string
	}{
		{"both dates set", resume.DateRange{Start: jan2020, End: &jun2021}, "Jan 2020 – Jun 2021"},
		{"end nil (present)", resume.DateRange{Start: jan2020, End: nil}, "Jan 2020 – Present"},
		{"end zero", resume.DateRange{Start: jan2020, End: func() *time.Time { t := time.Time{}; return &t }()}, "Jan 2020 – Present"},
		{"same month", resume.DateRange{Start: jan2020, End: &jan2020dup}, "Jan 2020"},
		{"both zero", resume.DateRange{}, ""},
		{"start zero end set", resume.DateRange{Start: time.Time{}, End: &jun2021}, "Jun 2021"},
		{"start set end same as start format", resume.DateRange{Start: jan2020, End: &jan2020}, "Jan 2020"},
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

func TestFormatOptionalDateRange(t *testing.T) {
	f := &baseFormatter{}

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		dr   *resume.DateRange
		want string
	}{
		{"nil", nil, ""},
		{"valid", &resume.DateRange{Start: jan2020, End: &jun2021}, "Jan 2020 – Jun 2021"},
		{"zero", &resume.DateRange{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatOptionalDateRange(tt.dr)
			if got != tt.want {
				t.Errorf("FormatOptionalDateRange() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatDates(t *testing.T) {
	f := &baseFormatter{}

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"string", "Jan 2020 – Jun 2021", "Jan 2020 – Jun 2021"},
		{"string with spaces", "  Jan 2020  ", "Jan 2020"},
		{"DateRange value", resume.DateRange{Start: jan2020, End: &jun2021}, "Jan 2020 – Jun 2021"},
		{"DateRange pointer nil", (*resume.DateRange)(nil), ""},
		{"DateRange pointer valid", &resume.DateRange{Start: jan2020, End: &jun2021}, "Jan 2020 – Jun 2021"},
		{"unsupported type", 42, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatDates(tt.value)
			if got != tt.want {
				t.Errorf("FormatDates() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatLocation(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name   string
		loc    *resume.Location
		escape func(string) string
		want   string
	}{
		{"nil", nil, nil, ""},
		{"city only", &resume.Location{City: "NYC"}, nil, "NYC"},
		{"full", &resume.Location{City: "NYC", State: "NY", Country: "USA"}, nil, "NYC, NY, USA"},
		{"city and country", &resume.Location{City: "Berlin", Country: "Germany"}, nil, "Berlin, Germany"},
		{"country only", &resume.Location{Country: "Canada"}, nil, "Canada"},
		{"dedup country matches state", &resume.Location{City: "Singapore", State: "Singapore", Country: "Singapore"}, nil, "Singapore, Singapore"},
		{"dedup case insensitive", &resume.Location{City: "Tokyo", Country: "tokyo"}, nil, "Tokyo"},
		{"with escape func", &resume.Location{City: "A&B"}, func(s string) string { return strings.ReplaceAll(s, "&", "&amp;") }, "A&amp;B"},
		{"whitespace trimmed", &resume.Location{City: "  NYC  ", State: "  NY  "}, nil, "NYC, NY"},
		{"all empty strings", &resume.Location{City: "", State: "", Country: ""}, nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatLocation(tt.loc, tt.escape)
			if got != tt.want {
				t.Errorf("FormatLocation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatGPA(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name     string
		gpa, max string
		want     string
	}{
		{"default max", "3.9", "4.0", "3.9"},
		{"custom max", "3.9", "5.0", "3.9 / 5.0"},
		{"empty gpa", "", "4.0", ""},
		{"empty max", "3.9", "", "3.9"},
		{"whitespace gpa", "  ", "4.0", ""},
		{"whitespace max treated as empty", "3.9", "  ", "3.9"},
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

func TestFormatGPAStruct(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name string
		gpa  *resume.GPA
		want string
	}{
		{"nil", nil, ""},
		{"default max", &resume.GPA{GPA: "3.9", MaxGPA: "4.0"}, "3.9"},
		{"custom max", &resume.GPA{GPA: "3.9", MaxGPA: "5.0"}, "3.9 / 5.0"},
		{"empty gpa field", &resume.GPA{GPA: "", MaxGPA: "4.0"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatGPAStruct(tt.gpa)
			if got != tt.want {
				t.Errorf("FormatGPAStruct() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSanitizePhone(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name  string
		phone string
		want  string
	}{
		{"formatted", "+1 (555) 123-4567", "+15551234567"},
		{"already clean", "+15551234567", "+15551234567"},
		{"empty", "", ""},
		{"letters removed", "abc123def", "123"},
		{"unicode removed", "☎️ 555-1234", "5551234"},
		{"plus preserved", "+44 20 7946 0958", "+442079460958"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.SanitizePhone(tt.phone)
			if got != tt.want {
				t.Errorf("SanitizePhone(%q) = %q, want %q", tt.phone, got, tt.want)
			}
		})
	}
}

func TestFormatList(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{"normal", []string{"Go", "Rust", "Python"}, "Go, Rust, Python"},
		{"filters empty", []string{"Go", "", "Python", "  "}, "Go, Python"},
		{"nil slice", nil, ""},
		{"single item", []string{"Go"}, "Go"},
		{"all empty", []string{"", " ", "  "}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.FormatList(tt.values)
			if got != tt.want {
				t.Errorf("FormatList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSkillNames(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name  string
		items []string
		want  int // expected length
	}{
		{"filters empty", []string{"Go", "", "Rust"}, 2},
		{"nil slice", nil, 0},
		{"all empty", []string{"", " "}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.SkillNames(tt.items)
			if len(got) != tt.want {
				t.Errorf("SkillNames() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name  string
		sep   string
		items []string
		want  string
	}{
		{"pipe sep", " | ", []string{"a", "b", "c"}, "a | b | c"},
		{"empty items", ", ", []string{}, ""},
		{"single", ", ", []string{"only"}, "only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.Join(tt.sep, tt.items)
			if got != tt.want {
				t.Errorf("Join() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCaseTransformations(t *testing.T) {
	f := &baseFormatter{}

	tests := []struct {
		name  string
		fn    func(string) string
		input string
		want  string
	}{
		{"lower mixed", f.Lower, "Hello World", "hello world"},
		{"lower empty", f.Lower, "", ""},
		{"upper mixed", f.Upper, "Hello World", "HELLO WORLD"},
		{"upper empty", f.Upper, "", ""},
		{"title lower", f.Title, "hello world", "Hello World"},
		{"title empty", f.Title, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", tt.name, tt.input, got, tt.want)
			}
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	f := &baseFormatter{}

	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end2y3m := time.Date(2022, time.April, 1, 0, 0, 0, 0, time.UTC)
	end6m := time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC)
	end10d := time.Date(2020, time.January, 11, 0, 0, 0, 0, time.UTC)
	end1y := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		start time.Time
		end   *time.Time
		want  string
	}{
		{"years and months", start, &end2y3m, "2 yr 3 mo"},
		{"months only", start, &end6m, "6 mo"},
		{"less than 1 month", start, &end10d, "< 1 mo"},
		{"exact year", start, &end1y, "1 yr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.CalculateDuration(tt.start, tt.end)
			if got != tt.want {
				t.Errorf("CalculateDuration() = %q, want %q", got, tt.want)
			}
		})
	}

	// Nil end uses time.Now(), so just check it returns a non-empty string
	t.Run("nil end uses now", func(t *testing.T) {
		recentStart := time.Now().AddDate(-1, -3, 0)
		got := f.CalculateDuration(recentStart, nil)
		if got == "" {
			t.Error("CalculateDuration with nil end should return non-empty string")
		}
		if !strings.Contains(got, "yr") && !strings.Contains(got, "mo") {
			t.Errorf("CalculateDuration with nil end returned unexpected format: %q", got)
		}
	})
}

func TestSortExperienceByDate(t *testing.T) {
	f := &baseFormatter{}

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	mar2022 := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC)

	t.Run("sorts descending by start date", func(t *testing.T) {
		input := []resume.Experience{
			{Title: "First", Dates: resume.DateRange{Start: jan2020}},
			{Title: "Third", Dates: resume.DateRange{Start: mar2022}},
			{Title: "Second", Dates: resume.DateRange{Start: jun2021}},
		}
		got := f.SortExperienceByDate(input)
		if got[0].Title != "Third" || got[1].Title != "Second" || got[2].Title != "First" {
			t.Errorf("unexpected order: %s, %s, %s", got[0].Title, got[1].Title, got[2].Title)
		}
	})

	t.Run("does not mutate input", func(t *testing.T) {
		input := []resume.Experience{
			{Title: "First", Dates: resume.DateRange{Start: jan2020}},
			{Title: "Second", Dates: resume.DateRange{Start: jun2021}},
		}
		_ = f.SortExperienceByDate(input)
		if input[0].Title != "First" {
			t.Error("SortExperienceByDate mutated input slice")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		got := f.SortExperienceByDate(nil)
		if len(got) != 0 {
			t.Errorf("expected empty, got %d", len(got))
		}
	})

	t.Run("single element", func(t *testing.T) {
		input := []resume.Experience{{Title: "Only"}}
		got := f.SortExperienceByDate(input)
		if len(got) != 1 || got[0].Title != "Only" {
			t.Error("unexpected result for single element")
		}
	})
}

func TestSortEducationByDate(t *testing.T) {
	f := &baseFormatter{}

	jan2018 := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	sep2020 := time.Date(2020, time.September, 1, 0, 0, 0, 0, time.UTC)

	t.Run("sorts descending", func(t *testing.T) {
		input := []resume.Education{
			{Institution: "Early", Dates: resume.DateRange{Start: jan2018}},
			{Institution: "Late", Dates: resume.DateRange{Start: sep2020}},
		}
		got := f.SortEducationByDate(input)
		if got[0].Institution != "Late" {
			t.Error("expected Late first")
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := f.SortEducationByDate(nil)
		if len(got) != 0 {
			t.Errorf("expected empty, got %d", len(got))
		}
	})

	t.Run("does not mutate input", func(t *testing.T) {
		input := []resume.Education{
			{Institution: "First", Dates: resume.DateRange{Start: jan2018}},
			{Institution: "Second", Dates: resume.DateRange{Start: sep2020}},
		}
		_ = f.SortEducationByDate(input)
		if input[0].Institution != "First" {
			t.Error("SortEducationByDate mutated input slice")
		}
	})
}

func TestSortProjectsByDate(t *testing.T) {
	f := &baseFormatter{}

	jan2020 := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	jun2021 := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	mar2022 := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC)

	t.Run("sorts descending", func(t *testing.T) {
		input := []resume.Project{
			{Name: "Early", Dates: &resume.DateRange{Start: jan2020}},
			{Name: "Late", Dates: &resume.DateRange{Start: mar2022}},
			{Name: "Mid", Dates: &resume.DateRange{Start: jun2021}},
		}
		got := f.SortProjectsByDate(input)
		if got[0].Name != "Late" || got[1].Name != "Mid" || got[2].Name != "Early" {
			t.Errorf("unexpected order: %s, %s, %s", got[0].Name, got[1].Name, got[2].Name)
		}
	})

	t.Run("nil dates sort to end", func(t *testing.T) {
		input := []resume.Project{
			{Name: "NoDates", Dates: nil},
			{Name: "HasDates", Dates: &resume.DateRange{Start: jan2020}},
		}
		got := f.SortProjectsByDate(input)
		if got[0].Name != "HasDates" || got[1].Name != "NoDates" {
			t.Errorf("nil dates should sort to end, got %s, %s", got[0].Name, got[1].Name)
		}
	})

	t.Run("all nil dates stable", func(t *testing.T) {
		input := []resume.Project{
			{Name: "A", Dates: nil},
			{Name: "B", Dates: nil},
		}
		got := f.SortProjectsByDate(input)
		if got[0].Name != "A" || got[1].Name != "B" {
			t.Error("nil dates should maintain original order")
		}
	})

	t.Run("mix of nil and non-nil", func(t *testing.T) {
		input := []resume.Project{
			{Name: "NoDates1", Dates: nil},
			{Name: "Late", Dates: &resume.DateRange{Start: mar2022}},
			{Name: "NoDates2", Dates: nil},
			{Name: "Early", Dates: &resume.DateRange{Start: jan2020}},
		}
		got := f.SortProjectsByDate(input)
		if got[0].Name != "Late" || got[1].Name != "Early" {
			t.Errorf("dated projects should come first in desc order, got %s, %s", got[0].Name, got[1].Name)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := f.SortProjectsByDate(nil)
		if len(got) != 0 {
			t.Errorf("expected empty, got %d", len(got))
		}
	})

	t.Run("does not mutate input", func(t *testing.T) {
		input := []resume.Project{
			{Name: "B", Dates: &resume.DateRange{Start: jan2020}},
			{Name: "A", Dates: &resume.DateRange{Start: mar2022}},
		}
		_ = f.SortProjectsByDate(input)
		if input[0].Name != "B" {
			t.Error("SortProjectsByDate mutated input slice")
		}
	})
}
