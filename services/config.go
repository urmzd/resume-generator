package services

import (
	"fmt"
	"strings"

	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
)

// ConfigInfo contains the current configuration state.
type ConfigInfo struct {
	ConfigPath string
	Defaults   templates.Defaults
	Templates  []TemplateStatus
	Warnings   []string
}

// TemplateStatus describes a registered template and whether it exists on disk.
type TemplateStatus struct {
	Name    string
	Version string
	Path    string
	Exists  bool
}

// ShowConfig loads and returns the current configuration state.
func ShowConfig() (*ConfigInfo, error) {
	configPath := templates.ConfigPath()
	if configPath == "" {
		return nil, fmt.Errorf("could not determine config path")
	}

	cfg, err := templates.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	info := &ConfigInfo{
		ConfigPath: configPath,
		Defaults:   cfg.Defaults,
		Warnings:   cfg.Validate(),
	}

	for _, t := range cfg.Templates {
		info.Templates = append(info.Templates, TemplateStatus{
			Name:    t.Name,
			Version: t.Version,
			Path:    t.Path,
			Exists:  utils.DirExists(t.Path),
		})
	}

	return info, nil
}

// SetConfig updates a configuration value. Supported keys:
// defaults.templates, defaults.latex_engine, defaults.output_dir.
func SetConfig(key, value string) error {
	cfg, err := templates.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "defaults.templates":
		if value == "" {
			cfg.Defaults.Templates = nil
		} else {
			cfg.Defaults.Templates = strings.Split(value, ",")
			for i := range cfg.Defaults.Templates {
				cfg.Defaults.Templates[i] = strings.TrimSpace(cfg.Defaults.Templates[i])
			}
		}
	case "defaults.latex_engine":
		cfg.Defaults.LaTeXEngine = value
	case "defaults.output_dir":
		if value != "" {
			resolved, err := utils.ResolvePath(value)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}
			value = resolved
		}
		cfg.Defaults.OutputDir = value
	default:
		return fmt.Errorf("unknown config key: %s (supported: defaults.templates, defaults.latex_engine, defaults.output_dir)", key)
	}

	if err := templates.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
