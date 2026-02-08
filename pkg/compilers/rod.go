package compilers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

// RodHTMLToPDFCompiler converts HTML to PDF using rod (headless Chromium).
// On first use rod auto-downloads a compatible Chromium binary.
// Set ROD_BROWSER_BIN to skip the download and use an existing browser.
type RodHTMLToPDFCompiler struct {
	logger *zap.SugaredLogger
}

// NewRodHTMLToPDFCompiler creates a new rod-based HTML-to-PDF compiler.
func NewRodHTMLToPDFCompiler(logger *zap.SugaredLogger) *RodHTMLToPDFCompiler {
	return &RodHTMLToPDFCompiler{logger: logger}
}

// Compile converts HTML content to a PDF file at outputPath.
func (c *RodHTMLToPDFCompiler) Compile(htmlContent, outputPath string) error {
	pdfBytes, err := c.CompileToBytes(htmlContent)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	c.logger.Infof("Successfully converted HTML to PDF: %s", outputPath)
	return nil
}

// CompileToBytes converts HTML content to PDF and returns the raw bytes.
func (c *RodHTMLToPDFCompiler) CompileToBytes(htmlContent string) ([]byte, error) {
	l := launcher.New()

	// Respect ROD_BROWSER_BIN for CI or pre-installed browsers
	if bin := os.Getenv("ROD_BROWSER_BIN"); bin != "" {
		l = l.Bin(bin)
		c.logger.Infof("Using browser from ROD_BROWSER_BIN: %s", bin)
	}

	u, err := l.Headless(true).Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	if err := page.SetDocumentContent(htmlContent); err != nil {
		return nil, fmt.Errorf("failed to set page content: %w", err)
	}

	if err := page.WaitStable(300 * time.Millisecond); err != nil {
		c.logger.Warnf("Page stability wait timed out, proceeding anyway: %v", err)
	}

	pdf, err := page.PDF(&proto.PagePrintToPDF{
		PrintBackground:   true,
		PreferCSSPageSize: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	pdfBytes, err := io.ReadAll(pdf)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF bytes: %w", err)
	}

	return pdfBytes, nil
}
