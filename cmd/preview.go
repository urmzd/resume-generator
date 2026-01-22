package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
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
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		// Resolve file path
		filePath, err := utils.ResolvePath(args[0])
		if err != nil {
			sugar.Fatalf("Error resolving file path: %v", err)
		}
		if !utils.FileExists(filePath) {
			sugar.Fatalf("File does not exist: %s", filePath)
		}

		fmt.Printf("Loading resume configuration from: %s\n\n", filePath)

		// Load using unified adapter
		inputData, err := resume.LoadResumeFromFile(filePath)
		if err != nil {
			sugar.Fatalf("Error loading resume: %v", err)
		}

		// Validate
		if err := inputData.Validate(); err != nil {
			sugar.Fatalf("Validation error: %v", err)
		}

		// Convert to the runtime resume structure for preview
		resume := inputData.ToResume()

		// Display preview
		fmt.Println("┌─────────────────────────────────────────────┐")
		fmt.Println("│         Resume Configuration Preview        │")
		fmt.Println("└─────────────────────────────────────────────┘")
		fmt.Println()

		fmt.Printf("Format Type: %s\n\n", inputData.GetFormat())

		// Contact Info
		fmt.Println("📧 Contact Information:")
		fmt.Printf("  Name:     %s\n", resume.Contact.Name)
		fmt.Printf("  Email:    %s\n", resume.Contact.Email)
		if resume.Contact.Phone != "" {
			fmt.Printf("  Phone:    %s\n", resume.Contact.Phone)
		}
		if resume.Contact.Location != nil {
			fmt.Printf("  Location: %s, %s\n", resume.Contact.Location.City, resume.Contact.Location.State)
		}
		if len(resume.Contact.Links) > 0 {
			fmt.Printf("  Links:    %d link(s)\n", len(resume.Contact.Links))
		}
		fmt.Println()

		// Skills
		if len(resume.Skills.Categories) > 0 {
			fmt.Println("🔧 Skills:")
			totalSkills := 0
			for _, cat := range resume.Skills.Categories {
				totalSkills += len(cat.Items)
			}
			fmt.Printf("  %d categories, %d total skills\n", len(resume.Skills.Categories), totalSkills)
			for _, cat := range resume.Skills.Categories {
				fmt.Printf("    - %s: %d skills\n", cat.Category, len(cat.Items))
			}
			fmt.Println()
		}

		// Experience
		if len(resume.Experience.Positions) > 0 {
			fmt.Println("💼 Experience:")
			fmt.Printf("  %d position(s)\n", len(resume.Experience.Positions))
			for _, exp := range resume.Experience.Positions {
				fmt.Printf("    - %s at %s\n", exp.Title, exp.Company)
			}
			fmt.Println()
		}

		// Projects
		if resume.Projects != nil && len(resume.Projects.Projects) > 0 {
			fmt.Println("🚀 Projects:")
			fmt.Printf("  %d project(s)\n", len(resume.Projects.Projects))
			for _, proj := range resume.Projects.Projects {
				fmt.Printf("    - %s\n", proj.Name)
			}
			fmt.Println()
		}

		// Education
		if len(resume.Education.Institutions) > 0 {
			fmt.Println("🎓 Education:")
			fmt.Printf("  %d institution(s)\n", len(resume.Education.Institutions))
			for _, edu := range resume.Education.Institutions {
				fmt.Printf("    - %s from %s\n", edu.Degree, edu.Institution)
			}
			fmt.Println()
		}

		fmt.Println("✓ Configuration is valid!")
		fmt.Println()
		fmt.Println("To see the full configuration in JSON format, add --json flag")
		fmt.Printf("To generate output, use: resume-generator run -i %s\n", filePath)

		// If --json flag is set, show full JSON
		if verbose, _ := cmd.Flags().GetBool("json"); verbose {
			fmt.Println("\n" + strings.Repeat("─", 50))
			fmt.Println("Full Configuration (JSON):")
			fmt.Println(strings.Repeat("─", 50))
			jsonData, _ := json.MarshalIndent(resume, "", "  ")
			fmt.Println(string(jsonData))
		}
	},
}

func init() {
	previewCmd.Flags().Bool("json", false, "Show full configuration in JSON format")
}
