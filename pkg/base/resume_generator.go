package base

import (
	"fmt"
	"log"
	"strings"

	"net/url"

	"errors"

	"github.com/thoas/go-funk"
	"github.com/urmzd/generate-resumes/pkg/standard"
)

type BaseResumeGenerator struct {
	code []string
}

func locationToStr(location standard.Location) string {
	return fmt.Sprintf("%s, %s", location.City, location.State)
}

func (generator *BaseResumeGenerator) AddEducation(education *[]standard.Education) {
	beforeCode := `\section*{education}`

	generator.write(beforeCode)

	for _, school := range *education {
		dateRanges := dateRangeToStr(school.Dates)
		degree := funk.Reduce(school.Suffixes, func(acc string, cur string) string {
			return fmt.Sprintf("%s (%s)", acc, cur)
		}, school.Degree)
		generator.addProject(school.School, "", dateRanges)
		generator.addSubProject(degree.(string))
		generator.addDescription(&school.Details)
	}

}

func NewPrefixedLink(text string, prefix string) *standard.Link {
	new_link := prefix + text
	return &standard.Link{
		Text: text,
		Ref:  new_link,
	}
}

func linkToStr(link standard.Link) string {
	if link.Text == "" {
		parsedLinked, err := url.Parse(link.Ref)

		if err != nil {
			panic(err)
		}

		urlWithoutSchema := fmt.Sprintf("%s%s", parsedLinked.Hostname(), parsedLinked.Path)

		link.Text = urlWithoutSchema
	}

	return fmt.Sprintf(`\href{%s}{%s}`, link.Ref, link.Text)
}

func dateRangeToStr(rng standard.DateRange) string {
	dateFmt := "Jan 2006"

	var end string
	if rng.End == nil {
		end = "Present"
	} else {
		end = rng.End.Format(dateFmt)
	}

	return fmt.Sprintf("%s - %s", rng.Start.Format(dateFmt), end)
}

func (generator *BaseResumeGenerator) AddProjects(projects *[]standard.Project) {
	beforeCode := `\section*{projects}`

	generator.write(beforeCode)

	for _, project := range *projects {
		generator.addProject(project.Name, project.Language, linkToStr(project.Link))
		generator.addAchievements(project.Details...)
	}
}

func (generator *BaseResumeGenerator) write(strs ...string) {
	for _, str := range strs {
		generator.code = append(generator.code, strings.TrimSpace(str))
	}
}

func (gen *BaseResumeGenerator) EndResume() string {
	gen.write(`\end{document}`)
	preTex := strings.Join(gen.code, "\n")
	// We need to need to escape amparcents.
	// We can move this to the compile step when finished.

	processedTex := strings.ReplaceAll(preTex, "&", `\&`)

	return processedTex
}

func (generator *BaseResumeGenerator) AddExperiences(experience *[]standard.Experience) {
	beforeCode := `\section*{experience}`

	generator.write(beforeCode)

	for _, xp := range *experience {
		dateRange := dateRangeToStr(xp.Dates)
		generator.addProject(xp.Title, xp.Company, dateRange)
		if xp.Location == nil {
			generator.addSubProject("Remote")
		} else {
			generator.addSubProject(locationToStr(*xp.Location))
		}
		generator.addAchievements(xp.Achievements...)
	}
}

func (generator *BaseResumeGenerator) addAchievements(achievements ...string) {
	beforeCode := `\achievements%`
	stringTemplate := `[%s]`

	generator.write(beforeCode)

	for _, achievement := range achievements {
		stringValue := fmt.Sprintf(stringTemplate, achievement)
		// escape % in stringValue
		stringValueEscaped := strings.ReplaceAll(stringValue, "%", `\%`)
		generator.write(stringValueEscaped)
	}
}

func (generator *BaseResumeGenerator) addProject(left string, middle string, right string) {
	stringTemplate := `\project{%s}{%s}{%s}`
	stringValue := fmt.Sprintf(stringTemplate, left, middle, right)
	generator.write(stringValue)
}

func (generator *BaseResumeGenerator) addSubProject(label string) {
	stringTemplate := `\subproject{%s}`
	stringValue := fmt.Sprintf(stringTemplate, label)
	generator.write(stringValue)
}

func (generator *BaseResumeGenerator) addDescription(skills *[]standard.Detail) {
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

func (generator *BaseResumeGenerator) AddSkills(skills *[]standard.Detail) {
	if len(*skills) > 0 {
		beforeCode := `\section*{skills}`
		generator.write(beforeCode)
		generator.addDescription(skills)
	}
}

func (generator *BaseResumeGenerator) StartResume(contact *standard.Contact) {
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

	basics := ""

	if len(contact.Links) == 1 {
		basics = fmt.Sprintf(`\contact{%s}{%s}{%s}`,
			linkToStr(*NewPrefixedLink(contact.Email, "mailto:")),
			linkToStr(*NewPrefixedLink(contact.Phone, "tel:")),
			linkToStr(contact.Links[0]),
		)
	} else if len(contact.Links) == 2 {
		basics = fmt.Sprintf(`\contact{%s}{%s}{%s}{%s}`,
			linkToStr(*NewPrefixedLink(contact.Email, "mailto:")),
			linkToStr(*NewPrefixedLink(contact.Phone, "tel:")),
			linkToStr(contact.Links[0]),
			linkToStr(contact.Links[1]),
		)
	} else {
		log.Fatal(errors.New("Contact length must be 1 or 2"))
	}

	generator.write(beforeCode...)
	generator.write(name, basics)
}
