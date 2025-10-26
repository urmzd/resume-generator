package compilers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestCanonicalToolName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ungoogled-chromium-browser", "ungoogled-chromium"},
		{"chromium-browser", "chromium"},
		{"chrome", "google-chrome"},
		{"chromium", "chromium"},
	}

	for _, tt := range tests {
		if got := canonicalToolName(tt.input); got != tt.want {
			t.Errorf("canonicalToolName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestContainsHelper(t *testing.T) {
	if contains("needle", []string{"hay", "stack"}) {
		t.Fatal("contains returned true for missing value")
	}
	if !contains("needle", []string{"hay", "needle", "stack"}) {
		t.Fatal("contains returned false for present value")
	}
}

func TestHTMLToPDFCompilerCompileNoTool(t *testing.T) {
	compiler := &HTMLToPDFCompiler{
		logger: zap.NewNop().Sugar(),
	}

	outputPath := filepath.Join(t.TempDir(), "resume.pdf")
	err := compiler.Compile("<html></html>", outputPath)
	if err == nil || !strings.Contains(err.Error(), "no HTML to PDF conversion tool found") {
		t.Fatalf("Compile() error = %v, want error about missing tool", err)
	}
}

func TestHTMLToPDFCompilerCompileWithChromium(t *testing.T) {
	tmpDir := t.TempDir()

	chromiumPath := filepath.Join(tmpDir, "chromium-mock.sh")
	chromiumScript := `#!/bin/sh
output=""
for arg in "$@"; do
  case "$arg" in
    --print-to-pdf=*)
      output=${arg#--print-to-pdf=}
      ;;
  esac
done

if [ -n "$MOCK_CHROMIUM_LOG" ]; then
  printf "%s" "$@" > "$MOCK_CHROMIUM_LOG"
fi

if [ -z "$output" ]; then
  echo "missing output flag" >&2
  exit 1
fi

touch "$output"
exit 0
`
	if err := os.WriteFile(chromiumPath, []byte(chromiumScript), 0755); err != nil {
		t.Fatalf("failed to write mock script: %v", err)
	}

	logPath := filepath.Join(tmpDir, "chromium.log")
	t.Setenv("MOCK_CHROMIUM_LOG", logPath)

	compiler := &HTMLToPDFCompiler{
		logger:   zap.NewNop().Sugar(),
		toolPath: chromiumPath,
		toolName: "chromium",
	}

	outputPath := filepath.Join(tmpDir, "resume.pdf")
	err := compiler.Compile("<html><body>test</body></html>", outputPath)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected PDF at %s, got error: %v", outputPath, err)
	}

	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	logContent := string(logData)
	if !strings.Contains(logContent, "--headless=new") {
		t.Errorf("chromium command missing headless flag: %s", logContent)
	}
	if !strings.Contains(logContent, "--print-to-pdf="+outputPath) {
		t.Errorf("chromium command missing print-to-pdf flag: %s", logContent)
	}
	if !strings.Contains(logContent, "--user-data-dir=") {
		t.Errorf("chromium command missing user-data-dir flag: %s", logContent)
	}
}

func TestHTMLToPDFCompilerCompileWithWKHTMLToPDF(t *testing.T) {
	tmpDir := t.TempDir()

	mockPath := filepath.Join(tmpDir, "wkhtmltopdf-mock.sh")
	mockScript := `#!/bin/sh
output=$4
if [ -z "$output" ]; then
  echo "missing output" >&2
  exit 1
fi
touch "$output"
if [ -n "$MOCK_WKHTMLTOPDF_LOG" ]; then
  printf "%s" "$@" > "$MOCK_WKHTMLTOPDF_LOG"
fi
exit 0
`
	if err := os.WriteFile(mockPath, []byte(mockScript), 0755); err != nil {
		t.Fatalf("failed to write mock wkhtmltopdf script: %v", err)
	}

	logPath := filepath.Join(tmpDir, "wkhtml.log")
	t.Setenv("MOCK_WKHTMLTOPDF_LOG", logPath)

	compiler := &HTMLToPDFCompiler{
		logger:   zap.NewNop().Sugar(),
		toolPath: mockPath,
		toolName: "wkhtmltopdf",
	}

	outputPath := filepath.Join(tmpDir, "resume.pdf")
	err := compiler.Compile("<html><body>test</body></html>", outputPath)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("expected PDF at %s, got error: %v", outputPath, err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	logContent := string(data)
	if !strings.Contains(logContent, "--enable-local-file-access") {
		t.Errorf("wkhtmltopdf command missing enable-local-file-access flag: %s", logContent)
	}
}
