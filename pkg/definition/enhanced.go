package definition

import (
	"strings"
	"time"
)

// EnhancedResume represents the new configuration format with ordering and self-contained templates
type EnhancedResume struct {
	Meta           ResumeMetadata             `json:"meta" yaml:"meta" toml:"meta"`
	Contact        EnhancedContact            `json:"contact" yaml:"contact" toml:"contact"`
	Skills         EnhancedSkills             `json:"skills" yaml:"skills" toml:"skills"`
	Experience     EnhancedExperienceList     `json:"experience" yaml:"experience" toml:"experience"`
	Projects       EnhancedProjectList        `json:"projects" yaml:"projects" toml:"projects"`
	Education      EnhancedEducationList      `json:"education" yaml:"education" toml:"education"`
	Certifications EnhancedCertificationList  `json:"certifications,omitempty" yaml:"certifications,omitempty" toml:"certifications,omitempty"`
}

// ResumeMetadata contains version info, templates, and output preferences
type ResumeMetadata struct {
	Version   string            `json:"version" yaml:"version" toml:"version"`
	Templates TemplateSet       `json:"templates,omitempty" yaml:"templates,omitempty" toml:"templates,omitempty"`
	Output    OutputPreferences `json:"output" yaml:"output" toml:"output"`
	Theme     string            `json:"theme,omitempty" yaml:"theme,omitempty" toml:"theme,omitempty"`
}

// TemplateSet contains self-contained template definitions
type TemplateSet struct {
	LaTeX string `json:"latex,omitempty" yaml:"latex,omitempty" toml:"latex,omitempty"`
	HTML  string `json:"html,omitempty" yaml:"html,omitempty" toml:"html,omitempty"`
	CSS   string `json:"css,omitempty" yaml:"css,omitempty" toml:"css,omitempty"`
}

// OutputPreferences defines output format and styling options
type OutputPreferences struct {
	Formats     []string          `json:"formats" yaml:"formats" toml:"formats"`
	Quality     string            `json:"quality,omitempty" yaml:"quality,omitempty" toml:"quality,omitempty"`
	Options     map[string]string `json:"options,omitempty" yaml:"options,omitempty" toml:"options,omitempty"`
	Destination string            `json:"destination,omitempty" yaml:"destination,omitempty" toml:"destination,omitempty"`
}

// EnhancedContact extends Contact with ordering and additional fields
type EnhancedContact struct {
	Order       int              `json:"order" yaml:"order" toml:"order"`
	Name        string           `json:"name" yaml:"name" toml:"name"`
	Title       string           `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Email       string           `json:"email" yaml:"email" toml:"email"`
	Phone       string           `json:"phone,omitempty" yaml:"phone,omitempty" toml:"phone,omitempty"`
	Website     string           `json:"website,omitempty" yaml:"website,omitempty" toml:"website,omitempty"`
	Location    *Location        `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Links       []EnhancedLink   `json:"links,omitempty" yaml:"links,omitempty" toml:"links,omitempty"`
	Summary     string           `json:"summary,omitempty" yaml:"summary,omitempty" toml:"summary,omitempty"`
	Visibility  VisibilityConfig `json:"visibility,omitempty" yaml:"visibility,omitempty" toml:"visibility,omitempty"`
}

