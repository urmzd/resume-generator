package services

import (
	"fmt"

	"github.com/urmzd/incipit/templates"
)

// InitOptions configures the initialization process.
type InitOptions struct {
	Version string
	Force   bool
}

// InitResult contains the outcome of initialization.
type InitResult struct {
	ConfigPath string
	Installed  []templates.InstalledTemplate
	Skipped    bool // true if already initialized and not forced
}

// Init initializes incipit by downloading templates and creating a config file.
func Init(opts InitOptions) (*InitResult, error) {
	version := opts.Version
	if version == "" {
		latest, err := templates.LatestVersion()
		if err != nil {
			return nil, fmt.Errorf("cannot determine version for template download: %w", err)
		}
		version = latest
	}

	configPath := templates.ConfigPath()
	if configPath == "" {
		return nil, fmt.Errorf("could not determine config directory")
	}

	cfg, err := templates.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load existing config: %w", err)
	}

	if len(cfg.Templates) > 0 && !opts.Force {
		return &InitResult{
			ConfigPath: configPath,
			Installed:  nil,
			Skipped:    true,
		}, nil
	}

	installed, err := templates.Install(templates.InstallOptions{
		Version: version,
		Force:   opts.Force,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to install templates: %w", err)
	}

	newCfg := &templates.Config{}
	for _, tmpl := range installed {
		_ = newCfg.Add(tmpl.Name, tmpl.Version, tmpl.Path)
	}

	newCfg.Defaults = templates.Defaults{
		Templates: []string{"modern-html"},
	}

	if err := templates.SaveConfig(newCfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return &InitResult{
		ConfigPath: configPath,
		Installed:  installed,
	}, nil
}
