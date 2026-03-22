package templates

import (
	"fmt"
	"os"
	"path/filepath"

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
	Templates   []string `yaml:"templates,omitempty"`    // template names to use when -t is omitted
	LaTeXEngine string   `yaml:"latex_engine,omitempty"` // preferred LaTeX engine
	OutputDir   string   `yaml:"output_dir,omitempty"`   // default output directory
}

// TemplateEntry is a reference to a template directory on the filesystem.
type TemplateEntry struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
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

// Add registers a template by name and filesystem path. The path is
// resolved to an absolute path before storing. Returns an error if the
// name is already registered or the path is invalid.
func (c *Config) Add(name, path string) error {
	for _, t := range c.Templates {
		if t.Name == name {
			return fmt.Errorf("template %q already registered (path: %s)", name, t.Path)
		}
	}

	absPath, err := utils.ResolvePath(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	metadataPath := filepath.Join(absPath, "metadata.yml")
	if !utils.FileExists(metadataPath) {
		return fmt.Errorf("template directory missing metadata.yml: %s", absPath)
	}

	c.Templates = append(c.Templates, TemplateEntry{Name: name, Path: absPath})
	return nil
}

// Remove unregisters a template by name. Returns an error if the name
// is not found.
func (c *Config) Remove(name string) error {
	for i, t := range c.Templates {
		if t.Name == name {
			c.Templates = append(c.Templates[:i], c.Templates[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("template %q not found in config", name)
}

// List returns all registered template entries.
func (c *Config) List() []TemplateEntry {
	return c.Templates
}

// Lookup finds a template entry by name, returning nil if not found.
func (c *Config) Lookup(name string) *TemplateEntry {
	for i := range c.Templates {
		if c.Templates[i].Name == name {
			return &c.Templates[i]
		}
	}
	return nil
}

// Validate checks config integrity and returns warnings for issues
// like stale paths or missing metadata files.
func (c *Config) Validate() []string {
	var warnings []string

	registered := make(map[string]bool)
	for _, t := range c.Templates {
		if registered[t.Name] {
			warnings = append(warnings, fmt.Sprintf("duplicate template name: %q", t.Name))
		}
		registered[t.Name] = true

		if !utils.DirExists(t.Path) {
			warnings = append(warnings, fmt.Sprintf("template %q: path does not exist: %s", t.Name, t.Path))
			continue
		}
		if !utils.FileExists(filepath.Join(t.Path, "metadata.yml")) {
			warnings = append(warnings, fmt.Sprintf("template %q: missing metadata.yml at %s", t.Name, t.Path))
		}
	}

	for _, name := range c.Defaults.Templates {
		if !registered[name] {
			warnings = append(warnings, fmt.Sprintf("default template %q is not registered", name))
		}
	}

	return warnings
}
