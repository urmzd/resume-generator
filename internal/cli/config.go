package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/templates"
	"github.com/urmzd/incipit/utils"
)

func initConfigCmd() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and modify configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit config show")

		configPath := templates.ConfigPath()
		if configPath == "" {
			ui.Error("Could not determine config path")
			os.Exit(1)
		}
		ui.Infof("Config: %s", configPath)

		cfg, err := templates.LoadConfig()
		if err != nil {
			ui.Errorf("Failed to load config: %v", err)
			os.Exit(1)
		}

		// Show defaults
		ui.Blank()
		ui.Section("Defaults")
		if len(cfg.Defaults.Templates) > 0 {
			ui.Detail(fmt.Sprintf("templates: %s", strings.Join(cfg.Defaults.Templates, ", ")))
		} else {
			ui.Detail("templates: (all)")
		}
		if cfg.Defaults.LaTeXEngine != "" {
			ui.Detail(fmt.Sprintf("latex_engine: %s", cfg.Defaults.LaTeXEngine))
		} else {
			ui.Detail("latex_engine: (auto-detect)")
		}
		if cfg.Defaults.OutputDir != "" {
			ui.Detail(fmt.Sprintf("output_dir: %s", cfg.Defaults.OutputDir))
		} else {
			ui.Detail("output_dir: (current directory)")
		}

		// Show registered templates
		ui.Blank()
		ui.Section("Templates")
		if len(cfg.Templates) == 0 {
			ui.Warn("No templates registered")
		} else {
			for _, t := range cfg.Templates {
				exists := "ok"
				if !utils.DirExists(t.Path) {
					exists = "missing"
				}
				ui.Detail(fmt.Sprintf("%s → %s [%s]", t.Name, t.Path, exists))
			}
		}

		// Show validation warnings
		warnings := cfg.Validate()
		if len(warnings) > 0 {
			ui.Blank()
			ui.Section("Warnings")
			for _, w := range warnings {
				ui.Warn(w)
			}
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Supported keys:
  defaults.templates       Comma-separated list of default template names
  defaults.latex_engine    Preferred LaTeX engine (xelatex, pdflatex, lualatex, latex)
  defaults.output_dir      Default output directory`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		cfg, err := templates.LoadConfig()
		if err != nil {
			ui.Errorf("Failed to load config: %v", err)
			os.Exit(1)
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
					ui.Errorf("Failed to resolve path: %v", err)
					os.Exit(1)
				}
				value = resolved
			}
			cfg.Defaults.OutputDir = value
		default:
			ui.Errorf("Unknown config key: %s", key)
			ui.Detail("Supported keys: defaults.templates, defaults.latex_engine, defaults.output_dir")
			os.Exit(1)
		}

		if err := templates.SaveConfig(cfg); err != nil {
			ui.Errorf("Failed to save config: %v", err)
			os.Exit(1)
		}

		ui.PhaseOk(fmt.Sprintf("Set %s", key), value)
	},
}
