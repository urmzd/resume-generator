# Incipit — here begins the new career.

example_input := "assets/example_resumes/software_engineer.yml"
outputs_dir := "outputs"

default:
    @just --list

# Install dependencies and tools
init:
    git config core.hooksPath .githooks
    go mod download && go mod tidy

# Install CLI binary to $GOPATH/bin
install:
    CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" ./cmd/incipit

# Build (if needed) and generate a resume
run input=example_input output=outputs_dir *args="": install
    @mkdir -p {{output}}
    incipit run -i {{input}} -o {{output}} {{args}}

# Generate PNG screenshots for each template
screenshots: install
    incipit screenshots --input {{example_input}}

