package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	agentsdk "github.com/urmzd/saige/agent"
	"github.com/urmzd/saige/agent/types"
	"github.com/urmzd/incipit/resume"
)

// OptimizeResult holds the optimized resume JSON.
type OptimizeResult struct {
	JSON string
}

// Optimize improves resume content, optionally tailoring it for a specific job description.
func Optimize(ctx context.Context, resumeJSON string, jobDesc string, opts ProviderOptions) (*OptimizeResult, error) {
	if strings.TrimSpace(resumeJSON) == "" {
		return nil, fmt.Errorf("resume content is empty")
	}

	provider, err := ResolveProvider(ctx, opts)
	if err != nil {
		return nil, err
	}

	schema := ResumeSchema()

	agent := agentsdk.NewAgent(agentsdk.AgentConfig{
		Name:           "resume-optimizer",
		Provider:       provider,
		MaxIter:        1,
		SystemPrompt:   optimizeSystemPrompt,
		ResponseSchema: schema,
	})

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	prompt := buildOptimizePrompt(resumeJSON, jobDesc)

	stream := agent.Invoke(ctx, []types.Message{
		types.NewUserMessage(prompt),
	})

	var output strings.Builder
	for delta := range stream.Deltas() {
		if td, ok := delta.(types.TextContentDelta); ok {
			output.WriteString(td.Content)
		}
	}

	if err := stream.Wait(); err != nil {
		return nil, fmt.Errorf("optimization failed: %w", err)
	}

	jsonStr, err := extractOptimizedJSON(output.String())
	if err != nil {
		return nil, fmt.Errorf("failed to extract optimized resume: %w\n\nRaw output:\n%s", err, output.String())
	}

	return &OptimizeResult{JSON: jsonStr}, nil
}

func extractOptimizedJSON(raw string) (string, error) {
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

	var r resume.Resume
	if err := json.Unmarshal([]byte(cleaned), &r); err != nil {
		return "", fmt.Errorf("output is not valid JSON: %w", err)
	}

	data, err := json.MarshalIndent(&r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}

	return string(data), nil
}

func buildOptimizePrompt(resumeJSON string, jobDesc string) string {
	var sb strings.Builder
	sb.WriteString("Optimize the following resume:\n\n## Current Resume (JSON)\n\n")
	sb.WriteString(resumeJSON)

	if strings.TrimSpace(jobDesc) != "" {
		sb.WriteString("\n\n## Target Job Description\n\n")
		sb.WriteString(jobDesc)
	}

	return sb.String()
}

const optimizeSystemPrompt = `You are a professional resume optimizer. Given a resume in JSON format and optionally a target job description, improve the resume content while preserving its structure.

Your improvements should focus on:
1. **Stronger bullet points**: Rewrite vague or weak highlights to be specific, quantified, and impact-driven. Use the XYZ formula: "Accomplished [X] as measured by [Y], by doing [Z]".
2. **Better metrics**: Add or improve quantitative measures (percentages, dollar amounts, team sizes, timeframes) where the data supports it. Do not fabricate numbers.
3. **Keyword optimization**: If a job description is provided, incorporate relevant keywords and skills naturally into the experience highlights and skills sections.
4. **Conciseness**: Tighten wordy bullet points. Each should be 1-2 lines max.
5. **Action verbs**: Start each bullet point with a strong action verb. Avoid repeating the same verb.
6. **Skill relevance**: If a job description is provided, prioritize skills and experiences that align with the target role.

Rules:
- Do NOT change contact information, dates, company names, job titles, or education details.
- Do NOT add fabricated experiences or achievements.
- Preserve all existing sections and their ordering.
- Output ONLY valid JSON matching the resume schema. Start with {.
- Output the complete improved resume.`
