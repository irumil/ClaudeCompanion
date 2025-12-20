#!/bin/bash

# Build script for macOS

echo "Building ClaudeCompanion for macOS..."

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

# Build for macOS (Intel)
echo "Compiling for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o "$OUTPUT_DIR/${APP_NAME}-amd64" "$MAIN_PACKAGE"

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "Build successful (Intel)!"
    echo "Output: $OUTPUT_DIR/${APP_NAME}-amd64"
    echo "========================================"
fi

# Build for macOS (Apple Silicon)
echo ""
echo "Compiling for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -ldflags="-s -w" -o "$OUTPUT_DIR/${APP_NAME}-arm64" "$MAIN_PACKAGE"

if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "Build successful (Apple Silicon)!"
    echo "Output: $OUTPUT_DIR/${APP_NAME}-arm64"
    echo "========================================"
fi

# Create universal binary
if [ -f "$OUTPUT_DIR/${APP_NAME}-amd64" ] && [ -f "$OUTPUT_DIR/${APP_NAME}-arm64" ]; then
    echo ""
    echo "Creating universal binary..."
    lipo -create -output "$OUTPUT_DIR/$APP_NAME" "$OUTPUT_DIR/${APP_NAME}-amd64" "$OUTPUT_DIR/${APP_NAME}-arm64"

    if [ $? -eq 0 ]; then
        echo ""
        echo "========================================"
        echo "Universal binary created!"
        echo "Output: $OUTPUT_DIR/$APP_NAME"
        echo "========================================"
        echo ""
        echo "To run the application:"
        echo "  ./$OUTPUT_DIR/$APP_NAME"
        echo ""
        echo "Config file will be created at:"
        echo "  ~/Library/Application Support/ClaudeCompanion/config.yaml"

        # Clean up architecture-specific binaries
        rm "$OUTPUT_DIR/${APP_NAME}-amd64" "$OUTPUT_DIR/${APP_NAME}-arm64"

        # Copy icon for notifications
        if [ -f "extension/icon96.png" ]; then
            cp extension/icon96.png "$OUTPUT_DIR/app-icon.png"
            echo ""
            echo "App icon copied to $OUTPUT_DIR/app-icon.png"
        fi
    fi
fi

echo ""
echo "Build complete!"
