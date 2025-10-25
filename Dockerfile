# ================================
# Resume Generator - Multistage Dockerfile
# ================================
# Builds a minimal resume generator with LaTeX support
# ================================

# ┌─────────────────────────────────────────────────────────┐
# │ Stage 1: Go Builder                                      │
# └─────────────────────────────────────────────────────────┘
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Cache Go modules for faster builds
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source and build optimized binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o resume-generator \
    main.go

# Verify binary works
RUN ./resume-generator --help

# ┌─────────────────────────────────────────────────────────┐
# │ Stage 2: Final Runtime Image                            │
# └─────────────────────────────────────────────────────────┘
FROM alpine:3.19

LABEL maintainer="urmzd"
LABEL description="Resume Generator - Generate professional resumes from YAML/JSON/TOML"
LABEL version="2.0"

# Install runtime dependencies
# - texlive: Full TeX distribution for LaTeX PDF generation
# - chromium: For HTML to PDF conversion
# - curl: For healthchecks
RUN apk add --no-cache \
    texlive \
    texlive-xetex \
    texlive-luatex \
    texmf-dist-latexextra \
    texmf-dist-fontsextra \
    chromium \
    curl \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Set Chromium path for HTML to PDF conversion
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/resume-generator /usr/local/bin/resume-generator

# Ensure binary is executable
RUN chmod +x /usr/local/bin/resume-generator

# Create directories for volumes
RUN mkdir -p /templates /examples /inputs /outputs /tmp/uploads /tmp/downloads

# Copy default templates
COPY templates /templates/

# Expose templates within the working directory for relative lookups
RUN ln -sfn /templates /app/templates

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD resume-generator --help > /dev/null || exit 1

# Environment variables
ENV PATH="/usr/local/bin:${PATH}" \
    GIN_MODE=release \
    PORT=8080

# Default entrypoint
ENTRYPOINT ["resume-generator"]
CMD ["--help"]

# ┌─────────────────────────────────────────────────────────┐
# │ Usage Examples                                           │
# └─────────────────────────────────────────────────────────┘
# Build:
#   docker build -t resume-generator .
#
# Generate HTML:
#   docker run --rm -v $(pwd):/work resume-generator run -i /work/resume.yml -f html -o /work
#
# Generate PDF:
#   docker run --rm -v $(pwd):/work resume-generator run -i /work/resume.yml -f pdf -o /work
#
# List templates:
#   docker run --rm -v $(pwd):/work resume-generator templates list
#
# Validate:
#   docker run --rm -v $(pwd):/work resume-generator validate /work/resume.yml
#
# Preview:
#   docker run --rm -v $(pwd):/work resume-generator preview /work/resume.yml
