package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"net/url"
)

type Resume struct {
	Contact     Contact
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

type Contact struct {
	Name  string
	Email string
	Phone string
	Links []Link
}

type Link struct {
	Text string
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
	StartResume(*Contact)
	AddSkills(*[]Detail)
	AddExperiences(*[]Experience)
	AddEducation(*[]Education)
	AddProjects(*[]Project)
	EndResume() string
}

type DefaultResumeGenerator struct {
	builder strings.Builder
}

func (generator *DefaultResumeGenerator) StartResume(contact *Contact) {

	// Ensure that all links have a text display.
	FillMissingLinkParts(&contact.Links[0])
	FillMissingLinkParts(&contact.Links[1])

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

	fmt.Println(contact)

	name := fmt.Sprintf(`\name{%s}`, contact.Name)
	basics := fmt.Sprintf(
		`\contact%%
			{\href{mailto://%s}{%s}}
			{\href{tel:%s}{%s}}
			{\href{%s}{%s}}
			{\href{%s}{%s}}
		`,
		contact.Email,
		contact.Email,
		contact.Phone,
		contact.Phone,
		contact.Links[0].Text,
		contact.Links[0].Link,
		contact.Links[1].Text,
		contact.Links[1].Link,
	)

	generator.builder.WriteString(beforeCode)
	generator.builder.WriteString(name)
	generator.builder.WriteString(basics)
}

func FillMissingLinkParts(link *Link) {
	if link.Text == "" {
		parsedLinked , err := url.Parse(link.Link)

		if err != nil {
			panic(err)
		}

		urlWithoutSchema := fmt.Sprintf("%s%s", parsedLinked.Hostname(), parsedLinked.Path)

		link.Text = urlWithoutSchema
	}
} 

func main() {
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	config := string(data)

	if err != nil {
		panic(err)
	}

	var resume Resume
	_, err = toml.Decode(config, &resume)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", &resume)

	resumeBuilder := &DefaultResumeGenerator{}

	resumeBuilder.StartResume(&resume.Contact)
	fmt.Println(resumeBuilder.builder.String())
}
