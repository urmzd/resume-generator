package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/thoas/go-funk"
	"net/url"
)

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

func (generator *DefaultResumeGenerator) AddEducation(education *[]Education) {
	beforeCode := `\section*{education}`

	generator.write(beforeCode)

	for _, school := range *education {
		dateRanges := school.Dates.toString()
		degree := funk.Reduce(school.Suffixes, func(acc string, cur string) string {
			return fmt.Sprintf("%s (%s)", acc, cur)
		}, school.Degree)
		generator.addProject(school.School, "", dateRanges)
		generator.addSubProject(degree.(string))
		generator.addDescription(&school.Details)
	}

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

func NewPrefixedLink(link string, prefix string) *Link {
	text := prefix + link
	return &Link{
		Text: text,
		Link: link,
	}
}

func (link *Link) toString() string {
	if link.Text == "" {
		parsedLinked, err := url.Parse(link.Link)

		if err != nil {
			panic(err)
		}

		urlWithoutSchema := fmt.Sprintf("%s%s", parsedLinked.Hostname(), parsedLinked.Path)

		link.Text = urlWithoutSchema
	}

	return fmt.Sprintf(`\href{%s}{%s}`, link.Link, link.Text)
}

type Location struct {
	City  string
	State string
}

func (loc *Location) toString() string {
	return fmt.Sprintf("%s, %s", loc.City, loc.State)
}

type DateRange struct {
	Start time.Time
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
	Company      string
	Title        string
	Achievements []string
	Dates        DateRange
	Location     Location
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

type ResumeGenerator interface {
	StartResume(*Contact)
	AddSkills(*[]Detail)
	AddExperiences(*[]Experience)
	AddEducation(*[]Education)
	AddProjects(*[]Project)
	EndResume() string
}

type DefaultResumeGenerator struct {
	code []string
}

func (generator *DefaultResumeGenerator) AddProjects(projects *[]Project) {
	beforeCode := `\section*{projects}`

	generator.write(beforeCode)

	for _, project := range *projects {
		generator.addProject(project.Name, project.Language, project.Link.toString())
		generator.addAchievements(project.Details...)
	}
}

func (generator *DefaultResumeGenerator) write(strs ...string) {
	for _, str := range strs {
		generator.code = append(generator.code, strings.TrimSpace(str))
	}
}

func (gen *DefaultResumeGenerator) EndResume() string {
	gen.write(`\end{document}`)
	return strings.Join(gen.code, "\n")
}

func (generator *DefaultResumeGenerator) AddExperiences(experience *[]Experience) {
	beforeCode := `\section*{experience}`

	generator.write(beforeCode)

	for _, xp := range *experience {
		dateRange := xp.Dates.toString()
		generator.addProject(xp.Title, xp.Company, dateRange)
		generator.addSubProject(xp.Location.toString())
		generator.addAchievements(xp.Achievements...)
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
	stringTemplate := `\item[%s:]{%s}`

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
	beforeCode := []string{
		`\documentclass{default}`,
		`\usepackage{geometry}`,
		`\usepackage{titlesec}`,
		`\usepackage[allcolors=blue]{hyperref}`,
		`\usepackage{helvet}`,
		`\geometry{a4paper,left=0.5in,right=0.5in,bottom=0.5in,top = 0.5in}`,
		`\hypersetup {colorlinks=true,linkcolor=blue }`,
		`\titleformat{\section}{\normalfont\Large\scshape\fontsize{12}{15}}{\thesection}{1em}{}[{\titlerule[0.8pt]}]`,
		`\begin{document}`,
		`\pagestyle{empty}`,
	}

	name := fmt.Sprintf(`\name{%s}`, contact.Name)
	basics := fmt.Sprintf(`\contact{%s}{%s}{%s}{%s}`,
		NewPrefixedLink(contact.Email, "mailto://").toString(),
		NewPrefixedLink(contact.Phone, "tel:").toString(),
		contact.Links[0].toString(),
		contact.Links[1].toString(),
	)

	generator.write(beforeCode...)
	generator.write(name, basics)
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

	//fmt.Printf("%+v\n", &resume)

	resumeBuilder := &DefaultResumeGenerator{}

	resumeBuilder.StartResume(&resume.Contact)
	resumeBuilder.AddSkills(&resume.Skills)
	resumeBuilder.AddExperiences(&resume.Experience)
	resumeBuilder.AddProjects(&resume.Projects)
	resumeBuilder.AddEducation(&resume.Education)

	fmt.Println(resumeBuilder.EndResume())
}
