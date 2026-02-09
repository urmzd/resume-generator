package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, got string)
	}{
		{
			name:    "empty path returns empty",
			input:   "",
			wantErr: false,
			check: func(t *testing.T, got string) {
				if got != "" {
					t.Errorf("ResolvePath() = %q, want empty string", got)
				}
			},
		},
		{
			name:    "absolute path returns cleaned",
			input:   "/tmp/test",
			wantErr: false,
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("ResolvePath() = %q, expected absolute path", got)
				}
			},
		},
		{
			name:    "relative path becomes absolute",
			input:   "./test",
			wantErr: false,
			check: func(t *testing.T, got string) {
				if !filepath.IsAbs(got) {
					t.Errorf("ResolvePath() = %q, expected absolute path", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolvePath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestResolvePathWithHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get home directory")
	}

	got, err := ResolvePath("~/test")
	if err != nil {
		t.Fatalf("ResolvePath() error = %v", err)
	}

	expected := filepath.Join(home, "test")
	if got != expected {
		t.Errorf("ResolvePath(~/test) = %q, want %q", got, expected)
	}
}

func TestResolveOutputPath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		input     string
		createDir bool
		setup     func(t *testing.T, tmpDir string) string
		wantErr   bool
	}{
		{
			name:      "existing directory",
			createDir: false,
			setup: func(t *testing.T, tmpDir string) string {
				return tmpDir
			},
			wantErr: false,
		},
		{
			name:      "new file path with createDir true",
			createDir: true,
			setup: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "subdir", "file.pdf")
			},
			wantErr: false,
		},
		{
			name:      "new file path with createDir false",
			createDir: false,
			setup: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "nonexistent", "file.pdf")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.setup(t, tmpDir)
			got, err := ResolveOutputPath(input, tt.createDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveOutputPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !filepath.IsAbs(got) {
				t.Errorf("ResolveOutputPath() = %q, expected absolute path", got)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "simple directory",
			path:    filepath.Join(tmpDir, "test"),
			wantErr: false,
		},
		{
			name:    "nested directory",
			path:    filepath.Join(tmpDir, "nested", "deep", "directory"),
			wantErr: false,
		},
		{
			name:    "already exists",
			path:    tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !DirExists(tt.path) {
				t.Errorf("EnsureDir() did not create directory %q", tt.path)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: testFile,
			want: true,
		},
		{
			name: "non-existent file",
			path: filepath.Join(tmpDir, "nonexistent.txt"),
			want: false,
		},
		{
			name: "directory instead of file",
			path: tmpDir,
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileExists(tt.path)
			if got != tt.want {
				t.Errorf("FileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing directory",
			path: tmpDir,
			want: true,
		},
		{
			name: "non-existent directory",
			path: filepath.Join(tmpDir, "nonexistent"),
			want: false,
		},
		{
			name: "file instead of directory",
			path: testFile,
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DirExists(tt.path)
			if got != tt.want {
				t.Errorf("DirExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestResolveAssetPath(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	assetPath := filepath.Join(tmpDir, "templates", "modern-html")
	if err := os.MkdirAll(assetPath, 0755); err != nil {
		t.Fatalf("Failed to create test asset: %v", err)
	}

	// Save current dir
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	// Change to temp dir
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantExist bool
	}{
		{
			name:      "existing relative path",
			input:     "templates/modern-html",
			wantErr:   false,
			wantExist: true,
		},
		{
			name:      "non-existent path returns clean path",
			input:     "nonexistent/path",
			wantErr:   false,
			wantExist: false,
		},
		{
			name:      "empty path",
			input:     "",
			wantErr:   false,
			wantExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveAssetPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveAssetPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.input == "" && got != "" {
				t.Errorf("ResolveAssetPath(\"\") = %q, want empty", got)
				return
			}

			if tt.wantExist && !DirExists(got) {
				t.Errorf("ResolveAssetPath() = %q, directory does not exist", got)
			}
		})
	}
}

func TestResolveAssetPathWithEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	assetPath := filepath.Join(tmpDir, "custom-templates")
	if err := os.MkdirAll(assetPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Set environment variable
	t.Setenv("RESUME_TEMPLATES_DIR", tmpDir)

	got, err := ResolveAssetPath("custom-templates")
	if err != nil {
		t.Fatalf("ResolveAssetPath() error = %v", err)
	}

	if !DirExists(got) {
		t.Errorf("ResolveAssetPath() with RESUME_TEMPLATES_DIR = %q, directory does not exist", got)
	}
}

func TestGetExecutableDir(t *testing.T) {
	got, err := GetExecutableDir()
	if err != nil {
		t.Errorf("GetExecutableDir() error = %v", err)
		return
	}

	if !filepath.IsAbs(got) {
		t.Errorf("GetExecutableDir() = %q, expected absolute path", got)
	}

	if !DirExists(got) {
		t.Errorf("GetExecutableDir() = %q, directory does not exist", got)
	}
}
