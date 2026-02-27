# AGENTS.md — AI Coding Agent Guide

## Project Overview

**resume-generator** is a Go CLI + Wails desktop application that converts structured resume data (YAML/JSON/TOML) into PDF, HTML, LaTeX, DOCX, and Markdown output formats. It uses Go's `text/template` and `html/template` engines with a formatter abstraction layer for output-specific escaping and rendering.

## Repository Structure

```
.
├── main.go                     # Entry point, embeds templates via //go:embed
├── app.go                      # Wails desktop app backend (context, file handling)
├── cmd/                        # Cobra CLI commands
│   ├── root.go                 # Root command, embedded FS setup
│   ├── run.go                  # `run` command: loads resume → generates output
│   ├── templates.go            # `templates list|validate|engines` subcommands
│   ├── validate.go             # `validate` command for resume data
│   ├── preview.go              # `preview` command (HTML live preview)
│   └── schema.go               # `schema` command (JSON Schema output)
├── pkg/
│   ├── generators/
│   │   ├── generator.go        # Core: template loading, type dispatch, embed FS
│   │   ├── formatter.go        # Formatter interface definition
│   │   ├── formatter_base.go   # Shared formatting logic (dates, location, GPA)
│   │   ├── formatter_html.go   # HTML-specific formatter
│   │   ├── formatter_latex.go  # LaTeX-specific formatter
│   │   ├── formatter_markdown.go # Markdown-specific formatter
│   │   ├── html.go             # HTMLGenerator (html/template)
│   │   ├── latex.go            # LaTeXGenerator (text/template)
│   │   ├── markdown.go         # MarkdownGenerator (text/template)
│   │   └── docx.go             # DOCXGenerator (programmatic)
│   ├── compilers/              # PDF compilation (LaTeX engines, Rod/Chromium)
│   ├── resume/                 # Resume data model and validation
│   └── utils/                  # Path resolution, file helpers
├── templates/
│   ├── modern-html/            # config.yml + template.html
│   ├── modern-latex/           # config.yml + template.tex + default.cls
│   ├── modern-cv/              # config.yml + template.tex
│   ├── modern-docx/            # config.yml (programmatic generation)
│   └── modern-markdown/        # config.yml + template.md
├── frontend/                   # React/TypeScript Wails frontend
├── assets/example_resumes/     # Example YAML resume files
└── justfile                    # Task runner (build, run, dev, demo)
```

## Architecture

### Data Flow

```
Input (YAML/JSON/TOML) → resume.LoadResumeFromFile() → Resume struct
    → Generator.GenerateWithTemplate(template, resume)
        → Formatter.TemplateFuncs() provides template helpers
        → text/template or html/template renders output
    → Compiler (LaTeX→PDF or HTML→PDF via Rod/Chromium)
    → Output file (.pdf, .html, .docx, .md)
```

### Formatter Pattern

Each output format has a formatter that embeds `baseFormatter` and overrides format-specific behavior:

- `baseFormatter` — shared date, location, GPA, list formatting
- `htmlFormatter` — HTML escaping, CSS layout helpers
- `latexFormatter` — LaTeX escaping, `\href{}{}` links, `\textendash` dates
- `markdownFormatter` — Markdown escaping, `[text](url)` links

All formatters implement the `Formatter` interface and expose `TemplateFuncs()` returning a `template.FuncMap`.

### Template System

Templates live in `templates/<name>/` with:
- `config.yml` — metadata (name, format, description, tags)
- Template file (`template.html`, `template.tex`, `template.md`)
- Optional support files (`.cls` for LaTeX)

Templates are embedded at build time via `//go:embed` in `main.go` and loaded through `generator.go`'s `LoadTemplate()` / `ListTemplates()` functions.

## How to Add a New Template

1. Create `templates/<name>/config.yml` with format, display_name, description
2. Create the template file (`template.html`, `template.tex`, or `template.md`)
3. Use Go template syntax with the formatter's `TemplateFuncs()` helpers
4. The template is auto-discovered — no code changes needed
5. Test: `go build && ./resume-generator run -i assets/example_resumes/software_engineer.yml -t <name>`

## How to Add a New Output Format

1. Create `pkg/generators/formatter_<format>.go` — embed `baseFormatter`, implement `Formatter`
2. Create `pkg/generators/<format>.go` — generator struct with `Generate()` method
3. In `generator.go`:
   - Add `TemplateType<Format>` constant
   - Add case in `parseTemplateType()`
   - Add case in `resolveTemplateFilename()`
   - Add case in `GenerateWithTemplate()` switch
   - Add `render<Format>()` method
4. In `cmd/run.go`: add output handling (direct file write or PDF compilation)
5. In `cmd/templates.go`: add display group and validation case
6. Create at least one template in `templates/modern-<format>/`
7. Write tests in `pkg/generators/<format>_test.go`

## Resume Data Model

See `pkg/resume/resume.go` for the full struct. Key types:

- `Resume` — top-level: Contact, Summary, Skills, Experience, Projects, Education, Languages, Certifications, Layout
- `Contact` — Name (required), Email (required), Phone, Location, Links
- `DateRange` — Start (time.Time), End (*time.Time, nil = Present)
- `Layout` — Sections ordering, density, typography, header style

Date format in YAML: RFC3339 (`2024-01-15T00:00:00Z`)

## Build & Test

```bash
# Build CLI binary
just install
# or: CGO_ENABLED=0 go build -trimpath -o resume-generator .

# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/generators/... -v

# Format check
gofmt -l .

# Lint (if installed)
golangci-lint run

# Fuzz testing
go test ./pkg/generators/ -fuzz=FuzzHTMLGenerate -fuzztime=10s
```

## CI Checks

- `gofmt` — all Go files must be formatted
- `go test ./...` — all tests must pass
- `golangci-lint run` — no lint errors
- Fuzz smoke test — 10-second fuzz run

## Commit Convention

Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, `ci:`

Examples:
- `feat: add markdown output format`
- `fix: handle nil date range in LaTeX formatter`
- `docs: update AGENTS.md with new format guide`

## Common Tasks

### Modify a template
Edit the template file directly. Use `{{escape .Field}}` for user content, formatter helpers for dates/locations/links.

### Add a resume field
1. Add field to the appropriate struct in `pkg/resume/resume.go`
2. Add to the input adapter in `pkg/resume/` (YAML/JSON/TOML mapping)
3. Update templates that should display the field
4. Add validation if required

### Debug template rendering
Set output dir and check the `_debug/` directory next to each output file. It contains intermediate files (rendered HTML, LaTeX source) before PDF compilation.
