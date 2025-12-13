# Claude.md - Project Context for AI Assistants

> This file provides context about the ClaudeCompanion project for AI assistants working on this codebase.

## Project Overview

**ClaudeCompanion** is a Windows system tray application that monitors Claude.ai API usage quota in real-time.

### Key Components

1. **Desktop Application (Go)** - System tray app that polls Claude.ai API
2. **Browser Extension (Firefox)** - Captures authentication and sends to desktop app
3. **Dynamic Icon Generation** - Creates tray icons with colored percentage numbers
4. **Toast Notifications** - Native Windows notifications for low/zero quota

### Tech Stack

- **Language**: Go 1.21+
- **Platform**: Windows (primary), cross-platform capable
- **Key Libraries**:
  - `github.com/getlantern/systray` - System tray integration
  - `github.com/go-toast/toast` - Windows Toast notifications
  - `gopkg.in/yaml.v3` - Configuration management
  - Standard library for image generation and HTTP

## Project Structure

```
ClaudeCompanion/
├── cmd/claudecompanion/
│   └── main.go                 # Application entry point, main loop
├── internal/
│   ├── api/
│   │   ├── client.go          # Claude.ai API client (uses curl)
│   │   ├── client_windows.go  # Windows-specific: hide curl window
│   │   └── client_other.go    # Unix: no-op for curl hiding
│   ├── config/
│   │   └── config.go          # YAML config with hot-reload (2s interval)
│   ├── icon/
│   │   ├── generator.go       # Dynamic 48x48 ICO generation
│   │   ├── digits.go          # Pixel-art digit definitions (6x6 blocks)
│   │   └── static.go          # Fallback static icons
│   ├── logger/
│   │   └── logger.go          # File/console logging with daily rotation
│   ├── notifier/
│   │   └── notifier.go        # Toast notifications with state tracking
│   ├── server/
│   │   └── server.go          # HTTP server (:8383) for extension
│   └── tray/
│       └── tray.go            # System tray manager, menu, click handlers
├── extension/
│   ├── manifest.json          # Firefox extension manifest
│   ├── background.js          # Extension logic: extract sessionKey, send to app
│   ├── options.html/js        # Extension settings page
│   └── icon48.png, icon96.png # Extension icons
├── build/
│   ├── create_extension_icons.go  # Icon generator script
│   └── package-extension.bat      # Extension packaging script
├── dist/
│   ├── claudecompanion.exe        # Compiled application
│   ├── config.yaml                # User configuration (gitignored)
│   └── claudecompanion-*.log      # Log files (gitignored)
└── schema/
    ├── architecture.md            # Detailed architecture with Mermaid
    └── architecture-simple.txt    # ASCII diagram
```

## Architecture Patterns

### 1. API Polling Loop

**File**: `cmd/claudecompanion/main.go` - `pollLoop()`

- Runs every N seconds (default: 30)
- Uses system `curl.exe` instead of Go HTTP client (better proxy support)
- Curl window hidden via Windows syscall flags
- Error threshold system: gray icon after 5 errors, notification after 10

### 2. Configuration Hot-Reload

**File**: `internal/config/config.go`

- Watches config file every 2 seconds
- Compares `ModTime` to detect changes
- Triggers callbacks on change (no restart needed)
- Applies defaults for missing values

### 3. Dynamic Icon Generation

**File**: `internal/icon/generator.go`

- Generates 48x48 pixel icons at runtime
- Transparent background (RGBA{0,0,0,0})
- Colored text: Green (>40%), Yellow (20-40%), Red (<20%), Gray (error)
- Uses 6x6 pixel blocks for each digit
- Converts PNG → ICO format for Windows compatibility

### 4. Notification State Machine

**File**: `internal/notifier/notifier.go`

- Tracks which notifications were shown (prevents duplicates)
- Three states: error, low value, zero value
- Resets on new context or value above threshold
- Uses `beeep.Alert` (changed from `Notify` for Windows compatibility)

