package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"go.uber.org/zap"
)

// HTMLToPDFCompiler converts HTML to PDF using available tools
type HTMLToPDFCompiler struct {
	logger   *zap.SugaredLogger
	toolPath string // Selected Chromium-based executable (resolved path)
	toolName string // Canonical tool name
}

// NewHTMLToPDFCompiler creates a new HTML to PDF compiler
// It auto-detects which tool is available
func NewHTMLToPDFCompiler(logger *zap.SugaredLogger) *HTMLToPDFCompiler {
	compiler := &HTMLToPDFCompiler{
		logger: logger,
	}

	// Try to find available tools
	tools := []string{
		"ungoogled-chromium",
		"ungoogled-chromium-browser",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"chrome",
	}

	for _, tool := range tools {
		if path, err := exec.LookPath(tool); err == nil {
			compiler.toolPath = path
			compiler.toolName = canonicalToolName(tool)
			logger.Infof("Using %s for HTML to PDF conversion", compiler.toolName)
			return compiler
		}
	}

	if runtime.GOOS == "darwin" {
		macAppPaths := map[string][]string{
			"chromium": {
				"/Applications/Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Chromium.app", "Contents", "MacOS", "Chromium"),
			},
			"ungoogled-chromium": {
				"/Applications/Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Chromium.app", "Contents", "MacOS", "Chromium"),
				"/Applications/Ungoogled Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Ungoogled Chromium.app", "Contents", "MacOS", "Chromium"),
			},
			"google-chrome": {
				"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
				filepath.Join(os.Getenv("HOME"), "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
			},
		}

		for _, tool := range tools {
			canonical := canonicalToolName(tool)
			paths, ok := macAppPaths[canonical]
			if !ok {
				continue
			}
			for _, appPath := range paths {
				if _, err := os.Stat(appPath); err == nil {
					compiler.toolPath = appPath
					compiler.toolName = canonical
					logger.Infof("Using %s app bundle for HTML to PDF conversion", compiler.toolName)
					return compiler
				}
			}
		}
	}

	logger.Warn("No HTML to PDF tool found (tried ungoogled-chromium, chromium, chrome)")
	return compiler
}

// Compile converts HTML content to PDF
func (c *HTMLToPDFCompiler) Compile(htmlContent, outputPath string) error {
	if c.toolPath == "" {
		return fmt.Errorf(`no HTML to PDF conversion tool found

Please install one of the following:
  - ungoogled-chromium: brew install ungoogled-chromium (macOS)
                        apt install ungoogled-chromium  (Debian/Ubuntu)
  - chromium:           apk add chromium               (Alpine/Docker)

  - Chrome:       brew install google-chrome (macOS)

Or use Docker which includes all dependencies:
  docker run --rm -v $(pwd):/work resume-generator run -i /work/resume.yml -t modern-html`)
	}

	// Create temporary HTML file
	tmpFile, err := os.CreateTemp("", "resume-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write HTML content
	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		return fmt.Errorf("failed to write HTML: %w", err)
	}
	tmpFile.Close()

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Convert based on available tool
	switch {
	case contains(c.toolName, []string{"ungoogled-chromium", "chromium", "google-chrome"}):
		return c.compileWithChromium(tmpFile.Name(), outputPath)
	default:
		return fmt.Errorf("unsupported tool: %s", c.toolName)
	}
}

// compileWithChromium uses a Chromium-based browser in headless mode to convert HTML to PDF
func (c *HTMLToPDFCompiler) compileWithChromium(htmlPath, outputPath string) error {
	c.logger.Infof("Converting HTML to PDF using %s", c.toolName)

	// Make paths absolute
	absHTMLPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	// Headless browser command
	cmd := exec.Command(c.toolPath,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--print-to-pdf="+absOutputPath,
		"file://"+absHTMLPath,
	)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Errorf("%s output: %s", c.toolName, string(output))
		return fmt.Errorf("%s failed: %w", c.toolName, err)
	}

	// Verify PDF was created
	if _, err := os.Stat(absOutputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF was not created at %s", absOutputPath)
	}

	c.logger.Infof("Successfully converted HTML to PDF: %s", outputPath)
	return nil
}

// contains checks if a string is in a slice
func contains(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// canonicalToolName normalizes tool names to a small canonical set
func canonicalToolName(tool string) string {
	switch tool {
	case "ungoogled-chromium-browser":
		return "ungoogled-chromium"
	case "chromium-browser":
		return "chromium"
	case "chrome":
		return "google-chrome"
	default:
		return tool
	}
}
