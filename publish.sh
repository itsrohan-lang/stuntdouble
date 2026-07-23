#!/bin/bash
set -e

echo "🚀 Initiating StuntDouble Marketplace Release..."

# 1. Publish NPM Wrapper
echo ">> Publishing stuntdouble-sandbox-cli to NPM registry..."
cd npm
npm publish --access public
cd ..

# 2. Publish VS Code Extension
echo ">> Publishing stuntdouble-vscode to Visual Studio Marketplace..."
cd vscode-extension
# Ensure you are logged in via `vsce login` before running this
npm run compile
npx vsce publish
cd ..

echo "✅ All packages successfully deployed to the developer ecosystem!"
