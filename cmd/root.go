package cmd

import (
	"embed"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	InputFile     string
	GeneratorType string
	LaTeXEngine   string

	EmbeddedTemplatesFS embed.FS

	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func initRootCmd() {
	initRunCmd()
	initValidateCmd()
	initTemplatesCmd()
	initPreviewCmd()
	initSchemaCmd()
	rootCmd.PersistentFlags().StringVarP(&GeneratorType, "generator", "g", "base", "The type of generator to use (e.g., base, json-resume)")
}

var rootCmd = &cobra.Command{
	Use:     "resume-generator",
	Short:   "Generate resumes from structured data using templates.",
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate),
}

func Execute() error {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
	initRootCmd()
	return rootCmd.Execute()
}
