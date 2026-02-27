package resume

import (
	"time"
)

type Resume struct {
	Contact        Contact         `json:"contact" yaml:"contact" toml:"contact"`
	Summary        string          `json:"summary,omitempty" yaml:"summary,omitempty" toml:"summary,omitempty"`
	Certifications *Certifications `json:"certifications,omitempty" yaml:"certifications,omitempty" toml:"certifications,omitempty"`
	Skills         Skills          `json:"skills" yaml:"skills" toml:"skills"`
	Experience     ExperienceList  `json:"experience" yaml:"experience" toml:"experience"`
	Projects       *ProjectList    `json:"projects,omitempty" yaml:"projects,omitempty" toml:"projects,omitempty"`
	Education      EducationList   `json:"education" yaml:"education" toml:"education"`
	Languages      *LanguageList   `json:"languages,omitempty" yaml:"languages,omitempty" toml:"languages,omitempty"`
	Layout         *Layout         `json:"layout,omitempty" yaml:"layout,omitempty" toml:"layout,omitempty"`
}

type Layout struct {
	Density      string   `json:"density,omitempty" yaml:"density,omitempty" toml:"density,omitempty"`
	Typography   string   `json:"typography,omitempty" yaml:"typography,omitempty" toml:"typography,omitempty"`
	Header       string   `json:"header,omitempty" yaml:"header,omitempty" toml:"header,omitempty"`
	Sections     []string `json:"sections,omitempty" yaml:"sections,omitempty" toml:"sections,omitempty"`
	SkillColumns int      `json:"skill_columns,omitempty" yaml:"skill_columns,omitempty" toml:"skill_columns,omitempty"`
}

type LanguageList struct {
	Title     string     `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Languages []Language `json:"languages" yaml:"languages" toml:"languages"`
}

type Language struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Proficiency string `json:"proficiency,omitempty" yaml:"proficiency,omitempty" toml:"proficiency,omitempty"`
}

type Certifications struct {
	Title string          `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Items []Certification `json:"items" yaml:"items" toml:"items"`
}

type Certification struct {
	Name   string     `json:"name" yaml:"name" toml:"name"`
	Issuer string     `json:"issuer,omitempty" yaml:"issuer,omitempty" toml:"issuer,omitempty"`
	Notes  string     `json:"notes,omitempty" yaml:"notes,omitempty" toml:"notes,omitempty"`
	Date   *time.Time `json:"date,omitempty" yaml:"date,omitempty" toml:"date,omitempty"`
}

type Location struct {
	City     string `json:"city" yaml:"city" toml:"city"`
	State    string `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`
	Province string `json:"province,omitempty" yaml:"province,omitempty" toml:"province,omitempty"`
	Country  string `json:"country,omitempty" yaml:"country,omitempty" toml:"country,omitempty"`
	Remote   bool   `json:"remote,omitempty" yaml:"remote,omitempty" toml:"remote,omitempty"`
}

type Link struct {
	URI   string `json:"uri" yaml:"uri" toml:"uri"`
	Label string `json:"label,omitempty" yaml:"label,omitempty" toml:"label,omitempty"`
}

type Contact struct {
	Name     string    `json:"name" yaml:"name" toml:"name"`
	Email    string    `json:"email" yaml:"email" toml:"email"`
	Phone    string    `json:"phone,omitempty" yaml:"phone,omitempty" toml:"phone,omitempty"`
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Links    []Link    `json:"links,omitempty" yaml:"links,omitempty" toml:"links,omitempty"`
}

type Skills struct {
	Title      string          `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Categories []SkillCategory `json:"categories" yaml:"categories" toml:"categories"`
}

type SkillCategory struct {
	Category string   `json:"category" yaml:"category" toml:"category"`
	Items    []string `json:"items" yaml:"items" toml:"items"`
}

type ExperienceList struct {
	Title     string       `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Positions []Experience `json:"positions" yaml:"positions" toml:"positions"`
}

