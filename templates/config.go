package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urmzd/incipit/utils"
	"gopkg.in/yaml.v3"
)

// Config represents the incipit configuration file that tracks
// registered templates, their filesystem locations, and user defaults.
type Config struct {
	Defaults  Defaults        `yaml:"defaults"`
	Templates []TemplateEntry `yaml:"templates"`
}

// Defaults holds user preferences applied when CLI flags are omitted.
type Defaults struct {
	Templates   []string `yaml:"templates,omitempty"`    // template refs (name or name:version) to use when -t is omitted
	LaTeXEngine string   `yaml:"latex_engine,omitempty"` // preferred LaTeX engine
	OutputDir   string   `yaml:"output_dir,omitempty"`   // default output directory
}

// TemplateEntry is a reference to a versioned template directory on the filesystem.
type TemplateEntry struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version,omitempty"`
	Path    string `yaml:"path"`
}

// ParseTemplateRef splits a "name:version" reference into its components.
// If no colon is present, version is returned as "".
func ParseTemplateRef(ref string) (name, version string) {
	ref = strings.TrimSpace(ref)
	if idx := strings.LastIndex(ref, ":"); idx > 0 && idx < len(ref)-1 {
		return ref[:idx], ref[idx+1:]
	}
	return ref, ""
}

// ConfigPath returns the path to the config file.
// e.g. ~/.config/incipit/config.yaml
func ConfigPath() string {
	base := utils.AppConfigDir()
	if base == "" {
		return ""
	}
	return filepath.Join(base, "config.yaml")
}

// LoadConfig reads and parses the config file. Returns an empty config
// if the file doesn't exist.
func LoadConfig() (*Config, error) {
	path := ConfigPath()
	if path == "" {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}

// SaveConfig writes the config to disk.
func SaveConfig(cfg *Config) error {
	path := ConfigPath()
	if path == "" {
		return fmt.Errorf("could not determine config path")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// Add registers a template by name, version, and filesystem path.
// If version is empty, it is read from the template's metadata.yml.
// Returns an error if the name+version pair is already registered or
// the path is invalid.
func (c *Config) Add(name, version, path string) error {
	absPath, err := utils.ResolvePath(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	metadataPath := filepath.Join(absPath, "metadata.yml")
	if !utils.FileExists(metadataPath) {
		return fmt.Errorf("template directory missing metadata.yml: %s", absPath)
	}

	for _, t := range c.Templates {
		if t.Name == name && t.Version == version {
			return fmt.Errorf("template %q version %q already registered (path: %s)", name, version, t.Path)
		}
	}

	c.Templates = append(c.Templates, TemplateEntry{Name: name, Version: version, Path: absPath})
	return nil
}

// Remove unregisters template entries. If version is "", all entries for the
// name are removed. If version is specified, only that exact entry is removed.
func (c *Config) Remove(name, version string) error {
	found := false
	var kept []TemplateEntry
	for _, t := range c.Templates {
		if t.Name == name && (version == "" || t.Version == version) {
			found = true
			continue
		}
		kept = append(kept, t)
	}
	if !found {
		if version == "" {
			return fmt.Errorf("template %q not found in config", name)
		}
		return fmt.Errorf("template %q version %q not found in config", name, version)
	}
	c.Templates = kept
	return nil
}

// List returns all registered template entries.
func (c *Config) List() []TemplateEntry {
	return c.Templates
}

// Lookup finds a template entry by name and version.
// If version is "", returns the entry with the highest semver for that name.
func (c *Config) Lookup(name, version string) *TemplateEntry {
	if version != "" {
		for i := range c.Templates {
			if c.Templates[i].Name == name && c.Templates[i].Version == version {
				return &c.Templates[i]
			}
		}
		return nil
	}

	// Find the entry with the highest version for this name.
	var best *TemplateEntry
	for i := range c.Templates {
		if c.Templates[i].Name != name {
			continue
		}
		if best == nil || compareVersions(c.Templates[i].Version, best.Version) > 0 {
			best = &c.Templates[i]
		}
	}
	return best
}

// Validate checks config integrity and returns warnings for issues
// like stale paths or missing metadata files.
func (c *Config) Validate() []string {
	var warnings []string

	type key struct{ name, version string }
	registered := make(map[key]bool)
	nameRegistered := make(map[string]bool)

	for _, t := range c.Templates {
		k := key{t.Name, t.Version}
		if registered[k] {
			warnings = append(warnings, fmt.Sprintf("duplicate template entry: %q version %q", t.Name, t.Version))
		}
		registered[k] = true
		nameRegistered[t.Name] = true

		if !utils.DirExists(t.Path) {
			warnings = append(warnings, fmt.Sprintf("template %q (v%s): path does not exist: %s", t.Name, t.Version, t.Path))
			continue
		}
		if !utils.FileExists(filepath.Join(t.Path, "metadata.yml")) {
			warnings = append(warnings, fmt.Sprintf("template %q (v%s): missing metadata.yml at %s", t.Name, t.Version, t.Path))
		}
	}

	for _, ref := range c.Defaults.Templates {
		refName, _ := ParseTemplateRef(ref)
		if !nameRegistered[refName] {
			warnings = append(warnings, fmt.Sprintf("default template %q is not registered", ref))
		}
	}

	return warnings
}

// compareVersions compares two semver-like version strings (X.Y.Z).
// Returns >0 if a > b, <0 if a < b, 0 if equal.
func compareVersions(a, b string) int {
	pa := parseVersionParts(a)
	pb := parseVersionParts(b)

	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			return pa[i] - pb[i]
		}
	}
	return 0
}

func parseVersionParts(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var result [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(parts[i])
		result[i] = n
	}
	return result
}
