#!/usr/bin/env bash
# install.sh — One-step bootstrap: clone (if needed), build, and deploy movie CLI.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/movie-cli-v3/main/install.sh | bash
#   bash install.sh
#   bash install.sh --dir ~/bin
#
# Requires: Git, Go 1.22+, Bash 4+

set -euo pipefail

# ── Config ────────────────────────────────────────────────────
REPO_NAME="movie-cli-v3"
REPO_URL="https://github.com/alimtvnetwork/movie-cli-v3.git"
BINARY_NAME="movie"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"
INSTALL_DIR=""

# ── Helpers ───────────────────────────────────────────────────
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
GRAY='\033[0;90m'
MAGENTA='\033[0;35m'
DCYAN='\033[0;36m'
NC='\033[0m'

ok()   { printf "  ${GREEN}OK${NC} %s\n" "$1"; }
info() { printf "  ${CYAN}->${NC} ${GRAY}%s${NC}\n" "$1"; }
err()  { printf "  ${RED}XX${NC} ${RED}%s${NC}\n" "$1"; }
die()  { err "$1"; [ -n "${2:-}" ] && info "$2"; exit 1; }

banner() {
    echo ""
    printf " ${DCYAN}+======================================+${NC}\n"
    printf " ${DCYAN}|${NC} ${CYAN}movie installer${NC}                  ${DCYAN}|${NC}\n"
    printf " ${DCYAN}+======================================+${NC}\n"
    echo ""
}

section() {
    echo ""
    printf " ${MAGENTA}%s${NC}\n" "$1"
    printf " ${GRAY}%s${NC}\n" "--------------------------------------------------"
}

# ── Parse arguments ───────────────────────────────────────────
while [[ $# -gt 0 ]]; do
    case "$1" in
        --dir)   INSTALL_DIR="$2"; shift 2 ;;
        --dir=*) INSTALL_DIR="${1#*=}"; shift ;;
        -h|--help)
            echo "Usage: install.sh [--dir <path>]"
            echo ""
            echo "Options:"
            echo "  --dir <path>   Install directory (default: ~/.local/bin)"
            echo ""
            exit 0
            ;;
        *) die "Unknown option: $1" "Use --help for usage" ;;
    esac
done

[ -z "$INSTALL_DIR" ] && INSTALL_DIR="$DEFAULT_INSTALL_DIR"

# ── Main ──────────────────────────────────────────────────────
banner

# [1/5] Prerequisites
section "[1/5] Checking prerequisites"

command -v git >/dev/null 2>&1 || die "Git is not installed" "Install from https://git-scm.com/downloads"
ok "Git: $(git --version)"

command -v go >/dev/null 2>&1 || die "Go is not installed" "Install from https://go.dev/dl/"
ok "Go: $(go version)"

# [2/5] Locate or clone repo
section "[2/5] Locating repository"

REPO_ROOT=""

if [ -f "go.mod" ] && [ -f "run.ps1" ]; then
    REPO_ROOT="$(pwd)"
    ok "Already inside repo: $REPO_ROOT"
elif [ -d "$REPO_NAME" ] && [ -f "$REPO_NAME/go.mod" ]; then
    REPO_ROOT="$(cd "$REPO_NAME" && pwd)"
    ok "Found repo: $REPO_ROOT"
else
    info "Cloning $REPO_URL ..."
    git clone "$REPO_URL" || die "Clone failed"
    REPO_ROOT="$(cd "$REPO_NAME" && pwd)"
    ok "Cloned to: $REPO_ROOT"
fi

# [3/5] Build
section "[3/5] Building binary"

cd "$REPO_ROOT"

go mod tidy
info "Dependencies resolved"

VERSION="$(git describe --tags --always 2>/dev/null || echo "dev")"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
BUILD_DATE="$(date -u +%Y-%m-%d)"

LDFLAGS="-s -w \
  -X 'github.com/mahin/movie-cli-v2/version.Version=$VERSION' \
  -X 'github.com/mahin/movie-cli-v2/version.Commit=$COMMIT' \
  -X 'github.com/mahin/movie-cli-v2/version.BuildDate=$BUILD_DATE'"

# Detect OS/arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
esac
info "Target: ${OS}/${ARCH}"

CGO_ENABLED=0 GOOS="$OS" GOARCH="$ARCH" go build -ldflags "$LDFLAGS" -o "$BINARY_NAME" .
ok "Built: $BINARY_NAME ($(du -h "$BINARY_NAME" | cut -f1))"

# [4/5] SHA256 verification & deploy
section "[4/5] Deploying to $INSTALL_DIR"

# Generate SHA256 checksum
SHA256="$(shasum -a 256 "$BINARY_NAME" 2>/dev/null || sha256sum "$BINARY_NAME" 2>/dev/null)"
SHA256="${SHA256%% *}"
info "SHA256: $SHA256"

# Create install dir
mkdir -p "$INSTALL_DIR"

# Backup existing binary
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    mv "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/${BINARY_NAME}.bak"
    info "Backed up existing binary"
fi

# Copy and set executable
cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"
ok "Deployed to: $INSTALL_DIR/$BINARY_NAME"

# Verify deployed binary checksum
DEPLOYED_SHA="$(shasum -a 256 "$INSTALL_DIR/$BINARY_NAME" 2>/dev/null || sha256sum "$INSTALL_DIR/$BINARY_NAME" 2>/dev/null)"
DEPLOYED_SHA="${DEPLOYED_SHA%% *}"

if [ "$SHA256" = "$DEPLOYED_SHA" ]; then
    ok "SHA256 verified ✓"
else
    die "SHA256 mismatch after deploy!" "Built: $SHA256 / Deployed: $DEPLOYED_SHA"
fi

# Clean up build artifact and backup
rm -f "$BINARY_NAME"
rm -f "$INSTALL_DIR/${BINARY_NAME}.bak"

# [5/5] Verify installation
section "[5/5] Verifying installation"

# Check if install dir is in PATH
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        info "$INSTALL_DIR is not in your PATH"
        info "Add it with:"
        echo ""
        SHELL_NAME="$(basename "$SHELL" 2>/dev/null || echo "bash")"
        case "$SHELL_NAME" in
            zsh)  echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc && source ~/.zshrc" ;;
            fish) echo "  fish_add_path $INSTALL_DIR" ;;
            *)    echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc && source ~/.bashrc" ;;
        esac
        echo ""
        export PATH="$INSTALL_DIR:$PATH"
        ;;
esac

if [[ -x "$INSTALL_DIR/$BINARY_NAME" ]]; then
    VER_OUTPUT="$("$INSTALL_DIR/$BINARY_NAME" version 2>&1 || true)"
    ok "$BINARY_NAME is ready"
    echo "$VER_OUTPUT" | while IFS= read -r line; do
        printf "       %s\n" "$line"
    done
else
    die "'$INSTALL_DIR/$BINARY_NAME' not found after install" "Check that the install directory is writable and on your PATH"
fi

# ── Done ──────────────────────────────────────────────────────
echo ""
printf " ${DCYAN}+======================================+${NC}\n"
printf " ${DCYAN}|${NC} ${GREEN}Installation complete${NC}            ${DCYAN}|${NC}\n"
printf " ${DCYAN}+======================================+${NC}\n"
echo ""
