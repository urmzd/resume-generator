---
name: resume
description: Generate polished resumes (PDF, DOCX, Markdown) from JSON data. Installs the incipit CLI, provides the resume schema, and uses templates for rendering. The agent writes JSON directly — no AI subcommands needed.
argument-hint: [file]
---

# Resume

Generate professional resumes with `incipit`. The agent writes structured JSON matching the schema, then the CLI renders it into polished output.

## Setup

Install the CLI if `incipit` is not on PATH:

```sh
curl -fsSL https://raw.githubusercontent.com/urmzd/incipit/main/install.sh | bash
export PATH="$HOME/.local/bin:$PATH"
```

Verify: `incipit --version`

## Schema

Get the full JSON Schema for resume data:

```sh
incipit schema
```

Use this schema to write or validate resume JSON files directly.

## Resume JSON Structure

Required fields: `contact.name`, `contact.email`, `skills`, `experience`, `education`.

Date formats: `"2024"` (year), `"2024-01"` (month), `"2024-01-15T00:00:00Z"` (full).

```json
{
  "contact": {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phone": "+1-555-0100",
    "location": { "city": "San Francisco", "state": "CA", "country": "USA" },
    "links": [
      { "uri": "https://linkedin.com/in/janedoe", "label": "LinkedIn" },
      { "uri": "https://github.com/janedoe", "label": "GitHub" }
    ]
  },
  "summary": "Senior software engineer with 8 years of experience...",
  "skills": {
    "title": "Technical Skills",
    "categories": [
      { "category": "Languages", "items": ["Go", "Python", "TypeScript"] },
      { "category": "Infrastructure", "items": ["AWS", "Kubernetes", "Terraform"] }
    ]
  },
  "experience": {
    "title": "Experience",
    "positions": [
      {
        "company": "Acme Corp",
        "title": "Senior Software Engineer",
        "employment_type": "Full-time",
        "highlights": [
          "Reduced API latency by 40% by migrating to a distributed cache layer serving 2M daily requests",
          "Led a team of 4 engineers to deliver a real-time analytics pipeline processing 500K events/sec"
        ],
        "dates": { "start": "2022-01" },
        "location": { "city": "San Francisco", "state": "CA" },
        "technologies": ["Go", "Redis", "Kafka"]
      }
    ]
  },
  "education": {
    "title": "Education",
    "institutions": [
      {
        "institution": "MIT",
        "degree": { "name": "B.S. Computer Science", "descriptions": ["Dean's List"] },
        "gpa": { "gpa": "3.8", "max_gpa": "4.0" },
        "dates": { "start": "2014-09", "end": "2018-05" },
        "location": { "city": "Cambridge", "state": "MA" }
      }
    ]
  },
  "projects": {
    "title": "Projects",
    "projects": [
      {
        "name": "OpenTracer",
        "highlights": ["Distributed tracing library with 2K GitHub stars"],
        "link": { "uri": "https://github.com/janedoe/opentracer" },
        "technologies": ["Go", "gRPC"]
      }
    ]
  },
  "certifications": {
    "title": "Certifications",
    "items": [
      { "name": "AWS Solutions Architect", "issuer": "Amazon", "date": "2023-06" }
    ]
  },
  "languages": {
    "title": "Languages",
    "languages": [
      { "name": "English", "proficiency": "Native" },
      { "name": "Spanish", "proficiency": "Professional" }
    ]
  }
}
```

## CLI Commands

### Generate output

```sh
incipit run resume.json -t modern-html          # HTML → PDF
incipit run resume.json -t modern-latex          # LaTeX → PDF
incipit run resume.json -t modern-docx           # Word document
incipit run resume.json -t modern-markdown       # Markdown file
incipit run resume.json                          # All templates
incipit run resume.json -t modern-html -o ./out  # Custom output dir
```

### Validate

```sh
incipit validate resume.json
```

### List templates

```sh
incipit templates list
```

## Available Templates

| Template | Format | Description |
|----------|--------|-------------|
| `modern-html` | HTML → PDF | Clean, modern design via Chromium |
| `modern-latex` | LaTeX → PDF | Classic academic style |
| `modern-cv` | LaTeX → PDF | Detailed CV format |
| `modern-docx` | DOCX | Microsoft Word document |
| `modern-markdown` | Markdown | Plain `.md` file |

## Writing Good Resume Content

When creating or improving resume highlights, follow these guidelines:

- **Quantify**: use numbers, percentages, dollar amounts, team sizes, timeframes
- **XYZ formula**: "Accomplished [X] as measured by [Y], by doing [Z]"
- **Action verbs**: start each bullet with a strong verb (Led, Built, Reduced, Designed, Migrated)
- **Concise**: 1-2 lines per bullet point max
- **Specific**: reference technologies, methodologies, and concrete outcomes
- **No fabrication**: only include data supported by the user's actual experience

## Workflow

1. Get the schema: `incipit schema`
2. Write `resume.json` matching the schema (the agent does this directly)
3. Validate: `incipit validate resume.json`
4. Generate: `incipit run resume.json -t modern-html`
