package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initVersionCmd() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("incipit %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	},
}
