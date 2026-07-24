#!/bin/bash
set -e

echo "Building StuntDouble WebAssembly port..."
GOOS=js GOARCH=wasm go build -o stuntdouble.wasm main.go

echo "Copying wasm_exec.js helper..."
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

echo "Build complete. To test in browser, serve this directory (e.g., python3 -m http.server)."