### 5. Browser Extension Flow

**File**: `extension/background.js`

```
1. User visits claude.ai
2. Extract sessionKey cookie
3. Call /api/organizations to get org UUID
4. Construct usage URL: /api/organizations/{uuid}/usage
5. POST to localhost:8383/set-context with cookies + URL
6. Desktop app stores and starts polling
```

## Critical Implementation Details

### Null Handling for `resets_at`

**Problem**: API can return `"resets_at": null` when quota is 0

**Solution**:
- Use `*time.Time` (pointer) instead of `time.Time`
- Check for nil before formatting
- Display `—` when null instead of zero time

```go
if usage.FiveHour.ResetsAt != nil {
    resetTime = usage.FiveHour.ResetsAt.Local().Format("15:04:05")
} else {
    resetTime = "—"
}
```

### Windows Console Window Hiding

**Files**: `internal/api/client_windows.go`, build with `-ldflags "-H windowsgui"`

Two places where console appears:
1. **Application**: Build with `-H windowsgui` flag
2. **Curl subprocess**: Set `CREATE_NO_WINDOW` flag via syscall

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    HideWindow:    true,
    CreationFlags: 0x08000000, // CREATE_NO_WINDOW
}
```

### Multiline Tooltips on Windows

Use `\r\n` instead of `\n` for line breaks in tooltips:

```go
return fmt.Sprintf("5 часов: %.0f%%\r\nНеделя: %.0f%%", ...)
```

### Icon Embedding

**Files**: `cmd/claudecompanion/rsrc_windows_amd64.syso`

1. Use `github.com/akavel/rsrc` to generate `.syso` file
2. Place in same directory as `main.go`
3. Go automatically includes `.syso` files during build
4. Icon appears in exe and can be used for notifications

```bash
rsrc -ico "../../icon.ico" -o rsrc_windows_amd64.syso
```

## Common Development Tasks

### Building the Application

```bash
# Debug build (with console)
go build -o dist/claudecompanion.exe ./cmd/claudecompanion

# Release build (no console)
go build -ldflags "-H windowsgui" -o dist/claudecompanion.exe ./cmd/claudecompanion
```

### Packaging the Extension

```bash
# Windows
cd build
package-extension.bat

# Or manually
cd extension
powershell -Command "Compress-Archive -Path manifest.json,background.js,options.html,options.js,icon48.png,icon96.png -DestinationPath ../dist/claudecompanion-extension.zip -Force"
```

### Testing Demo Mode

Enable in `config.yaml`:
```yaml
demo_mode:
  enabled: true
  duration_seconds: 60  # Full cycle: 100% → 0%
```

- Runs every 2 seconds (not 30s poll interval)
- Infinite loop, resets after duration
- Force zero in last 3 seconds
- Tests all notification thresholds

### Adding New Configuration Options

1. Add field to `Config` struct in `internal/config/config.go`
2. Add YAML tag: `yaml:"field_name"`
3. Add to `defaultConfig` in `createDefaultConfig()`
4. Use in code via `cfgMgr.Get().FieldName`
5. Hot-reload works automatically

### Modifying Icon Colors/Sizes

**File**: `internal/icon/generator.go`

- Size: Change `iconSize` constant (currently 48)
- Block size: Change pixel block drawing logic (currently 6x6)
- Colors: Modify `getColor()` method
- Digits: Edit definitions in `digits.go`

## Known Limitations and Quirks

1. **Firefox Only**: Extension uses `browser.*` API (Firefox-specific)
   - Chrome requires `chrome.*` and manifest v3

2. **Windows Primary**: While Go code is cross-platform, some features are Windows-specific:
   - Toast notifications (go-toast)
   - Systray icon embedding
   - Curl window hiding

3. **Systray Library Limitation**: `getlantern/systray` doesn't support click events on icon
   - Workaround: Added "Открыть Claude.ai" menu item

4. **Proxy Authentication**: Currently doesn't support proxy auth prompts
   - User must configure auth in proxy URL if needed

5. **Icon Format**: Must use ICO on Windows, PNG doesn't work with systray
   - Generator creates PNG, then converts to ICO

## API Response Format

**Endpoint**: `GET https://claude.ai/api/organizations/{uuid}/usage`