// EnhancedLink extends Link with ordering and metadata
type EnhancedLink struct {
	Order       int    `json:"order" yaml:"order" toml:"order"`
	Text        string `json:"text" yaml:"text" toml:"text"`
	URL         string `json:"url" yaml:"url" toml:"url"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty" toml:"type,omitempty"` // linkedin, github, portfolio, etc.
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty" toml:"icon,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

// VisibilityConfig controls what information to show/hide
type VisibilityConfig struct {
	ShowPhone    bool `json:"show_phone" yaml:"show_phone" toml:"show_phone"`
	ShowEmail    bool `json:"show_email" yaml:"show_email" toml:"show_email"`
	ShowLocation bool `json:"show_location" yaml:"show_location" toml:"show_location"`
	ShowSummary  bool `json:"show_summary" yaml:"show_summary" toml:"show_summary"`
}

// EnhancedSkills with categorization and ordering
type EnhancedSkills struct {
	Order      int                    `json:"order" yaml:"order" toml:"order"`
	Title      string                 `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Categories []EnhancedSkillCategory `json:"categories" yaml:"categories" toml:"categories"`
	Layout     string                 `json:"layout,omitempty" yaml:"layout,omitempty" toml:"layout,omitempty"` // grid, list, compact
}

// EnhancedSkillCategory represents a group of related skills
type EnhancedSkillCategory struct {
	Order       int                `json:"order" yaml:"order" toml:"order"`
	Name        string             `json:"name" yaml:"name" toml:"name"`
	Items       []EnhancedSkillItem `json:"items" yaml:"items" toml:"items"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Icon        string             `json:"icon,omitempty" yaml:"icon,omitempty" toml:"icon,omitempty"`
}

// EnhancedSkillItem represents an individual skill with proficiency
type EnhancedSkillItem struct {
	Order              int      `json:"order" yaml:"order" toml:"order"`
	Name               string   `json:"name" yaml:"name" toml:"name"`
	Level              string   `json:"level,omitempty" yaml:"level,omitempty" toml:"level,omitempty"` // beginner, intermediate, advanced, expert
	Years              int      `json:"years,omitempty" yaml:"years,omitempty" toml:"years,omitempty"`
	YearsOfExperience  int      `json:"yearsOfExperience,omitempty" yaml:"yearsOfExperience,omitempty" toml:"yearsOfExperience,omitempty"`
	Certification      string   `json:"certification,omitempty" yaml:"certification,omitempty" toml:"certification,omitempty"`
	Keywords           []string `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
}

// EnhancedExperienceList with ordering and grouping
type EnhancedExperienceList struct {
	Order     int                     `json:"order" yaml:"order" toml:"order"`
	Title     string                  `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Positions []EnhancedExperience    `json:"positions" yaml:"positions" toml:"positions"`
	Groups    []EnhancedExperienceGroup `json:"groups,omitempty" yaml:"groups,omitempty" toml:"groups,omitempty"`
}

// EnhancedExperience extends Experience with ordering and enhanced metadata
type EnhancedExperience struct {
	Order          int                   `json:"order" yaml:"order" toml:"order"`
	Company        string                `json:"company" yaml:"company" toml:"company"`
	Title          string                `json:"title" yaml:"title" toml:"title"`
	Department     string                `json:"department,omitempty" yaml:"department,omitempty" toml:"department,omitempty"`
	Type           string                `json:"type,omitempty" yaml:"type,omitempty" toml:"type,omitempty"` // full-time, part-time, contract, internship
	EmploymentType string                `json:"employment_type,omitempty" yaml:"employment_type,omitempty" toml:"employment_type,omitempty"` // alias for Type
	Description    []string              `json:"description" yaml:"description" toml:"description"`
	Highlights     []string              `json:"highlights,omitempty" yaml:"highlights,omitempty" toml:"highlights,omitempty"`
	Achievements   []EnhancedAchievement `json:"achievements,omitempty" yaml:"achievements,omitempty" toml:"achievements,omitempty"`
	Technologies   []string              `json:"technologies,omitempty" yaml:"technologies,omitempty" toml:"technologies,omitempty"`
	Dates          EnhancedDateRange     `json:"dates" yaml:"dates" toml:"dates"`
	Location       *Location             `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Website        string                `json:"website,omitempty" yaml:"website,omitempty" toml:"website,omitempty"`
	Manager        string                `json:"manager,omitempty" yaml:"manager,omitempty" toml:"manager,omitempty"`
	TeamSize       int                   `json:"team_size,omitempty" yaml:"team_size,omitempty" toml:"team_size,omitempty"`
	Keywords       []string              `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
	Metrics        []EnhancedMetric      `json:"metrics,omitempty" yaml:"metrics,omitempty" toml:"metrics,omitempty"`
}

// EnhancedExperienceGroup allows grouping experiences by theme
type EnhancedExperienceGroup struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Positions   []int  `json:"positions" yaml:"positions" toml:"positions"` // References to position orders
}

