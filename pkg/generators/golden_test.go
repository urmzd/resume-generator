package generators

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

var update = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
	// Resolve project root (tests run from pkg/generators/)
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve project root: %v", err)
	}

	// Point template resolution at the project root
	t.Setenv("RESUME_TEMPLATES_DIR", projectRoot)

	logger := zap.NewNop().Sugar()
	gen := NewGenerator(logger)

	testdataDir := filepath.Join("testdata")
	inputDir := filepath.Join(testdataDir, "input")
	goldenDir := filepath.Join(testdataDir, "golden")

	if *update {
		if err := os.MkdirAll(goldenDir, 0755); err != nil {
			t.Fatalf("failed to create golden dir: %v", err)
		}
	}

	inputs := []struct {
		name string
		file string
	}{
		{"software_engineer", "software_engineer.yml"},
		{"minimal", "minimal.yml"},
	}

	templates := []struct {
		name  string
		ttype TemplateType
		ext   string
	}{
		{"modern-html", TemplateTypeHTML, ""},
		{"modern-latex", TemplateTypeLaTeX, ""},
		{"modern-cv", TemplateTypeLaTeX, ""},
		{"modern-markdown", TemplateTypeMarkdown, ""},
		{"modern-docx", TemplateTypeDOCX, ""},
	}

	for _, input := range inputs {
		inputPath := filepath.Join(inputDir, input.file)
		inputData, err := resume.LoadResumeFromFile(inputPath)
		if err != nil {
			t.Fatalf("failed to load %s: %v", input.file, err)
		}
		resumeData := inputData.ToResume()

		for _, tmpl := range templates {
			testName := fmt.Sprintf("%s/%s", input.name, tmpl.name)

			t.Run(testName, func(t *testing.T) {
				tmplObj, err := LoadTemplate(tmpl.name)
				if err != nil {
					t.Fatalf("failed to load template %s: %v", tmpl.name, err)
				}

				var actual []byte

				if tmplObj.Type == TemplateTypeDOCX {
					docxBytes, err := gen.GenerateDOCX(resumeData)
					if err != nil {
						t.Fatalf("GenerateDOCX() error: %v", err)
					}
					// Extract word/document.xml from the ZIP
					actual, err = extractDocumentXML(docxBytes)
					if err != nil {
						t.Fatalf("failed to extract document.xml: %v", err)
					}
				} else {
					content, err := gen.GenerateWithTemplate(tmplObj, resumeData)
					if err != nil {
						t.Fatalf("GenerateWithTemplate() error: %v", err)
					}
					actual = []byte(content)
				}

				goldenExt := ".golden"
				if tmplObj.Type == TemplateTypeDOCX {
					goldenExt = ".golden.xml"
				}
				goldenFile := filepath.Join(goldenDir, input.name+"."+tmpl.name+goldenExt)

				if *update {
					if err := os.WriteFile(goldenFile, actual, 0644); err != nil {
						t.Fatalf("failed to write golden file: %v", err)
					}
					t.Logf("updated golden file: %s", goldenFile)
					return
				}

				expected, err := os.ReadFile(goldenFile)
				if err != nil {
					t.Fatalf("failed to read golden file %s (run with -update to generate): %v", goldenFile, err)
				}

				if !bytes.Equal(actual, expected) {
					// Write actual output for diffing
					actualFile := strings.TrimSuffix(goldenFile, filepath.Ext(goldenFile)) + ".actual" + filepath.Ext(goldenFile)
					_ = os.WriteFile(actualFile, actual, 0644)
					t.Errorf("output differs from golden file.\n  golden: %s\n  actual: %s\n  diff with: diff %s %s",
						goldenFile, actualFile, goldenFile, actualFile)
				}
			})
		}
	}
}

// extractDocumentXML reads word/document.xml from a DOCX (ZIP) byte slice.
func extractDocumentXML(docxBytes []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open docx as zip: %w", err)
	}

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open document.xml: %w", err)
			}
			defer func() { _ = rc.Close() }()
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(rc); err != nil {
				return nil, fmt.Errorf("failed to read document.xml: %w", err)
			}
			return buf.Bytes(), nil
		}
	}
	return nil, fmt.Errorf("word/document.xml not found in docx")
}
