package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/ai"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/utils"
)

var (
	aiOptimizeJob    string
	aiOptimizeOutput string
)

func initAIOptimizeCmd() {
	aiCmd.AddCommand(aiOptimizeCmd)
	aiOptimizeCmd.Flags().StringVarP(&aiOptimizeJob, "job", "j", "", "Job description text or path to a file containing the job description")
	aiOptimizeCmd.Flags().StringVarP(&aiOptimizeOutput, "output", "o", "", "Output JSON file path (default: <input>-optimized.json)")
}

var aiOptimizeCmd = &cobra.Command{
	Use:   "optimize [file]",
	Short: "Optimize resume content using an LLM, optionally targeting a job description",
	Args:  cobra.ExactArgs(1),
	Long: `Optimize a resume by improving bullet points, adding metrics, and tailoring
content for a specific role.

If a job description is provided (via --job), the optimizer will incorporate
relevant keywords and emphasize matching experience.

Examples:
  incipit ai optimize resume.json
  incipit ai optimize resume.json --job "Senior Go developer with 5+ years..."
  incipit ai optimize resume.json --job job-description.txt -o optimized.json`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit ai optimize")

		inputPath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving input path: %s", err)
			os.Exit(1)
		}
		if !utils.FileExists(inputPath) {
			ui.Errorf("Input file does not exist: %s", inputPath)
			os.Exit(1)
		}

		resumeBytes, err := os.ReadFile(inputPath)
		if err != nil {
			ui.Errorf("Error reading input file: %s", err)
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

		opts := aiProviderOpts()
		if jobDesc != "" {
			ui.Infof("Optimizing resume for target job description...")
		} else {
			ui.Infof("Optimizing resume content...")
		}

		ctx := context.Background()
		result, err := ai.Optimize(ctx, string(resumeBytes), jobDesc, opts)
		if err != nil {
			ui.Errorf("Optimization failed: %s", err)
			os.Exit(1)
		}

		outPath := aiOptimizeOutput
		if outPath == "" {
			dir := filepath.Dir(inputPath)
			base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
			outPath = filepath.Join(dir, base+"-optimized.json")
		}

		if err := os.WriteFile(outPath, []byte(result.JSON+"\n"), 0644); err != nil {
			ui.Errorf("Failed to write output file: %s", err)
			os.Exit(1)
		}

		ui.PhaseOk("Optimized", fmt.Sprintf("%s → %s", filepath.Base(inputPath), outPath))
		ui.Infof("Review the optimized JSON and run: incipit run %s", outPath)
	},
}
