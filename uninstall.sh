#!/usr/bin/env bash
#
# StuntDouble CLI uninstaller script.
#
# Removes installed StuntDouble binaries (/usr/local/bin/stuntdouble and /usr/local/bin/sd)
# and cleans up global npm installations.
#
# Options:
#   --purge     Remove local workspace .stuntdouble configuration and telemetry files.
#
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
PURGE=false

for arg in "$@"; do
    case "$arg" in
        --purge) PURGE=true ;;
        *) ;;
    esac
done

echo "🗑️ Uninstalling StuntDouble CLI..."

# 1. Remove binary files from INSTALL_DIR
for binary in "stuntdouble" "sd"; do
    target="${INSTALL_DIR}/${binary}"
    if [ -f "$target" ] || [ -L "$target" ]; then
        echo ">> Removing ${target}..."
        if [ -w "$INSTALL_DIR" ] || [ -w "$target" ]; then
            rm -f "$target"
        else
            echo "   ${target} requires superuser permissions; escalating with sudo."
            sudo rm -f "$target"
        fi
    fi
done

# 2. Check npm global package removal
if command -v npm >/dev/null 2>&1; then
    if npm list -g stuntdouble-sandbox-cli >/dev/null 2>&1; then
        echo ">> Uninstalling global npm package (stuntdouble-sandbox-cli)..."
        npm uninstall -g stuntdouble-sandbox-cli || true
    fi
fi

# 3. Purge workspace files if requested
if [ "$PURGE" = true ]; then
    echo ">> Purging local .stuntdouble configuration directory..."
    rm -rf .stuntdouble .stuntdouble.yaml .stuntdouble.telemetry.json
fi

echo ""
echo "✅ StuntDouble has been successfully uninstalled from your system."
