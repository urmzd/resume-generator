# ┌─────────────── Stage 1: Build Go binary ───────────────┐
FROM golang:1.22-rc-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY pkg/ pkg/
COPY cmd/ cmd/
COPY main.go .
RUN go build -o resume-generator main.go

# └─────────────────────────────────────────────────────────┘


# ┌─────────────── Stage 2: Install TeX Live ───────────────┐
FROM debian:bookworm-slim AS tex

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
      libfontconfig1 \
      fonts-dejavu-core \
      wget \
      perl \
      xz-utils \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp
RUN wget https://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz \
 && tar -xzf install-tl-unx.tar.gz \
 && cd install-tl-*/ \
 && perl install-tl --no-interaction --scheme=small \
 && rm -rf /tmp/install-tl-*

# Add TeX Live binaries to PATH
ENV PATH=/usr/local/texlive/2024/bin/x86_64-linux:$PATH

# Install extra packages via tlmgr
RUN tlmgr install enumitem titlesec \
 && tlmgr option autobackup 0

# Clean up package manager caches again
RUN apt-get purge -y --auto-remove wget perl xz-utils \
 && rm -rf /var/lib/apt/lists/*

# └──────────────────────────────────────────────────────────┘


# ┌─────────────── Stage 3: Final runtime image ───────────────┐
FROM tex

WORKDIR /app
COPY --from=builder /app/resume-generator .

ENTRYPOINT ["./resume-generator"]
# └────────────────────────────────────────────────────────────┘
