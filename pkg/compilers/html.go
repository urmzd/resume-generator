package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

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

	if overridePath, overrideName := resolveToolOverride(logger, os.Getenv("RESUME_HTML_TO_PDF_TOOL")); overridePath != "" {
		compiler.toolPath = overridePath
		compiler.toolName = overrideName
		logger.Infof("Using override HTML to PDF tool: %s", overrideName)
		return compiler
	}

	// Try to find available tools
	tools := []string{
		"wkhtmltopdf",
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
				"/Applications/ungoogled-chromium.app/Contents/MacOS/ungoogled-chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "ungoogled-chromium.app", "Contents", "MacOS", "ungoogled-chromium"),
				"/Applications/ungoogled-chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "ungoogled-chromium.app", "Contents", "MacOS", "Chromium"),
				"/Applications/Eloston-Ungoogled-Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Eloston-Ungoogled-Chromium.app", "Contents", "MacOS", "Chromium"),
				"/Applications/Eloston Ungoogled Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Eloston Ungoogled Chromium.app", "Contents", "MacOS", "Chromium"),
				"/Applications/Chromium.app/Contents/MacOS/Chromium",
				filepath.Join(os.Getenv("HOME"), "Applications", "Chromium.app", "Contents", "MacOS", "Chromium"),
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
	case c.toolName == "wkhtmltopdf":
		return c.compileWithWKHTMLToPDF(tmpFile.Name(), outputPath)
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
	args := []string{
		"--headless=new",
		"--disable-gpu",
		"--disable-dev-shm-usage",
		"--print-to-pdf=" + absOutputPath,
	}

	// Use a temporary user data directory to avoid touching default profiles (prevents crashes on some Chromium builds)
	userDataDir, err := os.MkdirTemp("", "resume-chromium-profile-*")
	if err != nil {
		return fmt.Errorf("failed to create temp user data dir: %w", err)
	}
	defer os.RemoveAll(userDataDir)
	args = append(args, "--user-data-dir="+userDataDir)

	if os.Geteuid() == 0 {
		args = append(args, "--no-sandbox")
	}

	// Allow users to inject custom flags for troubleshooting
	if extra := os.Getenv("RESUME_CHROMIUM_FLAGS"); extra != "" {
		args = append(args, splitArgs(extra)...)
	}

	args = append(args, "file://"+absHTMLPath)

	cmd := exec.Command(c.toolPath, args...)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Errorf("%s output: %s", c.toolName, string(output))

		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
				msg := fmt.Sprintf("%s exited with signal %s", c.toolName, status.Signal())
				if runtime.GOOS == "darwin" {
					msg += ". Headless Chromium from macOS app bundles often fails due to sandbox restrictions. Install wkhtmltopdf (brew install wkhtmltopdf) or run via Docker for reliable HTML→PDF conversion."
				}
				return fmt.Errorf("%s", msg)
			}
		}

		return fmt.Errorf("%s failed: %w", c.toolName, err)
	}

	// Verify PDF was created
	if _, err := os.Stat(absOutputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF was not created at %s", absOutputPath)
	}

	c.logger.Infof("Successfully converted HTML to PDF: %s", outputPath)
	return nil
}

func (c *HTMLToPDFCompiler) compileWithWKHTMLToPDF(htmlPath, outputPath string) error {
	c.logger.Infof("Converting HTML to PDF using wkhtmltopdf")

	absHTMLPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	args := []string{
		"--enable-local-file-access",
		"--quiet",
		absHTMLPath,
		absOutputPath,
	}

	if extra := os.Getenv("RESUME_WKHTMLTOPDF_FLAGS"); extra != "" {
		args = append(splitArgs(extra), args...)
	}

	cmd := exec.Command(c.toolPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Errorf("wkhtmltopdf output: %s", string(output))
		return fmt.Errorf("wkhtmltopdf failed: %w", err)
	}

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
	normalized := strings.ToLower(tool)
	normalized = strings.TrimSuffix(normalized, ".app")
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.ReplaceAll(normalized, " ", "-")

	switch normalized {
	case "ungoogled-chromium-browser":
		return "ungoogled-chromium"
	case "chromium-browser":
		return "chromium"
	case "chrome", "google-chrome.app", "google-chrome-stable", "google-chrome-mac":
		return "google-chrome"
	case "eloston-ungoogled-chromium", "eloston-ungoogled-chromium-mac":
		return "ungoogled-chromium"
	default:
		return normalized
	}
}

func splitArgs(value string) []string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return nil
	}
	return fields
}

func resolveToolOverride(logger *zap.SugaredLogger, override string) (string, string) {
	override = strings.TrimSpace(override)
	if override == "" {
		return "", ""
	}

	// Attempt to resolve named executables on PATH first
	if !strings.Contains(override, "/") {
		if resolved, err := exec.LookPath(override); err == nil {
			name := canonicalToolName(filepath.Base(resolved))
			return resolved, name
		}
	}

	candidate := override
	if !filepath.IsAbs(candidate) {
		if abs, err := filepath.Abs(candidate); err == nil {
			candidate = abs
		}
	}

	if info, err := os.Stat(candidate); err == nil {
		if info.IsDir() {
			if strings.HasSuffix(strings.ToLower(candidate), ".app") {
				if execPath := findMacAppExecutable(candidate); execPath != "" {
					return execPath, detectToolName(execPath)
				}
			}
			logger.Warnf("RESUME_HTML_TO_PDF_TOOL points to a directory without a known executable: %s", candidate)
			return "", ""
		}
		return candidate, detectToolName(candidate)
	}

	return "", ""
}

func findMacAppExecutable(appPath string) string {
	binaryHints := []string{
		filepath.Base(strings.TrimSuffix(appPath, ".app")),
		"Chromium",
		"Google Chrome",
		"chrome",
		"chromium",
		"ungoogled-chromium",
	}

	for _, hint := range binaryHints {
		path := filepath.Join(appPath, "Contents", "MacOS", hint)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}

func detectToolName(path string) string {
	lowerPath := strings.ToLower(path)
	switch {
	case strings.Contains(lowerPath, "wkhtmlto"):
		return "wkhtmltopdf"
	case strings.Contains(lowerPath, "ungoogled"):
		return "ungoogled-chromium"
	case strings.Contains(lowerPath, "chromium"):
		return "chromium"
	case strings.Contains(lowerPath, "chrome"):
		return "google-chrome"
	default:
		return canonicalToolName(filepath.Base(path))
	}
}
