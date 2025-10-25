# Resume Generator

## Introduction

Resume Generator is a CLI-focused toolkit for turning structured configuration files into polished resumes. It supports multiple output formats and ships with a flexible template system so you can generate professional PDFs or HTML resumes from YAML, JSON, or TOML data.

## Features

### Core Capabilities
- **Multiple Output Formats**: Generate resumes in PDF (LaTeX) or HTML with the same template system
- **Enhanced Configuration**: Support for YAML, JSON, and TOML with advanced ordering and template embedding
- **Template System**: Modular templates with embedded assets; customize or create new templates per project
- **Robust Path Resolution**: Works from any directory, supports `~`, relative paths, and timestamped output workspaces
- **Docker Support**: Containerized build that includes LaTeX and Chromium tooling for consistent output

### Technical Features
- **Docker Workflow**: Build the Go binary and supporting toolchain in a single container
- **CLI Commands**: Validate inputs, preview data, list templates, and generate outputs from the terminal
- **JSON Resume Schema**: Full compatibility with the JSON Resume standard

## Prerequisites

### For CLI Usage
- Go 1.21+
- Chromium/Chrome (for HTML → PDF conversion)
- TeX Live or compatible LaTeX engine (for LaTeX → PDF conversion)

### Optional
- Docker (to run the bundled image)
- [just](https://github.com/casey/just) for helper commands

## Getting Started

### Quick Start

1. **Clone the Repository**
   ```bash
   git clone https://github.com/urmzd/resume-generator.git
   cd resume-generator
   ```

2. **Build the CLI**
   ```bash
   go build -o resume-generator ./...
   ```

3. **Generate a Resume**
   ```bash
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -t modern-html
   ```

If you prefer Docker, build the bundled image and run commands inside the container:

```bash
docker build -t resume-generator .
docker run --rm -v "$(pwd)":/work resume-generator run -i /work/assets/example_inputs/sample-enhanced.yml -t modern-html
```

### CLI Usage

1. **Build the Application**
   ```bash
   go build -o resume-generator ./...
   ```

2. **Generate Resume (Enhanced CLI)**
   ```bash
   # Generate HTML resume with custom output path
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o outputs/sample.pdf -t modern-html

   # Generate with default filename (first_last_resume_timestamp.pdf)
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -t modern-html

   # Generate LaTeX PDF
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o ~/Documents/my_resume.pdf -t base-latex

   # Validate configuration
   ./resume-generator validate assets/example_inputs/sample-enhanced.yml

   # Preview without generation
   ./resume-generator preview assets/example_inputs/sample-enhanced.yml

   # List available templates
   ./resume-generator templates list
   ```

3. **Path Resolution**

   The CLI now supports flexible path resolution:
   - **Relative paths**: `./assets/example_inputs/sample-enhanced.yml`, `../data/resume.yml`
   - **Absolute paths**: `/Users/name/Documents/resume.yml`
   - **Home directory**: `~/Documents/resume.yml`
   - **Custom output locations**: Specify any file path for output
  - **Directory output**: Provide a directory, and a timestamped workspace will be created (with PDF + debug artifacts)

   Examples:
   ```bash
   # Relative input, custom output
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output/my_resume.pdf

   # Absolute paths
   ./resume-generator run -i /path/to/resume.yml -o /path/to/output.pdf

   # Home directory paths
   ./resume-generator run -i ~/resumes/resume.yml -o ~/Documents/resume.pdf

   # Output to directory (creates timestamped folder)
   ./resume-generator run -i resume.yml -o ~/Documents/
   ```

   Each run results in a directory named `first[_middle]_last_<timestamp>/` containing the generated `resume.pdf` (or your custom filename) along with a `debug/` subfolder that preserves the rendered `.tex`/`.html`, `.log`, `.aux`, and supporting class files.

## Showcase

Here are some examples of resumes generated with our tool:

### Sample Resume 1

![Sample Resume 1](assets/example_results/example.jpg)

A clean, professional layout suitable for various industries.

## Generators and Templates

This tool supports different generators and templates to customize your resume.

## Configuration Formats

### Enhanced Configuration (v2.0)
Supports advanced features like ordering, template embedding, and multiple output formats:

```yaml
meta:
  version: "2.0"
  output:
    formats: ["pdf", "html"]
    theme: "modern"

contact:
  order: 1
  name: "John Doe"
  email: "john@example.com"

skills:
  order: 2
  categories:
    - name: "Programming Languages"
      order: 1
      items: ["Python", "Go", "JavaScript"]

experience:
  order: 3
  positions:
    - company: "Tech Corp"
      order: 1
      title: "Software Engineer"
      # ... rest of experience
```

### Legacy Configuration
Still supports original YAML, JSON, and TOML formats for backward compatibility.

### Generators and Templates

#### Generators
- `base`: Default generator with clean layout
- `json-resume`: [JSON Resume](https://jsonresume.org/) schema support
- `html`: Modern HTML generator with responsive design

#### Templates
- **PDF Templates**: LaTeX-based templates in `templates/*-latex/`
- **HTML Templates**: Modern responsive templates in `templates/*-html/`
- **Custom Templates**: Create your own templates following the provided patterns

```bash
# List available templates
just templates

# Use specific template
just generate config.yml output custom-template
```

## Customization

To customize your resume, edit the source file (e.g., `example.yml`) with your personal information, experiences, and skills. The tool supports various file formats like TOML, YAML, and JSON.

## Advanced Usage

For more advanced users, the following options are available for interacting with the CLI tool:

1. **Direct CLI Interaction with Specific Commands**

   Execute specific commands within the Docker container using `just exec`. This allows you to pass custom arguments and commands directly to the CLI tool. For example, to execute a command `command-name` with arguments `arg1 arg2`, use:

   ```bash
   just exec "command-name arg1 arg2"
   ```

   Replace `command-name`, `arg1`, and `arg2` with your actual command and arguments. This method is useful for executing specific operations without starting an interactive shell.

2. **Interactive Shell Session**

   Start an interactive shell session within the Docker container for a more hands-on approach:

   ```bash
   just shell
   ```

   This command opens a `/bin/bash` session in the Docker container, allowing you to interact directly with the tool and the file system.

3. **Quick Examples**

   Run a specific example quickly with:

   ```bash
   just example example.yml
   ```

   Or run all examples at once:

   ```bash
   just examples
   ```

Use these advanced options for more control over the tool or for tasks that require direct interaction with the CLI environment. 

## Templates

Built-in templates live in the `templates/` directory, one folder per template (for example, `templates/modern-html/template.html` or `templates/base-latex/template.tex`). Each template ships with a `config.yml` describing its format (`html` or `latex`) and any supporting metadata. LaTeX templates bundle their required `.cls` or helper files directly alongside the template, so no additional classes directory is needed.

## Contributing

Contributions to the Generate Resumes project are welcome. Please read our contributing guidelines and submit pull requests for any enhancements, bug fixes, or documentation improvements.

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code as per the license terms.

## Acknowledgments

Thanks to all the contributors who have helped in building and maintaining this tool. Special thanks to the LaTeX community for the underlying typesetting system.