// EnhancedAchievement represents a notable accomplishment
type EnhancedAchievement struct {
	Order       int              `json:"order" yaml:"order" toml:"order"`
	Title       string           `json:"title" yaml:"title" toml:"title"`
	Description string           `json:"description" yaml:"description" toml:"description"`
	Impact      string           `json:"impact,omitempty" yaml:"impact,omitempty" toml:"impact,omitempty"`
	Metrics     []EnhancedMetric `json:"metrics,omitempty" yaml:"metrics,omitempty" toml:"metrics,omitempty"`
	Date        *time.Time       `json:"date,omitempty" yaml:"date,omitempty" toml:"date,omitempty"`
}

// EnhancedMetric represents quantifiable achievements
type EnhancedMetric struct {
	Name        string  `json:"name" yaml:"name" toml:"name"`
	Value       float64 `json:"value" yaml:"value" toml:"value"`
	Unit        string  `json:"unit" yaml:"unit" toml:"unit"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Improvement bool    `json:"improvement,omitempty" yaml:"improvement,omitempty" toml:"improvement,omitempty"`
}

// EnhancedDateRange extends DateRange with more precision
type EnhancedDateRange struct {
	Start   time.Time  `json:"start" yaml:"start" toml:"start"`
	End     *time.Time `json:"end,omitempty" yaml:"end,omitempty" toml:"end,omitempty"`
	Current bool       `json:"current,omitempty" yaml:"current,omitempty" toml:"current,omitempty"`
	Duration string    `json:"duration,omitempty" yaml:"duration,omitempty" toml:"duration,omitempty"`
}

// EnhancedProjectList with ordering and categorization
type EnhancedProjectList struct {
	Order      int                       `json:"order" yaml:"order" toml:"order"`
	Title      string                    `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Projects   []EnhancedProject         `json:"projects" yaml:"projects" toml:"projects"`
	Categories []EnhancedProjectCategory `json:"categories,omitempty" yaml:"categories,omitempty" toml:"categories,omitempty"`
}

// EnhancedProject extends Project with enhanced metadata
type EnhancedProject struct {
	Order        int                 `json:"order" yaml:"order" toml:"order"`
	Name         string              `json:"name" yaml:"name" toml:"name"`
	Category     string              `json:"category,omitempty" yaml:"category,omitempty" toml:"category,omitempty"`
	Type         string              `json:"type,omitempty" yaml:"type,omitempty" toml:"type,omitempty"` // personal, professional, open-source
	Status       string              `json:"status,omitempty" yaml:"status,omitempty" toml:"status,omitempty"` // completed, in-progress, maintenance
	Description  []string            `json:"description" yaml:"description" toml:"description"`
	Technologies []string            `json:"technologies" yaml:"technologies" toml:"technologies"`
	Links        []EnhancedLink      `json:"links,omitempty" yaml:"links,omitempty" toml:"links,omitempty"`
	Dates        *EnhancedDateRange  `json:"dates,omitempty" yaml:"dates,omitempty" toml:"dates,omitempty"`
	Role         string              `json:"role,omitempty" yaml:"role,omitempty" toml:"role,omitempty"`
	TeamSize     int                 `json:"team_size,omitempty" yaml:"team_size,omitempty" toml:"team_size,omitempty"`
	Achievements []EnhancedAchievement `json:"achievements,omitempty" yaml:"achievements,omitempty" toml:"achievements,omitempty"`
	Keywords     []string            `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
	Images       []string            `json:"images,omitempty" yaml:"images,omitempty" toml:"images,omitempty"`
}

// EnhancedProjectCategory for grouping projects
type EnhancedProjectCategory struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty" toml:"icon,omitempty"`
	Projects    []int  `json:"projects" yaml:"projects" toml:"projects"` // References to project orders
}

