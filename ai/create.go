package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/urmzd/incipit/resume"
	agentsdk "github.com/urmzd/saige/agent"
	"github.com/urmzd/saige/agent/types"
)

// CreateResult holds the converted JSON output.
type CreateResult struct {
	JSON string
}

// Create converts freeform plain text resume content to structured JSON
// using an LLM with structured output when available.
func Create(ctx context.Context, plainText string, opts ProviderOptions) (*CreateResult, error) {
	if strings.TrimSpace(plainText) == "" {
		return nil, fmt.Errorf("input text is empty")
	}

	provider, err := ResolveProvider(ctx, opts)
	if err != nil {
		return nil, err
	}

	schema := ResumeSchema()

	agent := agentsdk.NewAgent(agentsdk.AgentConfig{
		Name:           "resume-creator",
		Provider:       provider,
		MaxIter:        1,
		SystemPrompt:   createSystemPrompt,
		ResponseSchema: schema,
	})

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	stream := agent.Invoke(ctx, []types.Message{
		types.NewUserMessage(buildCreatePrompt(plainText)),
	})

	var output strings.Builder
	for delta := range stream.Deltas() {
		if td, ok := delta.(types.TextContentDelta); ok {
			output.WriteString(td.Content)
		}
	}

	if err := stream.Wait(); err != nil {
		return nil, fmt.Errorf("creation failed: %w", err)
	}

	jsonStr, err := extractResumeJSON(output.String())
	if err != nil {
		return nil, fmt.Errorf("failed to extract valid resume from LLM response: %w\n\nRaw output:\n%s", err, output.String())
	}

	return &CreateResult{JSON: jsonStr}, nil
}

// extractResumeJSON tries to parse the output as JSON, stripping markdown fences if needed.
func extractResumeJSON(raw string) (string, error) {
	cleaned := strings.TrimSpace(raw)

	// Strip markdown code fences
	if strings.HasPrefix(cleaned, "```") {
		lines := strings.SplitN(cleaned, "\n", 2)
		if len(lines) == 2 {
			cleaned = lines[1]
		}
		if idx := strings.LastIndex(cleaned, "```"); idx >= 0 {
			cleaned = cleaned[:idx]
		}
		cleaned = strings.TrimSpace(cleaned)
	}

	// Validate it parses as a Resume
	var r resume.Resume
	if err := json.Unmarshal([]byte(cleaned), &r); err != nil {
		return "", fmt.Errorf("output is not valid JSON: %w", err)
	}

	if r.Contact.Name == "" {
		return "", fmt.Errorf("output JSON is missing required 'contact.name' field")
	}

	// Re-marshal for consistent formatting
	data, err := json.MarshalIndent(&r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}

	return string(data), nil
}

func buildCreatePrompt(plainText string) string {
	return fmt.Sprintf(`Convert this plain-text resume into structured JSON matching the resume schema.

## Resume Text

%s`, plainText)
}

const createSystemPrompt = `You are a resume data extraction assistant. Convert freeform resume text into structured JSON matching the resume schema.

Rules:
1. Output ONLY valid JSON — no explanations, no markdown fences, no commentary.
2. Every required field (contact.name, contact.email, skills, experience, education) MUST be present. If the text does not contain an email, use a placeholder like "update@me.com".
3. Dates: use "YYYY-MM" format when the exact day is unknown, "YYYY" when the month is unknown.
4. Never fabricate data that is not implied by the text.
5. Map freeform bullet points to "highlights" arrays within experience positions.
6. Group skills into sensible categories (e.g., "Programming Languages", "Frameworks", "Tools").
7. Preserve the section ordering from the input text where possible.
8. For date ranges, use: start and end fields. Omit end for current/ongoing roles.
9. For locations, include only fields present in the text (city, state, country).
10. Start your output directly with {. Do not wrap in code fences.`
