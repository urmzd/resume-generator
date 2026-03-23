# Contributing to Incipit

## Prerequisites

- Go 1.25+
- Git
- LaTeX distribution (TeX Live, MacTeX, or MiKTeX) for LaTeX template development
- Chromium/Chrome for HTML template testing
- [just](https://github.com/casey/just) (optional, for helper commands)

## Development Setup

```bash
git clone https://github.com/urmzd/incipit.git
cd incipit
go install ./cmd/incipit
go test ./...
```

## Testing

```bash
# Test with JSON resume
incipit run -i assets/example_resumes/software_engineer.json -t modern-html

# Test AI review (requires LLM provider)
incipit ai review assets/example_resumes/software_engineer.json

# Run all tests
go test ./...
```

## Commit Convention

Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, `ci:`

## Adding Templates

1. Create `templates/<name>/metadata.yml` with format, display_name, description
2. Create the template file (`template.html`, `template.tex`, or `template.md`)
3. Test: `incipit run -i assets/example_resumes/software_engineer.json -t <name>`
