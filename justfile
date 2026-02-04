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
examples_dir := "assets/example_resumes"
templates_dir := "templates"

# Default recipe to display help
default:
    @just --list

# Run the CLI (local binary if present, otherwise go run)
cli *args:
    @if [ -x "{{cli_binary}}" ]; then \
        "{{cli_binary}}" {{args}}; \
    else \
        go run main.go {{args}}; \
    fi

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

# Build the Docker image (multistage build with Go + TeX)
docker-build:
    @echo "Building Docker image {{image_tag}}"
    docker build --tag {{image_tag}} .

# Run the resume-generator inside Docker
docker-run input_file=(examples_dir + "/software_engineer.yml") output_dir=outputs_dir template="modern-html": docker-build
    @echo "Running resume-generator in Docker"
    docker run --rm \
      -v "{{justfile_directory()}}:/work" \
      -v "{{justfile_directory()}}/{{output_dir}}:/outputs" \
      -v "{{justfile_directory()}}/{{templates_dir}}:/templates" \
      {{image_tag}} \
      run -i /work/{{input_file}} -o /outputs -t {{template}}

# Full build and test workflow
test: build-cli
    @echo "Running Go unit tests"
    just go-test

# Generate a resume using the CLI (local Go build)
generate input_file=(examples_dir + "/software_engineer.yml") output_dir=outputs_dir template="":
    @echo "Generating resume from {{input_file}} into {{output_dir}}"
    @if [ -n "{{template}}" ]; then \
        just cli run -i {{input_file}} -o {{output_dir}} -t {{template}}; \
    else \
        just cli run -i {{input_file}} -o {{output_dir}}; \
    fi

# Generate the JSON schema for the resume input format
schema output=(justfile_directory() + "/assets/schema/resume.schema.json"):
    @echo "Generating JSON schema -> {{output}}"
    just cli schema -o {{output}}

# Generate README preview images for the bundled templates
readme-previews:
    @scripts/generate-readme-previews.sh

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

# Clean generated outputs & CLI binary
clean:
    @echo "Cleaning up {{outputs_dir}} and {{cli_output}}"
    rm -rf {{outputs_dir}}
    @rm -f {{cli_output}}

# Clean everything including Docker images
clean-all: clean
    @echo "Removing Docker image"
    @-docker rmi {{image_tag}}
    docker system prune -f
