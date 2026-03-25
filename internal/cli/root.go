package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	LaTeXEngine string

	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func initRootCmd() {
	initInitCmd()
	initGenerateCmd()
	initTemplatesCmd()
	initConfigCmd()
	initAICmd()
}

var rootCmd = &cobra.Command{
	Use:     "incipit",
	Short:   "Here begins the new career. Generate resumes from structured data using templates.",
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate),
}

func Execute() error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
	initRootCmd()
	return rootCmd.Execute()
}
