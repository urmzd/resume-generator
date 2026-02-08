# Resume Generator

cli_binary := "resume-generator"
example_input := "assets/example_resumes/software_engineer.yml"
outputs_dir := "outputs"

default:
    @just --list

# Install dependencies and tools
init:
    go mod download && go mod tidy
    go install github.com/goreleaser/goreleaser/v2@latest

# Build the CLI binary for the current architecture
install:
    go build -trimpath -ldflags="-s -w" -o {{cli_binary}} .

# Build (if needed) and generate a resume
run input=example_input output=outputs_dir *args="": install
    @mkdir -p {{output}}
    ./{{cli_binary}} run -i {{input}} -o {{output}} {{args}}
