package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/internal/ui"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
)

func initPreviewCmd() {
	rootCmd.AddCommand(previewCmd)
}

var previewCmd = &cobra.Command{
	Use:   "preview [file]",
	Short: "Preview a resume configuration without generating output",
	Long: `Preview command loads and validates a resume configuration file,
then displays a summary of the contents without generating any output files.
This is useful for quickly checking if your configuration is valid.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.Header("resume-generator preview")

		// Resolve file path
		filePath, err := utils.ResolvePath(args[0])
		if err != nil {
			ui.Errorf("Error resolving file path: %v", err)
			os.Exit(1)
		}
		if !utils.FileExists(filePath) {
			ui.Errorf("File does not exist: %s", filePath)
			os.Exit(1)
		}

		ui.Infof("Loading %s", filePath)

		// Load using unified adapter
		inputData, err := resume.LoadResumeFromFile(filePath)
		if err != nil {
			ui.Errorf("Error loading resume: %v", err)
			os.Exit(1)
		}

		// Validate
		if err := inputData.Validate(); err != nil {
			ui.Errorf("Validation error: %v", err)
			os.Exit(1)
		}

		ui.PhaseOk("Loaded and validated", fmt.Sprintf("format: %s", inputData.GetFormat()))

		// Convert to the runtime resume structure for preview
		r := inputData.ToResume()

		ui.Blank()

		// Contact Info
		ui.Section("Contact Information")
		ui.Detail(fmt.Sprintf("Name:     %s", r.Contact.Name))
		ui.Detail(fmt.Sprintf("Email:    %s", r.Contact.Email))
		if r.Contact.Phone != "" {
			ui.Detail(fmt.Sprintf("Phone:    %s", r.Contact.Phone))
		}
		if r.Contact.Location != nil {
			ui.Detail(fmt.Sprintf("Location: %s, %s", r.Contact.Location.City, r.Contact.Location.State))
		}
		if len(r.Contact.Links) > 0 {
			ui.Detail(fmt.Sprintf("Links:    %d link(s)", len(r.Contact.Links)))
		}

		// Skills
		if len(r.Skills.Categories) > 0 {
			ui.Blank()
			ui.Section("Skills")
			totalSkills := 0
			for _, cat := range r.Skills.Categories {
				totalSkills += len(cat.Items)
			}
			ui.Detail(fmt.Sprintf("%d categories, %d total skills", len(r.Skills.Categories), totalSkills))
			for _, cat := range r.Skills.Categories {
				ui.Detail(fmt.Sprintf("- %s: %d skills", cat.Category, len(cat.Items)))
			}
		}

		// Experience
		if len(r.Experience.Positions) > 0 {
			ui.Blank()
			ui.Section("Experience")
			ui.Detail(fmt.Sprintf("%d position(s)", len(r.Experience.Positions)))
			for _, exp := range r.Experience.Positions {
				ui.Detail(fmt.Sprintf("- %s at %s", exp.Title, exp.Company))
			}
		}

		// Projects
		if r.Projects != nil && len(r.Projects.Projects) > 0 {
			ui.Blank()
			ui.Section("Projects")
			ui.Detail(fmt.Sprintf("%d project(s)", len(r.Projects.Projects)))
			for _, proj := range r.Projects.Projects {
				ui.Detail(fmt.Sprintf("- %s", proj.Name))
			}
		}

		// Education
		if len(r.Education.Institutions) > 0 {
			ui.Blank()
			ui.Section("Education")
			ui.Detail(fmt.Sprintf("%d institution(s)", len(r.Education.Institutions)))
			for _, edu := range r.Education.Institutions {
				ui.Detail(fmt.Sprintf("- %s from %s", edu.Degree, edu.Institution))
			}
		}

		ui.Blank()
		ui.PhaseOk("Configuration is valid", "")
		ui.Infof("To generate output: resume-generator run -i %s", filePath)

		// If --json flag is set, show full JSON
		if verbose, _ := cmd.Flags().GetBool("json"); verbose {
			ui.Blank()
			fmt.Fprintf(os.Stderr, "  %s\n", strings.Repeat("─", 50))
			ui.Section("Full Configuration (JSON)")
			fmt.Fprintf(os.Stderr, "  %s\n", strings.Repeat("─", 50))
			jsonData, _ := json.MarshalIndent(r, "", "  ")
			// JSON data goes to stdout
			fmt.Println(string(jsonData))
		}
	},
}

func init() {
	previewCmd.Flags().Bool("json", false, "Show full configuration in JSON format")
}
