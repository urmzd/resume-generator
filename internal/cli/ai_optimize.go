package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/ai"
	"github.com/urmzd/incipit/utils"
)

var aiOptimizeJob string

func initAIOptimizeCmd() {
	aiCmd.AddCommand(aiOptimizeCmd)
	aiOptimizeCmd.Flags().StringVarP(&aiOptimizeJob, "job", "j", "", "Job description text or path to a file containing the job description")
}

var aiOptimizeCmd = &cobra.Command{
	Use:   "optimize [file]",
	Short: "Print an LLM prompt that optimizes resume content, optionally for a job description",
	Args:  cobra.ExactArgs(1),
	Long: `Print a self-contained optimization prompt with the resume embedded.

The prompt asks the LLM to strengthen bullet points, add metrics, and
tighten wording while preserving the JSON structure. If a job description
is provided (via --job), it also asks the LLM to incorporate relevant
keywords and emphasize matching experience.

Save the LLM's JSON output to a file and run: incipit generate <file>

Examples:
  incipit ai optimize resume.json
  incipit ai optimize resume.json --job "Senior Go developer with 5+ years..."
  incipit ai optimize resume.json --job job-description.txt | claude -p > optimized.json`,
	Run: func(cmd *cobra.Command, args []string) {
		resumeJSON, err := readInputFile(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Resolve job description: could be inline text or a file path
		jobDesc := aiOptimizeJob
		if jobDesc != "" {
			if resolved, err := utils.ResolvePath(jobDesc); err == nil && utils.FileExists(resolved) {
				if data, err := os.ReadFile(resolved); err == nil {
					jobDesc = string(data)
				}
			}
		}

		fmt.Println(ai.OptimizePrompt(resumeJSON, jobDesc))
	},
}
