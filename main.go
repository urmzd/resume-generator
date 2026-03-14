package main

import (
	"os"

	"github.com/urmzd/resume-generator/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.EmbeddedTemplatesFS = EmbeddedTemplates
	cmd.Version = version
	cmd.Commit = commit
	cmd.BuildDate = date

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
