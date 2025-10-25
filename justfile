# Resume Generator justfile
# Run `just --list` to see all available commands

# Variables
organization := "urmzd"
version := env_var_or_default("VERSION", "latest")
image_tag := organization + "/" + "resume-generator" + ":" + version

# Directories
outputs_dir := "outputs"
inputs_dir := "inputs"
examples_dir := "examples"
assets_dir := "assets"

# Default recipe to display help
default:
    @just --list

# Initialize directories & copy examples
init:
    @echo "Initializing {{inputs_dir}} and {{outputs_dir}}"
    mkdir -p {{inputs_dir}} {{outputs_dir}}
    cp -r {{examples_dir}}/* {{inputs_dir}}/

# Build the Docker image (multistage build with Go + TeX)
build:
    @echo "Building Docker image {{image_tag}}"
    docker build --tag {{image_tag}} .

# Push the image to Docker Hub
push: build
    @echo "Pushing image to Docker Hub"
    docker push {{image_tag}}


# Run the resume-generator inside Docker
docker-run filename format="pdf":
    @echo "Running resume-generator in Docker"
    docker run --rm \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{outputs_dir}}:/outputs" \
      -v "{{justfile_directory()}}/{{assets_dir}}:/assets" \
      {{image_tag}} \
      run -i /inputs/{{filename}} -o /outputs -f {{format}}

# Exec an arbitrary command in the Docker container
exec cmd:
    @echo "Executing in Docker container"
    docker run --rm -it \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{assets_dir}}:/assets" \
      {{image_tag}} \
      {{cmd}}

# Start an interactive shell in Docker
shell:
    @echo "Launching shell in Docker container"
    docker run --rm -it \
      -v "{{justfile_directory()}}/{{inputs_dir}}:/inputs" \
      -v "{{justfile_directory()}}/{{outputs_dir}}:/outputs" \
      -v "{{justfile_directory()}}/{{assets_dir}}:/assets" \
      --entrypoint /bin/sh \
      {{image_tag}}

# Clean generated outputs & inputs
clean:
    @echo "Cleaning up {{inputs_dir}} and {{outputs_dir}}"
    rm -rf {{inputs_dir}} {{outputs_dir}}

# Run all examples through the pipeline
examples: init
    @echo "Running all example resumes"
    #!/usr/bin/env bash
    for f in {{examples_dir}}/*; do \
        filename=$(basename "$f"); \
        just docker-run "$filename" "pdf"; \
    done

# Build and run a specific example
example filename: init build
    just docker-run {{filename}} "pdf"

# Full build and test workflow
test: build
    just init
    just example "sample-enhanced.yml"
    @echo "Build and test completed successfully!"

# Generate a resume using the CLI (local Go build)
generate input_file output_file="resume" format="html" template="modern":
    @echo "Generating resume from {{input_file}}"
    go run main.go run -i {{input_file}} -o {{output_file}}.{{format}} -f {{format}} -t {{template}}

# Validate a resume configuration file
validate input_file:
    @echo "Validating resume configuration: {{input_file}}"
    go run main.go validate {{input_file}}

# Preview resume configuration (no compilation)
preview input_file:
    @echo "Previewing resume configuration: {{input_file}}"
    go run main.go preview {{input_file}}

# List available templates
templates:
    @echo "Available templates:"
    go run main.go templates list

# Install Go dependencies
deps:
    @echo "Installing Go dependencies"
    go mod tidy

# Clean everything including Docker images
clean-all: clean
    @echo "Removing Docker image"
    @-docker rmi {{image_tag}}
    docker system prune -f