package templates

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urmzd/incipit/utils"
)

// ConfigDir returns the platform-native templates directory.
// e.g. ~/.config/incipit/templates on Linux,
// ~/Library/Application Support/incipit/templates on macOS.
func ConfigDir() string {
	base := utils.AppConfigDir()
	if base == "" {
		return ""
	}
	return filepath.Join(base, "templates")
}

// InstallOptions configures template installation.
type InstallOptions struct {
	// Version is the release tag to download (e.g. "v1.2.0").
	Version string
	// Force overwrites existing templates if true.
	Force bool
}

// Install downloads templates from a GitHub release and extracts them
// to the platform config directory. Returns the installed path.
func Install(opts InstallOptions) (string, error) {
	if opts.Version == "" {
		return "", fmt.Errorf("version is required")
	}

	version := opts.Version
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	targetDir := ConfigDir()
	if targetDir == "" {
		return "", fmt.Errorf("could not determine config directory")
	}

	if utils.DirExists(targetDir) && !opts.Force {
		return targetDir, fmt.Errorf("templates already installed at %s (use Force to overwrite)", targetDir)
	}

	url := fmt.Sprintf("https://github.com/urmzd/incipit/releases/download/%s/templates.tar.gz", version)

	resp, err := http.Get(url) //nolint:gosec // URL constructed from known pattern
	if err != nil {
		return "", fmt.Errorf("failed to download templates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download templates: HTTP %d (url: %s)", resp.StatusCode, url)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create templates directory: %w", err)
	}

	if err := extractTarGz(resp.Body, targetDir); err != nil {
		return "", fmt.Errorf("failed to extract templates: %w", err)
	}

	return targetDir, nil
}

func extractTarGz(r io.Reader, destDir string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		target := filepath.Join(destDir, filepath.Clean(header.Name))

		// Prevent path traversal
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) && target != filepath.Clean(destDir) {
			return fmt.Errorf("invalid tar entry path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}
			if _, err := io.Copy(f, io.LimitReader(tr, 50<<20)); err != nil { // 50MB limit per file
				f.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
			f.Close()
		}
	}
	return nil
}
