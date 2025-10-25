# Resume Generator Usage Guide

A comprehensive platform for generating professional resumes from configuration files with support for multiple output formats, web interface, and REST API access.

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

### Fastest Way: Docker Compose (Recommended)

1. **Clone and Start**
   ```bash
   git clone https://github.com/urmzd/resume-generator.git
   cd resume-generator
   docker-compose up
   ```

2. **Access Web Interface**
   - Open: http://localhost:3000
   - Upload your resume config file
   - Choose format (HTML/PDF) and template
   - Download generated resume

3. **Or Use CLI Directly**
   ```bash
   # Copy example config
   just init

   # Generate resume
   just generate examples/sample-enhanced.yml resume html modern
   ```

---

## Installation Methods

### Docker-Based (Recommended)

#### Full Platform with Web Interface
```bash
# Start complete platform (API + Frontend + Nginx)
docker-compose up -d

# Access:
# - Web UI: http://localhost:3000
# - API: http://localhost:8080/api/v1
# - Nginx: http://localhost (production)
```

#### CLI-Only Docker Usage
```bash
# Build the CLI image
just build

# Generate resume using Docker
just run example.yml

# Interactive shell access
just shell
```

### Local Development Setup

#### Prerequisites
- Go 1.21+
- Node.js 18+ (for frontend)
- Docker (for TeX tools)

#### Setup Steps
```bash
# Install dependencies
just deps
just frontend-install

# Start development environment
just dev  # Starts both API (8080) and frontend (3000)

# Or start separately
just serve 8080     # API only
just frontend       # Frontend only
```

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
just validate examples/sample-enhanced.yml

# Or using Go directly
go run main.go validate config.yml

# Preview configuration (no generation)
just preview examples/sample-enhanced.yml
```

---

## Usage Methods

### 1. CLI Commands

#### Basic Generation
```bash
# Generate HTML resume
just generate config.yml output html modern

# Generate PDF resume
just generate config.yml output pdf base

# Generate multiple formats
go run main.go run -i config.yml -f pdf,html -o resume
```

#### Available Commands
```bash
# Core commands
resume-generator run -i config.yml -o output.pdf     # Generate resume (stores results in timestamped directory)
resume-generator validate config.yml                 # Validate config
resume-generator preview config.yml                  # Preview config
resume-generator serve -p 8080                      # Start API server

# Template management
resume-generator templates list                      # List templates
go run main.go templates list                       # Alternative

# Using just shortcuts
just generate config.yml resume html modern         # Generate with template
just validate config.yml                           # Validate
just serve 8080                                    # Start API
just templates                                     # List templates
```

Each CLI run creates a dedicated output directory named after the contact (`first[_middle]_last_<timestamp>`). The main PDF (or custom filename from `--output`) resides at the root of that folder, and supporting artifacts such as rendered `.tex`/`.html`, `.log`, and `.aux` files are stored under the `debug/` subdirectory.

### 2. Web Interface

1. **Upload Configuration**
   - Navigate to http://localhost:3000
   - Upload YAML, JSON, or TOML config file
   - Or paste configuration directly

2. **Select Options**
   - Choose output format (PDF/HTML)
   - Select template theme
   - Preview configuration if needed

3. **Generate & Download**
   - Click generate button
   - Download generated resume
   - Files expire after 24 hours

### 3. REST API

#### Start API Server
```bash
# Using just
just serve 8080

# Using Go directly
go run main.go serve -p 8080

# Using Docker
docker-compose up api
```

#### API Endpoints

##### Generate Resume
```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -F "config=@examples/sample-enhanced.yml" \
  -F "format=html" \
  -F "template=modern"
```

##### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

##### List Available Formats
```bash
curl http://localhost:8080/api/v1/formats
```

##### List Available Templates
```bash
curl http://localhost:8080/api/v1/templates
```

##### Test All Endpoints
```bash
just api-test  # Requires running API server
```

### 4. Docker Usage

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

#### Docker Compose Workflows
```bash
# Development mode
docker-compose -f docker-compose.yaml -f docker-compose.override.yaml up

# Production mode
docker-compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d

# View logs
just logs api        # API logs only
just logs frontend   # Frontend logs only
just logs            # All services
```

---

## Output Formats & Templates

### Supported Output Formats

#### PDF Generation
- **Engine**: XeLaTeX compiler
- **Templates**: LaTeX-based templates in `assets/templates/`
- **Features**: Professional typesetting, print-ready
- **File**: `resume.pdf`

#### HTML Generation
- **Engine**: Go html/template
- **Templates**: Modern responsive templates in `assets/templates/html/`
- **Features**: Responsive design, web-optimized, SEO-friendly
- **File**: `resume.html`

### Available Templates

#### PDF Templates
```bash
# List available templates
just templates

