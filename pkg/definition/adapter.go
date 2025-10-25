package definition

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// InputData is the unified interface for all resume input formats
// This allows generators and processors to work with any input type
type InputData interface {
	// ToEnhanced converts any input format to the EnhancedResume format
	ToEnhanced() *EnhancedResume

	// GetFormat returns the original format type
	GetFormat() string

	// Validate performs basic validation on the input data
	Validate() error
}

// EnhancedResumeAdapter wraps EnhancedResume to implement InputData
type EnhancedResumeAdapter struct {
	Resume *EnhancedResume
}

func (a *EnhancedResumeAdapter) ToEnhanced() *EnhancedResume {
	return a.Resume
}

func (a *EnhancedResumeAdapter) GetFormat() string {
	return "enhanced"
}

func (a *EnhancedResumeAdapter) Validate() error {
	validator := &ConfigurationValidator{StrictMode: false}
	errors := validator.ValidateEnhancedResume(a.Resume)
	if len(errors) > 0 {
		return fmt.Errorf("validation failed with %d errors: %v", len(errors), errors[0].Message)
	}
	return nil
}

// LegacyResumeAdapter wraps Resume to implement InputData
type LegacyResumeAdapter struct {
	Resume *Resume
}

func (a *LegacyResumeAdapter) ToEnhanced() *EnhancedResume {
	migrator := &MigrationUtility{
		PreserveOrder: true,
		DefaultTheme:  "modern",
	}
	return migrator.MigrateFromLegacy(a.Resume)
}

func (a *LegacyResumeAdapter) GetFormat() string {
	return "legacy"
}

func (a *LegacyResumeAdapter) Validate() error {
	if a.Resume.Contact.Name == "" {
		return fmt.Errorf("contact name is required")
	}
	if a.Resume.Contact.Email == "" {
		return fmt.Errorf("contact email is required")
	}
	return nil
}

// JSONResumeAdapter wraps JSONResume to implement InputData
type JSONResumeAdapter struct {
	Resume *JSONResume
}

func (a *JSONResumeAdapter) ToEnhanced() *EnhancedResume {
	// Convert JSON Resume to Enhanced Resume
	enhanced := &EnhancedResume{
		Meta: ResumeMetadata{
			Version: "2.0",
			Output: OutputPreferences{
				Formats: []string{"html"},
			},
		},
		Contact: EnhancedContact{
			Order:    1,
			Name:     a.Resume.Basics.Name,
			Title:    a.Resume.Basics.Label,
			Email:    a.Resume.Basics.Email,
			Phone:    a.Resume.Basics.Phone,
			Website:  a.Resume.Basics.URL,
			Summary:  a.Resume.Basics.Summary,
			Location: &a.Resume.Basics.Location,
			Visibility: VisibilityConfig{
				ShowPhone:    true,
				ShowEmail:    true,
				ShowLocation: true,
				ShowSummary:  true,
			},
		},
	}

	// Convert profiles to links
	for i, profile := range a.Resume.Basics.Profiles {
		enhanced.Contact.Links = append(enhanced.Contact.Links, EnhancedLink{
			Order: i + 1,
			Text:  profile.Username,
			URL:   profile.URL,
			Type:  strings.ToLower(profile.Network),
		})
	}

	// Convert skills
	enhanced.Skills = EnhancedSkills{
		Order: 2,
		Title: "Skills",
	}
	for i, skill := range a.Resume.Skills {
		category := EnhancedSkillCategory{
			Order: i + 1,
			Name:  skill.Name,
		}
		for j, keyword := range skill.Keywords {
			category.Items = append(category.Items, EnhancedSkillItem{
				Order: j + 1,
				Name:  keyword,
				Level: skill.Level,
			})
		}
		enhanced.Skills.Categories = append(enhanced.Skills.Categories, category)
	}

	// Convert work experience
	enhanced.Experience = EnhancedExperienceList{
		Order: 3,
		Title: "Experience",
	}
	for i, work := range a.Resume.Work {
		exp := EnhancedExperience{
			Order:       i + 1,
			Company:     work.Name,
			Title:       work.Position,
			Description: work.Highlights,
			Website:     work.URL,
		}
		// Parse dates
		if work.StartDate != "" {
			// Simple date parsing - in production you'd want better error handling
			exp.Dates.Start, _ = parseDate(work.StartDate)
		}
		if work.EndDate != "" {
			endDate, _ := parseDate(work.EndDate)
			exp.Dates.End = &endDate
		} else {
			exp.Dates.Current = true
		}
		enhanced.Experience.Positions = append(enhanced.Experience.Positions, exp)
	}

	// Convert education
	enhanced.Education = EnhancedEducationList{
		Order: 4,
		Title: "Education",
	}
	for i, edu := range a.Resume.Education {
		institution := EnhancedEducation{
			Order:       i + 1,
			Institution: edu.Institution,
			Degree:      edu.StudyType,
			Field:       edu.Area,
			GPA:         edu.GPA,
		}
		// Parse dates
		if edu.StartDate != "" {
			institution.Dates.Start, _ = parseDate(edu.StartDate)
		}
		if edu.EndDate != "" {
			endDate, _ := parseDate(edu.EndDate)
			institution.Dates.End = &endDate
		}
		enhanced.Education.Institutions = append(enhanced.Education.Institutions, institution)
	}

	// Convert projects
	enhanced.Projects = EnhancedProjectList{
		Order: 5,
		Title: "Projects",
	}
	for i, proj := range a.Resume.Projects {
		project := EnhancedProject{
			Order:        i + 1,
			Name:         proj.Name,
			Description:  proj.Highlights,
			Technologies: proj.Keywords,
		}
		if proj.URL != "" {
			project.Links = []EnhancedLink{{
				Order: 1,
				Text:  "View Project",
				URL:   proj.URL,
				Type:  "website",
			}}
		}
		// Parse dates
		if proj.StartDate != "" {
			startDate, _ := parseDate(proj.StartDate)
			project.Dates = &EnhancedDateRange{Start: startDate}
			if proj.EndDate != "" {
				endDate, _ := parseDate(proj.EndDate)
				project.Dates.End = &endDate
			}
		}
		enhanced.Projects.Projects = append(enhanced.Projects.Projects, project)
	}

	return enhanced
}

