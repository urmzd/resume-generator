---
name: add-template
description: >
  Guides the user through adding a new resume template to the resume-generator project.
  Use this skill when the user asks to create, add, or scaffold a new template variant
  (e.g., "add a compact HTML template", "create a minimal LaTeX resume template",
  "add a new markdown template"). Also use when they want to add an entirely new output
  format (e.g., "add JSON output support").
user_invocable: true
---

# Add Template Skill

## Overview

A **template** in this project is a `templates/<name>/` directory containing:

- `config.yml` — metadata (name, format, description, tags)
- A template file using Go template syntax (`template.html`, `template.tex`, or `template.md`)
- Optional support files (e.g., `.cls` for LaTeX)

Templates are auto-discovered at build time via `//go:embed` in `main.go`. No registration code is needed.

## Decision: New Template vs New Format

Ask the user which path applies:

| Path | When | Effort |
|------|------|--------|
| **New template** | Adding a variant for an existing format (HTML, LaTeX, Markdown) | Files only — zero Go code changes |
| **New output format** | Adding an entirely new format (e.g., JSON, AsciiDoc) | Go code + formatter + generator + template files |

Most users want a new template. Only proceed with "new format" if they explicitly need a format that doesn't exist yet.

## Workflow: Adding a Template (Existing Format)

### 1. Choose a name

Convention: `<style>-<format>` (e.g., `minimal-html`, `compact-latex`, `academic-markdown`).

### 2. Create `templates/<name>/config.yml`

```yaml
name: <name>
display_name: <Human Readable Name>
description: <One-line description of the template>
format: <html|latex|markdown>
version: "1.0.0"
author: <author>
tags:
  - <format>
  - <style keywords>
```

Required fields: `name`, `format`. All others are recommended.

See existing configs for reference:
- `templates/modern-html/config.yml`
- `templates/modern-latex/config.yml`
- `templates/modern-markdown/config.yml`

### 3. Create the template file

| Format | File | Go template engine |
|--------|------|--------------------|
| HTML | `template.html` | `html/template` (auto-escapes HTML) |
| LaTeX | `template.tex` | `text/template` (use `escape` for LaTeX chars) |
| Markdown | `template.md` | `text/template` (use `escape` for Markdown chars) |

**Start from an existing template.** Read the corresponding file on disk:
- `templates/modern-html/template.html`
- `templates/modern-latex/template.tex`
- `templates/modern-markdown/template.md`

Adapt the structure and styling to the user's desired look.

### 4. Test

```bash
go build -o resume-generator .
./resume-generator run -i assets/example_resumes/software_engineer.yml -t <name>
./resume-generator templates list
```

## Workflow: Adding a New Output Format

This requires Go code changes in multiple files. Follow this sequence:

### 1. Formatter — `pkg/generators/formatter_<fmt>.go`

- Embed `baseFormatter`
- Implement the `Formatter` interface (defined in `pkg/generators/formatter.go`)
- Key methods: `EscapeText`, `FormatLocation`, `FormatLink`, `TemplateFuncs`
- `TemplateFuncs()` returns a `template.FuncMap` with all helpers available in templates

### 2. Generator — `pkg/generators/<fmt>.go`

- Struct with `Generate(templateContent string, resume *resume.Resume) (string, error)`
- Parse template, apply `TemplateFuncs()`, execute against resume data

### 3. Wiring — `pkg/generators/generator.go`

Five touch points:

1. **Constant**: Add `TemplateType<Fmt> TemplateType = "<fmt>"` (line ~20)
2. **`parseTemplateType()`**: Add case for the new format string (line ~380)
3. **`resolveTemplateFilename()`**: Add case returning default filename (line ~394)
4. **`GenerateWithTemplate()` switch**: Add case calling `g.render<Fmt>()` (line ~231)
5. **`render<Fmt>()` method**: New method instantiating the generator (after line ~258)

### 4. CLI output — `cmd/run.go`

Add output handling in the format switch. Non-PDF formats write directly to file. PDF formats invoke a compiler.

### 5. Template listing — `cmd/templates.go`

- Add display group for the new format
- Add validation case so `templates validate` recognizes it

### 6. Template files

Create `templates/modern-<fmt>/` with `config.yml` + template file.

### 7. Tests — `pkg/generators/<fmt>_test.go`

Write unit tests for the generator. See `pkg/generators/markdown_test.go` as a reference.

## Template Functions Reference

Functions available in `{{ }}` blocks, organized by format:

### All formats (from `baseFormatter`)

