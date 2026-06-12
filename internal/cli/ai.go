package cli

import (
	"github.com/spf13/cobra"
)

func initAICmd() {
	rootCmd.AddCommand(aiCmd)

	initAIReviewCmd()
	initAIOptimizeCmd()
	initAICreateCmd()
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Generate LLM prompts for resume review, optimization, and creation",
	Long: `Generate ready-to-use LLM prompts for working with resumes.

Each subcommand prints a self-contained prompt to stdout with your resume
content embedded. Paste the prompt into any LLM, or pipe it directly to a
CLI agent:

  incipit ai review resume.json | claude -p

No API keys or network access are required.`,
}
