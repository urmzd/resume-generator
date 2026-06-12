// Package ai builds ready-to-use LLM prompts for resume review, optimization,
// and creation. The `incipit ai` commands print these prompts to stdout so they
// can be pasted into — or piped to — any LLM or coding agent. No API keys or
// network access are required.
package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/urmzd/incipit/resume"
)

// ReviewPrompt returns a self-contained prompt asking an LLM to review and
// score the given resume across four dimensions.
func ReviewPrompt(resumeContent string) string {
	return fmt.Sprintf(`You are a senior resume reviewer. Assess the resume below across four dimensions, score each 1-10, and synthesize a final report.

## Dimensions

1. **Content** — achievement quantity (3-5 strong bullets per role), metrics & numbers (percentages, dollar amounts, team sizes, timeframes), specificity of accomplishments, and demonstrated impact (results, not responsibilities).
2. **Writing** — succinctness (1-2 lines per bullet), clarity (understandable in under 5 seconds), readability (varied structure, consistent action verbs, parallel structure), grammar & mechanics, and professional tone.
3. **Industry fit** — first identify the candidate's target industry/role, then assess relevant keywords (ATS and recruiter visibility), industry conventions, role alignment, skill relevance, and competitive positioning against a typical applicant pool.
4. **Structure** — section ordering by relevance, length appropriate to experience level, information density (no redundancy or gaps), section completeness (contact, experience, education, skills), and logical career-story flow. Do NOT evaluate visual formatting (fonts, spacing, bullet styles) — that is handled by a generator.

## Report format

For each dimension, output:

CONTENT SCORE: X/10 (likewise WRITING SCORE, INDUSTRY FIT SCORE, STRUCTURE SCORE)

followed by Strengths, Weaknesses, and Suggestions as bullet points. Be direct and specific — quote actual phrases and bullet points from the resume.

Then conclude with:
- Target industry/role identified
- OVERALL SCORE: X/10 (weighted: content 30%%, industry 25%%, writing 25%%, format 20%%)
- Top 3 priority improvements (the most impactful changes across all dimensions)

## Resume

%s`, resumeContent)
}

// OptimizePrompt returns a self-contained prompt asking an LLM to improve the
// given resume JSON, optionally tailoring it to a job description.
func OptimizePrompt(resumeJSON, jobDesc string) string {
	var sb strings.Builder
	sb.WriteString(`You are a professional resume optimizer. Improve the resume JSON below while preserving its structure.

Focus on:
1. **Stronger bullet points**: Rewrite vague or weak highlights to be specific, quantified, and impact-driven. Use the XYZ formula: "Accomplished [X] as measured by [Y], by doing [Z]".
2. **Better metrics**: Add or improve quantitative measures (percentages, dollar amounts, team sizes, timeframes) where the data supports it. Do not fabricate numbers.
3. **Keyword optimization**: If a job description is provided, incorporate relevant keywords and skills naturally into the experience highlights and skills sections.
4. **Conciseness**: Tighten wordy bullet points. Each should be 1-2 lines max.
5. **Action verbs**: Start each bullet point with a strong action verb. Avoid repeating the same verb.
6. **Skill relevance**: If a job description is provided, prioritize skills and experiences that align with the target role.

Rules:
- Do NOT change contact information, dates, company names, job titles, or education details.
- Do NOT add fabricated experiences or achievements.
- Preserve all existing sections, their ordering, and the JSON structure.
- Output ONLY the complete improved resume as valid JSON. Start with {. No markdown fences, no commentary.

## Current Resume (JSON)

`)
	sb.WriteString(resumeJSON)

	if strings.TrimSpace(jobDesc) != "" {
		sb.WriteString("\n\n## Target Job Description\n\n")
		sb.WriteString(jobDesc)
	}

	return sb.String()
}

// CreatePrompt returns a self-contained prompt asking an LLM to convert
// freeform resume text into structured JSON matching the resume schema.
func CreatePrompt(plainText string) (string, error) {
	schemaJSON, err := resumeSchemaJSON()
	if err != nil {
		return "", fmt.Errorf("failed to build resume schema: %w", err)
	}

	return fmt.Sprintf(`You are a resume data extraction assistant. Convert the plain-text resume below into structured JSON matching the JSON Schema that follows.

Rules:
1. Output ONLY valid JSON — no explanations, no markdown fences, no commentary. Start your output directly with {.
2. Every required field (contact.name, contact.email, skills, experience, education) MUST be present. If the text does not contain an email, use a placeholder like "update@me.com".
3. Dates: use "YYYY-MM" format when the exact day is unknown, "YYYY" when the month is unknown.
4. Never fabricate data that is not implied by the text.
5. Map freeform bullet points to "highlights" arrays within experience positions.
6. Group skills into sensible categories (e.g., "Programming Languages", "Frameworks", "Tools").
7. Preserve the section ordering from the input text where possible.
8. For date ranges, use start and end fields. Omit end for current/ongoing roles.
9. For locations, include only fields present in the text (city, state, country).

## Resume JSON Schema

%s

## Resume Text

%s`, schemaJSON, plainText), nil
}

// resumeSchemaJSON reflects the Resume struct into an indented JSON Schema string.
func resumeSchemaJSON() (string, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&resume.Resume{})
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
