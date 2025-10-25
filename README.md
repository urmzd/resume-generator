# Resume Generator

## Introduction

Resume Generator is a comprehensive platform designed to create elegant resumes from structured configuration files. This modern tool supports multiple output formats, provides a web-based interface, and offers both CLI and REST API access for maximum flexibility.

## Features

### Core Capabilities
- **Multiple Output Formats**: Generate resumes in PDF (LaTeX), HTML, and other formats
- **Enhanced Configuration**: Support for YAML, JSON, and TOML with advanced ordering and template embedding
- **Web Interface**: Next.js-based frontend for easy resume generation and preview
- **REST API**: Full-featured API for programmatic resume generation
- **Template System**: Flexible template system supporting both LaTeX and HTML output

### Technical Features
- **Docker Optimization**: Containerized environment with optimized build process
- **CLI Tools**: Enhanced command-line interface with validation, preview, and generation commands
- **JSON Resume Schema**: Full compatibility with the JSON Resume standard
- **Responsive Design**: Web interface works seamlessly on desktop and mobile devices
- **Real-time Preview**: Live preview of generated resumes in the web interface

## Prerequisites

### For Docker-based Usage
- Docker and Docker Compose
- [just](https://github.com/casey/just) command runner

### For Local Development
- Go 1.21+ (for API and CLI development)
- Node.js 18+ and npm (for frontend development)
- Docker (for containerized builds)

## Getting Started

### Quick Start with Docker Compose

1. **Clone the Repository**
   ```bash
   git clone https://github.com/urmzd/resume-generator.git
   cd resume-generator
   ```

2. **Start the Full Platform**
   ```bash
   docker-compose up
   ```
   This starts both the API server (port 8080) and web frontend (port 3000).

3. **Access the Web Interface**
   Open your browser and navigate to `http://localhost:3000`

### Development Setup

1. **Initialize the Project**
   ```bash
   just init
   just deps
   ```

2. **Start Development Environment**
   ```bash
   just dev
   ```
   This runs both the API server and frontend in development mode.

3. **Frontend Development Only**
   ```bash
   just frontend-install
   just frontend
   ```

### CLI Usage

1. **Build the Application**
   ```bash
   just build
   ```

2. **Generate Resume (Enhanced CLI)**
   ```bash
   # Generate HTML resume with custom output path
   ./resume-generator run -i examples/sample-enhanced.yml -o examples/output.pdf -t modern-html

   # Generate with default filename (first_last_resume_timestamp.pdf)
   ./resume-generator run -i examples/sample-enhanced.yml -t modern-html

   # Generate LaTeX PDF
   ./resume-generator run -i examples/sample-enhanced.yml -o ~/Documents/my_resume.pdf -t base-latex

   # Validate configuration
   ./resume-generator validate examples/sample-enhanced.yml

   # Preview without generation
   ./resume-generator preview examples/sample-enhanced.yml

   # List available templates
   ./resume-generator templates list
   ```

3. **Path Resolution**

   The CLI now supports flexible path resolution:
   - **Relative paths**: `./examples/resume.yml`, `../data/resume.yml`
   - **Absolute paths**: `/Users/name/Documents/resume.yml`
   - **Home directory**: `~/Documents/resume.yml`
   - **Custom output locations**: Specify any file path for output
  - **Directory output**: Provide a directory, and a timestamped workspace will be created (with PDF + debug artifacts)

   Examples:
   ```bash
   # Relative input, custom output
   ./resume-generator run -i examples/sample.yml -o output/my_resume.pdf

   # Absolute paths
   ./resume-generator run -i /path/to/resume.yml -o /path/to/output.pdf

   # Home directory paths
   ./resume-generator run -i ~/resumes/resume.yml -o ~/Documents/resume.pdf

   # Output to directory (creates timestamped folder)
   ./resume-generator run -i resume.yml -o ~/Documents/
   ```

   Each run results in a directory named `first[_middle]_last_<timestamp>/` containing the generated `resume.pdf` (or your custom filename) along with a `debug/` subfolder that preserves the rendered `.tex`/`.html`, `.log`, `.aux`, and supporting class files.

4. **Custom Asset Locations**
   ```bash
   # Custom LaTeX classes folder
   ./resume-generator run -i resume.yml -c /path/to/classes -t base-latex
   ```

## Showcase

Here are some examples of resumes generated with our tool:

### Sample Resume 1

![Sample Resume 1](assets/example_results/example.jpg)

A clean, professional layout suitable for various industries.

## Generators and Templates

This tool supports different generators and templates to customize your resume.

## API Usage

The resume generator provides a REST API for programmatic access:

### Start API Server
```bash
just serve 8080
```

### API Endpoints

- **POST** `/api/v1/generate` - Generate resume from configuration
- **GET** `/api/v1/health` - Health check
- **GET** `/api/v1/formats` - List available output formats
- **GET** `/api/v1/templates` - List available templates

### Example API Usage
```bash
# Test all endpoints
just api-test

# Generate resume via API
curl -X POST http://localhost:8080/api/v1/generate \
  -F "config=@examples/sample-enhanced.yml" \
  -F "format=html" \
  -F "template=modern"
```

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
- **PDF Templates**: LaTeX-based templates in `assets/templates/`
- **HTML Templates**: Modern responsive templates in `assets/templates/html/`
- **Custom Templates**: Create your own templates following the provided patterns

```bash
# List available templates
just templates

# Use specific template
just generate config.yml output pdf custom-template
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

## Assets/Helpers Directory

This directory is specifically designed to assist you in creating a more effective and professional resume. Here's what you can expect to find inside:

### Key Features of the `assets` Folder:

1. **Keywords and Phrases**: This section contains a curated list of powerful keywords and phrases relevant to various industries and job roles. Incorporating these into your resume can significantly enhance its impact and help in catching the attention of recruiters and resume screening software.

2. **Formatting Templates**: We provide a selection of formatting templates to give your resume a polished and professional look. These templates are designed to be easily customizable to fit your personal style and the requirements of your industry.

3. **Action Verbs**: A comprehensive list of action verbs is included to help you describe your experiences and achievements in a dynamic and compelling way. These verbs are crucial for making your resume more engaging and effective.

4. **Examples and Samples**: To give you a better idea of how to craft your resume, we've included examples and sample resumes. These can serve as a great starting point or source of inspiration for your own resume.

## Contributing

Contributions to the Generate Resumes project are welcome. Please read our contributing guidelines and submit pull requests for any enhancements, bug fixes, or documentation improvements.

## License

This project is licensed under the [MIT License](LICENSE). Feel free to use, modify, and distribute the code as per the license terms.

## Acknowledgments

Thanks to all the contributors who have helped in building and maintaining this tool. Special thanks to the LaTeX community for the underlying typesetting system.
