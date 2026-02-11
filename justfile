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

# Record desktop demo (video + screenshot). Starts and stops wails dev automatically.
demo-desktop:
    #!/usr/bin/env bash
    set -euo pipefail

    # Start Wails dev server in background
    wails dev &
    WAILS_PID=$!
    trap 'kill $WAILS_PID 2>/dev/null; wait $WAILS_PID 2>/dev/null' EXIT

    # Wait for dev server to become ready (up to 120s for first-time builds)
    echo "Waiting for Wails dev server on :34115..."
    for i in $(seq 1 120); do
        if curl -sf http://localhost:34115 > /dev/null 2>&1; then
            echo "Dev server ready."
            break
        fi
        if [ "$i" -eq 120 ]; then
            echo "Timeout: dev server did not start within 120s."
            exit 1
        fi
        sleep 1
    done

    # Run Playwright test (produces video + screenshot)
    cd e2e/desktop && npx playwright test

    # Copy recorded video
    VIDEO=$(find ../../assets/playwright-results -name 'video.webm' | head -1)
    if [ -n "$VIDEO" ]; then
        cp "$VIDEO" ../../assets/demo-desktop.webm
    fi

# Record all demos
demo: demo-cli demo-desktop
