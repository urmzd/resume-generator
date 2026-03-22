---
name: generate
description: Generate polished resumes from YAML/JSON/TOML data using LaTeX, HTML, or DOCX templates. Use when creating resumes, adding templates, or working with resume generation.
argument-hint: [input-file] [template]
---

# Generate

Generate resumes using `incipit` — here begins the new career.

## Quick Start

```sh
# Build the CLI
just install

# Generate with default example
just run

# Generate with specific input and template
./incipit run -i resume.yml -t modern-html

# Generate DOCX
./incipit run -i resume.yml -t modern-docx

# Generate with all templates
./incipit run -i resume.yml

# Validate input
./incipit validate resume.yml

# List templates
./incipit templates list
```

## Templates

Five built-in templates ship with incipit. Each lives in `templates/<name>/` with a `metadata.yml` and a template file.

| Template | Format | Output | Best For |
|----------|--------|--------|----------|
| `modern-html` | HTML → PDF | Chromium render | Web-friendly, print-ready PDFs with CSS control |
| `modern-latex` | LaTeX → PDF | TeX Live | Typographically precise, classic academic style |
| `modern-cv` | LaTeX → PDF | TeX Live | Detailed CVs with comprehensive employment history |
| `modern-docx` | XML → DOCX | go-docx | Editable Word documents for recruiters and ATS |
| `modern-markdown` | Markdown | `.md` file | GitHub READMEs, plain-text contexts, portability |

### Selecting Templates

```sh
# Single template
./incipit run -i resume.yml -t modern-html

# Multiple templates
./incipit run -i resume.yml -t modern-html -t modern-latex

# Comma-separated
./incipit run -i resume.yml -t modern-html,modern-latex

# All templates (omit -t)
./incipit run -i resume.yml
```

### HTML Template Features

The `modern-html` template supports layout and typography options via the resume data:

- **Density**: `standard` (default), `compact`, or `detailed` — controls spacing, font sizes, and margins
- **Typography**: `classic` (serif), `modern` (sans-serif), or `elegant` (mixed)
- **Header style**: `centered`, `split` (name left / contact right), or `minimal`
- **Section ordering**: customizable via `section_order` in resume data
- **References footer**: optional "References available upon request"

### Adding Custom Templates

Create `templates/<name>/` with:
1. `metadata.yml` — name, format, description, author, tags
2. Template file — `.html`, `.tex`, `.xml`, or `.md` using Go template syntax

### Installing Community Templates

```sh
# Register a local template directory
./incipit templates add my-template ./path/to/template

# Install templates from a release
./incipit templates install --version v1.0.0
```

## Output Structure

```
outputs/<name>/<date>/<template>/
├── <name>_resume.pdf
└── <name>_resume_debug/
    └── <name>_resume.{html,tex}
```

## Web Usage

Install incipit as an agent skill for use in Claude Code, Cursor, or any compatible agent:

```sh
npx skills add urmzd/incipit
```

Once installed, invoke `/generate` to create resumes directly from your agent. Provide your resume data as YAML, JSON, or TOML, and specify a template — the skill handles rendering and output.

For a full GUI experience, use [incipit-app](https://github.com/urmzd/incipit-app) — a native desktop app with live preview, template gallery, and inline editing.
