package cmd

import (
	"github.com/spf13/cobra"
)

var (
	InputFile     string
	OutputFile    string
	TemplateFile  string
	KeepTex       bool
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
	Short: "Generate beautiful LaTeX resumes with one command.",
}

func Execute() error {
	initRootCmd()
	return rootCmd.Execute()
}
