package generators

import (
	"fmt"
	"net/url"
	"strings"

	"text/template"

	"github.com/urmzd/generate-resumes/pkg/definition"
	"go.uber.org/zap"
)

type BaseResumeGenerator struct {
	logger *zap.SugaredLogger
	// No need to store the template as a field anymore
}

func NewBaseResumeGenerator(logger *zap.SugaredLogger) *BaseResumeGenerator {
	return &BaseResumeGenerator{
		logger: logger,
	}
}

var _ definition.Generator = &BaseResumeGenerator{}

func (generator *BaseResumeGenerator) Generate(tmpl string, resume *definition.Resume) string {
	// debug resume
	generator.logger.Debugf("Resume: %+v", resume)

	var latexBuilder strings.Builder

	// Create a new template to encompass the entire document
	fullTemplate := template.New("fullTemplate")
	fullTemplate = fullTemplate.Funcs(template.FuncMap{
		"escapeLatexChars": generator.escapeLatexChars,
		"fmtDates":         generator.fmtDates,
		"fmtLink":          generator.fmtLink,
		"fmtLocation":      generator.fmtLocation,
	})

	// Parse the template
	fullTemplate = template.Must(fullTemplate.Parse(tmpl))

	// Execute the full template to generate the entire document
	fullTemplate.Execute(&latexBuilder, resume)

	return latexBuilder.String()
}

func (generator *BaseResumeGenerator) fmtDates(date definition.DateRange) string {
	if date.End != nil {
		return fmt.Sprintf("%s - %s", date.Start.Format("Jan 2006"), date.End.Format("Jan 2006"))
	}

	return fmt.Sprintf("%s - Present", date.Start.Format("Jan 2006"))
}

func (generator *BaseResumeGenerator) fmtLink(link definition.Link) string {
	generator.logger.Debugf("Formatting link: %s", link)

	if link.Text == "" {
		parsedLinked, err := url.Parse(link.Ref)

		if err != nil {
			generator.logger.Fatalf("Unable to parse link: %s", link.Ref)
		}

		urlWithoutSchema := fmt.Sprintf("%s%s", parsedLinked.Host, parsedLinked.Path)

		link.Text = urlWithoutSchema

	}

	return fmt.Sprintf(`\href{%s}{%s}`, link.Ref, link.Text)
}

func (generator *BaseResumeGenerator) fmtLocation(location *definition.Location) string {
	generator.logger.Debugf("Formatting location: %s", location)

	if location != nil {
		return fmt.Sprintf("%s, %s", location.City, location.State)
	}

	return "Remote"
}

func (generator *BaseResumeGenerator) escapeLatexChars(str string) string {
	replacer := strings.NewReplacer(
		"%", `\%`,
		"{", `\{`,
		"}", `\}`,
		"$", `\$`,
		"&", `\&`,
		"#", `\#`,
		"_", `\_`,
		"^", `\textasciicircum`,
		"~", `\textasciitilde`,
		"\\", `\textbackslash`,
	)
	return replacer.Replace(str)
}
