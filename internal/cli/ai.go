package cli

import (
	"github.com/spf13/cobra"
	"github.com/urmzd/incipit/ai"
)

var (
	aiProvider   string
	aiModel      string
	aiOllamaHost string
)

func initAICmd() {
	rootCmd.AddCommand(aiCmd)
	aiCmd.PersistentFlags().StringVarP(&aiProvider, "provider", "p", "", "LLM provider (anthropic|openai|google|ollama; auto-detected from API keys)")
	aiCmd.PersistentFlags().StringVarP(&aiModel, "model", "m", "", "Model name (uses provider default if omitted)")
	aiCmd.PersistentFlags().StringVar(&aiOllamaHost, "ollama-host", "http://localhost:11434", "Ollama host URL")

	initAIReviewCmd()
	initAIOptimizeCmd()
	initAICreateCmd()
}

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI-powered resume tools (review, optimize, create)",
	Long: `AI-powered resume tools that use LLMs to review, optimize, and create resumes.

Supports multiple providers: Anthropic (Claude), OpenAI (GPT), Google (Gemini), and Ollama (local).
Provider is auto-detected from API keys (ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY)
or falls back to Ollama.`,
}

// aiProviderOpts builds ProviderOptions from the shared CLI flags.
func aiProviderOpts() ai.ProviderOptions {
	return ai.ProviderOptions{
		Provider:   aiProvider,
		Model:      aiModel,
		OllamaHost: aiOllamaHost,
	}
}
