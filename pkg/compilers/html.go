package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
)

// HTMLToPDFCompiler converts HTML to PDF using available tools
type HTMLToPDFCompiler struct {
	logger *zap.SugaredLogger
	tool   string // "chromium", "wkhtmltopdf", or "chrome"
}

// NewHTMLToPDFCompiler creates a new HTML to PDF compiler
// It auto-detects which tool is available
func NewHTMLToPDFCompiler(logger *zap.SugaredLogger) *HTMLToPDFCompiler {
	compiler := &HTMLToPDFCompiler{
		logger: logger,
	}

	// Try to find available tools
	tools := []string{
		"chromium",
		"chromium-browser",
		"google-chrome",
		"chrome",
		"wkhtmltopdf",
	}

	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			compiler.tool = tool
			logger.Infof("Using %s for HTML to PDF conversion", tool)
			return compiler
		}
	}

	logger.Warn("No HTML to PDF tool found (tried chromium, chrome, wkhtmltopdf)")
	return compiler
}

// Compile converts HTML content to PDF
func (c *HTMLToPDFCompiler) Compile(htmlContent, outputPath string) error {
	if c.tool == "" {
		return fmt.Errorf(`no HTML to PDF conversion tool found

Please install one of the following:
  - Chromium:     brew install chromium     (macOS)
                  apt install chromium      (Debian/Ubuntu)
                  apk add chromium          (Alpine/Docker)

  - Chrome:       brew install google-chrome (macOS)

  - wkhtmltopdf:  brew install wkhtmltopdf   (macOS)
                  apt install wkhtmltopdf   (Debian/Ubuntu)

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
	case contains(c.tool, []string{"chromium", "chromium-browser", "google-chrome", "chrome"}):
		return c.compileWithChrome(tmpFile.Name(), outputPath)
	case c.tool == "wkhtmltopdf":
		return c.compileWithWkhtmltopdf(tmpFile.Name(), outputPath)
	default:
		return fmt.Errorf("unsupported tool: %s", c.tool)
	}
}

// compileWithChrome uses Chromium/Chrome headless to convert HTML to PDF
func (c *HTMLToPDFCompiler) compileWithChrome(htmlPath, outputPath string) error {
	c.logger.Infof("Converting HTML to PDF using %s", c.tool)

	// Make paths absolute
	absHTMLPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	// Chrome headless command
	cmd := exec.Command(c.tool,
		"--headless",
		"--disable-gpu",
		"--no-sandbox",
		"--print-to-pdf="+absOutputPath,
		"file://"+absHTMLPath,
	)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Errorf("Chrome output: %s", string(output))
		return fmt.Errorf("chrome failed: %w", err)
	}

	// Verify PDF was created
	if _, err := os.Stat(absOutputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF was not created at %s", absOutputPath)
	}

	c.logger.Infof("Successfully converted HTML to PDF: %s", outputPath)
	return nil
}

// compileWithWkhtmltopdf uses wkhtmltopdf to convert HTML to PDF
func (c *HTMLToPDFCompiler) compileWithWkhtmltopdf(htmlPath, outputPath string) error {
	c.logger.Info("Converting HTML to PDF using wkhtmltopdf")

	cmd := exec.Command("wkhtmltopdf",
		"--quiet",
		"--enable-local-file-access",
		htmlPath,
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Errorf("wkhtmltopdf output: %s", string(output))
		return fmt.Errorf("wkhtmltopdf failed: %w", err)
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
