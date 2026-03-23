package ai

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	agentsdk "github.com/urmzd/saige/agent"
	"github.com/urmzd/saige/agent/types"
)

// ReviewResult is the structured output of a resume review.
type ReviewResult struct {
	ContentScore  float64          `json:"contentScore"`
	WritingScore  float64          `json:"writingScore"`
	IndustryScore float64          `json:"industryScore"`
	FormatScore   float64          `json:"formatScore"`
	OverallScore  float64          `json:"overallScore"`
	Report        string           `json:"report"`
	Categories    []CategoryResult `json:"categories"`
}

// CategoryResult holds the score and feedback for one assessment dimension.
type CategoryResult struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

// Review runs a multi-agent resume review and returns structured results.
func Review(ctx context.Context, resumeYAML string, opts ProviderOptions) (*ReviewResult, error) {
	provider, err := ResolveProvider(ctx, opts)
	if err != nil {
		return nil, err
	}

	agent := agentsdk.NewAgent(agentsdk.AgentConfig{
		Name:         "resume-coordinator",
		Provider:     provider,
		MaxIter:      10,
		SystemPrompt: coordinatorPrompt,
		SubAgents:    buildReviewSubAgents(provider),
	})

	prompt := fmt.Sprintf("Assess the following resume (in YAML format):\n\n---\n%s\n---", resumeYAML)

	stream := agent.Invoke(ctx, []types.Message{
		types.NewUserMessage(prompt),
	})

	var report strings.Builder
	for delta := range stream.Deltas() {
		if td, ok := delta.(types.TextContentDelta); ok {
			report.WriteString(td.Content)
		}
	}

	if err := stream.Wait(); err != nil {
		return nil, fmt.Errorf("review failed: %w", err)
	}

	return ParseReport(report.String()), nil
}

// ParseReport extracts structured scores from the raw assessment report text.
func ParseReport(report string) *ReviewResult {
	result := &ReviewResult{Report: report}

	scores := make([]float64, 4)
	for i, pattern := range scorePatterns {
		if m := pattern.FindStringSubmatch(report); len(m) > 1 {
			if v, err := strconv.ParseFloat(m[1], 64); err == nil {
				scores[i] = v
			}
		}
	}

	result.ContentScore = scores[0]
	result.WritingScore = scores[1]
	result.IndustryScore = scores[2]
	result.FormatScore = scores[3]

	// Weighted average: content 30%, industry 25%, writing 25%, format 20%
	result.OverallScore = scores[0]*0.30 + scores[2]*0.25 + scores[1]*0.25 + scores[3]*0.20

	for i, label := range categoryLabels {
		result.Categories = append(result.Categories, CategoryResult{
			Label: label,
			Score: scores[i],
		})
	}

	return result
}

var scorePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)CONTENT\s+SCORE:\s*(\d+(?:\.\d+)?)\s*/\s*10`),
	regexp.MustCompile(`(?i)WRITING\s+SCORE:\s*(\d+(?:\.\d+)?)\s*/\s*10`),
	regexp.MustCompile(`(?i)INDUSTRY\s+(?:FIT\s+)?SCORE:\s*(\d+(?:\.\d+)?)\s*/\s*10`),
	regexp.MustCompile(`(?i)STRUCTURE\s+SCORE:\s*(\d+(?:\.\d+)?)\s*/\s*10`),
}

var categoryLabels = []string{"Content Depth", "Writing Quality", "Industry Alignment", "Format & Flow"}

const coordinatorPrompt = `You are a senior resume review coordinator. You have four specialist analysts available.

Your process:
1. Read the resume carefully and identify the candidate's target industry/role.
2. Delegate to ALL FOUR analysts — content, writing, industry, and format — by calling each delegate tool. Pass the full resume text as the task to each one, prefixed with the target industry/role you identified.
3. After receiving all four reports, synthesize a final assessment that includes:
   - Target industry/role identified
   - Individual dimension scores (from each analyst)
   - Overall score (weighted average: content 30%, industry 25%, writing 25%, format 20%)
   - Top 3 priority improvements (the most impactful changes across all dimensions)

Always delegate to all four analysts. Do not skip any. Present the final report in a clean, readable format.`

func buildReviewSubAgents(provider types.Provider) []agentsdk.SubAgentDef {
	return []agentsdk.SubAgentDef{
		{
			Name:     "content_analyst",
			Provider: provider,
			MaxIter:  1,
			Description: "Analyzes resume content quality: achievement quantity, use of metrics/numbers, " +
				"specificity of accomplishments, and demonstrated impact. Delegate the full resume text to this agent.",
			SystemPrompt: `You are a resume content analyst. Score the resume on CONTENT (1-10) based on:

