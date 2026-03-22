package main

import (
	"os"

	"github.com/urmzd/incipit/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Version = version
	cli.Commit = commit
	cli.BuildDate = date

	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
