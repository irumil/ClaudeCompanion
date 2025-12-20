# ClaudeCompanion

🇬🇧 English | 🇷🇺 [Русский](README.ru.md) | ❤️ [Татарча](README.tt.md)

---

System tray application for monitoring Claude.ai API usage quota in real-time.

![Architecture](schema/architecture-simple.txt)

## Features

- 🎯 **Real-time monitoring** - Updates every 60 seconds
- 🎨 **Dynamic tray icon** - Shows remaining quota percentage with color coding
- 🔔 **Smart notifications** - Alerts when quota is low or exhausted
- ☀️ **Morning Greeting to Claude** - Automatic scheduled messages to optimize 5-hour limit
- 🌐 **Browser integration** - Firefox extension for seamless authentication
- ⚙️ **Hot-reload config** - No restart needed for configuration changes
- 🔒 **Proxy support** - Works with corporate proxies
- 📊 **Detailed tooltips** - Shows 5-hour and 7-day quota information
- 🍎 **Cross-platform** - Supports Windows, macOS (Intel & Apple Silicon), and Linux

## Color Coding

- 🟢 **Green** (> 40%) - Plenty of quota remaining
- 🟡 **Yellow** (20-40%) - Moderate usage
- 🔴 **Red** (< 20%) - Low quota
- ⚪ **Gray** - Connection error

## ⚠️ Disclaimer and Risks

### Important Information Before Using

**ClaudeCompanion uses unofficial Claude.ai API.** This means:

1. **Potential risk of account suspension**
   - The app accesses internal API not intended for public use
   - Claude.ai may consider this a violation of Terms of Service
   - While risk is low with reasonable usage, we cannot fully exclude it

2. **Automation may violate ToS**
   - Automatic requests every 30-120 seconds
   - Automatic message sending (greeting feature)
   - May be classified as bot activity

### ✅ Recommendations for Safe Usage

1. **Increase polling interval**
   ```yaml
   poll_interval_seconds: 60  # Minimum 60 seconds
   # or for better safety:
   poll_interval_seconds: 120  # 2 minutes
   ```

2. **Use greeting feature carefully**
   ```yaml
   greeting:
     greeting_cron: "0 8 * * *"  # ✅ Once a day - safe
     # greeting_cron: "*/5 * * * *"  # ❌ Every 5 minutes - DANGEROUS
   ```

3. **Personal use only**
   - One account only
   - Don't use for mass automation
   - Don't run on multiple accounts simultaneously

4. **Don't run 24/7**
   - Run only during work hours
   - Stop overnight and on weekends

### 📊 Risk Assessment

**LOW risk** if:
- ✅ Polling interval ≥ 60 seconds
- ✅ Greeting maximum 1-2 times per day
- ✅ Personal use on single account
- ✅ Running only during work hours

**HIGH risk** if:
- ❌ Frequent requests every 5-10 seconds
- ❌ Frequent automatic messages
- ❌ Usage on multiple accounts
- ❌ Running 24/7 non-stop

### 🎯 Disclaimer

**Use at your own risk.** Developers are not responsible for possible account suspension or other consequences of using this application. The official and safe way to work with Claude.ai is using only the web interface at claude.ai.

If Claude.ai is critical for your work, consider using only the official web interface.

## Installation

### 1. Download

Download the latest release from [Releases](../../releases) or build from source.

### 2. Configure

Copy `config.yaml.example` to `config.yaml` and adjust settings:

```yaml
proxy: "http://your-proxy:port"  # Optional: your proxy server
browser_path: "C:\\Program Files\\Mozilla Firefox\\firefox.exe"  # Optional: custom browser
enable_file_logging: false  # true = log to file, false = console only
```

### 3. Install Browser Extension

1. Open Firefox
2. Go to `about:debugging#/runtime/this-firefox`
3. Click "Load Temporary Add-on"
4. Select `extension/manifest.json`

### 4. Run

Double-click `claudecompanion.exe` - it will start in system tray (no console window).

