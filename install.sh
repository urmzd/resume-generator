#!/bin/sh
# install.sh — Installs the incipit binary from GitHub releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/urmzd/incipit/main/install.sh | sh
#
# Environment variables:
#   INCIPIT_VERSION     — version to install (e.g. "v1.0.0"); defaults to latest
#   INCIPIT_INSTALL_DIR — installation directory; defaults to $HOME/.local/bin

set -eu

REPO="urmzd/incipit"

# curl with optional auth — uses GH_TOKEN or GITHUB_TOKEN if set.
gh_curl() {
    token="${GH_TOKEN:-${GITHUB_TOKEN:-}}"
    if [ -n "$token" ]; then
        curl -fsSL -H "Authorization: token $token" "$@"
    else
        curl -fsSL "$@"
    fi
}

main() {
    os=$(uname -s)
    arch=$(uname -m)

    case "$os" in
        Linux)
            case "$arch" in
                x86_64)          target="linux-amd64" ;;
                aarch64|arm64)   target="linux-arm64" ;;
                *)               err "Unsupported Linux architecture: $arch" ;;
            esac
            ;;
        Darwin)
            case "$arch" in
                x86_64)  target="darwin-amd64" ;;
                arm64)   target="darwin-arm64" ;;
                *)       err "Unsupported macOS architecture: $arch" ;;
            esac
            ;;
        MINGW*|MSYS*|CYGWIN*|Windows_NT)
            err "Windows is not supported by this installer. Download a binary from https://github.com/$REPO/releases/latest"
            ;;
        *)
            err "Unsupported operating system: $os"
            ;;
    esac

    if [ -n "${INCIPIT_VERSION:-}" ]; then
        tag="$INCIPIT_VERSION"
    else
        tag=$(gh_curl "https://api.github.com/repos/$REPO/releases/latest" \
            | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
        if [ -z "$tag" ]; then
            err "Failed to fetch latest release tag"
        fi
    fi

    artifact="incipit-${target}"
    url="https://github.com/$REPO/releases/download/${tag}/${artifact}"

    install_dir="${INCIPIT_INSTALL_DIR:-$HOME/.local/bin}"
    mkdir -p "$install_dir"

    echo "Downloading incipit $tag for $target..."
    gh_curl "$url" -o "$install_dir/incipit"

    chmod +x "$install_dir/incipit"

    echo "Installed incipit to $install_dir/incipit"

    # Initialize: download bundled templates and create config
    echo "Initializing templates..."
    "$install_dir/incipit" init --version "$tag"

    case ":$PATH:" in
        *":$install_dir:"*) ;;
        *) add_to_path "$install_dir" ;;
    esac
}

add_to_path() {
    install_dir="$1"

    case "$(basename "$SHELL")" in
        zsh)  profile="$HOME/.zshrc" ;;
        bash)
            if [ -f "$HOME/.bashrc" ]; then
                profile="$HOME/.bashrc"
            else
                profile="$HOME/.profile"
            fi
            ;;
        fish) profile="$HOME/.config/fish/config.fish" ;;
        *)    profile="$HOME/.profile" ;;
    esac

    if [ "$(basename "$SHELL")" = "fish" ]; then
        if ! grep -q "$install_dir" "$profile" 2>/dev/null; then
            mkdir -p "$(dirname "$profile")"
            echo "" >> "$profile"
            echo "# Added by incipit installer" >> "$profile"
            echo "set -Ux fish_user_paths $install_dir \$fish_user_paths" >> "$profile"
            echo "Added $install_dir to $profile"
            echo "Restart your shell or run: source $profile"
        fi
    elif [ -n "$profile" ] && ! grep -q "$install_dir" "$profile" 2>/dev/null; then
        echo "" >> "$profile"
        echo "# Added by incipit installer" >> "$profile"
        echo "export PATH=\"$install_dir:\$PATH\"" >> "$profile"
        echo "Added $install_dir to $profile"
        echo "Restart your shell or run: source $profile"
    fi
}

err() {
    echo "Error: $1" >&2
    exit 1
}

main
