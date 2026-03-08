package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/compilers"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

var screenshotsOutputDir string

func initScreenshotsCmd() {
	rootCmd.AddCommand(screenshotsCmd)
	screenshotsCmd.Flags().StringVarP(&InputFile, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	screenshotsCmd.Flags().StringVarP(&screenshotsOutputDir, "output-dir", "o", "assets/example_results", "Directory to save screenshot PNGs")

	_ = screenshotsCmd.MarkFlagRequired("input")

	generators.SetEmbeddedFS(EmbeddedTemplatesFS)
}

var screenshotsCmd = &cobra.Command{
	Use:   "screenshots",
	Short: "Generate PNG screenshots of each HTML-capable template",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		inputPath, err := utils.ResolvePath(InputFile)
		if err != nil {
			sugar.Fatalf("Error resolving input path: %s", err)
		}
		if !utils.FileExists(inputPath) {
			sugar.Fatalf("Input file does not exist: %s", inputPath)
		}

		inputData, err := resume.LoadResumeFromFile(inputPath)
		if err != nil {
			sugar.Fatalf("Error loading resume data: %s", err)
		}
		if err := inputData.Validate(); err != nil {
			sugar.Fatalf("Validation error: %s", err)
		}
		resumeData := inputData.ToResume()

		allTemplates, err := generators.ListTemplates()
		if err != nil {
			sugar.Fatalf("Failed to list templates: %v", err)
		}

		htmlFallback, err := generators.LoadTemplate("modern-html")
		if err != nil {
			sugar.Fatalf("Failed to load HTML fallback template: %v", err)
		}

		generator := generators.NewGenerator(sugar)

		outputDir, err := utils.ResolvePath(screenshotsOutputDir)
		if err != nil {
			sugar.Fatalf("Error resolving output directory: %s", err)
		}
		if err := utils.EnsureDir(outputDir); err != nil {
			sugar.Fatalf("Error creating output directory: %s", err)
		}

		for _, tmpl := range allTemplates {
			tmplPtr := &tmpl
			// For non-HTML templates, use the HTML fallback to render a screenshot
			if tmpl.Type != generators.TemplateTypeHTML {
				tmplPtr = htmlFallback
			}

			htmlContent, err := generator.GenerateWithTemplate(tmplPtr, resumeData)
			if err != nil {
				sugar.Errorf("Failed to generate HTML for template %s: %v", tmpl.Name, err)
				continue
			}

			outputPath := filepath.Join(outputDir, tmpl.Name+".png")
			if err := compilers.ScreenshotHTML(sugar, htmlContent, outputPath, 1200); err != nil {
				sugar.Errorf("Failed to screenshot template %s: %v", tmpl.Name, err)
				continue
			}

			sugar.Infof("Generated screenshot for %s: %s", tmpl.Name, outputPath)
		}
	},
}
