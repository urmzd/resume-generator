package templates

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
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

// VersionedTemplateDir returns the path for a specific template version.
// e.g. ~/.config/incipit/templates/modern-html/1.0.0
func VersionedTemplateDir(name, version string) string {
	base := ConfigDir()
	if base == "" {
		return ""
	}
	return filepath.Join(base, name, version)
}

// LatestVersion queries the GitHub API for the latest release tag of incipit.
func LatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/urmzd/incipit/releases/latest") //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to query latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query latest release: HTTP %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse release response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no release tag found")
	}

	return release.TagName, nil
}

// InstallOptions configures template installation.
type InstallOptions struct {
	// Version is the release tag to download (e.g. "v1.2.0").
	Version string
	// Force overwrites existing template versions if true.
	Force bool
}

// InstalledTemplate describes a single template that was installed.
type InstalledTemplate struct {
	Name    string
	Version string
	Path    string
}

// Install downloads templates from a GitHub release, extracts them into
// versioned subdirectories (templates/<name>/<version>/), and returns
// the list of installed templates.
func Install(opts InstallOptions) ([]InstalledTemplate, error) {
	if opts.Version == "" {
		return nil, fmt.Errorf("version is required")
	}

	version := opts.Version
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	baseDir := ConfigDir()
	if baseDir == "" {
		return nil, fmt.Errorf("could not determine config directory")
	}

	url := fmt.Sprintf("https://github.com/urmzd/incipit/releases/download/%s/templates.tar.gz", version)

	resp, err := http.Get(url) //nolint:gosec // URL constructed from known pattern
	if err != nil {
		return nil, fmt.Errorf("failed to download templates: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download templates: HTTP %d (url: %s)", resp.StatusCode, url)
	}

	// Extract to a temp directory first, then move into versioned layout.
	tmpDir, err := os.MkdirTemp("", "incipit-templates-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	if err := extractTarGz(resp.Body, tmpDir); err != nil {
		return nil, fmt.Errorf("failed to extract templates: %w", err)
	}

	// Move each template subdirectory into templates/<name>/<version>/
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read extracted templates: %w", err)
	}

	var installed []InstalledTemplate
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		tmplName := entry.Name()
		srcDir := filepath.Join(tmpDir, tmplName)

		// Version comes from the release tag, not metadata.yml
		tmplVersion := strings.TrimPrefix(version, "v")

		destDir := VersionedTemplateDir(tmplName, tmplVersion)

		if utils.DirExists(destDir) {
			if !opts.Force {
				continue
			}
			if err := os.RemoveAll(destDir); err != nil {
				return nil, fmt.Errorf("failed to remove existing template %s:%s: %w", tmplName, tmplVersion, err)
			}
		}

		if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
			return nil, fmt.Errorf("failed to create template directory: %w", err)
		}

		if err := os.Rename(srcDir, destDir); err != nil {
			return nil, fmt.Errorf("failed to move template %s to %s: %w", tmplName, destDir, err)
		}

		installed = append(installed, InstalledTemplate{
			Name:    tmplName,
			Version: tmplVersion,
			Path:    destDir,
		})
	}

	return installed, nil
}

func extractTarGz(r io.Reader, destDir string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() { _ = gz.Close() }()

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
				_ = f.Close()
				return fmt.Errorf("failed to write file %s: %w", target, err)
			}
			if err := f.Close(); err != nil {
				return fmt.Errorf("failed to close file %s: %w", target, err)
			}
		}
	}
	return nil
}
