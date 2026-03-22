package cli

import (
	"fmt"
	"os"

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
	initCmd.Flags().StringVar(&initVersion, "version", "", "Template version to install (default: latest GitHub release)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing config and templates")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize incipit: create config and install default templates",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit init")

		version := initVersion
		if version == "" {
			ui.Info("Checking for latest release...")
			latest, err := templates.LatestVersion()
			if err != nil {
				ui.Error("Cannot determine version for template download.")
				ui.Detail("Specify a version with --version, e.g.: incipit init --version 1.0.0")
				os.Exit(1)
			}
			version = latest
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

		// Download and install templates into versioned subdirectories
		ui.Infof("Downloading templates (%s)...", version)
		installed, err := templates.Install(templates.InstallOptions{
			Version: version,
			Force:   initForce,
		})
		if err != nil {
			ui.Errorf("Failed to install templates: %v", err)
			os.Exit(1)
		}

		// Register installed templates
		newCfg := &templates.Config{}
		for _, tmpl := range installed {
			_ = newCfg.Add(tmpl.Name, tmpl.Version, tmpl.Path)
			ui.PhaseOk(fmt.Sprintf("Installed %s:%s", tmpl.Name, tmpl.Version), tmpl.Path)
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
