# AGENTS.md

## Project Overview

**incipit** is a Go CLI tool that converts structured resume data (JSON/Markdown) into PDF, HTML, LaTeX, DOCX, and Markdown output formats. It includes AI-powered commands for reviewing, optimizing, and creating resumes using multiple LLM providers (Anthropic, OpenAI, Google, Ollama).

## Repository Structure

```
.
├── cmd/incipit/main.go             # Entry point
├── internal/cli/                   # Cobra CLI commands
│   ├── root.go                     # Root command setup
│   ├── generate.go                 # `generate` command: JSON output, dry-run, schema
│   ├── ai.go                       # `ai` parent command + shared provider flags
│   ├── ai_review.go                # `ai review`: multi-agent resume assessment
│   ├── ai_optimize.go              # `ai optimize`: resume optimization for a role
│   ├── ai_create.go                # `ai create`: plain text to structured JSON
│   └── templates.go                # `templates list|validate|engines` subcommands
├── ai/                             # AI agent logic
│   ├── provider.go                 # Multi-provider resolution (Anthropic/OpenAI/Google/Ollama)
│   ├── schema.go                   # Resume JSON Schema to saige ParameterSchema converter
│   ├── review.go                   # Coordinator + 4 sub-agent review architecture
│   ├── optimize.go                 # Single-agent resume optimizer with structured output
│   └── create.go                   # Single-agent text to JSON converter with structured output
├── generators/                     # Template loading, formatters, HTML/LaTeX/MD/DOCX generators
├── compilers/                      # PDF compilation (LaTeX engines, Rod/Chromium)
├── resume/                         # Resume data model, validation, JSON/Markdown parsing
├── services/                       # High-level service layer
├── templates/                      # Built-in templates (modern-html, modern-latex, etc.)
├── assets/example_resumes/         # Example JSON resume files
├── skills/resume/                  # Agent skill definition
└── justfile                        # Task runner
```

## Architecture

### Data Flow

```
Input (JSON/Markdown) -> resume.LoadResumeFromFile() -> Resume struct
    -> Generator.GenerateWithTemplate(template, resume)
        -> Formatter.TemplateFuncs() provides template helpers
        -> text/template or html/template renders output
    -> Compiler (LaTeX->PDF or HTML->PDF via Rod/Chromium)
    -> Output file (.pdf, .html, .docx, .md)

AI review flow (multi-agent via saige):
    Input -> Resume JSON -> Coordinator agent
        -> delegate_to_content_analyst  (quantity, metrics, specificity, impact)
        -> delegate_to_writing_analyst  (succinctness, clarity, readability, grammar)
        -> delegate_to_industry_analyst (keywords, conventions, role fit, ATS)
        -> delegate_to_format_analyst   (structure, ordering, length, density)
    -> Coordinator synthesizes final scored report -> stdout

AI create/optimize flow (structured output via saige):
    Input -> plain text or resume JSON -> Agent with ResponseSchema
    -> LLM produces valid Resume JSON (constrained by schema)
    -> Output JSON file
```

### AI Provider Resolution

The `ai/` package supports multiple LLM providers, auto-detected from environment:

1. `ANTHROPIC_API_KEY` -> Anthropic (Claude)
2. `OPENAI_API_KEY` -> OpenAI (GPT)
3. `GOOGLE_API_KEY` -> Google (Gemini)
4. Fallback -> Ollama (local, no API key needed)

Override with `--provider` / `--model` flags on the `ai` parent command.

### Input Formats

Resume data is accepted as **JSON** or **Markdown**. Unrecognized file extensions fall through to the Markdown parser.

### Template System

Templates live in `templates/<name>/` with:
- `metadata.yml` -- metadata (name, format, description, tags)
- Template file (`template.html`, `template.tex`, `template.md`)
- Optional support files (`.cls` for LaTeX)

## Resume Data Model

See `resume/resume.go` for the full struct. Key types:

- `Resume` -- Contact, Summary, Skills, Experience, Projects, Education, Languages, Certifications, Layout
- `PartialDate` -- date with precision (year, month, or full)
- `DateRange` -- Start (PartialDate), End (*PartialDate, nil = Present)

Date formats in JSON: `"2024"`, `"2024-01"`, or `"2024-01-15T00:00:00Z"`

## Build & Test

```bash
just install
go test ./...
gofmt -l .
golangci-lint run
```

## Commit Convention

Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, `ci:`

## Dependencies

- **saige** (`github.com/urmzd/saige`) -- streaming AI agent framework with multi-provider support and structured output. Used by the `ai` commands.