# Available templates:
# - base: Clean, professional layout
# - json-resume: JSON Resume schema compatible
```

#### HTML Templates
```bash
# Available HTML themes:
# - modern: Clean, responsive design with modern styling
# - minimal: Simplified layout focused on content
# - creative: Enhanced visual design with more colors
```

### Template Selection

```bash
# CLI: Specify template with -t flag
go run main.go run -i config.yml -t modern -f html

# API: Include template parameter
curl -X POST http://localhost:8080/api/v1/generate \
  -F "config=@config.yml" \
  -F "template=modern"

# Config: Set in configuration file
meta:
  output:
    theme: "modern"
```

### Multi-Format Generation

```bash
# Generate both PDF and HTML
go run main.go run -i config.yml -f pdf,html

# Configuration-based
meta:
  output:
    formats: ["pdf", "html"]
```

---

## Advanced Usage

### Custom Templates

#### Creating HTML Templates
1. **Create Template File**
   ```html
   <!-- assets/templates/html/custom.html -->
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
   go run main.go run -i config.yml -t custom -f html
   ```

#### Creating LaTeX Templates
1. **Template Structure**
   ```latex
   % assets/templates/custom.txt
   \documentclass{article}
   \begin{document}
   \section*{ {{.Name}} }
   % Template content
   \end{document}
   ```

### API Integration Examples

#### Python Integration
```python
import requests

# Generate resume
with open('config.yml', 'rb') as f:
    response = requests.post(
        'http://localhost:8080/api/v1/generate',
        files={'config': f},
        data={'format': 'html', 'template': 'modern'}
    )

if response.status_code == 200:
    with open('resume.html', 'wb') as f:
        f.write(response.content)
```

#### JavaScript/Node.js Integration
```javascript
const FormData = require('form-data');
const fs = require('fs');

const form = new FormData();
form.append('config', fs.createReadStream('config.yml'));
form.append('format', 'html');
form.append('template', 'modern');

fetch('http://localhost:8080/api/v1/generate', {
    method: 'POST',
    body: form
})
.then(response => response.buffer())
.then(buffer => fs.writeFileSync('resume.html', buffer));
```

### Batch Processing

#### Multiple Resume Generation
```bash
#!/bin/bash
# generate_all.sh
for config in configs/*.yml; do
    basename=$(basename "$config" .yml)
    just generate "$config" "output/$basename" html modern
    just generate "$config" "output/$basename" pdf base
done
```

#### Automated CI/CD Integration
```yaml
# .github/workflows/resume.yml
name: Generate Resume
on: [push]
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Generate Resume
        run: |
          docker-compose up --build -d api
          curl -X POST http://localhost:8080/api/v1/generate \
            -F "config=@resume.yml" -F "format=pdf" \
            --output resume.pdf
```

### Development Workflows

#### Hot Reload Development
```bash
# Start with auto-reload
just dev

# Or manually
go install github.com/cosmtrek/air@latest
air  # Auto-reloads on Go file changes

# Frontend development
cd frontend && npm run dev
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
# Issue: Container fails to start
# Solution: Check Docker daemon and rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up

# Issue: Port already in use
# Solution: Change ports or stop conflicting services
docker-compose down
lsof -ti:8080 | xargs kill -9  # Kill process on port 8080
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
cp examples/sample-enhanced.yml my-config.yml
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

#### API Response Time
```bash
# Monitor API performance
curl -w "%{time_total}s\n" -X POST http://localhost:8080/api/v1/generate \
  -F "config=@config.yml" -F "format=html"
```

### Security Considerations

#### File Upload Security
- Maximum file size: 10MB
- Allowed file types: .yml, .yaml, .json, .toml
- Content validation before processing
- Temporary file cleanup after processing

#### API Security
```bash
# Rate limiting in place
# CORS configuration for web requests
# Input sanitization
# No sensitive data logging
```

### Getting Help

#### Debug Information
```bash
# Enable verbose logging
go run main.go run -i config.yml --verbose

# Check service health
curl http://localhost:8080/api/v1/health

# View container logs
docker-compose logs api
docker-compose logs frontend
```

#### Common Commands Reference
```bash
# Quick reference for just commands
just --list

# Essential commands:
just init          # Initialize project
just build         # Build Docker images
just dev           # Start development environment
just generate      # Generate resume
just validate      # Validate configuration
just serve         # Start API server
just frontend      # Start frontend
just test          # Run tests
just clean         # Clean outputs
```

---

## Additional Resources

- **Examples**: Check `examples/` directory for sample configurations
- **Templates**: Explore `assets/templates/` for template customization
- **API Documentation**: Visit http://localhost:8080/docs when server is running
- **Contributing**: See contribution guidelines in the main README
- **Issues**: Report bugs and feature requests on GitHub

This usage guide covers all major use cases and should get you started with the Resume Generator platform quickly and effectively.
