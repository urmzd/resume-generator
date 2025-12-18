package generators

import (
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/definition"
	"go.uber.org/zap"
)

func TestBuildHeaderView(t *testing.T) {
	formatter := newHTMLFormatter()

	resume := &definition.Resume{
		Contact: definition.Contact{
			Name:     " Alice Example ",
			Title:    " Senior Engineer ",
			Email:    "alice@example.com",
			Phone:    "+1 (555) 123-4567",
			Website:  " https://alice.dev ",
			Location: &definition.Location{City: "New York", State: "NY", Country: "USA"},
			Links: []definition.Link{
				{Order: 2, Text: "", URL: "https://linkedin.com/in/alice"},
				{Order: 1, Text: "GitHub", URL: "https://github.com/alice"},
			},
			Visibility: definition.VisibilityConfig{
				ShowEmail:    true,
				ShowPhone:    true,
				ShowLocation: true,
			},
		},
	}

	header := buildHeaderView(resume, formatter)

	if header.Name != "Alice Example" {
		t.Fatalf("Name = %q, want %q", header.Name, "Alice Example")
	}
	if header.Title != "Senior Engineer" {
		t.Fatalf("Title = %q, want %q", header.Title, "Senior Engineer")
	}

	wantItems := []htmlInlineItem{
		{Text: "alice@example.com", URL: "mailto:alice@example.com"},
		{Text: "+1 (555) 123-4567", URL: "tel:+15551234567"},
		{Text: "New York, NY, USA"},
		{Text: "https://alice.dev", URL: "https://alice.dev"},
		{Text: "GitHub", URL: "https://github.com/alice"},
		{Text: "https://linkedin.com/in/alice", URL: "https://linkedin.com/in/alice"},
	}

	if len(header.ContactItems) != len(wantItems) {
		t.Fatalf("ContactItems len = %d, want %d", len(header.ContactItems), len(wantItems))
	}

	for i, item := range header.ContactItems {
		if item != wantItems[i] {
			t.Errorf("ContactItems[%d] = %+v, want %+v", i, item, wantItems[i])
		}
	}

	if !header.HasContactRow {
		t.Error("HasContactRow = false, want true")
	}
}

func TestBuildSkillsView(t *testing.T) {
	formatter := newHTMLFormatter()

	resume := &definition.Resume{
		Skills: definition.Skills{
			Title: "Capabilities",
			Categories: []definition.SkillCategory{
				{
					Order: 2,
					Name:  "Languages",
					Items: []definition.SkillItem{
						{Order: 1, Name: "Go"},
						{Order: 3, Name: " "},
						{Order: 2, Name: " Python "},
					},
				},
				{
					Order: 1,
					Name:  "Tools",
					Items: []definition.SkillItem{
						{Order: 2, Name: "Docker"},
						{Order: 1, Name: "Git"},
					},
				},
			},
		},
	}

	view := buildSkillsView(resume, formatter)
	if view == nil {
		t.Fatal("buildSkillsView returned nil, want view")
	}

	if view.Title != "Capabilities" {
		t.Fatalf("view.Title = %q, want %q", view.Title, "Capabilities")
	}

	if len(view.Categories) != 2 {
		t.Fatalf("Categories len = %d, want 2", len(view.Categories))
	}

	if view.Categories[0].Name != "Tools" || view.Categories[0].Display != "Git, Docker" {
		t.Errorf("Categories[0] = %+v, want Name=Tools Display=Git, Docker", view.Categories[0])
	}

	if view.Categories[1].Name != "Languages" || view.Categories[1].Display != "Go, Python" {
		t.Errorf("Categories[1] = %+v, want Name=Languages Display=Go, Python", view.Categories[1])
	}
}