// EnhancedEducationList with ordering and additional metadata
type EnhancedEducationList struct {
	Order       int                 `json:"order" yaml:"order" toml:"order"`
	Title       string              `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Institutions []EnhancedEducation `json:"institutions" yaml:"institutions" toml:"institutions"`
}

// EnhancedEducation extends Education with enhanced metadata
type EnhancedEducation struct {
	Order        int                     `json:"order" yaml:"order" toml:"order"`
	Institution  string                  `json:"institution" yaml:"institution" toml:"institution"`
	Degree       string                  `json:"degree" yaml:"degree" toml:"degree"`
	Field        string                  `json:"field,omitempty" yaml:"field,omitempty" toml:"field,omitempty"`
	Level        string                  `json:"level,omitempty" yaml:"level,omitempty" toml:"level,omitempty"` // bachelor, master, phd, certificate
	GPA          string                  `json:"gpa,omitempty" yaml:"gpa,omitempty" toml:"gpa,omitempty"`
	MaxGPA       string                  `json:"max_gpa,omitempty" yaml:"max_gpa,omitempty" toml:"max_gpa,omitempty"`
	Honors       []string                `json:"honors,omitempty" yaml:"honors,omitempty" toml:"honors,omitempty"`
	Description  []CategoryValuePair     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Coursework   []EnhancedCoursework    `json:"coursework,omitempty" yaml:"coursework,omitempty" toml:"coursework,omitempty"`
	Projects     []string                `json:"projects,omitempty" yaml:"projects,omitempty" toml:"projects,omitempty"`
	Activities   []string                `json:"activities,omitempty" yaml:"activities,omitempty" toml:"activities,omitempty"`
	Location     *Location               `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Dates        EnhancedDateRange       `json:"dates" yaml:"dates" toml:"dates"`
	Website      string                  `json:"website,omitempty" yaml:"website,omitempty" toml:"website,omitempty"`
	Thesis       *EnhancedThesis         `json:"thesis,omitempty" yaml:"thesis,omitempty" toml:"thesis,omitempty"`
	Keywords     []string                `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
}

// EnhancedCoursework represents relevant coursework
type EnhancedCoursework struct {
	Order       int    `json:"order" yaml:"order" toml:"order"`
	Name        string `json:"name" yaml:"name" toml:"name"`
	Code        string `json:"code,omitempty" yaml:"code,omitempty" toml:"code,omitempty"`
	Grade       string `json:"grade,omitempty" yaml:"grade,omitempty" toml:"grade,omitempty"`
	Credits     int    `json:"credits,omitempty" yaml:"credits,omitempty" toml:"credits,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

// EnhancedThesis represents thesis or dissertation information
type EnhancedThesis struct {
	Title       string   `json:"title" yaml:"title" toml:"title"`
	Advisor     string   `json:"advisor,omitempty" yaml:"advisor,omitempty" toml:"advisor,omitempty"`
	Committee   []string `json:"committee,omitempty" yaml:"committee,omitempty" toml:"committee,omitempty"`
	Abstract    string   `json:"abstract,omitempty" yaml:"abstract,omitempty" toml:"abstract,omitempty"`
	Keywords    []string `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
	Publication string   `json:"publication,omitempty" yaml:"publication,omitempty" toml:"publication,omitempty"`
	DOI         string   `json:"doi,omitempty" yaml:"doi,omitempty" toml:"doi,omitempty"`
}

// EnhancedCertificationList with ordering
type EnhancedCertificationList struct {
	Order          int                      `json:"order" yaml:"order" toml:"order"`
	Title          string                   `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Certifications []EnhancedCertification  `json:"certifications" yaml:"certifications" toml:"certifications"`
}

// EnhancedCertification represents a professional certification or credential
type EnhancedCertification struct {
	Order           int        `json:"order" yaml:"order" toml:"order"`
	Name            string     `json:"name" yaml:"name" toml:"name"`
	Issuer          string     `json:"issuer" yaml:"issuer" toml:"issuer"`
	IssueDate       *time.Time `json:"issueDate,omitempty" yaml:"issueDate,omitempty" toml:"issueDate,omitempty"`
	ExpirationDate  *time.Time `json:"expirationDate,omitempty" yaml:"expirationDate,omitempty" toml:"expirationDate,omitempty"`
	CredentialID    string     `json:"credentialId,omitempty" yaml:"credentialId,omitempty" toml:"credentialId,omitempty"`
	CredentialURL   string     `json:"credentialUrl,omitempty" yaml:"credentialUrl,omitempty" toml:"credentialUrl,omitempty"`
	VerificationURL string     `json:"verificationUrl,omitempty" yaml:"verificationUrl,omitempty" toml:"verificationUrl,omitempty"`
	Description     string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Skills          []string   `json:"skills,omitempty" yaml:"skills,omitempty" toml:"skills,omitempty"`
}

// ConfigurationValidator provides validation for enhanced configurations
type ConfigurationValidator struct {
	StrictMode bool
	MinVersion string
}

// ValidateEnhancedResume validates the enhanced resume configuration
func (cv *ConfigurationValidator) ValidateEnhancedResume(resume *EnhancedResume) []ValidationError {
	var errors []ValidationError

	// Validate metadata
	if resume.Meta.Version == "" {
		errors = append(errors, ValidationError{
			Field:   "meta.version",
			Message: "Version is required",
			Type:    "required",
		})
	}

	// Validate output formats
	if len(resume.Meta.Output.Formats) == 0 {
		errors = append(errors, ValidationError{
			Field:   "meta.output.formats",
			Message: "At least one output format must be specified",
			Type:    "required",
		})
	}

	// Validate ordering (no duplicates)
	errors = append(errors, cv.validateOrdering(resume)...)

	// Validate contact information
	if resume.Contact.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "contact.name",
			Message: "Name is required",
			Type:    "required",
		})
	}

	if resume.Contact.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "contact.email",
			Message: "Email is required",
			Type:    "required",
		})
	}

	return errors
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Value   interface{} `json:"value,omitempty"`
}

