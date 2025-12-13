# Publishing ClaudeCompanion Extension to Mozilla Add-ons

## Pre-requisites

1. Create a Mozilla Add-ons account at https://addons.mozilla.org
2. Read Mozilla's [Extension Workshop](https://extensionworkshop.com/) guidelines

## Package the Extension

### Using the build script (Windows)

```bash
cd build
package-extension.bat
```

### Manual packaging

```bash
cd extension
# Create ZIP with only required files
zip -r ../dist/claudecompanion-extension.zip \
  manifest.json \
  background.js \
  options.html \
  options.js \
  icon48.png \
  icon96.png
```

Or using PowerShell:
```powershell
cd extension
Compress-Archive -Path manifest.json,background.js,options.html,options.js,icon48.png,icon96.png -DestinationPath ../dist/claudecompanion-extension.zip -Force
```

## Submit to Mozilla Add-ons

1. Go to https://addons.mozilla.org/developers/
2. Click "Submit a New Add-on"
3. Choose "On this site" (for listed add-on) or "On your own" (for unlisted)
4. Upload `dist/claudecompanion-extension.zip`
5. Fill in the required information:
   - **Name**: ClaudeCompanion Context Provider
   - **Category**: Productivity or Privacy & Security
   - **Summary**: Sends cookies and URL to ClaudeCompanion desktop application
   - **Description**: See below

### Recommended Description

```
ClaudeCompanion Context Provider is a companion extension for the ClaudeCompanion desktop application.

This extension automatically extracts authentication information from Claude.ai and sends it to the local desktop application for monitoring API usage quota.

Features:
- Seamless integration with ClaudeCompanion desktop app
- Automatic authentication handling
- No manual cookie copying required
- Privacy-focused: all data stays on your computer

Requirements:
- ClaudeCompanion desktop application must be installed and running
- Works only with claude.ai website

Source code: https://github.com/[your-username]/ClaudeCompanion
```

## Important Notes

### Permissions Explanation

When submitting, Mozilla will ask why you need each permission:

- **cookies**: To extract sessionKey cookie from claude.ai for authentication
- **tabs**: To detect when user visits claude.ai
- **storage**: To store extension settings (server port)
- **notifications**: To notify user about connection status
- **host permission (*://claude.ai/*)**: To access claude.ai cookies and API

### Privacy Policy

You may need to provide a privacy policy. Example:

```
ClaudeCompanion Extension Privacy Policy

Data Collection:
- This extension does NOT collect any personal data
- This extension does NOT send data to any external servers
- All data stays on your local computer

What the extension does:
- Reads cookies from claude.ai (only when you visit the site)
- Sends cookies to LOCAL desktop application (localhost:8383)
- No data is transmitted over the internet

Open Source:
- Full source code is available at: [GitHub URL]
```

### Mozilla Review Process

1. **First submission**: Usually takes 1-2 weeks
2. **Updates**: Faster, typically 1-3 days
3. **Be responsive**: Answer reviewer questions promptly

## Testing Before Submission

1. Load extension as temporary in Firefox (`about:debugging`)
2. Test all functionality
3. Check for JavaScript errors in console
4. Verify permissions are minimal and justified

## Version Management

When updating:

1. Update `version` in `manifest.json`
2. Re-package the extension
3. Submit update with changelog

### Changelog Template

```
Version 1.0.1
- Bug fixes
- Improved error handling

Version 1.0.0
- Initial release
```

## Common Rejection Reasons

- **Too broad permissions**: Only request what you need
- **Missing privacy policy**: Required if you handle user data
- **Unclear description**: Explain clearly what the extension does
- **Missing source code**: Provide GitHub link if using bundlers/minifiers

## After Approval

1. Extension will be available at: `https://addons.mozilla.org/firefox/addon/[addon-name]/`
2. Users can install directly from Mozilla Add-ons
3. Automatic updates will work

## Resources

- [Extension Workshop](https://extensionworkshop.com/)
- [Submission Guidelines](https://extensionworkshop.com/documentation/publish/submitting-an-add-on/)
- [Review Policies](https://extensionworkshop.com/documentation/publish/add-on-policies/)