- **Quantity of achievements**: Does each role have 3-5 strong bullet points? Are there enough concrete accomplishments?
- **Metrics & numbers**: Are achievements quantified (percentages, dollar amounts, team sizes, timeframes)?
- **Specificity**: Are bullet points specific to this person's contribution, or generic/vague?
- **Impact**: Do bullet points show results and outcomes, not just responsibilities?

Output format:
CONTENT SCORE: X/10

Strengths:
- ...

Weaknesses:
- ...

Suggestions:
- ...

Be direct and specific. Reference actual bullet points from the resume.`,
		},
		{
			Name:     "writing_analyst",
			Provider: provider,
			MaxIter:  1,
			Description: "Analyzes resume writing quality: succinctness, clarity, readability, grammar, " +
				"and professional tone. Delegate the full resume text to this agent.",
			SystemPrompt: `You are a resume writing analyst. Score the resume on WRITING QUALITY (1-10) based on:

- **Succinctness**: Are bullet points concise (ideally 1-2 lines)? Is there unnecessary wordiness or filler?
- **Clarity**: Can a recruiter understand each bullet point in under 5 seconds? Is the language unambiguous?
- **Readability**: Is sentence structure varied? Are action verbs used consistently? Is parallel structure maintained?
- **Grammar & mechanics**: Any spelling errors, grammatical issues, or inconsistent punctuation/formatting?
- **Professional tone**: Is the language professional without being stiff or overly casual?

Output format:
WRITING SCORE: X/10

Strengths:
- ...

Weaknesses:
- ...

Suggestions:
- ...

Be direct and specific. Quote actual phrases from the resume that could be improved.`,
		},
		{
			Name:     "industry_analyst",
			Provider: provider,
			MaxIter:  1,
			Description: "Analyzes resume industry fit: relevant keywords, industry conventions, " +
				"role-specific expectations, and ATS compatibility. Delegate the full resume text with the target industry/role.",
			SystemPrompt: `You are a resume industry analyst. The task will include the target industry/role and the resume text. Score on INDUSTRY FIT (1-10) based on:

- **Keywords**: Does the resume include relevant industry/role keywords that ATS systems and recruiters look for?
- **Conventions**: Does the resume follow the norms for this industry (e.g., tech resumes emphasize projects and skills; sales resumes emphasize revenue and quotas; academic CVs emphasize publications)?
- **Role alignment**: Do the experiences and skills clearly map to the target role?
- **Skill relevance**: Are the listed skills current and valued in this industry? Are outdated or irrelevant skills cluttering the resume?
- **Competitive positioning**: How would this resume compare to a typical applicant pool for this role?

Output format:
INDUSTRY FIT SCORE: X/10

Target role/industry analyzed: ...

Strengths:
- ...

Weaknesses:
- ...

Missing keywords/skills:
- ...

Suggestions:
- ...

Be direct and specific to the industry identified.`,
		},
		{
			Name:     "format_analyst",
			Provider: provider,
			MaxIter:  1,
			Description: "Analyzes resume content structure: section ordering, information density, " +
				"completeness, and length appropriateness. Delegate the full resume text to this agent.",
			SystemPrompt: `You are a resume structure analyst. The visual formatting is handled automatically by a generator — do NOT evaluate fonts, spacing, bullet styles, or visual hierarchy. Instead, score the resume on STRUCTURE (1-10) based on its content organization:

- **Section ordering**: Are sections ordered by relevance to the target role? (Most impactful sections first)
- **Length**: Is the amount of content appropriate for the candidate's experience level? (Concise for <10 years, more detail acceptable for senior)
- **Information density**: Is there redundant or filler content that could be condensed or removed? Are there gaps where more detail is needed?
- **Section completeness**: Are expected sections present (contact, experience, education, skills)? Are any critical sections missing?
- **Logical flow**: Does the resume tell a coherent career story? Do sections build on each other logically?

Output format:
STRUCTURE SCORE: X/10

Strengths:
- ...

Weaknesses:
- ...

Suggestions:
- ...

Be direct and specific about structural improvements. Do not comment on visual formatting — only content organization.`,
		},
	}
}
