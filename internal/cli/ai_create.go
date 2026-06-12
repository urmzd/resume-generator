package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/ai"
)

func initAICreateCmd() {
	aiCmd.AddCommand(aiCreateCmd)
}

var aiCreateCmd = &cobra.Command{
	Use:   "create [file.txt]",
	Short: "Print an LLM prompt that converts plain text to structured resume JSON",
	Args:  cobra.ExactArgs(1),
	Long: `Print a self-contained conversion prompt with the plain-text resume and
the resume JSON Schema embedded.

The prompt asks the LLM to extract structured data from the text and map
it to the schema. Save the LLM's JSON output to a file, review it, and
run: incipit generate <file>

Examples:
  incipit ai create resume.txt
  incipit ai create resume.txt | claude -p > resume.json`,
	Run: func(cmd *cobra.Command, args []string) {
		plainText, err := readInputFile(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		prompt, err := ai.CreatePrompt(plainText)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(prompt)
	},
}
