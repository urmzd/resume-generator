package definition

import (
	"time"
)

type Compiler interface {
	LoadClasses(string)
	AddOutputFolder(string)
	Compile(string, string) string
}

type Resume struct {
	Contact    Contact
	Skills     []CategoryValuePair
	Experience []Experience
	Projects   []Project
	Education  []Education
}

type Generator interface {
	Generate(string, *Resume) string
}

type Education struct {
	School      string
	Degree      string
	Suffixes    []string
	Description []CategoryValuePair
	Location    Location
	Dates       DateRange
}

type Contact struct {
	Name  string
	Email string
	Phone string
	Links []Link
}

type Link struct {
	Text string
	Ref  string
}

type Location struct {
	Address     string `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`
	PostalCode  string `json:"postalCode,omitempty" yaml:"postalCode,omitempty" toml:"postalCode,omitempty"`
	City        string `json:"city" yaml:"city" toml:"city"`
	State       string `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`
	Country     string `json:"country,omitempty" yaml:"country,omitempty" toml:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty" yaml:"countryCode,omitempty" toml:"countryCode,omitempty"`
	Region      string `json:"region,omitempty" yaml:"region,omitempty" toml:"region,omitempty"`
}

type Experience struct {
	Company      string    `yaml:"company" json:"company" toml:"company"`
	Title        string    `yaml:"title" json:"title" toml:"title"`
	Description  []string  `yaml:"description,omitempty" json:"description,omitempty" toml:"description,omitempty"`
	Achievements []string  `yaml:"achievements,omitempty" json:"achievements,omitempty" toml:"achievements,omitempty"`
	Dates        DateRange `yaml:"dates" json:"dates" toml:"dates"`
	Location     *Location `yaml:"location,omitempty" json:"location,omitempty" toml:"location,omitempty"`
}

type Project struct {
	Name        string
	Language    string
	Description []string
	Link        Link
}

type CategoryValuePair struct {
	Category string
	Value    string
}

type DateRange struct {
	Start time.Time
	End   *time.Time
}