## Usage

1. **Start the application** - Run `claudecompanion.exe`
2. **Open Claude.ai** in Firefox - The extension will automatically send authentication
3. **Check the tray icon** - Shows remaining quota percentage
4. **Right-click the icon** for menu:
   - "Открыть Claude.ai" - Open Claude.ai in browser
   - "Открыть настройки" - Edit configuration
   - "Выход" - Exit application

## Configuration

All settings are in `config.yaml`:

### Basic Settings

```yaml
server_port: 8383              # Port for browser extension connection
poll_interval_seconds: 60      # How often to check quota
proxy: ""                      # HTTP proxy (leave empty if not needed)
browser_path: ""               # Custom browser path (leave empty for default)
enable_file_logging: false     # Enable file logging
```

### Notification Settings

```yaml
low_value_notifications:
  enabled: true
  threshold: 20                # Notify when quota <= 20%
  phrases:                     # Random phrases for low quota
    - "Time to go home! 🏡"
  zero_phrases:                # Random phrases for zero quota
    - "Game over! 🎮"
```

### Demo Mode

For testing all features and notifications:

```yaml
demo_mode:
  enabled: true
  duration_seconds: 60         # Simulates 100% → 0% in 60 seconds
```

Demo mode demonstrates:
- Icon changes from 100% to 0% (green → yellow → red)
- "Morning Greeting to Claude" notification at the start of each cycle ☀️
- Low quota notifications (when threshold is reached)
- Zero quota notifications

### Morning Greeting to Claude

Automatically send messages to a chat on schedule to optimize 5-hour limit:

```yaml
greeting:
  greeting_cron: "0 8 * * *"   # Cron schedule: 8 AM every day
  greeting_text: "Ok"          # Message text
  greeting_chat_id: ""         # Chat UUID (required for greeting to work)
```

**How to get chat UUID:**
1. Open the desired chat on claude.ai
2. UUID is in the URL: `https://claude.ai/chat/{UUID}`
3. Copy the UUID to `greeting_chat_id` setting

**Cron schedule examples:**
- `"0 8 * * *"` - every day at 8:00 AM
- `"0 9 * * 1-5"` - at 9:00 AM on weekdays (Mon-Fri)
- `"30 7 * * *"` - every day at 7:30 AM
- `"*/5 * * * *"` - every 5 minutes (for testing)

When greeting is sent, you'll see "Morning Greeting to Claude" notification ☀️

### Work Hours

Limit API polling to specific time ranges (e.g., only during work hours):

```yaml
work_hours:
  enabled: true               # Enable to limit polling to work hours only
  start: "08:00"              # Start time (HH:MM format)
  end: "20:00"                # End time (HH:MM format)
```

**How it works:**
- When `enabled: true`, API polling only happens between `start` and `end` times
- Uses 24-hour format (HH:MM)
- Supports overnight ranges (e.g., `start: "20:00"`, `end: "08:00"` for night shift)
- Outside work hours, polling is skipped (no API requests made)
- Helps reduce unnecessary API usage and potential detection risk

**Examples:**
- `start: "08:00"`, `end: "20:00"` - typical work day (8 AM to 8 PM)
- `start: "09:00"`, `end: "17:00"` - standard office hours (9 AM to 5 PM)
- `start: "20:00"`, `end: "08:00"` - overnight shift

## Architecture

The application consists of two parts:

### 1. Browser Extension (Firefox)
- Extracts `sessionKey` cookie from Claude.ai
- Fetches organization UUID
- Sends authentication data and organizationId to desktop app

### 2. Desktop Application (Go)
- **HTTP Server** - Receives data from browser extension
- **API Client** - Polls Claude.ai API every 60 seconds
- **Cron Scheduler** - Sends greeting messages on schedule
- **Tray Manager** - Shows dynamic icon with percentage
- **Icon Generator** - Creates 48x48 icons with colored numbers
- **Notifier** - Windows Toast notifications
- **Config Manager** - Hot-reload configuration changes
- **Logger** - Optional file logging

