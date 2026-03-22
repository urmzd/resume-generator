---
name: resume-generate
description: Generate polished resumes from YAML/JSON/TOML data using LaTeX, HTML, or DOCX templates. Use when creating resumes, adding templates, or working with resume generation.
argument-hint: [input-file] [template]
---

# Resume Generation

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

## Template Types

| Type | Output | Engine |
|------|--------|--------|
| `*-html` | PDF | Rod/Chromium |
| `*-latex` | PDF | TeX Live |
| `*-docx` | Word | go-docx |

## Output Structure

```
outputs/<name>/<date>/<template>/
├── <name>_resume.pdf
└── <name>_resume_debug/
    └── <name>_resume.{html,tex}
```

## Adding a Template

Create `templates/<name>/` with `config.yml` + template file. See existing templates for patterns.
