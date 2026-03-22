#!/usr/bin/env bash
set -euo pipefail

REPO="urmzd/incipit"
INSTALL_DIR="$HOME/.local/bin"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Darwin) PLATFORM="darwin" ;;
  Linux)  PLATFORM="linux" ;;
  *)
    echo "Error: Unsupported OS '$OS'. This script supports macOS and Linux." >&2
    exit 1
    ;;
esac

# Check dependencies
if ! command -v curl &>/dev/null; then
  echo "Error: curl is required but not installed." >&2
  exit 1
fi

# Fetch latest release tag
echo "Fetching latest release..."
RELEASE_JSON="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")"
TAG="$(echo "$RELEASE_JSON" | grep '"tag_name"' | head -1 | sed 's/.*: *"//;s/".*//')"

if [ -z "$TAG" ]; then
  echo "Error: Could not determine latest release tag." >&2
  exit 1
fi

echo "Latest release: $TAG"

# Download and install
TMPDIR_INSTALL="$(mktemp -d)"
trap 'rm -rf "$TMPDIR_INSTALL"' EXIT

if [ "$PLATFORM" = "darwin" ]; then
  ASSET_NAME="incipit-darwin"
else
  ASSET_NAME="incipit-${INCIPIT_VARIANT:-linux}"
fi

ASSET_URL="https://github.com/$REPO/releases/download/$TAG/$ASSET_NAME"
echo "Downloading $ASSET_URL..."
curl -fsSL -o "$TMPDIR_INSTALL/incipit" "$ASSET_URL"
BINARY="$TMPDIR_INSTALL/incipit"

if [ ! -f "$BINARY" ]; then
  echo "Error: Binary not found after download." >&2
  exit 1
fi

# Install
mkdir -p "$INSTALL_DIR"
cp "$BINARY" "$INSTALL_DIR/incipit"
chmod +x "$INSTALL_DIR/incipit"

echo "Installed incipit ($TAG) to $INSTALL_DIR/incipit"

# Check PATH and ensure binary is accessible for init
INCIPIT_BIN="$INSTALL_DIR/incipit"
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "WARNING: $INSTALL_DIR is not in your PATH."
    echo "Add it by appending this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo ""
    echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    echo ""
    ;;
esac

# Initialize: download bundled templates and create config
echo "Initializing templates..."
"$INCIPIT_BIN" init --version "$TAG"
