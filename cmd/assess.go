package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	agentsdk "github.com/urmzd/adk"
	"github.com/urmzd/adk/core"
	"github.com/urmzd/adk/provider/ollama"
	"github.com/urmzd/adk/tui"
	"golang.org/x/term"

	"github.com/spf13/cobra"
	"github.com/urmzd/resume-generator/pkg/utils"
	"go.uber.org/zap"
)

var (
	assessInput     string
	assessModel     string
	assessOllamaURL string
	assessVerbose   bool
)

func initAssessCmd() {
	rootCmd.AddCommand(assessCmd)
	assessCmd.Flags().StringVarP(&assessInput, "input", "i", "", "Path to the resume data file (e.g., resume.yml)")
	assessCmd.Flags().StringVarP(&assessModel, "model", "m", "qwen3.5:4b", "Ollama model to use for assessment")
	assessCmd.Flags().StringVar(&assessOllamaURL, "ollama-url", "http://localhost:11434", "Ollama server URL")
	assessCmd.Flags().BoolVarP(&assessVerbose, "verbose", "v", false, "Show full streaming output from all agents")

	_ = assessCmd.MarkFlagRequired("input")
}

var assessCmd = &cobra.Command{
	Use:   "assess",
	Short: "Rate and review a resume using specialized LLM agents via Ollama",
	Long: `Assess a resume by delegating to four specialist sub-agents:

  - content-analyst:  achievement quantity, metrics, specificity, impact
  - writing-analyst:  succinctness, clarity, readability, grammar
  - industry-analyst: industry-specific keywords, conventions, relevance
  - format-analyst:   structure, section ordering, length, visual hierarchy

Each agent scores its dimension 1-10 with bullet-point feedback.
A coordinator synthesizes the results into a final report.

Requires Ollama running locally (https://ollama.com).`,
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

		yamlBytes, err := os.ReadFile(inputPath)
		if err != nil {
			sugar.Fatalf("Error reading input file: %s", err)
		}
		resumeText := string(yamlBytes)

		// Check Ollama is reachable
		httpClient := &http.Client{Timeout: 5 * time.Second}
		resp, err := httpClient.Get(assessOllamaURL)
		if err != nil {
			sugar.Fatalf("Ollama is not available at %s. Install Ollama (https://ollama.com) and start it with 'ollama serve'.\n  Error: %v", assessOllamaURL, err)
		}
		_ = resp.Body.Close()

		client := ollama.NewClient(assessOllamaURL, assessModel, "")
		adapter := ollama.NewAdapter(client)

		agent := agentsdk.NewAgent(agentsdk.AgentConfig{
			Name:     "resume-coordinator",
			Provider: adapter,
			MaxIter:  10,
			SystemPrompt: `You are a senior resume review coordinator. You have four specialist analysts available.

Your process:
1. Read the resume carefully and identify the candidate's target industry/role.
2. Delegate to ALL FOUR analysts — content, writing, industry, and format — by calling each delegate tool. Pass the full resume text as the task to each one, prefixed with the target industry/role you identified.
3. After receiving all four reports, synthesize a final assessment that includes:
   - Target industry/role identified
   - Individual dimension scores (from each analyst)
   - Overall score (weighted average: content 30%, industry 25%, writing 25%, format 20%)
   - Top 3 priority improvements (the most impactful changes across all dimensions)

Always delegate to all four analysts. Do not skip any. Present the final report in a clean, readable format.`,
			SubAgents: buildAssessSubAgents(adapter),
		})

		prompt := fmt.Sprintf("Assess the following resume (in YAML format):\n\n---\n%s\n---", resumeText)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		stream := agent.Invoke(ctx, []core.Message{
			core.NewUserMessage(prompt),
		})

		header := buildAssessHeader(agent)
		isTTY := term.IsTerminal(int(os.Stdout.Fd()))

		if assessVerbose || !isTTY {
			runVerbose(stream, header, cancel)
		} else {
			runTUI(stream, header, cancel, sugar)
		}
	},
}

func buildAssessHeader(agent *agentsdk.Agent) tui.AgentHeader {
	info := agent.Info()
	header := tui.AgentHeader{
		Name:      info.Name,
		Provider:  info.Provider,
		Tools:     info.Tools,
		SubAgents: info.SubAgents,
	}
	tui.PopulateEnv(&header)
	return header
}

// runVerbose streams all agent output with colored prefixes using the SDK's StreamVerbose.
func runVerbose(stream *agentsdk.EventStream, header tui.AgentHeader, cancel context.CancelFunc) {
	result := tui.StreamVerbose(header, stream.Deltas(), os.Stdout)
	if result.Err != nil {
		cancel()
		fmt.Fprintf(os.Stderr, "Assessment error: %v\n", result.Err)
		os.Exit(1)
	}

	if err := stream.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Assessment failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
}

// runTUI runs the bubbletea progress UI, then prints the final report.
func runTUI(stream *agentsdk.EventStream, header tui.AgentHeader, cancel context.CancelFunc, sugar *zap.SugaredLogger) {
	model := tui.NewStreamModel(header, stream.Deltas())
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		cancel()
		sugar.Fatalf("TUI error: %v", err)
	}

	m := finalModel.(tui.StreamModel)
	if m.Err() != nil {
		cancel()
		sugar.Fatalf("Assessment error: %v", m.Err())
	}

	report := m.FinalReport()
	if report != "" {
		fmt.Println(tui.RenderReport("Resume Assessment", report))
	}

	if err := stream.Wait(); err != nil {
		sugar.Fatalf("Assessment failed: %v", err)
	}
}

func buildAssessSubAgents(provider core.Provider) []agentsdk.SubAgentDef {
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