// validateOrdering checks for duplicate order values within sections
func (cv *ConfigurationValidator) validateOrdering(resume *EnhancedResume) []ValidationError {
	var errors []ValidationError

	// Check skill category ordering
	categoryOrders := make(map[int]bool)
	for _, category := range resume.Skills.Categories {
		if categoryOrders[category.Order] {
			errors = append(errors, ValidationError{
				Field:   "skills.categories.order",
				Message: "Duplicate order value found in skill categories",
				Type:    "duplicate",
				Value:   category.Order,
			})
		}
		categoryOrders[category.Order] = true
	}

	// Check experience ordering
	expOrders := make(map[int]bool)
	for _, exp := range resume.Experience.Positions {
		if expOrders[exp.Order] {
			errors = append(errors, ValidationError{
				Field:   "experience.positions.order",
				Message: "Duplicate order value found in experience positions",
				Type:    "duplicate",
				Value:   exp.Order,
			})
		}
		expOrders[exp.Order] = true
	}

	return errors
}

// MigrationUtility provides utilities for migrating from old to new format
type MigrationUtility struct {
	PreserveOrder bool
	DefaultTheme  string
}

// MigrateFromLegacy converts a legacy Resume to EnhancedResume
func (mu *MigrationUtility) MigrateFromLegacy(legacy *Resume) *EnhancedResume {
	enhanced := &EnhancedResume{
		Meta: ResumeMetadata{
			Version: "2.0",
			Output: OutputPreferences{
				Formats: []string{"pdf"},
			},
			Theme: mu.DefaultTheme,
		},
		Contact: EnhancedContact{
			Order: 1,
			Name:  legacy.Contact.Name,
			Email: legacy.Contact.Email,
			Phone: legacy.Contact.Phone,
			Visibility: VisibilityConfig{
				ShowPhone:    true,
				ShowEmail:    true,
				ShowLocation: true,
				ShowSummary:  true,
			},
		},
	}

	// Migrate links
	for i, link := range legacy.Contact.Links {
		enhanced.Contact.Links = append(enhanced.Contact.Links, EnhancedLink{
			Order: i + 1,
			Text:  link.Text,
			URL:   link.Ref,
		})
	}

	// Migrate skills
	enhanced.Skills = EnhancedSkills{
		Order: 2,
		Title: "Skills",
	}

	for i, skill := range legacy.Skills {
		category := EnhancedSkillCategory{
			Order: i + 1,
			Name:  skill.Category,
		}

		// Parse comma-separated values
		skillItems := parseCommaSeparated(skill.Value)
		for j, item := range skillItems {
			category.Items = append(category.Items, EnhancedSkillItem{
				Order: j + 1,
				Name:  item,
			})
		}

		enhanced.Skills.Categories = append(enhanced.Skills.Categories, category)
	}

	// Migrate experience
	enhanced.Experience = EnhancedExperienceList{
		Order: 3,
		Title: "Experience",
	}

	for i, exp := range legacy.Experience {
		// Use achievements if available, otherwise use description
		highlights := exp.Description
		if len(exp.Achievements) > 0 {
			highlights = exp.Achievements
		}

		enhanced.Experience.Positions = append(enhanced.Experience.Positions, EnhancedExperience{
			Order:       i + 1,
			Company:     exp.Company,
			Title:       exp.Title,
			Description: exp.Description,
			Highlights:  highlights,
			Dates: EnhancedDateRange{
				Start:   exp.Dates.Start,
				End:     exp.Dates.End,
				Current: exp.Dates.End == nil,
			},
			Location: exp.Location,
		})
	}

	// Migrate projects
	enhanced.Projects = EnhancedProjectList{
		Order: 4,
		Title: "Projects",
	}

	for i, proj := range legacy.Projects {
		enhanced.Projects.Projects = append(enhanced.Projects.Projects, EnhancedProject{
			Order:        i + 1,
			Name:         proj.Name,
			Description:  proj.Description,
			Technologies: parseCommaSeparated(proj.Language),
			Links: []EnhancedLink{{
				Order: 1,
				Text:  "Repository",
				URL:   proj.Link.Ref,
				Type:  "github",
			}},
		})
	}

	// Migrate education
	enhanced.Education = EnhancedEducationList{
		Order: 5,
		Title: "Education",
	}

	for i, edu := range legacy.Education {
		enhanced.Education.Institutions = append(enhanced.Education.Institutions, EnhancedEducation{
			Order:       i + 1,
			Institution: edu.School,
			Degree:      edu.Degree,
			Honors:      edu.Suffixes,
			Description: edu.Description,
			Location:    &edu.Location,
			Dates: EnhancedDateRange{
				Start: edu.Dates.Start,
				End:   edu.Dates.End,
			},
		})
	}

	return enhanced
}

