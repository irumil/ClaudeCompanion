#!/bin/bash

# Script to package macOS .app bundle from compiled binary
# Usage: ./package-macos-app.sh <binary-path> <output-dir> <arch>
# Example: ./package-macos-app.sh dist/darwin-amd64/claudecompanion dist/macos-app intel

set -e

BINARY_PATH="$1"
OUTPUT_DIR="$2"
ARCH="$3"  # intel or apple-silicon

if [ -z "$BINARY_PATH" ] || [ -z "$OUTPUT_DIR" ] || [ -z "$ARCH" ]; then
    echo "Usage: $0 <binary-path> <output-dir> <arch>"
    echo "Example: $0 dist/darwin-amd64/claudecompanion dist/macos-app intel"
    exit 1
fi

APP_NAME="ClaudeCompanion"
BUNDLE_ID="com.github.claudecompanion"
VERSION="1.0.0"

# Create .app bundle structure
APP_DIR="$OUTPUT_DIR/${APP_NAME}.app"
CONTENTS_DIR="$APP_DIR/Contents"
MACOS_DIR="$CONTENTS_DIR/MacOS"
RESOURCES_DIR="$CONTENTS_DIR/Resources"

echo "Creating .app bundle structure..."
mkdir -p "$MACOS_DIR"
mkdir -p "$RESOURCES_DIR"

# Copy binary
echo "Copying binary..."
cp "$BINARY_PATH" "$MACOS_DIR/$APP_NAME"
chmod +x "$MACOS_DIR/$APP_NAME"

# Create Info.plist
echo "Creating Info.plist..."
cat > "$CONTENTS_DIR/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>en</string>
    <key>CFBundleExecutable</key>
    <string>$APP_NAME</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

# Create .icns icon from PNG
echo "Creating application icon..."
ICONSET_DIR="$OUTPUT_DIR/AppIcon.iconset"
mkdir -p "$ICONSET_DIR"

# Check if we have icon96.png
if [ -f "extension/icon96.png" ]; then
    # Use sips to create different sizes
    sips -z 16 16     extension/icon96.png --out "$ICONSET_DIR/icon_16x16.png" > /dev/null 2>&1
    sips -z 32 32     extension/icon96.png --out "$ICONSET_DIR/icon_16x16@2x.png" > /dev/null 2>&1
    sips -z 32 32     extension/icon96.png --out "$ICONSET_DIR/icon_32x32.png" > /dev/null 2>&1
    sips -z 64 64     extension/icon96.png --out "$ICONSET_DIR/icon_32x32@2x.png" > /dev/null 2>&1
    sips -z 128 128   extension/icon96.png --out "$ICONSET_DIR/icon_128x128.png" > /dev/null 2>&1
    sips -z 256 256   extension/icon96.png --out "$ICONSET_DIR/icon_128x128@2x.png" > /dev/null 2>&1
    sips -z 256 256   extension/icon96.png --out "$ICONSET_DIR/icon_256x256.png" > /dev/null 2>&1
    sips -z 512 512   extension/icon96.png --out "$ICONSET_DIR/icon_256x256@2x.png" > /dev/null 2>&1
    sips -z 512 512   extension/icon96.png --out "$ICONSET_DIR/icon_512x512.png" > /dev/null 2>&1
    cp extension/icon96.png "$ICONSET_DIR/icon_512x512@2x.png"

    # Convert iconset to icns
    iconutil -c icns "$ICONSET_DIR" -o "$RESOURCES_DIR/AppIcon.icns"
    rm -rf "$ICONSET_DIR"
    echo "Icon created successfully"
else
    echo "Warning: extension/icon96.png not found, skipping icon creation"
fi

# Copy additional resources
if [ -f "extension/icon96.png" ]; then
    cp extension/icon96.png "$RESOURCES_DIR/app-icon.png"
fi

# Create install script
cat > "$OUTPUT_DIR/install.sh" << 'INSTALLEOF'
#!/bin/bash

