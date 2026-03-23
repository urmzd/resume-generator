package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/urmzd/saige/agent/provider/anthropic"
	"github.com/urmzd/saige/agent/provider/google"
	"github.com/urmzd/saige/agent/provider/ollama"
	"github.com/urmzd/saige/agent/provider/openai"
	"github.com/urmzd/saige/agent/types"
)

var defaultModels = map[string]string{
	"anthropic": "claude-sonnet-4-6-20250514",
	"openai":    "gpt-4o",
	"google":    "gemini-2.0-flash",
	"ollama":    "qwen3.5:4b",
}

// ProviderOptions configures LLM provider resolution.
type ProviderOptions struct {
	Provider   string // "anthropic", "openai", "google", "ollama", or "" for auto-detect
	Model      string // model override; uses provider default if empty
	OllamaHost string // Ollama endpoint (default http://localhost:11434)
}

// resolvedProvider returns the provider name, falling back to env-based auto-detection.
func (o ProviderOptions) resolvedProvider() string {
	if o.Provider != "" {
		return o.Provider
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		return "anthropic"
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		return "openai"
	}
	if os.Getenv("GOOGLE_API_KEY") != "" {
		return "google"
	}
	return "ollama"
}

// resolvedModel returns the model name, falling back to the provider default.
func (o ProviderOptions) resolvedModel() string {
	if o.Model != "" {
		return o.Model
	}
	return defaultModels[o.resolvedProvider()]
}

// ResolveProvider creates a types.Provider from the resolved options.
func ResolveProvider(ctx context.Context, opts ProviderOptions) (types.Provider, error) {
	name := opts.resolvedProvider()
	model := opts.resolvedModel()

	switch name {
	case "ollama":
		host := opts.OllamaHost
		if host == "" {
			host = "http://localhost:11434"
		}
		client := ollama.NewClient(host, model, "")
		return ollama.NewAdapter(client), nil

	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for openai provider")
		}
		return openai.NewAdapter(apiKey, model), nil

	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY is required for anthropic provider")
		}
		return anthropic.NewAdapter(apiKey, model), nil

	case "google":
		apiKey := os.Getenv("GOOGLE_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GOOGLE_API_KEY is required for google provider")
		}
		return google.NewAdapter(ctx, apiKey, model)

	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
