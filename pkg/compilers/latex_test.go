package compilers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestDetectLaTeXEngine(t *testing.T) {
	tmpDir := t.TempDir()
	enginePath := filepath.Join(tmpDir, "xelatex")
	if err := os.WriteFile(enginePath, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("failed to create fake engine: %v", err)
	}

	t.Setenv("PATH", tmpDir)

	if got := DetectLaTeXEngine(); got != "xelatex" {
		t.Fatalf("DetectLaTeXEngine() = %q, want xelatex", got)
	}

	available := GetAvailableLaTeXEngines()
	if len(available) != 1 || available[0] != "xelatex" {
		t.Fatalf("GetAvailableLaTeXEngines() = %v, want [xelatex]", available)
	}
}

func TestDetectLaTeXEngineNone(t *testing.T) {
	t.Setenv("PATH", "")
	if got := DetectLaTeXEngine(); got != "" {
		t.Fatalf("DetectLaTeXEngine() with empty PATH = %q, want empty string", got)
	}
}

func TestCopyFileAndDir(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		t.Fatalf("failed to create dst dir: %v", err)
	}

	srcFile := filepath.Join(srcDir, "example.cls")
	content := []byte("class content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("failed to write src file: %v", err)
	}

	if err := copyFile(srcFile, filepath.Join(dstDir, "copied.cls")); err != nil {
		t.Fatalf("copyFile error: %v", err)
	}

	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir error: %v", err)
	}

	copied, err := os.ReadFile(filepath.Join(dstDir, "example.cls"))
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(copied) != string(content) {
		t.Fatalf("copied content = %q, want %q", copied, content)
	}
}

func TestLaTeXCompilerCompile(t *testing.T) {
	tmpDir := t.TempDir()

	scriptPath := filepath.Join(tmpDir, "latex-mock.sh")
	script := `#!/bin/sh
if [ -n "$MOCK_LATEX_LOG" ]; then
  printf "%s" "$@" > "$MOCK_LATEX_LOG"
fi
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("failed to write mock latex script: %v", err)
	}

	classesDir := filepath.Join(tmpDir, "classes")
	if err := os.MkdirAll(classesDir, 0755); err != nil {
		t.Fatalf("failed to create classes dir: %v", err)
	}
	classFile := filepath.Join(classesDir, "resume.cls")
	if err := os.WriteFile(classFile, []byte("class data"), 0644); err != nil {
		t.Fatalf("failed to write class file: %v", err)
	}

	logFile := filepath.Join(tmpDir, "latex.log")
	t.Setenv("MOCK_LATEX_LOG", logFile)

	logger := zap.NewNop().Sugar()
	compIface := NewLaTeXCompiler(scriptPath, logger)
	compiler, ok := compIface.(*LaTeXCompiler)
	if !ok {
		t.Fatalf("NewLaTeXCompiler did not return *LaTeXCompiler")
	}

	compiler.LoadClasses(classesDir)
	compiler.AddOutputFolder(tmpDir)

	output := compiler.Compile("test content", "resume")
	if !strings.HasSuffix(output, "resume.tex") {
		t.Fatalf("Compile() output = %q, want suffix resume.tex", output)
	}

	if _, err := os.Stat(output); err != nil {
		t.Fatalf("expected .tex file at %s, got error: %v", output, err)
	}

	copiedClass := filepath.Join(tmpDir, "resume.cls")
	if _, err := os.Stat(copiedClass); err != nil {
		t.Fatalf("expected class file copied to output: %v", err)
	}

	logData, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read latex log: %v", err)
	}
	if !strings.Contains(string(logData), output) {
		t.Fatalf("latex command did not receive tex path, log: %s", string(logData))
	}
}

func TestNewAutoLaTeXCompilerNoEngine(t *testing.T) {
	t.Setenv("PATH", "")
	logger := zap.NewNop().Sugar()
	compiler, err := NewAutoLaTeXCompiler(logger)
	if err == nil || compiler != nil {
		t.Fatalf("NewAutoLaTeXCompiler expected error when no engine available")
	}
}
