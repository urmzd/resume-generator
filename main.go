package main

import "fmt"
import "github.com/BurntSushi/toml"

/**
  The pipeline for construction will go as follows:

  CreateResume(Basics) -> AddExperience -> AddSkill -> AddProject -> AddContact
*/

type Resume struct {
	Basics      Basics
	Skills      []Detail
	Experiences []Experience
	Projects    []Project
	Education   []Education
}

type Education struct {
	School   string
	Degree   string
	Suffixes []string
	Details  []Detail
	Location Location
	Dates    DateRange
}

type Basics struct {
	Name  string
	Email string
	Phone string
	Links []Link
}

type Link struct {
	Link string
}

type Location struct {
	City  string
	State string
}

type DateRange struct {
	Start string
	End   string
}

type Experience struct {
	CompanyName string
	Title       string
	Achievement []string
	Dates       DateRange
	Location    Location
}

type Project struct {
	Name         string
	LanguageUsed string
	Details      []string
}

type Detail struct {
	Category string
	Value    string
}

type ResumeGenerator interface {
	CreateResume(*Basics) Resume
	AddProject(*Resume, Project) *Resume
	AddExperience(*Resume, Experience) *Resume
	AddSkill(*Resume, Detail) *Resume
	GenerateResume(Resume) string
}

func main() {
  resume := Resume {}
	fmt.Println(resume)
}
