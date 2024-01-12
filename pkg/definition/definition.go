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
	City  string
	State string
}

type Experience struct {
	Company     string
	Title       string
	Description []string
	Dates       DateRange
	Location    *Location
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
