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

var aiCreateOutput string

func initAICreateCmd() {
	aiCmd.AddCommand(aiCreateCmd)
	aiCreateCmd.Flags().StringVarP(&aiCreateOutput, "output", "o", "", "Output JSON file path (default: same directory, .json extension)")
}

var aiCreateCmd = &cobra.Command{
	Use:   "create [file.txt]",
	Short: "Convert plain text or Markdown to structured resume JSON via LLM",
	Args:  cobra.ExactArgs(1),
	Long: `Convert a freeform plain-text resume into structured JSON that can be
used with 'incipit run' to generate polished PDFs.

The conversion uses an LLM to extract structured data from your plain text
and map it to the resume schema. The generated JSON file can be reviewed
and edited before generating output.

Examples:
  incipit ai create resume.txt
  incipit ai create resume.txt -o my-resume.json
  incipit ai create resume.txt -p anthropic -m claude-sonnet-4-6-20250514`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("incipit ai create")

		inputPath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving input path: %s", err)
			os.Exit(1)
		}
		if !utils.FileExists(inputPath) {
			ui.Errorf("Input file does not exist: %s", inputPath)
			os.Exit(1)
		}

		textBytes, err := os.ReadFile(inputPath)
		if err != nil {
			ui.Errorf("Error reading input file: %s", err)
			os.Exit(1)
		}

		opts := aiProviderOpts()
		ui.Infof("Creating structured JSON from plain text...")

		ctx := context.Background()
		result, err := ai.Create(ctx, string(textBytes), opts)
		if err != nil {
			ui.Errorf("Creation failed: %s", err)
			os.Exit(1)
		}

		outPath := aiCreateOutput
		if outPath == "" {
			dir := filepath.Dir(inputPath)
			base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
			outPath = filepath.Join(dir, base+".json")
		}

		if err := os.WriteFile(outPath, []byte(result.JSON+"\n"), 0644); err != nil {
			ui.Errorf("Failed to write output file: %s", err)
			os.Exit(1)
		}

		ui.PhaseOk("Created", fmt.Sprintf("%s → %s", filepath.Base(inputPath), outPath))
		ui.Infof("Review the generated JSON and run: incipit run %s", outPath)
	},
}
