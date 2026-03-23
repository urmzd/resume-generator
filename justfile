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

# Generate showcase assets (template previews + demo GIF) via teasr
showcase: install
    @mkdir -p showcase/pdfs
    incipit run -i {{example_input}} -o showcase/build
    @find showcase/build -name '*.pdf' -exec sh -c 'cp "$1" showcase/pdfs/"$(echo "$1" | sed "s/.*\.\(.*\)\.pdf/\1.pdf/")"' _ {} \;
    teasr showme
