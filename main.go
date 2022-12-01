package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

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
	Start time.Time
	End   time.Time
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
	StartResume(*Basics)
	AddSkills(*[]Detail)
	AddExperiences(*[]Experience)
	AddEducation(*[]Education)
	AddProjects(*[]Project)
	EndResume() string
}

type DefaultResumeGenerator struct {
	builder strings.Builder
}

func (generator *DefaultResumeGenerator) StartResume(basics *Basics) {
	//generator.builder.Write(
	beforeCode := `
		\documentclass{resume}
		\usepackage{geometry}
		\usepackage{titlesec}
		\usepackage[allcolors=blue]{hyperref}
		\usepackage{helvet}

		\geometry{
				a4paper,
				left=0.5in,
				right=0.5in,
				bottom=0.5in,
				top = 0.5in,
		}

		\hypersetup {
			colorlinks=true,
			linkcolor=blue
		}

		\titleformat{\section}
			{\normalfont\Large\scshape\fontsize{12}{15}}
			{\thesection}{1em}
			{}[{\titlerule[0.8pt]}]

		\begin{document}
		\pagestyle{empty}
	`

	name := fmt.Sprintf(`\name{%s}`, basics.Name)
	contact := fmt.Sprintf(
		`\contact%
			{\href{mailto://%s}{%s}}
			{\href{tel:%s}{%s}}
			{\href{%s}{%s}}
			{\href{%s}{%s}}
		`,
	)

}

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	config := string(data)

	fmt.Println(config)

	if err != nil {
		panic(err)
	}

	var resume Resume
	_, err = toml.Decode(config, &resume)

	if err != nil {
		panic(err)
	}

	resumeJson, err := json.MarshalIndent(resume, "", "\t")
	fmt.Println(string(resumeJson))
}
