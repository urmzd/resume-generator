# Resume Generator

cli_binary := "resume-generator"
example_input := "assets/example_resumes/software_engineer.yml"
outputs_dir := "outputs"

default:
    @just --list

# Install dependencies and tools
init:
    go mod download && go mod tidy
    cd frontend && npm install
    brew install vhs
    cd e2e/desktop && npm install && npx playwright install chromium

# Install Air live-reload tool
install-air:
    go install github.com/air-verse/air@latest

# Build CLI-only binary (no GUI, no CGO)
install:
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o {{cli_binary}} .

# Build desktop app (GUI + CLI, requires wails CLI)
build-desktop:
    wails build -trimpath -ldflags="-s -w"

# Build (if needed) and generate a resume
run input=example_input output=outputs_dir *args="": install
    @mkdir -p {{output}}
    ./{{cli_binary}} run -i {{input}} -o {{output}} {{args}}

# Dev: Wails dev mode with hot reload
dev:
    wails dev

# Dev: Clean frontend cache, rebuild, and start dev mode
dev-clean:
    rm -rf frontend/dist frontend/node_modules/.vite
    cd frontend && npm run build
    wails dev

# Record CLI demo GIF (requires: brew install vhs)
demo-cli:
    vhs e2e/demo.tape

# Record desktop demo video (requires: wails dev running, npx playwright installed)
demo-desktop:
    cd e2e/desktop && npx playwright test

# Record all demos
demo: demo-cli demo-desktop
