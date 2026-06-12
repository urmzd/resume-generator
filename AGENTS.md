# AGENTS.md

## Project Overview

**incipit** is a Go CLI tool that converts structured resume data (JSON/Markdown) into PDF, HTML, LaTeX, DOCX, and Markdown output formats. It includes `ai` commands that print self-contained LLM prompts for reviewing, optimizing, and creating resumes — usable with any LLM or agent, no API keys required.

## Repository Structure

```
.
├── cmd/incipit/main.go             # Entry point
├── internal/cli/                   # Cobra CLI commands
│   ├── root.go                     # Root command setup
│   ├── generate.go                 # `generate` command: JSON output, dry-run, schema
│   ├── ai.go                       # `ai` parent command
│   ├── ai_review.go                # `ai review`: prints a resume review prompt
│   ├── ai_optimize.go              # `ai optimize`: prints an optimization prompt
│   ├── ai_create.go                # `ai create`: prints a text-to-JSON conversion prompt
│   └── templates.go                # `templates list|validate|engines` subcommands
├── ai/                             # LLM prompt builders
│   └── prompts.go                  # Review/optimize/create prompts + resume JSON Schema embed
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

AI prompt flow (no LLM calls — prompts only):
    Input file -> embedded into a self-contained prompt -> stdout
        ai review   -> four-dimension scoring rubric (content, writing, industry, structure)
        ai optimize -> bullet/metric/keyword improvement instructions (+ optional job description)
        ai create   -> extraction rules + resume JSON Schema (reflected from resume.Resume)
    User pipes the prompt to any LLM (e.g. `| claude -p`) and saves the output.
```

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

### Golden files

`generators/golden_test.go` compares generator output against snapshots in
`generators/testdata/golden/`. Any intentional change to template output
(separators, headings, layout) MUST regenerate the golden files in the same
commit, or CI fails:

```bash
just golden   # go test ./generators -run TestGolden -update
```

## Commit Convention

Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, `ci:`
