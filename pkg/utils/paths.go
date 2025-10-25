package utils

import (
	"os"
	"path/filepath"
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
// 1. The path as provided (if absolute or relative to cwd)
// 2. Relative to the executable directory
// 3. Relative to the current working directory
func ResolveAssetPath(relativePath string) (string, error) {
	// Try as-is first
	if filepath.IsAbs(relativePath) {
		if _, err := os.Stat(relativePath); err == nil {
			return filepath.Clean(relativePath), nil
		}
	}

	// Try relative to current working directory
	cwd, err := os.Getwd()
	if err == nil {
		cwdPath := filepath.Join(cwd, relativePath)
		if _, err := os.Stat(cwdPath); err == nil {
			return filepath.Clean(cwdPath), nil
		}
	}

	// Try relative to executable
	exDir, err := GetExecutableDir()
	if err == nil {
		exPath := filepath.Join(exDir, relativePath)
		if _, err := os.Stat(exPath); err == nil {
			return filepath.Clean(exPath), nil
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
