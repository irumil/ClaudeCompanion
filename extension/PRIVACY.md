# Privacy Policy for ClaudeCompanion Context Provider

**Last Updated:** December 18, 2024

## Overview

ClaudeCompanion Context Provider is a browser extension designed to work with the ClaudeCompanion desktop application. It monitors your Claude.ai API usage quota by extracting authentication data and sending it to your local desktop application.

## Data Collection

This extension collects the following data from claude.ai:

1. **Session Cookie (sessionKey)**
   - Purpose: Authentication with Claude.ai API
   - Source: Browser cookies from claude.ai domain
   - Storage: Not stored by extension, only transmitted

2. **Organization UUID**
   - Purpose: Construct API endpoint for quota monitoring
   - Source: Claude.ai API response
   - Storage: Not stored by extension, only transmitted

3. **API Usage URL**
   - Purpose: Endpoint for monitoring API quota
   - Source: Constructed from Organization UUID
   - Storage: Not stored by extension, only transmitted

## Data Usage

**All collected data is sent ONLY to:**
- `http://localhost:8383` (your local computer)
- The ClaudeCompanion desktop application running on your machine

**Data is used for:**
- Monitoring your Claude.ai API usage quota
- Displaying quota information in system tray
- Sending notifications when quota is low

## Data Sharing

**NO data is shared with:**
- ❌ External servers
- ❌ Third-party services
- ❌ Analytics services
- ❌ The extension developer
- ❌ Any other parties

**All data stays on your computer.**

## Data Storage

This extension does NOT store any data persistently:
- No cookies are stored
- No local storage is used
- No indexedDB is used
- Data is only kept in memory temporarily during transmission

## Third-Party Services

This extension does NOT use any third-party services:
- No analytics (Google Analytics, etc.)
- No crash reporting
- No advertising
- No tracking pixels

## Permissions Used

The extension requests the following permissions:

1. **`cookies`**
   - Purpose: Read sessionKey from claude.ai
   - Used only when you visit claude.ai

2. **`host_permissions: ["*://*.claude.ai/*"]`**
   - Purpose: Access Claude.ai cookies and API
   - Limited to claude.ai domain only

3. **`webRequest`, `webRequestBlocking`**
   - Purpose: Monitor network requests to detect claude.ai visits
   - Does not modify requests

## User Control

You have full control:
- Disable the extension at any time (about:addons)
- Remove the extension to stop all data collection
- Close the desktop application to stop quota monitoring

## Updates to This Policy

We may update this privacy policy. Changes will be reflected in the "Last Updated" date.

## Open Source

This extension is open source. You can review the complete source code at:
[https://github.com/irumil/ClaudeCompanion]

## Contact

For questions or concerns about privacy:
- GitHub Issues: https://github.com/irumil/ClaudeCompanion/issues
- Email: rizmailov1983@gmail.com

## Compliance

This extension:
- ✅ Complies with Mozilla Add-on Policies
- ✅ Does not collect personal information
- ✅ Does not transmit data to external servers
- ✅ Is transparent about data usage

---

**Summary:** This extension helps you monitor your Claude.ai quota by sending authentication data to a desktop application running on YOUR computer. No data leaves your machine.
