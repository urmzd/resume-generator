# Resume Generator justfile
# Run `just --list` to see all available commands

# Variables
organization := "urmzd"
version := env_var_or_default("VERSION", "latest")
image_tag := organization + "/" + "resume-generator" + ":" + version
cli_output := env_var_or_default("CLI_OUTPUT", "resume-generator")
cli_binary := justfile_directory() + "/" + cli_output
release_ldflags := env_var_or_default("GO_LDFLAGS", "-s -w")

# Directories
outputs_dir := "../outputs"
inputs_dir := "../inputs"
examples_dir := "assets/example_resumes"
templates_dir := "templates"

# Build the CLI binary for local development or release
build-cli mode="dev":
    @echo "Building CLI binary (mode={{mode}}) -> {{cli_output}}"
    @mkdir -p "$(dirname {{cli_binary}})"
    @if [ "{{mode}}" = "release" ]; then \
        CGO_ENABLED=${CGO_ENABLED:-0} go build -trimpath -ldflags "{{release_ldflags}}" -o {{cli_output}} .; \
    else \
        go build -o {{cli_output}} .; \
    fi

# Run Go unit tests
go-test pattern="./...":
    @echo "Running Go tests: {{pattern}}"
    go test {{pattern}}

# Default recipe to display help
default:
    @just --list

# Initialize directories & copy examples
init:
    @echo "Initializing {{inputs_dir}} and {{outputs_dir}}"
    mkdir -p {{inputs_dir}} {{outputs_dir}}
    cp -r {{examples_dir}}/* {{inputs_dir}}/

# Build the Docker image (multistage build with Go + TeX)
docker-build:
    @echo "Building Docker image {{image_tag}}"
    docker build --tag {{image_tag}} .

# Push the image to Docker Hub
push: docker-build
    @echo "Pushing image to Docker Hub"
    docker push {{image_tag}}

# Build both the local CLI binary and Docker image
build: build-cli docker-build
    @echo "Built CLI binary ({{cli_output}}) and Docker image ({{image_tag}})"

# Run the resume-generator inside Docker
docker-run filename template="modern-html":
    @echo "Running resume-generator in Docker"
    docker run --rm \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{outputs_dir}}:/outputs" \
      -v "{{justfile_directory()}}/{{templates_dir}}:/templates" \
      {{image_tag}} \
      run -i /inputs/{{filename}} -o /outputs -t {{template}}

# Exec an arbitrary command in the Docker container
exec cmd:
    @echo "Executing in Docker container"
    docker run --rm -it \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{templates_dir}}:/templates" \
      {{image_tag}} \
      {{cmd}}

# Start an interactive shell in Docker
shell:
    @echo "Launching shell in Docker container"
    docker run --rm -it \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{outputs_dir}}:/outputs" \
      -v "{{justfile_directory()}}/{{templates_dir}}:/templates" \
      --entrypoint /bin/sh \
      {{image_tag}}

# Clean generated outputs & inputs
clean:
    @echo "Cleaning up {{inputs_dir}}, {{outputs_dir}}, and {{cli_output}}"
    rm -rf {{inputs_dir}} {{outputs_dir}}
    @rm -f {{cli_output}}

# Run all examples through the pipeline
examples: init docker-build
    @echo "Running all example resumes"
    #!/usr/bin/env bash
    for f in {{examples_dir}}/*; do \
        filename=$(basename "$f"); \
        just docker-run "$filename"; \
    done

# Build and run a specific example
example filename: init docker-build
    just docker-run {{filename}}

# Convert a generated PDF into a README-friendly JPEG preview
pdf-to-jpeg pdf output="docs/readme-preview.jpg":
    @echo "Converting {{pdf}} -> {{output}}"
    @pdf_path="{{pdf}}"; \
     output_path="{{output}}"; \
     mkdir -p "$(dirname "$output_path")"; \
     if command -v magick >/dev/null 2>&1; then \
        magick -density 300 "${pdf_path}[0]" -quality 85 "$output_path"; \
     elif command -v convert >/dev/null 2>&1; then \
        convert -density 300 "${pdf_path}[0]" -quality 85 "$output_path"; \
     elif command -v pdftoppm >/dev/null 2>&1; then \
        output_base="${output_path%.*}"; \
        pdftoppm -jpeg -singlefile -r 300 "$pdf_path" "$output_base"; \
     else \
        echo "Error: install ImageMagick (magick/convert) or poppler-utils (pdftoppm) for PDF → JPEG conversion." >&2; \
        exit 1; \
     fi

# Full build and test workflow
test: build-cli docker-build
    @echo "Running Go unit tests"
    just go-test
    just init
    just example "example.yml"
    @echo "Build and test completed successfully!"

# Generate a resume using the CLI (local Go build)
generate input_file="./assets/example_resumes/example.yml" output_dir="outputs":
    @echo "Generating resume from {{input_file}} into {{output_dir}}"
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" run -i {{input_file}} -o {{output_dir}}; \
    else \
        go run main.go run -i {{input_file}} -o {{output_dir}}; \
    fi

# Validate a resume configuration file
validate input_file:
    @echo "Validating resume configuration: {{input_file}}"
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" validate {{input_file}}; \
    else \
        go run main.go validate {{input_file}}; \
    fi

# Preview resume configuration (no compilation)
preview input_file:
    @echo "Previewing resume configuration: {{input_file}}"
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" preview {{input_file}}; \
    else \
        go run main.go preview {{input_file}}; \
    fi

# List available templates
templates:
    @echo "Available templates:"
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" templates list; \
    else \
        go run main.go templates list; \
    fi

# Check available LaTeX engines
latex-engines:
    @echo "Checking available LaTeX engines:"
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" templates engines; \
    else \
        go run main.go templates engines; \
    fi

# Generate JSON schema for resume format
schema output_file="":
    @echo "Generating JSON schema"
    @if [ -z "{{output_file}}" ]; then \
        if [ -x "{{cli_binary}}" ]; then \
            "{{cli_binary}}" schema; \
        else \
            go run main.go schema; \
        fi; \
    else \
        if [ -x "{{cli_binary}}" ]; then \
            "{{cli_binary}}" schema -o {{output_file}}; \
        else \
            go run main.go schema -o {{output_file}}; \
        fi; \
    fi

# Format Go code
fmt:
    @echo "Formatting Go code"
    gofmt -w .

# Run Go linter
lint:
    @echo "Running Go linter"
    @if command -v golangci-lint &> /dev/null; then \
        golangci-lint run; \
    else \
        echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
    fi

# Install Go dependencies
deps:
    @echo "Installing Go dependencies"
    go mod tidy

# Clean everything including Docker images
clean-all: clean
    @echo "Removing Docker image"
    @-docker rmi {{image_tag}}
    docker system prune -f

check-latest-template output_path="./outputs" template="modern-latex":
    @echo "Checking latest template: {{template}}"
    @zsh -c 't="{{template}}"; slug="${t//-/}"; ls -t {{output_path}}/**/$slug/**/*.pdf | head -n 1 | xargs open'
