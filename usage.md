# Resume Generator Usage Guide

A CLI-focused toolkit for generating professional resumes from structured configuration files with support for multiple output formats and flexible templates.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Installation Methods](#installation-methods)
3. [Configuration Guide](#configuration-guide)
4. [Usage Methods](#usage-methods)
5. [Output Formats & Templates](#output-formats--templates)
6. [Advanced Usage](#advanced-usage)
7. [Troubleshooting](#troubleshooting)

---

## Quick Start

### Fastest Way: CLI

1. **Clone the repository**
   ```bash
   git clone https://github.com/urmzd/resume-generator.git
   cd resume-generator
   ```

2. **Build the binary**
   ```bash
   go build -o resume-generator ./...
   ```

3. **Generate a resume**
   ```bash
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -t modern-html
   ```

Prefer Docker? Build the bundled image and run the CLI inside it:

```bash
docker build -t resume-generator .
docker run --rm -v "$(pwd)":/work resume-generator run -i /work/assets/example_inputs/sample-enhanced.yml -t modern-html
```

---

## Installation Methods

### Go Toolchain

```bash
go build -o resume-generator ./...
./resume-generator --help
```

### Docker Image

```bash
# Build container with Go binary, LaTeX, and Chromium
docker build -t resume-generator .

# Run CLI inside container
docker run --rm -v "$(pwd)":/work resume-generator run -i /work/assets/example_inputs/sample-enhanced.yml -t base-latex
```

### Optional Helpers

Install [just](https://github.com/casey/just) if you want shorthand commands for building images or copying examples. The provided `justfile` includes targets such as `just build` (Docker image), `just init` (copy examples), and `just docker-run` (execute the CLI inside Docker).

### Using Just Command Runner

Install [just](https://github.com/casey/just):
```bash
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash

# See all available commands
just --list
```

---

## Configuration Guide

### Enhanced Configuration Format (v2.0)

The new format supports ordering, multiple output formats, and template embedding:

```yaml
meta:
  version: "2.0"
  output:
    formats: ["pdf", "html"]  # Supported: pdf, html
    theme: "modern"           # Template theme

contact:
  order: 1                    # Section ordering
  name: "John Doe"
  title: "Senior Software Engineer"
  email: "john.doe@example.com"
  phone: "+1 (555) 123-4567"
  website: "https://johndoe.dev"
  location:
    city: "San Francisco"
    state: "CA"
    zipcode: "94102"
    country: "USA"
  summary: "Experienced software engineer..."
  links:
    - order: 1
      type: "github"
      text: "github.com/johndoe"
      url: "https://github.com/johndoe"

skills:
  order: 2
  title: "Technical Skills"
  categories:
    - name: "Programming Languages"
      order: 1
      items:
        - name: "JavaScript/TypeScript"
          level: "Advanced"
          yearsOfExperience: 6

experience:
  order: 3
  positions:
    - order: 1
      company: "TechCorp Inc."
      title: "Senior Software Engineer"
      dates:
        start: "2021-03-15T00:00:00Z"
        end: null  # null for current position
      description:
        - "Led development of microservices architecture..."
      technologies: ["React", "Node.js", "AWS"]

education:
  order: 4
  institutions:
    - order: 1
      institution: "University of California, Berkeley"
      degree: "Bachelor of Science in Computer Science"
      dates:
        start: "2013-08-15T00:00:00Z"
        end: "2017-05-15T00:00:00Z"
      gpa: "3.8"
```

### Legacy Format Support

Still supports original YAML/JSON/TOML formats:

```yaml
# Legacy format (automatically detected)
name: "John Doe"
email: "john@example.com"
experience:
  - company: "Tech Corp"
    position: "Engineer"
    # ... rest of config
```

### Configuration Validation

```bash
# Validate configuration file
just validate assets/example_inputs/sample-enhanced.yml

# Or using Go directly
go run main.go validate config.yml

# Preview configuration (no generation)
just preview assets/example_inputs/sample-enhanced.yml
```

---

## Usage Methods

### 1. CLI Commands

#### Basic Generation
```bash
# Generate resume with Modern HTML template (renders to PDF)
just generate config.yml output modern-html

# Generate resume with Base LaTeX template
just generate config.yml output base-latex
```

#### Available Commands
```bash
# Core commands
resume-generator run -i config.yml -o output.pdf     # Generate resume (stores results in timestamped directory)
resume-generator validate config.yml                 # Validate config
resume-generator preview config.yml                  # Preview config

# Template management
resume-generator templates list                      # List templates
go run main.go templates list                       # Alternative

# Using just shortcuts
just generate config.yml resume modern-html         # Generate with template
just validate config.yml                           # Validate
just templates                                     # List templates
```

Each CLI run creates a dedicated output directory named after the contact (`first[_middle]_last_<timestamp>`). The main PDF (or custom filename from `--output`) resides at the root of that folder, and supporting artifacts such as rendered `.tex`/`.html`, `.log`, and `.aux` files are stored under the `debug/` subdirectory.

### 2. Docker Usage

#### Direct Docker Commands
```bash
# Run specific example
just example example.yml

# Run all examples
just examples

# Execute custom command
just exec "resume-generator --help"

# Interactive shell
just shell
```

## Output Formats & Templates

### Supported Output Formats

#### PDF Generation
- **Engine**: XeLaTeX compiler
- **Templates**: LaTeX-based templates in `templates/*-latex/`
- **Features**: Professional typesetting, print-ready
- **File**: `resume.pdf`

#### HTML Generation
- **Engine**: Go html/template
- **Templates**: Modern responsive templates in `templates/*-html/`
- **Features**: Responsive design, web-optimized, SEO-friendly
- **File**: `resume.html`

### Available Templates

#### PDF Templates
```bash
# List available templates
just templates

# Available templates:
# - base-latex: Clean, professional layout
# - json-resume-latex: JSON Resume schema compatible
```

#### HTML Templates
```bash
# Available HTML themes:
# - modern-html: Clean, responsive design with modern styling
```

### Template Selection

```bash
# CLI: Specify template with -t flag
go run main.go run -i config.yml -t modern-html

# Config: Set in configuration file
meta:
  output:
    theme: "modern"
```

### Template Metadata

Each template directory includes a `config.yml` file that declares its format (`html` or `latex`) along with display metadata. The CLI uses that configuration to choose the correct rendering pipeline, so you no longer need to pass a `--formats`/`-f` flag. To build a custom template, provide both the markup file (`template.html` or `template.tex`) and a `config.yml` similar to:

```yaml
name: custom-html
display_name: Custom HTML
description: Lightweight HTML template with accent colors.
format: html
tags:
  - html
  - custom
```

---

## Advanced Usage

### Custom Templates

#### Creating HTML Templates
1. **Create Template File**
   ```html
   <!-- templates/custom-html/template.html -->
   <!DOCTYPE html>
   <html>
   <head>
       <title>{{.Contact.Name}} - Resume</title>
       <style>
           /* Custom CSS */
       </style>
   </head>
   <body>
       <h1>{{.Contact.Name}}</h1>
       <!-- Template content using Go template syntax -->
   </body>
   </html>
   ```

2. **Use Custom Template**
   ```bash
   go run main.go run -i config.yml -t custom-html
   ```

#### Creating LaTeX Templates
1. **Template Structure**
   ```latex
   % templates/custom-latex/template.tex
   \documentclass{article}
   \begin{document}
   \section*{ {{.Name}} }
   % Template content
   \end{document}
   ```

### Batch Processing

#### Multiple Resume Generation
```bash
#!/bin/bash
# generate_all.sh
for config in configs/*.yml; do
    basename=$(basename "$config" .yml)
    just generate "$config" "output/$basename" modern-html
    just generate "$config" "output/$basename" base-latex
done
```

### Development Workflows

#### Hot Reload Development
```bash
# Use Air for auto-reload during CLI development
go install github.com/cosmtrek/air@latest
air  # Auto-reloads on Go file changes
```

#### Testing
```bash
# Run all tests
go test ./...

# Test specific package
go test ./pkg/generators

# Test with coverage
go test -cover ./...

# Integration tests
just test
```

---

## Troubleshooting

### Common Issues

#### Docker Issues
```bash
# Issue: Image is outdated or missing dependencies
# Solution: Rebuild the container image
docker build -t resume-generator .

# Issue: Container exits immediately
# Solution: Run with interactive shell to inspect
docker run --rm -it -v "$(pwd)":/work resume-generator /bin/sh
```

#### LaTeX/PDF Generation Issues
```bash
# Issue: LaTeX compilation fails
# Solution: Check template syntax and rebuild base image
just build-base
just build

# Issue: Missing fonts or packages
# Solution: Update Dockerfile.base with required packages
```

#### Configuration Issues
```bash
# Issue: Invalid YAML syntax
# Solution: Validate YAML syntax
just validate config.yml

# Issue: Missing required fields
# Solution: Check against example configurations
cp assets/example_inputs/sample-enhanced.yml my-config.yml
```

### Performance Optimization

#### Docker Image Size
```bash
# Current optimizations in place:
# - Multi-stage builds
# - Minimal TeX Live installation
# - Alpine base images

# Check image size
docker images | grep resume-generator
```

### Security Considerations

#### File Upload Security
- Maximum file size: 10MB
- Allowed file types: .yml, .yaml, .json, .toml
- Content validation before processing
- Temporary file cleanup after processing

### Getting Help

#### Debug Information
```bash
# Enable verbose logging (set environment variable before running)
LOG_LEVEL=debug ./resume-generator run -i config.yml

# Inspect generated debug artifacts
ls -R output_directory/debug
```

#### Common Commands Reference
```bash
# Quick reference for just commands
just --list

# Essential commands:
just init          # Initialize project
just build         # Build Docker images
just generate      # Generate resume
just validate      # Validate configuration
just test          # Run tests
just clean         # Clean outputs
```

---

## Additional Resources

- **Examples**: Check `assets/example_inputs/` for sample configurations
- **Templates**: Explore `templates/` for template customization
- **Contributing**: See contribution guidelines in the main README
- **Issues**: Report bugs and feature requests on GitHub

This usage guide covers all major use cases and should get you started with the Resume Generator platform quickly and effectively.