func (a *JSONResumeAdapter) GetFormat() string {
	return "jsonresume"
}

func (a *JSONResumeAdapter) Validate() error {
	if a.Resume.Basics.Name == "" {
		return fmt.Errorf("basics.name is required")
	}
	if a.Resume.Basics.Email == "" {
		return fmt.Errorf("basics.email is required")
	}
	return nil
}

// LoadResumeFromFile automatically detects the format and returns an InputData adapter
func LoadResumeFromFile(filePath string) (InputData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fileExt := filepath.Ext(filePath)

	// Try to detect if it's a JSON Resume by looking for "basics" field
	if fileExt == ".json" {
		var jsonResume JSONResume
		if err := json.Unmarshal(data, &jsonResume); err == nil && jsonResume.Basics.Name != "" {
			return &JSONResumeAdapter{Resume: &jsonResume}, nil
		}
	}

	// Try enhanced format first
	var enhanced EnhancedResume
	switch fileExt {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &enhanced); err == nil && enhanced.Meta.Version != "" {
			return &EnhancedResumeAdapter{Resume: &enhanced}, nil
		}
	case ".json":
		if err := json.Unmarshal(data, &enhanced); err == nil && enhanced.Meta.Version != "" {
			return &EnhancedResumeAdapter{Resume: &enhanced}, nil
		}
	case ".toml":
		if _, err := toml.Decode(string(data), &enhanced); err == nil && enhanced.Meta.Version != "" {
			return &EnhancedResumeAdapter{Resume: &enhanced}, nil
		}
	}

	// Fall back to legacy format
	var legacy Resume
	switch fileExt {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &legacy); err == nil {
			return &LegacyResumeAdapter{Resume: &legacy}, nil
		}
	case ".json":
		if err := json.Unmarshal(data, &legacy); err == nil {
			return &LegacyResumeAdapter{Resume: &legacy}, nil
		}
	case ".toml":
		if _, err := toml.Decode(string(data), &legacy); err == nil {
			return &LegacyResumeAdapter{Resume: &legacy}, nil
		}
	}

	return nil, fmt.Errorf("failed to parse resume file in any supported format")
}

// Helper function to parse date strings
func parseDate(dateStr string) (time.Time, error) {
	// Try multiple date formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"2006-01",
		"2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
