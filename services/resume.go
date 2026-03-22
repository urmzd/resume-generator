package services

import (
	"encoding/json"
	"fmt"

	"github.com/urmzd/incipit/resume"
	"github.com/urmzd/incipit/utils"
)

// ResumeData holds the loaded and validated resume alongside its metadata.
type ResumeData struct {
	Resume       *resume.Resume
	SectionOrder []string
	Format       string
}

// PreviewData contains structured resume information for display.
type PreviewData struct {
	ContactName     string
	ContactEmail    string
	ContactPhone    string
	ContactLocation string
	LinksCount      int
	SkillCategories []SkillCategorySummary
	TotalSkills     int
	Positions       []PositionSummary
	Projects        []string
	Institutions    []InstitutionSummary
	Format          string
	JSON            string
}

// SkillCategorySummary is a summary of a skill category.
type SkillCategorySummary struct {
	Category string
	Count    int
}

// PositionSummary is a summary of a work position.
type PositionSummary struct {
	Title   string
	Company string
}

// InstitutionSummary is a summary of an education entry.
type InstitutionSummary struct {
	Degree      string
	Institution string
}

// LoadResume loads and validates a resume from a file path, returning
// the runtime Resume, section order, and detected format.
func LoadResume(filePath string) (*ResumeData, error) {
	resolved, err := utils.ResolvePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("error resolving path: %w", err)
	}
	if !utils.FileExists(resolved) {
		return nil, fmt.Errorf("file does not exist: %s", resolved)
	}

	inputData, err := resume.LoadResumeFromFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("error loading resume data: %w", err)
	}

	if err := inputData.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return &ResumeData{
		Resume:       inputData.ToResume(),
		SectionOrder: inputData.GetSectionOrder(),
		Format:       inputData.GetFormat(),
	}, nil
}

// ValidationError mirrors resume.ValidationError for public consumption.
type ValidationError struct {
	Field   string
	Message string
}

// ValidateResume loads a resume and returns structured validation errors.
// Returns nil if the resume is valid.
func ValidateResume(filePath string) ([]ValidationError, error) {
	resolved, err := utils.ResolvePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("error resolving path: %w", err)
	}
	if !utils.FileExists(resolved) {
		return nil, fmt.Errorf("file does not exist: %s", resolved)
	}

	inputData, err := resume.LoadResumeFromFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("failed to load resume data: %w", err)
	}

	resumeData := inputData.ToResume()
	errors := resume.Validate(resumeData)
	if len(errors) == 0 {
		return nil, nil
	}

	result := make([]ValidationError, len(errors))
	for i, e := range errors {
		result[i] = ValidationError{Field: e.Field, Message: e.Message}
	}
	return result, nil
}

// PreviewResume loads a resume and returns structured preview data.
func PreviewResume(filePath string, includeJSON bool) (*PreviewData, error) {
	data, err := LoadResume(filePath)
	if err != nil {
		return nil, err
	}

	r := data.Resume
	preview := &PreviewData{
		ContactName:  r.Contact.Name,
		ContactEmail: r.Contact.Email,
		ContactPhone: r.Contact.Phone,
		Format:       data.Format,
	}

	if r.Contact.Location != nil {
		preview.ContactLocation = fmt.Sprintf("%s, %s", r.Contact.Location.City, r.Contact.Location.State)
	}
	preview.LinksCount = len(r.Contact.Links)

	for _, cat := range r.Skills.Categories {
		preview.SkillCategories = append(preview.SkillCategories, SkillCategorySummary{
			Category: cat.Category,
			Count:    len(cat.Items),
		})
		preview.TotalSkills += len(cat.Items)
	}

	for _, exp := range r.Experience.Positions {
		preview.Positions = append(preview.Positions, PositionSummary{
			Title:   exp.Title,
			Company: exp.Company,
		})
	}

	if r.Projects != nil {
		for _, proj := range r.Projects.Projects {
			preview.Projects = append(preview.Projects, proj.Name)
		}
	}

	for _, edu := range r.Education.Institutions {
		preview.Institutions = append(preview.Institutions, InstitutionSummary{
			Degree:      edu.Degree.Name,
			Institution: edu.Institution,
		})
	}

	if includeJSON {
		jsonData, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resume to JSON: %w", err)
		}
		preview.JSON = string(jsonData)
	}

	return preview, nil
}
