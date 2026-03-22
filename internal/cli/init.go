package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/templates"
)

var (
	initVersion string
	initForce   bool
)

func initInitCmd() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&initVersion, "version", "", "Template version to install (default: binary version)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing config and templates")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize incipit: create config and install bundled templates",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit init")

		version := initVersion
		if version == "" {
			version = Version
		}
		if version == "" || version == "dev" {
			ui.Error("Cannot determine version for template download.")
			ui.Detail("Specify a version with --version, e.g.: incipit init --version 1.0.0")
			os.Exit(1)
		}

		// Ensure version has v prefix
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		configPath := templates.ConfigPath()
		if configPath == "" {
			ui.Error("Could not determine config directory")
			os.Exit(1)
		}

		// Check if already initialized
		cfg, err := templates.LoadConfig()
		if err != nil {
			ui.Errorf("Failed to load existing config: %v", err)
			os.Exit(1)
		}
		if len(cfg.Templates) > 0 && !initForce {
			ui.Infof("Already initialized (%d template(s) registered)", len(cfg.Templates))
			ui.Detail(fmt.Sprintf("Config: %s", configPath))
			ui.Detail("Use --force to re-initialize")
			return
		}

		// Download bundled templates
		ui.Infof("Downloading bundled templates (%s)...", version)
		installedDir, err := templates.Install(templates.InstallOptions{
			Version: version,
			Force:   initForce,
		})
		if err != nil {
			ui.Errorf("Failed to install templates: %v", err)
			os.Exit(1)
		}
		ui.PhaseOk("Templates downloaded", installedDir)

		// Discover and register installed templates
		entries, err := os.ReadDir(installedDir)
		if err != nil {
			ui.Errorf("Failed to read templates directory: %v", err)
			os.Exit(1)
		}

		newCfg := &templates.Config{}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			tmplPath := filepath.Join(installedDir, entry.Name())
			_ = newCfg.Add(entry.Name(), tmplPath)
		}

		// Set sensible defaults
		newCfg.Defaults = templates.Defaults{
			Templates: []string{"modern-html"},
		}

		if err := templates.SaveConfig(newCfg); err != nil {
			ui.Errorf("Failed to save config: %v", err)
			os.Exit(1)
		}

		ui.PhaseOk("Config created", configPath)
		ui.Infof("Registered %d template(s), default: modern-html", len(newCfg.Templates))
	},
}
