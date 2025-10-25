package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath resolves a path to an absolute path.
// It supports:
// - Absolute paths (returned as-is after cleaning)
// - Relative paths (resolved from current working directory)
// - Paths with ~ (expanded to user home directory)
func ResolvePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	// Expand home directory
	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absPath), nil
}

// ResolveOutputPath resolves an output path, creating the directory if it doesn't exist.
// If the path is a directory, it returns the directory.
// If the path is a file, it ensures the parent directory exists and returns the file path.
func ResolveOutputPath(path string, createDir bool) (string, error) {
	resolved, err := ResolvePath(path)
	if err != nil {
		return "", err
	}

	if resolved == "" {
		return "", nil
	}

	// Check if path exists and is a directory
	info, err := os.Stat(resolved)
	if err == nil && info.IsDir() {
		return resolved, nil
	}

	// Path doesn't exist or is a file - ensure parent directory exists
	if createDir {
		dir := filepath.Dir(resolved)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	return resolved, nil
}

// GetExecutableDir returns the directory containing the executable.
// This is useful for finding assets relative to the binary location.
func GetExecutableDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	// Resolve symlinks
	exReal, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exReal), nil
}

// ResolveAssetPath tries to resolve an asset path by checking multiple locations:
// 1. Absolute path as provided
// 2. RESUME_TEMPLATES_DIR environment variable (if set, falls back to legacy RESUME_ASSETS_DIR)
// 3. Current working directory and its parents
// 4. Executable directory and its parents
// 5. Falls back to the cwd-relative path for clearer error messages
func ResolveAssetPath(relativePath string) (string, error) {
	if relativePath == "" {
		return "", nil
	}

	relativePath = filepath.Clean(relativePath)

	// Try as-is first
	if filepath.IsAbs(relativePath) {
		if _, err := os.Stat(relativePath); err == nil {
			return filepath.Clean(relativePath), nil
		}
	}

	var candidates []string

	// Allow users to override resource root via environment variable
	if assetsRoot := strings.TrimSpace(os.Getenv("RESUME_TEMPLATES_DIR")); assetsRoot != "" {
		if resolvedRoot, err := ResolvePath(assetsRoot); err == nil && resolvedRoot != "" {
			candidates = append(candidates, filepath.Join(resolvedRoot, relativePath))
		}
	} else if legacyRoot := strings.TrimSpace(os.Getenv("RESUME_ASSETS_DIR")); legacyRoot != "" {
		if resolvedRoot, err := ResolvePath(legacyRoot); err == nil && resolvedRoot != "" {
			candidates = append(candidates, filepath.Join(resolvedRoot, relativePath))
		}
	}

	// Try relative to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		for dir := cwd; dir != ""; dir = filepath.Dir(dir) {
			cwdPath := filepath.Join(dir, relativePath)
			candidates = append(candidates, cwdPath)

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	// Try relative to executable
	exDir, err := GetExecutableDir()
	if err == nil {
		for dir := exDir; dir != ""; dir = filepath.Dir(dir) {
			exPath := filepath.Join(dir, relativePath)
			candidates = append(candidates, exPath)

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	// Deduplicate and return the first path that exists
	seen := make(map[string]struct{})
	for _, candidate := range candidates {
		cleanCandidate := filepath.Clean(candidate)
		if _, ok := seen[cleanCandidate]; ok {
			continue
		}
		seen[cleanCandidate] = struct{}{}

		if _, err := os.Stat(cleanCandidate); err == nil {
			return cleanCandidate, nil
		}
	}

	// Return the cwd-relative path as fallback (will fail later if doesn't exist)
	if cwd != "" {
		return filepath.Clean(filepath.Join(cwd, relativePath)), nil
	}

	return filepath.Clean(relativePath), nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists checks if a file exists and is not a directory
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
