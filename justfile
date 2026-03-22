# Incipit — here begins the new career.

cli_binary := "incipit"
example_input := "assets/example_resumes/software_engineer.yml"
outputs_dir := "outputs"

default:
    @just --list

# Install dependencies and tools
init:
    git config core.hooksPath .githooks
    go mod download && go mod tidy

# Build CLI binary
install:
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o {{cli_binary}} ./cmd/incipit

# Build (if needed) and generate a resume
run input=example_input output=outputs_dir *args="": install
    @mkdir -p {{output}}
    ./{{cli_binary}} run -i {{input}} -o {{output}} {{args}}

# Generate PNG screenshots for each template
screenshots: install
    ./{{cli_binary}} screenshots --input {{example_input}}