type Experience struct {
	Company        string    `json:"company" yaml:"company" toml:"company"`
	Title          string    `json:"title" yaml:"title" toml:"title"`
	EmploymentType string    `json:"employment_type,omitempty" yaml:"employment_type,omitempty" toml:"employment_type,omitempty"`
	Highlights     []string  `json:"highlights,omitempty" yaml:"highlights,omitempty" toml:"highlights,omitempty"`
	Duties         []string  `json:"duties,omitempty" yaml:"duties,omitempty" toml:"duties,omitempty"`
	Notes          string    `json:"notes,omitempty" yaml:"notes,omitempty" toml:"notes,omitempty"`
	Dates          DateRange `json:"dates" yaml:"dates" toml:"dates"`
	Location       *Location `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Technologies   []string  `json:"technologies,omitempty" yaml:"technologies,omitempty" toml:"technologies,omitempty"`
}
type ExperienceGroup struct {
	Name      string `json:"name" yaml:"name" toml:"name"`
	Positions []int  `json:"positions" yaml:"positions" toml:"positions"` // References to position orders
}

type DateRange struct {
	Start time.Time  `json:"start" yaml:"start" toml:"start"`
	End   *time.Time `json:"end,omitempty" yaml:"end,omitempty" toml:"end,omitempty"`
}

type ProjectList struct {
	Title    string    `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Projects []Project `json:"projects" yaml:"projects" toml:"projects"`
}
type Project struct {
	Name         string     `json:"name" yaml:"name" toml:"name"`
	Link         Link       `json:"link,omitempty" yaml:"link,omitempty" toml:"link,omitempty"`
	Highlights   []string   `json:"highlights,omitempty" yaml:"highlights,omitempty" toml:"highlights,omitempty"`
	Dates        *DateRange `json:"dates,omitempty" yaml:"dates,omitempty" toml:"dates,omitempty"`
	Technologies []string   `json:"technologies,omitempty" yaml:"technologies,omitempty" toml:"technologies,omitempty"`
}

type EducationList struct {
	Title        string      `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Institutions []Education `json:"institutions" yaml:"institutions" toml:"institutions"`
}

type Degree struct {
	Name         string   `json:"name" yaml:"name" toml:"name"`
	Descriptions []string `json:"descriptions,omitempty" yaml:"descriptions,omitempty" toml:"descriptions,omitempty"`
}

type GPA struct {
	GPA    string `json:"gpa,omitempty" yaml:"gpa,omitempty" toml:"gpa,omitempty"`
	MaxGPA string `json:"max_gpa,omitempty" yaml:"max_gpa,omitempty" toml:"max_gpa,omitempty"`
}

type Education struct {
	Institution     string    `json:"institution" yaml:"institution" toml:"institution"`
	Degree          Degree    `json:"degree" yaml:"degree" toml:"degree"`
	Specializations []string  `json:"specializations,omitempty" yaml:"specializations,omitempty" toml:"specializations,omitempty"`
	GPA             *GPA      `json:"gpa,omitempty" yaml:"gpa,omitempty" toml:"gpa,omitempty"`
	Awards          []Award   `json:"awards,omitempty" yaml:"awards,omitempty" toml:"awards,omitempty"`
	Dates           DateRange `json:"dates" yaml:"dates" toml:"dates"`
	Location        *Location `json:"location,omitempty" yaml:"location,omitempty" toml:"location,omitempty"`
	Thesis          *Thesis   `json:"thesis,omitempty" yaml:"thesis,omitempty" toml:"thesis,omitempty"`
}

type Thesis struct {
	Title       string   `json:"title" yaml:"title" toml:"title"`
	Highlights  []string `json:"highlights,omitempty" yaml:"highlights,omitempty" toml:"highlights,omitempty"`
	Link        Link     `json:"link" yaml:"link" toml:"link"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
}

type Award struct {
	Name  string     `json:"name" yaml:"name" toml:"name"`
	Date  *time.Time `json:"date,omitempty" yaml:"date,omitempty" toml:"date,omitempty"`
	Notes string     `json:"notes,omitempty" yaml:"notes,omitempty" toml:"notes,omitempty"`
}

func Validate(resume *Resume) []ValidationError {
	var errors []ValidationError

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
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Value   interface{} `json:"value,omitempty"`
}
