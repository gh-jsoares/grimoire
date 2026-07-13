#!/bin/sh
set -e

REPO="gh-jsoares/grimoire"
INSTALL_DIR="${GRIMOIRE_INSTALL_DIR:-/usr/local/bin}"

# Detect OS and arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS" && exit 1 ;;
esac

# Get latest version
VERSION="$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)"
if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest version"
    exit 1
fi

# Check if already up-to-date
if command -v grimoire >/dev/null 2>&1; then
    CURRENT="$(grimoire --version 2>/dev/null | awk '{print $2}')"
    if [ "$CURRENT" = "$VERSION" ]; then
        echo "grimoire is already up-to-date ($VERSION)"
        exit 0
    fi
fi

echo "Installing grimoire ${VERSION} (${OS}/${ARCH})..."

# Download and extract
URL="https://github.com/${REPO}/releases/download/${VERSION}/grimoire_${OS}_${ARCH}.tar.gz"
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

curl -sL "$URL" -o "$TMPDIR/grimoire.tar.gz"
tar -xzf "$TMPDIR/grimoire.tar.gz" -C "$TMPDIR"

# Install binary
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMPDIR/grimoire" "$INSTALL_DIR/grimoire"
else
    sudo mv "$TMPDIR/grimoire" "$INSTALL_DIR/grimoire"
fi

# Install man page
MAN_DIR="${GRIMOIRE_MAN_DIR:-/usr/local/share/man/man1}"
if [ -f "$TMPDIR/grimoire.1" ]; then
    if [ -w "$MAN_DIR" ] 2>/dev/null || mkdir -p "$MAN_DIR" 2>/dev/null; then
        mv "$TMPDIR/grimoire.1" "$MAN_DIR/grimoire.1"
    elif sudo mkdir -p "$MAN_DIR"; then
        sudo mv "$TMPDIR/grimoire.1" "$MAN_DIR/grimoire.1"
    fi
    echo "Installed man page to ${MAN_DIR}/grimoire.1"
fi

echo "Installed grimoire to ${INSTALL_DIR}/grimoire"
grimoire --version
