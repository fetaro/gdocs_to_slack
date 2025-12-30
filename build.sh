#!/bin/bash
mkdir -p dist

# Clean up
rm -rf dist/*

NAME="gdocs_to_slack"

# Build for arm64
mkdir -p dist/arm64
echo "Building for darwin/arm64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o dist/arm64/$NAME
chmod +x dist/arm64/$NAME

# Build for x86_64 (amd64)
mkdir -p dist/x86_64
echo "Building for darwin/amd64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o dist/x86_64/$NAME
chmod +x dist/x86_64/$NAME

echo "Build complete."
ls -R dist/

