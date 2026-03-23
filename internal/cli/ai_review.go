package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/ai"
	"github.com/urmzd/incipit/utils"
)

func initAIReviewCmd() {
	aiCmd.AddCommand(aiReviewCmd)
}

var aiReviewCmd = &cobra.Command{
	Use:   "review [file]",
	Short: "Review and score a resume using specialist LLM agents",
	Args:  cobra.ExactArgs(1),
	Long: `Review a resume by delegating to four specialist sub-agents:

  - content-analyst:  achievement quantity, metrics, specificity, impact
  - writing-analyst:  succinctness, clarity, readability, grammar
  - industry-analyst: industry-specific keywords, conventions, relevance
  - format-analyst:   structure, section ordering, length, visual hierarchy

Each agent scores its dimension 1-10 with bullet-point feedback.
A coordinator synthesizes the results into a final report.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, err := utils.ResolvePath(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving input path: %s\n", err)
			os.Exit(1)
		}
		if !utils.FileExists(inputPath) {
			fmt.Fprintf(os.Stderr, "Input file does not exist: %s\n", inputPath)
			os.Exit(1)
		}

		yamlBytes, err := os.ReadFile(inputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input file: %s\n", err)
			os.Exit(1)
		}

		ctx := context.Background()
		result, err := ai.Review(ctx, string(yamlBytes), aiProviderOpts())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Review failed: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(result.Report)
	},
}