func TestBuildExperienceView(t *testing.T) {
	formatter := newHTMLFormatter()
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	resume := &definition.Resume{
		Experience: definition.ExperienceList{
			Positions: []definition.Experience{
				{
					Order:       2,
					Company:     "Beta Corp",
					Title:       "Developer",
					Description: []string{},
					Highlights:  []string{"  Improved reliability  ", ""},
					Dates: definition.DateRange{
						Start: start,
						End:   &end,
					},
				},
				{
					Order:       1,
					Company:     "Acme Inc.",
					Title:       "Senior Developer",
					Description: []string{"  Built APIs  ", ""},
					Location:    &definition.Location{City: "Seattle", State: "WA"},
					Dates: definition.DateRange{
						Start:   start,
						Current: true,
					},
				},
			},
		},
	}

	view := buildExperienceView(resume, formatter)
	if view == nil {
		t.Fatal("buildExperienceView returned nil, want view")
	}

	if view.Title != "Experience" {
		t.Fatalf("view.Title = %q, want Experience", view.Title)
	}

	if len(view.Positions) != 2 {
		t.Fatalf("Positions len = %d, want 2", len(view.Positions))
	}

	first := view.Positions[0]
	if first.Company != "Acme Inc." {
		t.Errorf("first.Company = %q, want Acme Inc.", first.Company)
	}
	if first.DateRange != "Jan 2020 – Present" {
		t.Errorf("first.DateRange = %q, want Jan 2020 – Present", first.DateRange)
	}
	if len(first.Highlights) != 1 || first.Highlights[0] != "Built APIs" {
		t.Errorf("first.Highlights = %v, want [Built APIs]", first.Highlights)
	}
	if first.Location != "Seattle, WA" {
		t.Errorf("first.Location = %q, want Seattle, WA", first.Location)
	}

	second := view.Positions[1]
	if second.DateRange != "Jan 2020 – Jun 2021" {
		t.Errorf("second.DateRange = %q, want Jan 2020 – Jun 2021", second.DateRange)
	}
	if len(second.Highlights) != 1 || second.Highlights[0] != "Improved reliability" {
		t.Errorf("second.Highlights = %v, want [Improved reliability]", second.Highlights)
	}
}

func TestFormatDateRange(t *testing.T) {
	formatter := newHTMLFormatter()
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	mid := time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		start   time.Time
		end     *time.Time
		current bool
		want    string
	}{
		{"empty", time.Time{}, nil, false, ""},
		{"start only", start, nil, false, "Jan 2020 – Present"},
		{"current true", start, nil, true, "Jan 2020 – Present"},
		{"range", start, &mid, false, "Jan 2020 – Mar 2021"},
		{"same month", start, &start, false, "Jan 2020"},
		{"no start", time.Time{}, &mid, false, "Mar 2021"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.formatDateRange(tt.start, tt.end, tt.current)
			if got != tt.want {
				t.Errorf("formatDateRange() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTMLGeneratorGenerateStandalone(t *testing.T) {
	logger := zap.NewNop().Sugar()
	gen := NewHTMLGenerator(logger)

	resume := &definition.Resume{
		Contact: definition.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Visibility: definition.VisibilityConfig{
				ShowEmail: true,
			},
		},
	}

	templateContent := `<div id="content">
<h1>{{.View.Header.Name}}</h1>
{{range .View.Header.ContactItems}}<span>{{.Text}}</span>{{end}}
</div>`
	cssContent := "body { color: #333; }"

	got, err := gen.GenerateStandalone(templateContent, cssContent, resume)
	if err != nil {
		t.Fatalf("GenerateStandalone() error = %v", err)
	}

	if !strings.Contains(got, "<!DOCTYPE html>") {
		t.Error("GenerateStandalone() missing <!DOCTYPE html>")
	}
	if !strings.Contains(got, cssContent) {
		t.Errorf("GenerateStandalone() missing CSS content %q", cssContent)
	}
	if !strings.Contains(got, "Jane Doe") {
		t.Error("GenerateStandalone() missing rendered resume content")
	}
	if !strings.Contains(got, "jane@example.com") {
		t.Error("GenerateStandalone() missing email address")
	}
}

func TestFormatGPAValue(t *testing.T) {
	formatter := newHTMLFormatter()

	tests := []struct {
		name string
		gpa  string
		max  string
		want string
	}{
		{"empty gpa", "", "4.0", ""},
		{"default max", "3.9", "4.0", "3.9"},
		{"custom max", "3.8", "5.0", "3.8 / 5.0"},
		{"trimmed", " 3.7 ", " 5.0 ", "3.7 / 5.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatter.FormatGPA(tt.gpa, tt.max)
			if got != tt.want {
				t.Errorf("formatGPAValue(%q, %q) = %q, want %q", tt.gpa, tt.max, got, tt.want)
			}
		})
	}
}
