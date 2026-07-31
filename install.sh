#!/usr/bin/env bash
#
# StuntDouble CLI installer.
#
# Downloads the release binary for this platform, verifies it against the
# SHA256SUMS file published with the release, and installs it to
# /usr/local/bin. Any failure aborts the install rather than leaving a
# partially-verified binary on the system.
#
# Environment:
#   STUNTDOUBLE_VERSION   install a specific tag (default: latest release)
#   INSTALL_DIR           install target (default: /usr/local/bin)
#
set -euo pipefail

REPO="itsrohan-lang/stuntdouble"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

die() {
    echo "❌ $*" >&2
    exit 1
}

need() {
    command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"
}

need curl

# Resolve a SHA-256 tool. macOS ships shasum; most Linux images ship sha256sum.
if command -v sha256sum >/dev/null 2>&1; then
    sha256() { sha256sum "$1" | cut -d' ' -f1; }
elif command -v shasum >/dev/null 2>&1; then
    sha256() { shasum -a 256 "$1" | cut -d' ' -f1; }
else
    die "need sha256sum or shasum to verify the download"
fi

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux|darwin) ;;
    *) die "unsupported operating system: $OS" ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) die "unsupported architecture: $ARCH" ;;
esac

echo "🚀 Installing StuntDouble CLI (${OS}/${ARCH})..."

# Resolve the version. Never fall back to a guessed tag: installing an
# arbitrary older release because an API call failed is worse than failing.
if [ -n "${STUNTDOUBLE_VERSION:-}" ]; then
    TAG="$STUNTDOUBLE_VERSION"
else
    echo ">> Resolving latest release..."
    TAG="$(curl --fail --silent --show-error --location \
        "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name":' \
        | sed -E 's/.*"([^"]+)".*/\1/')" \
        || die "could not resolve the latest release. Set STUNTDOUBLE_VERSION=vX.Y.Z to pin a version."
    [ -n "$TAG" ] || die "could not parse a release tag from the GitHub API response."
fi

BINARY_NAME="stuntdouble-${OS}-${ARCH}"
BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"

TMPDIR_SD="$(mktemp -d)"
trap 'rm -rf "$TMPDIR_SD"' EXIT

echo ">> Downloading ${BINARY_NAME} (${TAG})..."
curl --fail --silent --show-error --location \
    -o "${TMPDIR_SD}/${BINARY_NAME}" "${BASE_URL}/${BINARY_NAME}" \
    || die "download failed: ${BASE_URL}/${BINARY_NAME}"

echo ">> Downloading checksums..."
curl --fail --silent --show-error --location \
    -o "${TMPDIR_SD}/SHA256SUMS" "${BASE_URL}/SHA256SUMS" \
    || die "could not download SHA256SUMS for ${TAG}. Refusing to install an unverified binary."

EXPECTED="$(grep -E "( |\*)${BINARY_NAME}\$" "${TMPDIR_SD}/SHA256SUMS" | awk '{print $1}' || true)"
[ -n "$EXPECTED" ] || die "no checksum for ${BINARY_NAME} in SHA256SUMS. Refusing to install."

ACTUAL="$(sha256 "${TMPDIR_SD}/${BINARY_NAME}")"
if [ "$EXPECTED" != "$ACTUAL" ]; then
    die "checksum mismatch for ${BINARY_NAME}
    expected: ${EXPECTED}
    actual:   ${ACTUAL}
The download may be corrupt or tampered with. Nothing was installed."
fi
echo "✅ Checksum verified."

echo ">> Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    install -m 0755 "${TMPDIR_SD}/${BINARY_NAME}" "${INSTALL_DIR}/stuntdouble"
else
    echo "   ${INSTALL_DIR} is not writable; escalating with sudo."
    need sudo
    sudo install -m 0755 "${TMPDIR_SD}/${BINARY_NAME}" "${INSTALL_DIR}/stuntdouble" \
        || die "failed to install to ${INSTALL_DIR}. Set INSTALL_DIR to a writable location."
fi

echo ""
echo "✅ StuntDouble ${TAG} installed to ${INSTALL_DIR}/stuntdouble"
echo "Run 'stuntdouble --help' to get started."