echo "ClaudeCompanion Installer"
echo "=========================="
echo ""

APP_PATH="$(dirname "$0")/ClaudeCompanion.app"
TARGET_PATH="/Applications/ClaudeCompanion.app"

# Check if app exists
if [ ! -d "$APP_PATH" ]; then
    echo "❌ Error: ClaudeCompanion.app not found in current directory"
    exit 1
fi

# Fix permissions
echo "🔧 Fixing permissions..."
chmod +x "$APP_PATH/Contents/MacOS/ClaudeCompanion"

# Copy to Applications
echo "📦 Installing to /Applications..."
if [ -d "$TARGET_PATH" ]; then
    echo "   Removing old version..."
    rm -rf "$TARGET_PATH"
fi

cp -R "$APP_PATH" "$TARGET_PATH"
chmod +x "$TARGET_PATH/Contents/MacOS/ClaudeCompanion"

# Remove quarantine
echo "🔓 Removing quarantine attribute..."
xattr -d com.apple.quarantine "$TARGET_PATH" 2>/dev/null || echo "   (no quarantine attribute found)"

echo ""
echo "✅ Installation complete!"
echo ""
echo "To start ClaudeCompanion:"
echo "  1. Open Spotlight (Cmd+Space)"
echo "  2. Type 'ClaudeCompanion'"
echo "  3. Press Enter"
echo ""
echo "Or run from terminal:"
echo "  open /Applications/ClaudeCompanion.app"
echo ""
echo "Config location:"
echo "  ~/Library/Application Support/ClaudeCompanion/config.yaml"
echo ""
INSTALLEOF

chmod +x "$OUTPUT_DIR/install.sh"

# Create README with instructions
cat > "$OUTPUT_DIR/README.txt" << 'EOF'
ClaudeCompanion для macOS
=========================

БЫСТРАЯ УСТАНОВКА (рекомендуется):
1. Откройте Терминал в этой папке
2. Запустите: ./install.sh
3. Готово! Приложение установлено в /Applications

РУЧНАЯ УСТАНОВКА:
1. Переместите ClaudeCompanion.app в папку /Applications
2. ВАЖНО: Установите права на выполнение (в Терминале):
   chmod +x /Applications/ClaudeCompanion.app/Contents/MacOS/ClaudeCompanion

3. При первом запуске macOS заблокирует приложение (неподписанное)

ОБХОД GATEKEEPER (выберите один способ):

Способ 1 (через Системные настройки):
1. Попробуйте открыть приложение двойным кликом
2. Откройте Системные настройки → Конфиденциальность и безопасность
3. Нажмите "Все равно открыть" рядом с ClaudeCompanion

Способ 2 (через контекстное меню):
1. Удерживайте Control и кликните на ClaudeCompanion.app
2. Выберите "Открыть" в меню
3. Подтвердите открытие в диалоговом окне

Способ 3 (через терминал - рекомендуется):
cd /Applications
chmod +x ClaudeCompanion.app/Contents/MacOS/ClaudeCompanion
xattr -d com.apple.quarantine ClaudeCompanion.app
open ClaudeCompanion.app

ЕСЛИ ПОЛУЧАЕТЕ "permission denied":
chmod +x /Applications/ClaudeCompanion.app/Contents/MacOS/ClaudeCompanion

РАСПОЛОЖЕНИЕ КОНФИГА:
~/Library/Application Support/ClaudeCompanion/config.yaml

ЗАПУСК ИЗ ТЕРМИНАЛА (для отладки):
/Applications/ClaudeCompanion.app/Contents/MacOS/ClaudeCompanion

Приложение запустится в трее (menu bar).
EOF

echo ""
echo "✅ macOS .app bundle created successfully!"
echo "   Location: $APP_DIR"
echo "   Architecture: $ARCH"
echo ""
echo "To test locally:"
echo "  open \"$APP_DIR\""
echo ""
echo "To remove quarantine attribute:"
echo "  xattr -d com.apple.quarantine \"$APP_DIR\""
