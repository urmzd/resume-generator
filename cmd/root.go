package cmd

import (
	"github.com/spf13/cobra"
)

var (
	InputFile     string
	GeneratorType string
	LaTeXEngine   string
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
	Use:   "resume-generator",
	Short: "Generate resumes from structured data using templates.",
}

func Execute() error {
	initRootCmd()
	return rootCmd.Execute()
}
