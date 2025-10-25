package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/definition"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func initValidateCmd() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a resume configuration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		// Resolve file path
		filePath, err := utils.ResolvePath(args[0])
		if err != nil {
			sugar.Fatalf("Error resolving file path: %v", err)
		}
		if !utils.FileExists(filePath) {
			sugar.Fatalf("File does not exist: %s", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			sugar.Fatalf("failed to read file: %v", err)
		}

		var resume definition.EnhancedResume
		if err := yaml.Unmarshal(content, &resume); err != nil {
			sugar.Fatalf("failed to unmarshal YAML: %v", err)
		}

		validator := &definition.ConfigurationValidator{
			StrictMode: true,
		}

		errors := validator.ValidateEnhancedResume(&resume)
		if len(errors) > 0 {
			sugar.Errorf("Validation failed with %d errors:", len(errors))
			for _, e := range errors {
				sugar.Errorf("  - Field: %s, Message: %s", e.Field, e.Message)
			}
			return
		}

		sugar.Info("Validation successful!")
	},
}
