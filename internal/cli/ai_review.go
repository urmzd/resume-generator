package cli

import (
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
	Short: "Print an LLM prompt that reviews and scores a resume",
	Args:  cobra.ExactArgs(1),
	Long: `Print a self-contained review prompt with the resume embedded.

The prompt asks the LLM to score four dimensions (content, writing,
industry fit, structure) 1-10 with bullet-point feedback, then synthesize
an overall score and top 3 priority improvements.

Examples:
  incipit ai review resume.json
  incipit ai review resume.json | claude -p`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := readInputFile(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(ai.ReviewPrompt(content))
	},
}

// readInputFile resolves and reads a CLI input path.
func readInputFile(path string) (string, error) {
	resolved, err := utils.ResolvePath(path)
	if err != nil {
		return "", fmt.Errorf("error resolving input path: %w", err)
	}
	if !utils.FileExists(resolved) {
		return "", fmt.Errorf("input file does not exist: %s", resolved)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %w", err)
	}
	return string(data), nil
}
