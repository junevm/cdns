#!/bin/sh
set -e

# CDNS Installer Script
# Usage: curl -sfL https://raw.githubusercontent.com/junevm/cdns/main/install.sh | sh

REPO="junevm/cdns"
BINARY="cdns"
INSTALL_DIR="/usr/local/bin"

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" != "linux" ]; then
    echo "This script only supports Linux."
    exit 1
fi

case $ARCH in
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo "Detected $OS/$ARCH..."

# Get latest release tag
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not find latest release."
    exit 1
fi

# Check if already installed and compare versions
if command -v "$BINARY" >/dev/null 2>&1; then
    CURRENT_VER=$($BINARY --version 2>/dev/null | awk '{print $NF}')
    if [ "$CURRENT_VER" = "${LATEST_TAG#v}" ]; then
        echo "$BINARY is already up to date (version $CURRENT_VER)."
        exit 0
    fi
    echo "Updating $BINARY from $CURRENT_VER to ${LATEST_TAG#v}..."
else
    echo "Installing $BINARY $LATEST_TAG..."
fi

# Construct download URL
# Example: https://github.com/junevm/cdns/releases/download/v1.0.0/cdns_1.0.0_Linux_amd64.tar.gz

VERSION=${LATEST_TAG#v}
OS_CAP="$(echo "$OS" | awk '{print toupper(substr($0,1,1))substr($0,2)}')"
# Correcting logic to handle GoReleaser naming convention: {{ .ProjectName }}_{{ .Version }}_{{ title .Os }}_{{ .Arch }}.tar.gz

DOWNLOAD_URL="https://github.com/junevm/cdns/releases/download/$LATEST_TAG/${BINARY}_${VERSION}_${OS_CAP}_${ARCH}.tar.gz"

# Download and extract
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/release.tar.gz"

tar -xzf "$TMP_DIR/release.tar.gz" -C "$TMP_DIR"

if [ ! -f "$TMP_DIR/$BINARY" ]; then
    echo "Error: Binary not found in archive."
    rm -rf "$TMP_DIR"
    exit 1
fi

# Install
echo "Installing to $INSTALL_DIR (requires sudo)..."
sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
sudo chmod +x "$INSTALL_DIR/$BINARY"

rm -rf "$TMP_DIR"

echo "Installation complete! Run '$BINARY --help' to get started."
