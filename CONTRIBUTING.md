# Contributing to Resume Generator

Thank you for your interest in contributing to Resume Generator! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [Project Structure](#project-structure)
5. [Making Changes](#making-changes)
6. [Testing](#testing)
7. [Submitting Changes](#submitting-changes)
8. [Style Guidelines](#style-guidelines)
9. [Adding Templates](#adding-templates)
10. [Documentation](#documentation)

## Code of Conduct

This project follows a standard code of conduct:

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other contributors

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- Go 1.21 or higher
- Git
- LaTeX distribution (TeX Live, MacTeX, or MiKTeX) for LaTeX template development
- Chromium/Chrome for HTML template testing
- [just](https://github.com/casey/just) (optional, for helper commands)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/resume-generator.git
   cd resume-generator
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/urmzd/resume-generator.git
   ```

## Development Setup

### Building the CLI

```bash
# Build the binary
go build -o resume-generator .

# Or use the justfile
just install
```

### Installing Dependencies

```bash
# Install Go dependencies
go mod tidy

# Or use the justfile
just init
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/generators

```

## Making Changes

### Creating a Branch

Create a feature branch for your changes:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-new-template` - New features
- `fix/latex-compilation-error` - Bug fixes
- `docs/update-readme` - Documentation changes
- `refactor/path-resolution` - Code refactoring

### Development Workflow

1. Make your changes in your feature branch
2. Test your changes thoroughly
3. Commit your changes with clear, descriptive messages
4. Push to your fork
5. Create a pull request

## Testing

### Manual Testing

Test your changes with example inputs:

```bash
# Test with resume format
./resume-generator run -i assets/example_resumes/software_engineer.yml -t modern-html

# Test with LaTeX template
./resume-generator run -i assets/example_resumes/software_engineer.yml -t modern-latex

# Test validation
./resume-generator validate assets/example_resumes/software_engineer.yml

# Test preview
./resume-generator preview assets/example_resumes/software_engineer.yml

# Test template listing
./resume-generator templates list

# Test LaTeX engine detection
./resume-generator templates engines
```

### Automated Testing

```bash
# Run unit tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Testing Templates

When adding or modifying templates:

1. Test with multiple example configurations
2. Verify PDF output quality
3. Test with edge cases (long text, special characters, missing fields)
4. Validate on different operating systems if possible

## Submitting Changes

### Commit Messages

Write clear, descriptive commit messages following this format:

```
<type>: <subject>

<body (optional)>

<footer (optional)>
```

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, no logic changes)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks

Examples:

```
feat: add support for custom fonts in LaTeX templates

Added configuration option for custom font selection in LaTeX
templates. Updated modern-latex template to support font overrides.

Closes #123
```

```
fix: resolve path resolution issue on Windows

Fixed bug where home directory expansion failed on Windows systems.
Updated path utility to use filepath.Join consistently.

Fixes #456
```

### Pull Request Process

1. **Update documentation** - Ensure README.md and relevant docs are updated
2. **Add tests** - Include tests for new features or bug fixes
3. **Update changelog** - Add entry to CHANGELOG.md if it exists
4. **Verify CI passes** - Ensure all automated checks pass
5. **Request review** - Tag relevant maintainers for review

### Pull Request Template

When creating a pull request, include:

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other (describe)

## Testing
Describe how you tested your changes:
- [ ] Tested with example configurations
- [ ] Added/updated unit tests
- [ ] Tested on multiple platforms
- [ ] Verified output quality

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests pass locally
```

## Style Guidelines

### Go Code Style

Follow standard Go conventions:

```go
// Use gofmt for formatting
gofmt -w .

// Use golint for linting
golint ./...

// Use go vet for static analysis
go vet ./...
```

**Best Practices:**

- Use meaningful variable names
- Keep functions small and focused
- Add comments for exported functions
- Handle errors explicitly
- Use structured logging (zap)

**Example:**

```go
// LoadTemplate loads a template by name from the templates directory.
// It returns the template metadata and an error if the template is not found.
func LoadTemplate(name string) (*Template, error) {
    templatePath := utils.ResolveAssetPath(filepath.Join("templates", name))

    if !utils.DirExists(templatePath) {
        return nil, fmt.Errorf("template not found: %s", name)
    }

    // Load template configuration...
    return template, nil
}
```

### Template Style

**LaTeX Templates:**

- Use consistent indentation (2 or 4 spaces)
- Include comments for complex sections
- Use template variables for all dynamic content
- Test with special characters and edge cases

**HTML Templates:**

- Follow semantic HTML5 standards
- Use CSS classes for styling (no inline styles)
- Ensure responsive design
- Test print media styles for PDF output

### Documentation Style

- Use clear, concise language
- Include code examples where appropriate
- Keep line length to 80-100 characters for readability
- Use proper Markdown formatting
- Add table of contents for long documents

## Adding Templates

### Creating a New Template

1. **Create template directory:**
   ```bash
   mkdir -p templates/your-template-name
   ```

2. **Add template file:**
   - For HTML: `templates/your-template-name/template.html`
   - For LaTeX: `templates/your-template-name/template.tex`

3. **Add configuration:**
   ```yaml
   # templates/your-template-name/config.yml
   name: your-template-name
   display_name: Your Template Name
   description: Brief description of your template
   format: html  # or latex
   tags:
     - modern
     - minimal
   ```

4. **Add supporting files:**
   - CSS files for HTML templates
   - LaTeX class files for LaTeX templates

5. **Test the template:**
   ```bash
   ./resume-generator run -i assets/example_resumes/software_engineer.yml -t your-template-name
   ```

6. **Add example output:**
   - Generate a sample PDF
   - Add screenshot to `assets/example_results/`

### Template Guidelines

**Required Features:**
- Support all standard resume sections (contact, experience, education, skills)
- Handle missing/optional fields gracefully
- Provide clear visual hierarchy
- Generate readable, professional output

**HTML Template Requirements:**
- Responsive design (mobile-friendly)
- Print-optimized CSS
- Cross-browser compatible
- Proper semantic HTML

**LaTeX Template Requirements:**
- Use standard LaTeX packages when possible
- Include all required class files
- Support multiple LaTeX engines (xelatex, pdflatex)
- Handle Unicode characters properly

## Documentation

### Documentation Standards

All documentation should:

1. **Be accurate** - Verify all commands and examples work
2. **Be complete** - Cover all features and use cases
3. **Be clear** - Use simple language and good examples
4. **Be up-to-date** - Update docs when code changes

### Documentation Files

- **README.md** - Main project documentation, quick start
- **usage.md** - Detailed usage instructions and examples
- **CONTRIBUTING.md** - This file, contribution guidelines
- **docs/** - Additional documentation (features, guides)

### Updating Documentation

When adding features:

1. Update README.md with new functionality
2. Add detailed examples to usage.md
3. Create feature-specific docs in `docs/` if needed
4. Update CLI help text in command files
5. Add schema documentation if data structures change

## Getting Help

If you need help:

- Open an issue on GitHub for bugs or feature requests
- Check existing issues and pull requests
- Review the documentation in `docs/`
- Ask questions in pull request comments

## Recognition

Contributors will be recognized in:

- GitHub contributors list
- CHANGELOG.md for significant contributions
- Release notes for major features

Thank you for contributing to Resume Generator!
