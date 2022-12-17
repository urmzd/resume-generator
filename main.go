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
	Contact     *Contact
	Skills      []Detail
	Experiences []Experience
	Projects    []Project
	Education   []Education
}

type Education struct {
	School   *string
	Degree   *string
	Suffixes []string
	Details  []Detail
	Location *Location
	Dates    *DateRange
}

type Contact struct {
	Name  *string
	Email *string
	Phone *string
	Links []Link
}

type Link struct {
	Text *string
	Link *string
}

func NewLink(text string, link string, prefix string) *Link {
	return &Link {
		Text: prefix + text,
		Link: link,
	}
}

func (link *Link) toString() string {
	return fmt.Sprintf(`\href{%s}{%s}`, *link.Text, *link.Link)
}

type Location struct {
	City  *string
	State *string
}

func (loc *Location) toString() string {
	return fmt.Sprintf("%s, %s", *loc.City, *loc.State)
}

type DateRange struct {
	Start *time.Time
	End   *time.Time
}

func (rng *DateRange) toString() string {
	dateFmt := "Jan 2006"

	var end string
	if rng.End == nil {
		end = "Present"
	} else {
		end = rng.End.Format(dateFmt)
	}

	return fmt.Sprintf("%s - %s", rng.Start.Format(dateFmt), end)
}

type Experience struct {
	CompanyName  string
	Title        string
	Achievements []string
	Dates        DateRange
	Location     Location
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

func (generator *DefaultResumeGenerator) write(strs ...string) {
	for _, str := range strs {
		generator.builder.WriteString(strings.TrimSpace(str))
		generator.builder.WriteString("\n")
	}
}

func (generator *DefaultResumeGenerator) AddExperiences(experiences *[]Experience) {
	beforeCode := `\section*{education}`

	generator.write(beforeCode)

	for _, experience := range *experiences {
		dateRange := experience.Dates.toString()
		generator.addProject(experience.Title, experience.CompanyName, dateRange)
		generator.addSubProject(experience.Location.toString())
		generator.addAchievements(experience.Achievements...)
	}
}

func (generator *DefaultResumeGenerator) addAchievements(achievements ...string) {
	beforeCode := `\achievements%`
	stringTemplate := `[%s]`

	generator.write(beforeCode)

	for _, achievement := range achievements {
		stringValue := fmt.Sprintf(stringTemplate, achievement)
		generator.write(stringValue)
	}
}

func (generator *DefaultResumeGenerator) addProject(left string, middle string, right string) {
	stringTemplate := `\project{%s}{%s}{%s}`
	stringValue := fmt.Sprintf(stringTemplate, left, middle, right)
	generator.write(stringValue)
}

func (generator *DefaultResumeGenerator) addSubProject(label string) {
	stringTemplate := `\subproject{%s}`
	stringValue := fmt.Sprintf(stringTemplate, label)
	generator.write(stringValue)
}

func (generator *DefaultResumeGenerator) addDescription(skills *[]Detail) {
	beforeCode := `\begin{description}`
	afterCode := `\end{description}`
	stringTemplate := `\item[%s]{%s}`

	generator.write(beforeCode)

	for _, skill := range *skills {
		skillTemplate := fmt.Sprintf(stringTemplate, skill.Category, skill.Value)
		generator.write(skillTemplate)
	}

	generator.write(afterCode)
}

func (generator *DefaultResumeGenerator) AddSkills(skills *[]Detail) {
	beforeCode := `\section*{skills}`

	generator.write(beforeCode)
	generator.addDescription(skills)
}

func (generator *DefaultResumeGenerator) StartResume(contact *Contact) {

	// Ensure that all links have a text display.
	FillMissingLinkParts(&contact.Links[0])
	FillMissingLinkParts(&contact.Links[1])

	beforeCode := []string{
		`\documentclass{resume}`,
		`\usepackage{geometry}`,
		`\usepackage{titlesec}`,
		`\usepackage[]{hyperref}`,
		`\usepackage{helvet}`,
		`\geometry{a4paper,left=0.5in,right=0.5in,bottom=0.5in,top = 0.5in}`,
		`\hypersetup {colorlinks=true,linkcolor=blue }`,
		`\titleformat{\section}{\normalfont\Large\scshape\fontsize{12}{15}}{\thesection}{1em}{}[{\titlerule[0.8pt]}]`,
		`\begin{document}`,
		`\pagestyle{empty}`,
	}

	name := fmt.Sprintf(`\name{%s}`, *contact.Name)
	basics := fmt.Sprintf(`\contact{}{}{}{}`,
		NewLink(*contact.Email, *contact.Email, "mailto://"),
		NewLink(*contact.Phone, *contact.Phone, "tel:"),
		contact.Links[0],
		contact.Links[1]
	)

	generator.write(beforeCode...)
	generator.write(name, basics)
}

func FillMissingLinkParts(link *Link) {
	if link.Text == nil {
		parsedLinked, err := url.Parse(*link.Link)

		if err != nil {
			panic(err)
		}

		urlWithoutSchema := fmt.Sprintf("%s%s", parsedLinked.Hostname(), parsedLinked.Path)

		link.Text = &urlWithoutSchema
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

	resumeBuilder.StartResume(resume.Contact)
	resumeBuilder.AddSkills(&resume.Skills)
	resumeBuilder.AddExperiences(&resume.Experiences)
	fmt.Println(resumeBuilder.builder.String())
}
