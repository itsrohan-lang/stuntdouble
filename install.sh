#!/usr/bin/env bash
set -e

echo "🚀 Installing StuntDouble Native CLI..."

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
else
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
fi

REPO="itsrohan-lang/stuntdouble"
LATEST_TAG=$(curl -sL https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "⚠️  Could not fetch latest release. Defaulting to v2.0.0"
    LATEST_TAG="v2.0.0"
fi

BINARY_NAME="stuntdouble-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${BINARY_NAME}"

echo ">> Downloading ${BINARY_NAME} (${LATEST_TAG})..."
curl -sSL -o stuntdouble "$DOWNLOAD_URL"

chmod +x stuntdouble

echo ">> Installing to /usr/local/bin (may require sudo)..."
sudo mv stuntdouble /usr/local/bin/stuntdouble || {
    echo "❌ Failed to move binary to /usr/local/bin. Try running the script with sudo."
    exit 1
}

echo ""
echo "✅ StuntDouble successfully installed!"
echo "Run 'stuntdouble --help' to get started."
