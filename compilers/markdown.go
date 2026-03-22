package compilers

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
)

// MarkdownToHTML converts Markdown content to a standalone HTML document
// suitable for PDF rendering.
func MarkdownToHTML(markdownContent string) (string, error) {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdownContent), &buf); err != nil {
		return "", fmt.Errorf("failed to convert Markdown to HTML: %w", err)
	}

	// Wrap in a minimal HTML document for PDF rendering
	html := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<style>
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
    font-size: 11pt;
    line-height: 1.4;
    color: #333;
    max-width: 8.5in;
    margin: 0 auto;
    padding: 0.5in;
  }
  h1 { font-size: 1.6em; margin: 0 0 0.3em 0; }
  h2 { font-size: 1.2em; margin: 1em 0 0.3em 0; border-bottom: 1px solid #ccc; padding-bottom: 0.2em; }
  h3 { font-size: 1em; margin: 0.8em 0 0.2em 0; }
  ul { margin: 0.2em 0; padding-left: 1.5em; }
  li { margin: 0.1em 0; }
  p { margin: 0.3em 0; }
  a { color: #0066cc; text-decoration: none; }
  hr { border: none; border-top: 1px solid #ccc; margin: 0.5em 0; }
  @media print {
    body { padding: 0; }
  }
</style>
</head>
<body>
` + buf.String() + `
</body>
</html>`

	return html, nil
}
