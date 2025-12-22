#!/bin/bash

# Build script for macOS - creates .app bundle ready to use

set -e

echo "Building ClaudeCompanion for macOS..."

# Change to project root
cd "$(dirname "$0")/.."

# Set variables
OUTPUT_DIR="dist"
APP_NAME="ClaudeCompanion"
MAIN_PACKAGE="./cmd/claudecompanion"

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf "$OUTPUT_DIR/macos-intel" "$OUTPUT_DIR/macos-apple-silicon"
rm -rf "$OUTPUT_DIR/ClaudeCompanion.app"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Copy icon for go:embed
echo "Copying icon..."
cp icon.ico cmd/claudecompanion/icon.ico

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Detect current architecture
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    echo "Detected Apple Silicon (M1/M2/M3)"
    BUILD_ARCH="arm64"
    BUILD_NAME="apple-silicon"
else
    echo "Detected Intel Mac"
    BUILD_ARCH="amd64"
    BUILD_NAME="intel"
fi

# Build for current architecture
echo ""
echo "Compiling for macOS ($BUILD_ARCH)..."
CGO_ENABLED=1 GOARCH=$BUILD_ARCH go build -ldflags="-s -w" -o "$OUTPUT_DIR/temp-build/$APP_NAME" "$MAIN_PACKAGE"

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "========================================"
echo "Build successful!"
echo "========================================"

# Package into .app bundle
echo ""
echo "Creating .app bundle..."
chmod +x build/package-macos-app.sh
./build/package-macos-app.sh "$OUTPUT_DIR/temp-build/$APP_NAME" "$OUTPUT_DIR" "$BUILD_NAME"

# Clean up temp files
rm -rf "$OUTPUT_DIR/temp-build"
rm -f cmd/claudecompanion/icon.ico

# Move .app to dist root for easy access
if [ -d "$OUTPUT_DIR/$APP_NAME.app" ]; then
    echo "Cleaning old .app in dist root..."
    rm -rf "$OUTPUT_DIR/$APP_NAME.app"
fi

echo "Moving .app bundle to dist root..."
mv "$OUTPUT_DIR/$APP_NAME.app" "$OUTPUT_DIR/"

# Copy config example
if [ ! -f "$OUTPUT_DIR/config.yaml" ]; then
    echo "Copying config example..."
    cp config.yaml.example "$OUTPUT_DIR/config.yaml"
fi

echo ""
echo "========================================"
echo "✅ Build complete!"
echo "========================================"
echo ""
echo "App bundle: $OUTPUT_DIR/$APP_NAME.app"
echo ""
echo "To install:"
echo "  1. Option A (Simple): Double-click ClaudeCompanion.app"
echo "  2. Option B (Recommended): Run the installer:"
echo "     cd $OUTPUT_DIR"
echo "     ./install.sh"
echo ""
echo "To run from terminal:"
echo "  open $OUTPUT_DIR/$APP_NAME.app"
echo ""
echo "Config location (after first run):"
echo "  ~/Library/Application Support/ClaudeCompanion/config.yaml"
echo ""
echo "To open config:"
echo "  open -a TextEdit ~/Library/Application\\ Support/ClaudeCompanion/config.yaml"
echo ""
