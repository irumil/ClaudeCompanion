# ClaudeCompanion Monitor

**[Татарча / Tatar](./README.tt.md)**

Browser extension for the ClaudeCompanion desktop application.

## What does it do?

This extension works together with the [ClaudeCompanion desktop application](https://github.com/irumil/ClaudeCompanion) to monitor your Claude.ai API usage quota in real-time.

### Features

- 🔍 **Automatic detection** - Detects when you visit claude.ai
- 🔐 **Authentication extraction** - Reads sessionKey cookie for API access
- 📊 **Quota monitoring** - Sends data to desktop app for quota tracking
- 🖥️ **Local-only** - All data stays on your computer (sent to localhost)
- 🔔 **Desktop notifications** - Get notified when quota is low

## How it works

1. You visit [claude.ai](https://claude.ai) in Firefox
2. Extension extracts your sessionKey cookie
3. Extension calls Claude API to get your organization UUID
4. Extension sends data to ClaudeCompanion desktop app (localhost:8383)
5. Desktop app monitors your API quota and displays it in system tray

## Installation

### From Mozilla Add-ons (AMO)
1. Visit the [extension page on AMO](LINK_WILL_BE_ADDED)
2. Click "Add to Firefox"
3. Install ClaudeCompanion desktop app
4. Visit claude.ai to activate

### Manual Installation (Developer)
1. Download the source code
2. Open Firefox
3. Go to `about:debugging#/runtime/this-firefox`
4. Click "Load Temporary Add-on"
5. Select `manifest.json` from the extension folder

## Requirements

- **Firefox 57+** (or Firefox Developer Edition)
- **ClaudeCompanion desktop application** running on port 8383
- **Active Claude.ai account**

## Privacy

**This extension does NOT:**
- ❌ Send data to external servers
- ❌ Track your browsing
- ❌ Collect personal information
- ❌ Store data permanently

**All collected data is sent ONLY to:**
- ✅ `http://localhost:8383` (your computer)
- ✅ ClaudeCompanion desktop app

See [PRIVACY.md](./PRIVACY.md) for detailed privacy policy.

## Configuration

Click the extension icon or go to `about:addons` → ClaudeCompanion → Preferences to configure:

- **Desktop App Port** - Port where desktop app is listening (default: 8383)
- **Auto-send** - Automatically send data when visiting claude.ai

## Permissions Explained

This extension requests the following permissions:

| Permission | Why we need it |
|------------|----------------|
| `cookies` | Read sessionKey from claude.ai for API authentication |
| `*://claude.ai/*` | Access Claude.ai cookies and API |
| `tabs` | Detect when you open claude.ai |
| `webRequest` | Monitor network requests to claude.ai |
| `storage` | Save settings (port number, auto-send preference) |
| `notifications` | Show desktop notifications (optional) |

## Troubleshooting

### Extension not working?

1. **Check desktop app is running**
   ```bash
   # Windows
   tasklist | findstr claudecompanion

   # Linux/Mac
   ps aux | grep claudecompanion
   ```

2. **Check port 8383 is listening**
   ```bash
   netstat -an | findstr :8383
   ```

3. **Check Browser Console for errors**
   - Press `Ctrl+Shift+J` (Firefox)
   - Look for errors related to ClaudeCompanion

4. **Verify you're logged into claude.ai**
   - Visit https://claude.ai
   - Make sure you're logged in

### Desktop app not receiving data?

- Check firewall settings (allow localhost connections)
- Restart both Firefox and desktop app
- Check desktop app logs

## Development

### Building from source

```bash
cd extension

# Create XPI package
zip -r ../claudecompanion.xpi manifest.json background.js options.html options.js icon48.png icon96.png

# Or use web-ext
npm install -g web-ext
web-ext build
```

### Testing locally

```bash
web-ext run --firefox-profile=default-release
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

[Add your license here - MIT recommended]

## Links

- 🏠 **Homepage:** [https://github.com/irumil/ClaudeCompanion]
- 🐛 **Bug Reports:** [https://github.com/irumil/ClaudeCompanion/issues]
- 📖 **Documentation:** [https://github.com/irumil/ClaudeCompanion/wiki]
- 💬 **Discussions:** [https://github.com/irumil/ClaudeCompanion/discussions]

## Credits

Developed by **rizmailov**

Built with ❤️ for the Claude.ai community.

---

**Note:** This extension is not affiliated with or endorsed by Anthropic PBC.
