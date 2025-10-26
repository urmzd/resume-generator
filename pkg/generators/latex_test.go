package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

func TestLaTeXTemplateFuncs(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewLaTeXGenerator(logger)
	funcs := gen.templateFuncs()

	escape := funcs["escape"].(func(string) string)
	if got := escape(`50% {value} #test`); !strings.Contains(got, `50\%`) || !strings.Contains(got, `\{value\}`) || !strings.Contains(got, `\#test`) {
		t.Errorf("escape returned %q, expected LaTeX-escaped characters", got)
	}

	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	fmtDateRange := funcs["fmtDateRange"].(func(definition.DateRange) string)
	if got := fmtDateRange(definition.DateRange{Start: start, End: &end}); got != "Jan 2020 - Jun 2021" {
		t.Errorf("fmtDateRange returned %q, want Jan 2020 - Jun 2021", got)
	}
	if got := fmtDateRange(definition.DateRange{Start: start, Current: true}); got != "Jan 2020 - Present" {
		t.Errorf("fmtDateRange current returned %q, want Jan 2020 - Present", got)
	}

	fmtLocation := funcs["fmtLocation"].(func(*definition.Location) string)
	if got := fmtLocation(nil); got != "Remote" {
		t.Errorf("fmtLocation(nil) = %q, want Remote", got)
	}
	if got := fmtLocation(&definition.Location{City: "Paris", State: "Île-de-France"}); got != "Paris, Île-de-France" {
		t.Errorf("fmtLocation returned %q, want Paris, Île-de-France", got)
	}

	skillNames := funcs["skillNames"].(func([]definition.SkillItem) []string)
	names := skillNames([]definition.SkillItem{{Name: "Go"}, {Name: "Rust"}})
	if strings.Join(names, ",") != "Go,Rust" {
		t.Errorf("skillNames returned %v, want [Go Rust]", names)
	}

	fmtLink := funcs["fmtLink"].(func(interface{}) string)
	if got := fmtLink(definition.Link{URL: "https://example.com", Text: "Example"}); got != `\href{https://example.com}{Example}` {
		t.Errorf("fmtLink returned %q, want \\href{https://example.com}{Example}", got)
	}

	fmtDates := funcs["fmtDates"].(func(interface{}) string)
	if got := fmtDates("Jan 2020"); got != "Jan 2020" {
		t.Errorf("fmtDates string returned %q, want Jan 2020", got)
	}
	if got := fmtDates(definition.DateRange{Start: start, End: &end}); got != "Jan 2020 - Jun 2021" {
		t.Errorf("fmtDates range returned %q, want Jan 2020 - Jun 2021", got)
	}

	lower := funcs["lower"].(func(string) string)
	if got := lower("GoLang"); got != "golang" {
		t.Errorf("lower returned %q, want golang", got)
	}
}

func TestLaTeXGeneratorGenerate(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewLaTeXGenerator(logger)

	start := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	resume := &definition.Resume{
		Contact: definition.Contact{
			Name: "John & Co.",
			Links: []definition.Link{
				{URL: "https://example.com", Text: "Website"},
			},
		},
		Skills: definition.Skills{
			Categories: []definition.SkillCategory{
				{
					Items: []definition.SkillItem{
						{Name: "Go"},
						{Name: "Rust"},
					},
				},
			},
		},
		Experience: definition.ExperienceList{
			Positions: []definition.Experience{
				{
					Location: &definition.Location{City: "Berlin"},
					Dates: definition.DateRange{
						Start:   start,
						Current: true,
					},
				},
			},
		},
	}

	templateContent := `
Name: {{escape .Contact.Name}}
Dates: {{fmtDateRange (index .Experience.Positions 0).Dates}}
Location: {{fmtLocation (index .Experience.Positions 0).Location}}
Link: {{fmtLink (index .Contact.Links 0)}}
Skills: {{join ", " (skillNames (index .Skills.Categories 0).Items)}}
`

	got, err := gen.Generate(templateContent, resume)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(got, `John \& Co.`) {
		t.Errorf("Generate() output missing escaped name: %q", got)
	}
	if !strings.Contains(got, "Jan 2022 - Present") {
		t.Errorf("Generate() output missing date range: %q", got)
	}
	if !strings.Contains(got, "Berlin") {
		t.Errorf("Generate() output missing location: %q", got)
	}
	if !strings.Contains(got, `\href{https://example.com}{Website}`) {
		t.Errorf("Generate() output missing link: %q", got)
	}
	if !strings.Contains(got, "Go, Rust") {
		t.Errorf("Generate() output missing skills: %q", got)
	}
}
