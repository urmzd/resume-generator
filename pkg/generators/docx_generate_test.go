package generators

import (
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

func TestDOCXGenerate(t *testing.T) {
	logger := zap.NewNop().Sugar()

	expStart := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	eduStart := time.Date(2018, time.September, 1, 0, 0, 0, 0, time.UTC)
	eduEnd := time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)

	minimalResume := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		},
	}

	fullResume := &resume.Resume{
		Contact: resume.Contact{
			Name:  "Jane Doe",
			Email: "jane@example.com",
			Phone: "+1 (555) 123-4567",
			Location: &resume.Location{
				City:    "NYC",
				State:   "NY",
				Country: "USA",
			},
			Links: []resume.Link{
				{URI: "https://github.com/janedoe"},
			},
		},
		Summary: "Experienced software engineer.",
		Certifications: &resume.Certifications{
			Items: []resume.Certification{
				{Name: "AWS Certified", Issuer: "Amazon", Notes: "2023"},
			},
		},
		Skills: resume.Skills{
			Categories: []resume.SkillCategory{
				{Category: "Languages", Items: []string{"Go", "Rust"}},
			},
		},
		Experience: resume.ExperienceList{
			Positions: []resume.Experience{
				{
					Title:      "Senior Engineer",
					Company:    "Acme Corp",
					Highlights: []string{"Led team of 5"},
					Dates:      resume.DateRange{Start: expStart},
					Location:   &resume.Location{City: "NYC"},
				},
			},
		},
		Education: resume.EducationList{
			Institutions: []resume.Education{
				{
					Institution: "MIT",
					Degree:      resume.Degree{Name: "B.Sc CS", Descriptions: []string{"Dean's List"}},
					GPA:         &resume.GPA{GPA: "3.9", MaxGPA: "4.0"},
					Dates:       resume.DateRange{Start: eduStart, End: &eduEnd},
				},
			},
		},
		Projects: &resume.ProjectList{
			Projects: []resume.Project{
				{
					Name:       "Tool",
					Highlights: []string{"Used by 1000+ devs"},
					Link:       resume.Link{URI: "https://github.com/tool"},
				},
			},
		},
		Languages: &resume.LanguageList{
			Languages: []resume.Language{
				{Name: "English", Proficiency: "Native"},
			},
		},
	}

	tests := []struct {
		name   string
		resume *resume.Resume
	}{
		{"minimal resume", minimalResume},
		{"full resume", fullResume},
		{"with layout section order", &resume.Resume{
			Contact: resume.Contact{Name: "Test", Email: "t@t.com"},
			Layout: &resume.Layout{
				Sections: []string{"experience", "education", "skills"},
			},
			Skills: resume.Skills{
				Categories: []resume.SkillCategory{{Category: "Lang", Items: []string{"Go"}}},
			},
			Experience: resume.ExperienceList{
				Positions: []resume.Experience{{Title: "Dev", Company: "Co", Dates: resume.DateRange{Start: expStart}}},
			},
			Education: resume.EducationList{
				Institutions: []resume.Education{{Institution: "Uni", Dates: resume.DateRange{Start: eduStart, End: &eduEnd}}},
			},
		}},
		{"with references", &resume.Resume{
			Contact: resume.Contact{Name: "Test", Email: "t@t.com"},
			Layout:  &resume.Layout{References: true},
		}},
		{"nil optional sections", &resume.Resume{
			Contact:        resume.Contact{Name: "Test", Email: "t@t.com"},
			Certifications: nil,
			Projects:       nil,
			Languages:      nil,
		}},
		{"empty lists", &resume.Resume{
			Contact:    resume.Contact{Name: "Test", Email: "t@t.com"},
			Skills:     resume.Skills{Categories: []resume.SkillCategory{}},
			Experience: resume.ExperienceList{Positions: []resume.Experience{}},
			Education:  resume.EducationList{Institutions: []resume.Education{}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewDOCXGenerator(logger)
			data, err := gen.Generate(tt.resume)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if len(data) == 0 {
				t.Error("Generate() returned empty bytes")
			}
		})
	}
}
