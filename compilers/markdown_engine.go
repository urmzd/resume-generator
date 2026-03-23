package compilers

// markdownEngine implements Engine by converting Markdown to HTML,
// then delegating to an HTML engine for PDF generation.
type markdownEngine struct {
	htmlEngine Engine
}

func (e *markdownEngine) Name() string { return "markdown+" + e.htmlEngine.Name() }

func (e *markdownEngine) Compile(content string, outputPath string) error {
	htmlContent, err := MarkdownToHTML(content)
	if err != nil {
		return err
	}
	return e.htmlEngine.Compile(htmlContent, outputPath)
}
