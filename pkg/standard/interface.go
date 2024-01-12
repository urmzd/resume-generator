package standard

import (
	"time"
)

type Compiler interface {
	LoadClasses(classes ...string)
	AddOutputFolder(string)
	Compile(string, string)
}

type Resume struct {
	Contact    Contact
	Skills     []Detail
	Experience []Experience
	Projects   []Project
	Education  []Education
}

type Education struct {
	School   string
	Degree   string
	Suffixes []string
	Details  []Detail
	Location Location
	Dates    DateRange
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
	Company      string
	Title        string
	Achievements []string
	Dates        DateRange
	Location     *Location
}

type Project struct {
	Name     string
	Language string
	Details  []string
	Link     Link
}

type Detail struct {
	Category string
	Value    string
}

type DateRange struct {
	Start time.Time
	End   *time.Time
}

type ResumeGenerator interface {
	StartResume(*Contact)
	AddSkills(*[]Detail)
	AddExperiences(*[]Experience)
	AddEducation(*[]Education)
	AddProjects(*[]Project)
	EndResume() string
}
