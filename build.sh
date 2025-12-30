#!/bin/bash
mkdir -p dist

# Clean up
rm -rf dist/*

VERSION="v1.0.0"
NAME="docs_to_slack"

# Build for arm64
DIR_ARM64="dist/${NAME}_darwin_arm64_${VERSION}"
mkdir -p $DIR_ARM64
echo "Building for darwin/arm64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o $DIR_ARM64/$NAME
chmod +x $DIR_ARM64/$NAME
tar -C dist -czf dist/${NAME}_darwin_arm64_${VERSION}.tar.gz ${NAME}_darwin_arm64_${VERSION}

# Build for amd64
DIR_AMD64="dist/${NAME}_darwin_amd64_${VERSION}"
mkdir -p $DIR_AMD64
echo "Building for darwin/amd64..."
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o $DIR_AMD64/$NAME
chmod +x $DIR_AMD64/$NAME
tar -C dist -czf dist/${NAME}_darwin_amd64_${VERSION}.tar.gz ${NAME}_darwin_amd64_${VERSION}

# Universal Binary (Optional, but good for local use)
echo "Creating universal binary..."
lipo -create -output dist/$NAME $DIR_ARM64/$NAME $DIR_AMD64/$NAME
chmod +x dist/$NAME

echo "Build complete."
ls -l dist/

