#!/bin/bash

# Build script for Linux (amd64 and arm64)
# Builds ClaudeCompanion for Linux using Docker (for CGO support)

set -e

echo "======================================"
echo "Building ClaudeCompanion for Linux"
echo "======================================"

# Get script directory and project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Create output directory
mkdir -p dist/linux-amd64 dist/linux-arm64

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo ""
    echo "❌ Error: Docker is not installed or not in PATH"
    echo ""
    echo "Docker is required for building Linux binaries with CGO support."
    echo "Please install Docker from: https://www.docker.com/"
    echo ""
    echo "Alternative: Build directly on a Linux machine with:"
    echo "  sudo apt-get install -y libayatana-appindicator3-dev libgtk-3-dev pkg-config"
    echo "  CGO_ENABLED=1 go build -ldflags \"-s -w\" -o dist/claudecompanion ./cmd/claudecompanion"
    exit 1
fi

echo ""
echo "Using Docker to build with CGO support..."
echo ""

# Build for amd64
echo "======================================"
echo "Building for Linux amd64..."
echo "======================================"

docker run --rm \
    -v "$PROJECT_ROOT":/workspace \
    -w /workspace \
    golang:1.21-bookworm \
    bash -c "
        set -e
        echo '📦 Installing system dependencies...'
        apt-get update -qq > /dev/null 2>&1
        apt-get install -y -qq libayatana-appindicator3-dev libgtk-3-dev pkg-config > /dev/null 2>&1

        echo '📥 Downloading Go dependencies...'
        go mod download > /dev/null 2>&1

        echo '🔨 Building for amd64...'
        CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
            go build -ldflags \"-s -w\" \
            -o dist/linux-amd64/claudecompanion \
            ./cmd/claudecompanion

        echo '✅ amd64 build completed!'
    "

# Build for arm64
echo ""
echo "======================================"
echo "Building for Linux arm64..."
echo "======================================"

docker run --rm \
    -v "$PROJECT_ROOT":/workspace \
    -w /workspace \
    golang:1.21-bookworm \
    bash -c "
        set -e
        echo '📦 Installing system dependencies...'
        dpkg --add-architecture arm64 > /dev/null 2>&1 || true
        apt-get update -qq > /dev/null 2>&1
        apt-get install -y -qq \
            gcc-aarch64-linux-gnu \
            libayatana-appindicator3-dev:arm64 \
            libgtk-3-dev:arm64 \
            pkg-config > /dev/null 2>&1 || {
            echo '⚠️  Warning: Could not install arm64 cross-compilation tools'
            echo 'Skipping arm64 build...'
            exit 0
        }

        echo '📥 Downloading Go dependencies...'
        go mod download > /dev/null 2>&1

        echo '🔨 Building for arm64...'
        CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
            CC=aarch64-linux-gnu-gcc \
            go build -ldflags \"-s -w\" \
            -o dist/linux-arm64/claudecompanion \
            ./cmd/claudecompanion || {
            echo '⚠️  Warning: arm64 build failed'
            exit 0
        }

        echo '✅ arm64 build completed!'
    "

# Create symlink to amd64 as default
if [ -f "dist/linux-amd64/claudecompanion" ]; then
    ln -sf linux-amd64/claudecompanion dist/claudecompanion-linux
fi

echo ""
echo "======================================"
echo "✅ Build completed successfully!"
echo "======================================"
echo ""
echo "Output files:"
if [ -f "dist/linux-amd64/claudecompanion" ]; then
    echo "  ✓ dist/linux-amd64/claudecompanion ($(du -h dist/linux-amd64/claudecompanion | cut -f1))"
fi
if [ -f "dist/linux-arm64/claudecompanion" ]; then
    echo "  ✓ dist/linux-arm64/claudecompanion ($(du -h dist/linux-arm64/claudecompanion | cut -f1))"
fi
echo ""
echo "To create a distribution package:"
echo "  1. Copy dist/linux-amd64/claudecompanion to target Linux system"
echo "  2. Copy config.yaml.example as config.yaml"
echo "  3. Install Firefox extension"
echo ""
echo "System requirements on target Linux:"
echo "  - libayatana-appindicator3-1 (or libappindicator3-1)"
echo "  - libgtk-3-0"
echo "  - curl"
echo "  - libnotify-bin"
echo ""
echo "Install on Ubuntu/Debian:"
echo "  sudo apt-get install libayatana-appindicator3-1 libgtk-3-0 curl libnotify-bin"
echo ""
echo "Config file location:"
echo "  ~/.config/claudecompanion/config.yaml"
echo ""
