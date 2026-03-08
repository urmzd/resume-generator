package compilers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"go.uber.org/zap"
)

// ScreenshotHTML renders HTML content in a headless browser and saves a full-page PNG screenshot.
func ScreenshotHTML(logger *zap.SugaredLogger, htmlContent, outputPath string, width int) error {
	l := launcher.New()

	if bin := os.Getenv("ROD_BROWSER_BIN"); bin != "" {
		l = l.Bin(bin)
		logger.Infof("Using browser from ROD_BROWSER_BIN: %s", bin)
	}

	u, err := l.Headless(true).Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.MustClose()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:             width,
		Height:            900,
		DeviceScaleFactor: 2,
	}); err != nil {
		return fmt.Errorf("failed to set viewport: %w", err)
	}

	if err := page.SetDocumentContent(htmlContent); err != nil {
		return fmt.Errorf("failed to set page content: %w", err)
	}

	if err := page.WaitStable(300 * time.Millisecond); err != nil {
		logger.Warnf("Page stability wait timed out, proceeding anyway: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	screenshotData, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
	})
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	if err := os.WriteFile(outputPath, screenshotData, 0644); err != nil {
		return fmt.Errorf("failed to write screenshot: %w", err)
	}

	logger.Infof("Saved screenshot: %s", outputPath)
	return nil
}
