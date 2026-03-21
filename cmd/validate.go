package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/internal/ui"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
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

		ui.Header("resume-generator validate")

		// Resolve file path
		filePath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving file path: %v", err)
			os.Exit(1)
		}
		if !utils.FileExists(filePath) {
			ui.Errorf("File does not exist: %s", filePath)
			os.Exit(1)
		}

		ui.Infof("Loading %s", filePath)

		inputData, err := resume.LoadResumeFromFile(filePath)
		if err != nil {
			ui.Errorf("Failed to load resume data: %v", err)
			os.Exit(1)
		}
		ui.PhaseOk("Loaded resume data", "")

		resumeData := inputData.ToResume()
		errors := resume.Validate(resumeData)
		if len(errors) > 0 {
			ui.Errorf("Validation failed with %d errors:", len(errors))
			for _, e := range errors {
				ui.Detail(fmt.Sprintf("Field: %s, Message: %s", e.Field, e.Message))
			}
			os.Exit(1)
		}

		ui.PhaseOk("Validation successful", "")

		_ = sugar // keep logger available for debug
	},
}
