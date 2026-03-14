package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	agentsdk "github.com/urmzd/agent-sdk"
	"github.com/urmzd/agent-sdk/core"
	"github.com/urmzd/agent-sdk/provider/ollama"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/generators"
	"github.com/urmzd/resume-generator/pkg/resume"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

var (
	assessInput     string
	assessModel     string
	assessOllamaURL string
)

func initAssessCmd() {
	rootCmd.AddCommand(assessCmd)
	assessCmd.Flags().StringVarP(&assessInput, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	assessCmd.Flags().StringVarP(&assessModel, "model", "m", "qwen3:4b", "Ollama model to use for assessment")
	assessCmd.Flags().StringVar(&assessOllamaURL, "ollama-url", "http://localhost:11434", "Ollama server URL")

	_ = assessCmd.MarkFlagRequired("input")

	generators.SetEmbeddedFS(EmbeddedTemplatesFS)
}

var assessCmd = &cobra.Command{
	Use:   "assess",
	Short: "Rate and review a resume using an LLM via Ollama",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()

		inputPath, err := utils.ResolvePath(assessInput)
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

		resumeData := inputData.ToResume()

		// Load the markdown template and render the resume
		mdTmpl, err := generators.LoadTemplate("modern-markdown")
		if err != nil {
			sugar.Fatalf("Failed to load markdown template: %v", err)
		}

		generator := generators.NewGenerator(sugar)
		markdownText, err := generator.GenerateWithTemplate(mdTmpl, resumeData)
		if err != nil {
			sugar.Fatalf("Failed to render resume to markdown: %v", err)
		}

		prompt := fmt.Sprintf(`Rate the following resume on a scale of 1-10. Provide:
1. Overall rating (1-10)
2. Key strengths (bullet points)
3. Key weaknesses (bullet points)
4. Specific suggestions for improvement

Resume:
---
%s
---`, markdownText)

		// Check Ollama is reachable before doing any LLM work
		httpClient := &http.Client{Timeout: 5 * time.Second}
		resp, err := httpClient.Get(assessOllamaURL)
		if err != nil {
			sugar.Fatalf("Ollama is not available at %s. Install Ollama (https://ollama.com) and start it with 'ollama serve'.\n  Error: %v", assessOllamaURL, err)
		}
		resp.Body.Close()

		client := ollama.NewClient(assessOllamaURL, assessModel, "")
		adapter := ollama.NewAdapter(client)

		agent := agentsdk.NewAgent(agentsdk.AgentConfig{
			Name:     "resume-assessor",
			Provider: adapter,
			MaxIter:  1,
		})

		stream := agent.Invoke(context.Background(), []core.Message{
			core.NewUserMessage(prompt),
		})

		for delta := range stream.Deltas() {
			switch d := delta.(type) {
			case core.TextContentDelta:
				fmt.Print(d.Content)
			case core.ErrorDelta:
				sugar.Fatalf("Assessment error: %v", d.Error)
			}
		}

		if err := stream.Wait(); err != nil {
			sugar.Fatalf("Assessment failed: %v", err)
		}

		fmt.Fprintln(os.Stdout)
	},
}
