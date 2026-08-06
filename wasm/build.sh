#!/bin/bash
set -euo pipefail

echo "Building StuntDouble WebAssembly port..."
GOOS=js GOARCH=wasm go build -o stuntdouble.wasm main.go

# wasm_exec.js moved from misc/wasm to lib/wasm in Go 1.24. Probe both so this
# works on the 1.25 toolchain the modules require as well as older ones.
echo "Copying wasm_exec.js helper..."
GOROOT_DIR="$(go env GOROOT)"
WASM_EXEC=""
for candidate in "$GOROOT_DIR/lib/wasm/wasm_exec.js" "$GOROOT_DIR/misc/wasm/wasm_exec.js"; do
    if [ -f "$candidate" ]; then
        WASM_EXEC="$candidate"
        break
    fi
done

if [ -z "$WASM_EXEC" ]; then
    echo "error: wasm_exec.js not found under $GOROOT_DIR (looked in lib/wasm and misc/wasm)" >&2
    exit 1
fi

cp "$WASM_EXEC" .

echo "Build complete. To test in browser, serve this directory (e.g., python3 -m http.server)."
