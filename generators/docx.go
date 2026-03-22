package generators

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"go.uber.org/zap"
)

// DOCXGenerator generates Word documents from resume data using a Go template
// that produces word/document.xml, then assembles the DOCX zip archive.
type DOCXGenerator struct {
	logger    *zap.SugaredLogger
	formatter *docxFormatter
}

// NewDOCXGenerator creates a new DOCX generator.
func NewDOCXGenerator(logger *zap.SugaredLogger) *DOCXGenerator {
	return &DOCXGenerator{
		logger:    logger,
		formatter: newDocxFormatter(),
	}
}

// Generate renders a DOCX template and assembles a .docx zip archive.
// templateContent is the Go template that produces word/document.xml.
// scaffoldFS is the filesystem containing the static DOCX scaffolding files.
// scaffoldDir is the directory within scaffoldFS containing the scaffold.
func (g *DOCXGenerator) Generate(templateContent string, td *TemplateData, scaffoldFS fs.FS, scaffoldDir string) ([]byte, error) {
	g.logger.Info("Rendering DOCX template")

	funcs := g.formatter.TemplateFuncs()

	tmpl, err := template.New("docx").Funcs(funcs).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DOCX template: %w", err)
	}

	var documentXML strings.Builder
	if err := tmpl.Execute(&documentXML, td); err != nil {
		return nil, fmt.Errorf("failed to execute DOCX template: %w", err)
	}

	// Assemble the .docx zip archive
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// Write static scaffold files
	if err := g.writeScaffold(zw, scaffoldFS, scaffoldDir); err != nil {
		return nil, fmt.Errorf("failed to write scaffold files: %w", err)
	}

	// Write the rendered document.xml
	w, err := zw.Create("word/document.xml")
	if err != nil {
		return nil, fmt.Errorf("failed to create word/document.xml in zip: %w", err)
	}
	if _, err := w.Write([]byte(documentXML.String())); err != nil {
		return nil, fmt.Errorf("failed to write word/document.xml: %w", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize docx zip: %w", err)
	}

	g.logger.Info("Successfully rendered DOCX template")
	return buf.Bytes(), nil
}

// writeScaffold walks the scaffold directory and writes all files into the zip.
func (g *DOCXGenerator) writeScaffold(zw *zip.Writer, scaffoldFS fs.FS, scaffoldDir string) error {
	return fs.WalkDir(scaffoldFS, scaffoldDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Compute the path within the zip (strip the scaffold directory prefix)
		relPath := strings.TrimPrefix(path, scaffoldDir+"/")
		if relPath == path {
			return nil
		}

		data, err := fs.ReadFile(scaffoldFS, path)
		if err != nil {
			return fmt.Errorf("failed to read scaffold file %s: %w", path, err)
		}

		w, err := zw.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create %s in zip: %w", relPath, err)
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("failed to write %s: %w", relPath, err)
		}

		return nil
	})
}