See [schema/architecture.md](schema/architecture.md) for detailed architecture diagram.

## Building from Source

### Prerequisites

- Go 1.21 or higher
- Windows (for Windows build), macOS (for macOS build), or Linux (for Linux build)
- Git (optional, for cloning the repository)
- **macOS only**: Homebrew curl (required, see installation instructions below)
- **Linux only**: System libraries (libayatana-appindicator3-dev, libgtk-3-dev, curl, libnotify-bin)

### Building for macOS

The application supports macOS (Intel and Apple Silicon) with native notifications and system tray integration.

```bash
# 1. Clone repository (or download ZIP)
git clone https://github.com/irumil/ClaudeCompanion.git
cd ClaudeCompanion

# 2. Run build script
cd build
./build-macos.sh
```

The script will:
- Download dependencies
- Build for Intel (amd64)
- Build for Apple Silicon (arm64)
- Create universal binary that works on both architectures

**Output files:**
- `dist/claudecompanion` - Universal binary

**Running the app:**
```bash
# From project root
./dist/claudecompanion

# Config will be created at:
# ~/Library/Application Support/ClaudeCompanion/config.yaml
```

**Notes for macOS:**
- Browser extension works in Firefox (same setup as Windows)
- Notifications use native macOS notification center
- System tray icon appears in menu bar
- CGO is required (enabled automatically by build script)
- **Requires Homebrew curl** - System curl may be blocked by Cloudflare. Install with:
  ```bash
  brew install curl
  ```
  The app uses `/opt/homebrew/opt/curl/bin/curl` by default. To use a different curl binary, set `curl_path` in `config.yaml`:
  ```yaml
  curl_path: "/path/to/your/curl"
  ```

### Building for Linux

The application supports Linux with native notifications and system tray integration.

**Prerequisites:**
```bash
# Install required system libraries
# For Ubuntu/Debian:
sudo apt-get install -y libayatana-appindicator3-dev libgtk-3-dev curl libnotify-bin

# For Fedora/RHEL:
sudo dnf install -y libappindicator-gtk3-devel gtk3-devel curl libnotify

# For Arch Linux:
sudo pacman -S libappindicator-gtk3 gtk3 curl libnotify
```

**Building:**
```bash
# 1. Clone repository (or download ZIP)
git clone https://github.com/irumil/ClaudeCompanion.git
cd ClaudeCompanion

# 2. Install Go dependencies
go mod download

# 3. Build for Linux
CGO_ENABLED=1 go build -ldflags "-s -w" -o dist/claudecompanion ./cmd/claudecompanion

# 4. Copy config
cp config.yaml.example dist/config.yaml

# 5. Run the app
./dist/claudecompanion

# Config will be created at:
# ~/.config/claudecompanion/config.yaml
```

**Notes for Linux:**
- Browser extension works in Firefox (same setup as Windows)
- Notifications use `notify-send` (libnotify)
- System tray icon requires AppIndicator library
- CGO is required (set `CGO_ENABLED=1`)
- Uses system `curl` by default (`/usr/bin/curl`)
- Desktop environment with system tray support required (GNOME, KDE, XFCE, etc.)

### Full Build (Application + Extension) for Windows

#### Option 1: Automated Build (Recommended)

The easiest way - use the build script:

```bash
# Build everything (application + extension)
cd build
build-all.bat

# Or just debug version with console
build-debug.bat
```

#### Option 2: Quick Manual Build

```bash
# 1. Clone repository (or download ZIP)
git clone https://github.com/irumil/ClaudeCompanion.git
cd ClaudeCompanion

# 2. Install dependencies
go mod download

# 3. Build application (release version without console)
go build -ldflags "-H windowsgui" -o dist/claudecompanion.exe ./cmd/claudecompanion

# 4. Copy required files
copy config.yaml.example dist\config.yaml

# 5. Build browser extension
cd build
package-extension.bat
cd ..
```