| Function | Signature | Description |
|----------|-----------|-------------|
| `escape` | `string → string` | Format-specific character escaping |
| `fmtDateRange` | `DateRange → string` | "Jan 2020 – Present" |
| `fmtOptDateRange` | `*DateRange → string` | Nil-safe date range |
| `fmtLocation` | `interface{} → string` | "City, State, Country" |
| `formatList` | `[]string → string` | Comma-separated, empty-filtered |
| `formatGPA` | `*GPA → string` | "3.8" or "3.8 / 4.0" |
| `skillNames` | `[]string → []string` | Filtered skill names |
| `join` | `(sep, []string) → string` | Join with separator |
| `filterEmpty` | `[]string → []string` | Remove blank strings |
| `title` | `string → string` | Title Case |
| `upper` | `string → string` | UPPER CASE |
| `lower` | `string → string` | lower case |
| `trim` | `string → string` | Trim whitespace |
| `sanitizePhone` | `string → string` | Keep digits and + only |
| `default` | `(default, value) → interface{}` | Fallback if nil/empty |

### HTML-specific

| Function | Description |
|----------|-------------|
| `safeHTML` | Mark string as trusted HTML (no escaping) |
| `layoutClass` | CSS classes from `*Layout` |
| `hasSection` | Check if a section has data |
| `containsSection` | Check if name is in section list |
| `formatDate` | `time.Time → "January 2006"` |
| `formatDateShort` | `time.Time → "Jan 2006"` |
| `calculateDuration` | `(start, end) → "X yr Y mo"` |

### LaTeX-specific

| Function | Description |
|----------|-------------|
| `escapeLatexChars` | Alias for `escape` |
| `fmtLinkWithDomain` | `\href{url}{domain}` |
| `extractDisplayURL` | Strip protocol/www from URL |
| `formatLocationFull` | Full location with country |
| `formatLocationShort` | City + state only |
| `fmtDateLegal` | `time.Time → "January 2, 2006"` |

### Markdown-specific

| Function | Description |
|----------|-------------|
| `bold` | Wrap in `**...**` |
| `italic` | Wrap in `*...*` |
| `extractDisplayURL` | Strip protocol/www from URL |
| `fmtLinkWithDomain` | `[domain](url)` |
| `add` | Integer addition |

## Resume Data Model Quick Reference

The template receives a `*resume.Resume` as `{{ . }}`. Key fields:

```
.Contact.Name           string (required)
.Contact.Email          string (required)
.Contact.Phone          string
.Contact.Location       *Location (.City, .State, .Country, .Remote)
.Contact.Links          []Link (.URI, .Label)

.Summary                string

.Skills.Title           string
.Skills.Categories      []SkillCategory
  .Category             string
  .Items                []string

.Experience.Title       string
.Experience.Positions   []Experience
  .Title                string
  .Company              string
  .EmploymentType       string
  .Dates                DateRange (.Start, .End)
  .Location             *Location
  .Highlights           []string
  .Technologies         []string

.Education.Title        string
.Education.Institutions []Education
  .Institution          string
  .Degree               Degree (.Name, .Descriptions)
  .Specializations      []string
  .GPA                  *GPA (.GPA, .MaxGPA)
  .Dates                DateRange
  .Location             *Location
  .Thesis               *Thesis (.Title, .Link, .Highlights)
  .Awards               []Award (.Name, .Date, .Notes)

.Projects               *ProjectList (may be nil)
.Projects.Title         string
.Projects.Projects      []Project
  .Name                 string
  .Link                 Link (.URI, .Label)
  .Highlights           []string
  .Dates                *DateRange
  .Technologies         []string

.Languages              *LanguageList (may be nil)
.Languages.Title        string
.Languages.Languages    []Language
  .Name                 string
  .Proficiency          string

.Certifications         *Certifications (may be nil)
.Certifications.Title   string
.Certifications.Items   []Certification
  .Name                 string
  .Issuer               string
  .Date                 *time.Time
  .Notes                string

.Layout                 *Layout (may be nil)
  .Sections             []string
  .Density              string
  .Typography           string
  .Header               string
  .SkillColumns         int
```

**Important:** `.Projects`, `.Languages`, `.Certifications`, and `.Layout` are pointers — always nil-check before accessing (e.g., `{{if .Projects}}`).

## config.yml Schema

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Template identifier (must match directory name) |
| `format` | Yes | One of: `html`, `latex`, `markdown`, `docx` |
| `display_name` | No | Human-readable name (defaults to `name`) |
| `description` | No | One-line description |
| `version` | No | Semver string |
| `author` | No | Author name |
| `tags` | No | List of keyword strings |
| `template_file` | No | Override default filename (e.g., `custom.html`) |

## Verification Checklist

After creating the template:

- [ ] `go build -o resume-generator .` succeeds
- [ ] `go test ./...` passes
- [ ] `./resume-generator templates list` shows the new template
- [ ] `./resume-generator run -i assets/example_resumes/software_engineer.yml -t <name>` produces output
- [ ] Output file looks correct (open HTML in browser, compile LaTeX, view Markdown)