// MigrateToLegacy converts an EnhancedResume to a legacy Resume
func (er *EnhancedResume) ToLegacy() *Resume {
	legacy := &Resume{}

	// Contact
	legacy.Contact.Name = er.Contact.Name
	legacy.Contact.Email = er.Contact.Email
	legacy.Contact.Phone = er.Contact.Phone
	for _, l := range er.Contact.Links {
		legacy.Contact.Links = append(legacy.Contact.Links, Link{Text: l.Text, Ref: l.URL})
	}

	// Skills
	for _, c := range er.Skills.Categories {
		var items []string
		for _, i := range c.Items {
			items = append(items, i.Name)
		}
		legacy.Skills = append(legacy.Skills, CategoryValuePair{Category: c.Name, Value: strings.Join(items, ", ")})
	}

	// Experience
	for _, p := range er.Experience.Positions {
		legacy.Experience = append(legacy.Experience, Experience{
			Company:     p.Company,
			Title:       p.Title,
			Description: p.Description,
			Dates:       DateRange{Start: p.Dates.Start, End: p.Dates.End},
			Location:    p.Location,
		})
	}

	// Projects
	for _, p := range er.Projects.Projects {
		var link Link
		if len(p.Links) > 0 {
			link = Link{Text: p.Links[0].Text, Ref: p.Links[0].URL}
		}
		legacy.Projects = append(legacy.Projects, Project{
			Name:        p.Name,
			Language:    strings.Join(p.Technologies, ", "),
			Description: p.Description,
			Link:        link,
		})
	}

	// Education
	for _, i := range er.Education.Institutions {
		legacy.Education = append(legacy.Education, Education{
			School:      i.Institution,
			Degree:      i.Degree,
			Suffixes:    i.Honors,
			Description: i.Description,
			Location:    *i.Location,
			Dates:       DateRange{Start: i.Dates.Start, End: i.Dates.End},
		})
	}

	return legacy
}

// parseCommaSeparated splits a comma-separated string into a slice
func parseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	var result []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}