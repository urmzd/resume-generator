package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/compilers"
	"github.com/urmzd/incipit/generators"
	"github.com/urmzd/incipit/internal/ui"
	"github.com/urmzd/incipit/resume"
	"github.com/urmzd/incipit/utils"
	"go.uber.org/zap"
)

var screenshotsOutputDir string

func initScreenshotsCmd() {
	rootCmd.AddCommand(screenshotsCmd)
	screenshotsCmd.Flags().StringVarP(&InputFile, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	screenshotsCmd.Flags().StringVarP(&screenshotsOutputDir, "output-dir", "o", "assets/example_results", "Directory to save screenshot PNGs")

	_ = screenshotsCmd.MarkFlagRequired("input")
}

var screenshotsCmd = &cobra.Command{
	Use:   "screenshots",
	Short: "Generate PNG screenshots of each HTML-capable template",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		ui.Header("incipit screenshots")

		inputPath, err := utils.ResolvePath(InputFile)
		if err != nil {
			ui.Errorf("Error resolving input path: %s", err)
			os.Exit(1)
		}
		if !utils.FileExists(inputPath) {
			ui.Errorf("Input file does not exist: %s", inputPath)
			os.Exit(1)
		}

		inputData, err := resume.LoadResumeFromFile(inputPath)
		if err != nil {
			ui.Errorf("Error loading resume data: %s", err)
			os.Exit(1)
		}
		if err := inputData.Validate(); err != nil {
			ui.Errorf("Validation error: %s", err)
			os.Exit(1)
		}
		resumeData := inputData.ToResume()
		sectionOrder := inputData.GetSectionOrder()
		td := generators.NewTemplateData(resumeData, sectionOrder)
		ui.PhaseOk("Loaded resume data", "")

		allTemplates, err := generators.ListTemplates()
		if err != nil {
			ui.Errorf("Failed to list templates: %v", err)
			os.Exit(1)
		}

		htmlFallback, err := generators.LoadTemplate("modern-html")
		if err != nil {
			ui.Errorf("Failed to load HTML fallback template: %v", err)
			os.Exit(1)
		}

		generator := generators.NewGenerator(sugar)

		outputDir, err := utils.ResolvePath(screenshotsOutputDir)
		if err != nil {
			ui.Errorf("Error resolving output directory: %s", err)
			os.Exit(1)
		}
		if err := utils.EnsureDir(outputDir); err != nil {
			ui.Errorf("Error creating output directory: %s", err)
			os.Exit(1)
		}

		ui.Infof("Generating screenshots for %d template(s)", len(allTemplates))

		for _, tmpl := range allTemplates {
			tmplPtr := &tmpl
			// For non-HTML templates, use the HTML fallback to render a screenshot
			if tmpl.Type != generators.TemplateTypeHTML {
				tmplPtr = htmlFallback
			}

			htmlContent, err := generator.GenerateWithTemplate(tmplPtr, td)
			if err != nil {
				ui.Errorf("Failed to generate HTML for template %s: %v", tmpl.Name, err)
				continue
			}

			outputPath := filepath.Join(outputDir, tmpl.Name+".png")
			if err := compilers.ScreenshotHTML(sugar, htmlContent, outputPath, 1200); err != nil {
				ui.Errorf("Failed to screenshot template %s: %v", tmpl.Name, err)
				continue
			}

			ui.PhaseOk(fmt.Sprintf("Screenshot: %s", tmpl.Name), outputPath)
		}

		ui.Blank()
		ui.PhaseOk("Screenshots complete", outputDir)
	},
}