Done! Application is in `dist/claudecompanion.exe`, extension in `dist/claudecompanion-extension.zip`

#### Option 3: Build with Embedded Icon

To embed the icon in the exe file (for better notifications):

```bash
# 1-2. Same steps as above

# 3. Install resource compiler (once)
go install github.com/akavel/rsrc@latest

# 4. Generate Windows resources
cd cmd/claudecompanion
rsrc -ico "..\..\icon.ico" -o rsrc_windows_amd64.syso
cd ..\..

# 5. Build with embedded icon
go build -ldflags "-H windowsgui" -o dist/claudecompanion.exe ./cmd/claudecompanion

# 6. Copy files
copy config.yaml.example dist\config.yaml

# 7. Build extension
cd build
package-extension.bat
```

### Debug Build (with Console)

For development and debugging (with visible console):

```bash
# Build with console for viewing logs
go build -o dist/claudecompanion-debug.exe ./cmd/claudecompanion

# Run in console
dist\claudecompanion-debug.exe
```

### Extension Only

To rebuild only the browser extension:

```bash
cd build
package-extension.bat
```

Or manually:

```bash
cd extension
powershell -Command "Compress-Archive -Path manifest.json,background.js,options.html,options.js,icon48.png,icon96.png -DestinationPath ../dist/claudecompanion-extension.zip -Force"
```

### Verify Build

After building, check:

```bash
# Check files exist
dir dist

# Should have:
# claudecompanion.exe (release) or claudecompanion-debug.exe (debug)
# claudecompanion-extension.zip
# config.yaml
```

## Project Structure

```
ClaudeCompanion/
├── cmd/
│   └── claudecompanion/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # Claude.ai API client
│   ├── config/                  # Configuration management
│   ├── icon/                    # Dynamic icon generator
│   ├── logger/                  # Logging system
│   ├── notifier/                # Toast notifications
│   ├── server/                  # HTTP server for extension
│   └── tray/                    # System tray manager
├── extension/
│   ├── manifest.json            # Firefox extension manifest
│   ├── background.js            # Extension logic
│   └── icon.png                 # Extension icon
├── dist/
│   ├── claudecompanion.exe      # Built executable
│   └── config.yaml              # User configuration (not in git)
├── schema/
│   ├── architecture.md          # Detailed architecture
│   └── architecture-simple.txt  # ASCII diagram
├── go.mod
├── go.sum
├── icon.ico                     # Application icon
├── config.yaml.example          # Example configuration
└── README.md
```

## Troubleshooting

### Extension not working
- Check that Firefox is running
- Reload the extension in `about:debugging`
- Open browser console for errors

### No icon updates
- Check that you visited Claude.ai in Firefox
- Verify extension sent data (check logs if enabled)
- Try restarting the application

### Proxy errors
- Verify proxy address in config.yaml
- Test proxy with curl manually
- Check proxy authentication if required

### Notifications not showing
- Enable notifications in Windows settings
- Check notification threshold settings
- Try demo mode to test notifications

## Dependencies

### Go Libraries
- [github.com/getlantern/systray](https://github.com/getlantern/systray) - System tray (cross-platform)
- [github.com/go-toast/toast](https://github.com/go-toast/toast) - Toast notifications (Windows)
- [github.com/gen2brain/beeep](https://github.com/gen2brain/beeep) - Notifications (macOS/Linux)
- [github.com/robfig/cron/v3](https://github.com/robfig/cron) - Cron scheduler for scheduled tasks
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing

### External Tools
- **Windows**: curl.exe - For API requests (included in Windows 10+)
- **macOS**: Homebrew curl - For API requests (install with `brew install curl`)
- **Linux**: curl - For API requests (install with package manager)
- **Linux**: notify-send - For notifications (install `libnotify-bin` or `libnotify`)
- **Windows**: notepad.exe - For opening config file

## License

[Your License Here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on GitHub.
