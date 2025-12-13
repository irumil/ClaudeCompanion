#!/bin/bash

# Build script for Linux

echo "Building ClaudeCompanion for Linux..."

# Change to project root
cd "$(dirname "$0")/.."

# Set variables
OUTPUT_DIR="dist"
APP_NAME="claudecompanion"
MAIN_PACKAGE="./cmd/claudecompanion"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Download dependencies
echo "Downloading dependencies..."
go mod tidy
go mod download

# Build for Linux (amd64)
echo "Compiling for Linux (amd64)..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o "$OUTPUT_DIR/${APP_NAME}-amd64" "$MAIN_PACKAGE"

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "Build successful (amd64)!"
    echo "Output: $OUTPUT_DIR/${APP_NAME}-amd64"
    echo "========================================"
fi

# Build for Linux (arm64)
echo ""
echo "Compiling for Linux (arm64)..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o "$OUTPUT_DIR/${APP_NAME}-arm64" "$MAIN_PACKAGE"

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "Build successful (arm64)!"
    echo "Output: $OUTPUT_DIR/${APP_NAME}-arm64"
    echo "========================================"
fi

# Create symlink to amd64 as default
if [ -f "$OUTPUT_DIR/${APP_NAME}-amd64" ]; then
    ln -sf "${APP_NAME}-amd64" "$OUTPUT_DIR/$APP_NAME"
    echo ""
    echo "Default binary (amd64) linked to: $OUTPUT_DIR/$APP_NAME"
fi

echo ""
echo "========================================"
echo "Build complete!"
echo "========================================"
echo ""
echo "To run the application:"
echo "  ./$OUTPUT_DIR/$APP_NAME"
echo ""
echo "Config file will be created at:"
echo "  ~/.config/claudecompanion/config.yaml"
echo ""
echo "Note: On Linux, you may need to install dependencies for system tray:"
echo "  Ubuntu/Debian: sudo apt-get install libayatana-appindicator3-dev"
echo "  Fedora: sudo dnf install libappindicator-gtk3-devel"