**Headers**: `Cookie: sessionKey={value}`

**Response**:
```json
{
  "five_hour": {
    "utilization": 75.5,
    "resets_at": "2025-12-12T14:30:00Z"  // Can be null!
  },
  "seven_day": {
    "utilization": 42.3,
    "resets_at": "2025-12-15T00:00:00Z"
  }
}
```

**Important**:
- `utilization` is 0-100 (percentage used)
- Display shows `100 - utilization` (remaining quota)
- `resets_at` can be `null` when quota is exhausted

## Configuration Schema

```yaml
server_port: 8383                    # HTTP server port for extension
poll_interval_seconds: 30            # API polling interval
use_curl_fallback: true              # Always true (HTTP client removed)
gray_mode_threshold: 5               # Errors before gray icon
notification_threshold: 10           # Errors before notification
proxy: ""                            # Optional: "http://host:port"
enable_file_logging: false           # true = file+console, false = console
browser_path: ""                     # Optional: path to browser exe

low_value_notifications:
  enabled: true
  threshold: 20                      # Notify when ≤ 20%
  phrases: ["Phrase 1", "..."]       # Random selection
  zero_phrases: ["Phrase 1", "..."]  # When quota = 0

demo_mode:
  enabled: false
  duration_seconds: 60               # Cycle duration for testing
```

## Testing Checklist

- [ ] Build without console window (`-H windowsgui`)
- [ ] Curl doesn't show console
- [ ] Icon updates every 30 seconds
- [ ] Tooltip shows correct data
- [ ] Notifications appear for low/zero quota
- [ ] Settings menu opens config.yaml
- [ ] "Открыть Claude.ai" opens browser
- [ ] Hot-reload works (edit config, see changes in 2s)
- [ ] Extension sends cookies on claude.ai visit
- [ ] Gray icon appears after 5 API errors
- [ ] Handles `resets_at: null` correctly (shows `—`)

## Git Workflow

**Ignored files** (`.gitignore`):
- `dist/config.yaml` - User configuration
- `dist/*.exe`, `dist/*.log` - Build artifacts
- `cmd/claudecompanion/rsrc_windows_amd64.syso` - Generated resource
- `.idea/`, `.claude/` - IDE settings

**Include in commits**:
- `config.yaml.example` - Template configuration
- `schema/` - Architecture documentation
- All source code
- `README.md`, `README.ru.md` - Documentation

## Security Considerations

1. **No Hardcoded Secrets**: Proxy address was removed from code
   - Use `config.yaml` (gitignored) for personal settings

2. **Local Communication Only**: Extension → Desktop app is localhost
   - No external data transmission

3. **Session Token**: sessionKey stays on user's machine
   - Never sent to any server except Claude.ai API

4. **Proxy Support**: All requests can go through corporate proxy
   - Supports HTTP proxies via curl

## Future Improvements (Ideas)

- [ ] Chrome extension support (manifest v3)
- [ ] macOS/Linux testing and packaging
- [ ] Custom notification sounds
- [ ] Historical quota tracking/graphs
- [ ] Multiple account support
- [ ] Configurable color thresholds
- [ ] Auto-update mechanism
- [ ] Tray icon click support (requires different library)

## Resources

- **Go Documentation**: https://golang.org/doc/
- **Systray Library**: https://github.com/getlantern/systray
- **Toast Notifications**: https://github.com/go-toast/toast
- **Firefox Extensions**: https://extensionworkshop.com/
- **Claude.ai API**: (No official docs, reverse-engineered)

---

**Last Updated**: 2025-12-12
**Claude AI**: Assisted in development and created this documentation
