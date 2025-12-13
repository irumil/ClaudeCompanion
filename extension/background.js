// Default settings
const DEFAULT_PORT = 8383;
let currentPort = DEFAULT_PORT;

// Storage for captured anthropic-client-sha from real browser requests
let capturedClientSha = null;

// Capture anthropic-client-sha from real API requests
browser.webRequest.onBeforeSendHeaders.addListener(
  (details) => {
    // Find anthropic-client-sha header in real browser request
    const shaHeader = details.requestHeaders.find(h => h.name.toLowerCase() === 'anthropic-client-sha');
    if (shaHeader && shaHeader.value) {
      capturedClientSha = shaHeader.value;
      console.log('[ClaudeCompanion] ✅ Captured anthropic-client-sha:', capturedClientSha);
    }
  },
  {urls: ["*://claude.ai/api/*"]},
  ["requestHeaders"]
);

// Load settings on startup
browser.storage.local.get(['port']).then((result) => {
  if (result.port) {
    currentPort = result.port;
  }
  console.log(`ClaudeCompanion Extension loaded. Port: ${currentPort}`);
});

// Listen for storage changes
browser.storage.onChanged.addListener((changes, area) => {
  if (area === 'local' && changes.port) {
    currentPort = changes.port.newValue || DEFAULT_PORT;
    console.log(`Port updated to: ${currentPort}`);
  }
});

// Function to get Organization UUID and build usage URL
async function getOrgData() {
  try {
    console.log('[ClaudeCompanion] Fetching organization info...');

    // Fetch organizations list
    const orgsResponse = await fetch("https://claude.ai/api/organizations", {
      credentials: "include"
    });

    if (!orgsResponse.ok) {
      console.error('[ClaudeCompanion] ❌ Failed to fetch organizations:', orgsResponse.status);
      return null;
    }

    const orgs = await orgsResponse.json();

    if (!orgs || orgs.length === 0) {
      console.error('[ClaudeCompanion] ❌ No organizations found');
      return null;
    }

    // Get UUID of first organization
    const orgUuid = orgs[0].uuid;
    const usageUrl = `https://claude.ai/api/organizations/${orgUuid}/usage`;

    console.log('[ClaudeCompanion] ✅ Organization UUID:', orgUuid);
    console.log('[ClaudeCompanion] ✅ Usage URL:', usageUrl);

    return {
      organizationId: orgUuid,
      usageUrl: usageUrl
    };

  } catch (error) {
    console.error('[ClaudeCompanion] ❌ Error fetching organization:', error);
    return null;
  }
}

// Function to send context to desktop app
async function sendContext() {
  const endpoint = `http://127.0.0.1:${currentPort}/set-context`;

  // Get organization data (UUID and usage URL)
  const orgData = await getOrgData();

  if (!orgData) {
    console.error('[ClaudeCompanion] ❌ Cannot get organization data');
    return false;
  }

  // Get ALL cookies for claude.ai to better emulate browser requests
  const allCookies = await browser.cookies.getAll({
    url: "https://claude.ai"
  });

  if (!allCookies || allCookies.length === 0) {
    console.error('[ClaudeCompanion] ❌ No cookies found for claude.ai');
    return false;
  }

  // Convert cookies array to cookie string format
  const cookieString = allCookies.map(c => `${c.name}=${c.value}`).join('; ');
  console.log('[ClaudeCompanion] ✅ Found cookies:', allCookies.map(c => c.name).join(', '));

  // Extract important values from cookies for anthropic-* headers
  const findCookie = (name) => allCookies.find(c => c.name === name)?.value || '';
  const anonymousId = findCookie('ajs_anonymous_id') || findCookie('anthropic-anonymous-id');
  const deviceId = findCookie('anthropic-device-id');

  console.log('[ClaudeCompanion] ✅ Anonymous ID:', anonymousId);
  console.log('[ClaudeCompanion] ✅ Device ID:', deviceId);

  // Build headers to match real browser API requests exactly
  // Based on actual Firefox request to /api/organizations/.../usage
  const headers = {
    'User-Agent': navigator.userAgent,  // Real browser User-Agent
    'Accept': '*/*',
    'Accept-Language': navigator.language || 'ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3',
    // Note: Accept-Encoding excluded - curl doesn't handle gzip automatically
    // Note: Connection and Host are added by curl automatically
    'Referer': 'https://claude.ai/settings/usage',
    'Content-Type': 'application/json',

    // Critical Anthropic headers - these are essential!
    'anthropic-client-platform': 'web_claude_ai',
    'anthropic-client-version': '1.0.0',

    // Sec-Fetch headers for CORS requests
    'Sec-Fetch-Dest': 'empty',
    'Sec-Fetch-Mode': 'cors',
    'Sec-Fetch-Site': 'same-origin',
    'Priority': 'u=4',
    'TE': 'trailers'
  };

  // Add dynamic anthropic headers
  if (anonymousId) {
    headers['anthropic-anonymous-id'] = anonymousId;
  }
  if (deviceId) {
    headers['anthropic-device-id'] = deviceId;
  }

  // Add captured client SHA from real browser requests
  // This is dynamically extracted to always match current Claude.ai version
  if (capturedClientSha) {
    headers['anthropic-client-sha'] = capturedClientSha;
    console.log('[ClaudeCompanion] ✅ Using captured client SHA:', capturedClientSha);
  } else {
    console.log('[ClaudeCompanion] ⚠️ No client SHA captured yet (will be captured on next API request)');
  }

  console.log('[ClaudeCompanion] ✅ Headers prepared:', Object.keys(headers).join(', '));

  const payload = {
    cookies: cookieString,  // Send ALL cookies
    targetUrl: orgData.usageUrl,
    organizationId: orgData.organizationId,
    headers: headers  // All browser headers including User-Agent
  };

  console.log('[ClaudeCompanion] Sending context to desktop app:', {
    usageUrl: orgData.usageUrl,
    organizationId: orgData.organizationId,
    cookies: allCookies.map(c => c.name).join(', '),
    headers: Object.keys(headers).join(', '),
    endpoint: endpoint
  });

  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload)
    });

    if (response.ok) {
      const data = await response.json();
      console.log('[ClaudeCompanion] ✅ Context sent successfully:', data);
      return true;
    } else {
      console.error('[ClaudeCompanion] ❌ Failed to send context:', response.status, response.statusText);
      return false;
    }
  } catch (error) {
    console.error('[ClaudeCompanion] ❌ Error sending context to desktop app:', error);
    console.error('[ClaudeCompanion] ⚠️ Make sure the ClaudeCompanion desktop application is running');

    // Show notification if desktop app is not running
    if (browser.notifications) {
      browser.notifications.create({
        type: 'basic',
        iconUrl: 'icon48.png',
        title: 'ClaudeCompanion не запущен',
        message: 'Запустите ClaudeCompanion.exe для отслеживания квоты.'
      });
    }

    return false;
  }
}

// Function to check if URL is claude.ai
function isClaudeAi(url) {
  return url && url.includes('claude.ai');
}

// Listen for tab updates
browser.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
  // Send context whenever a claude.ai page is loaded
  if (changeInfo.status === 'complete' && isClaudeAi(tab.url)) {
    console.log('[ClaudeCompanion] Claude.ai page detected:', tab.url);

    // Wait a bit for the page to make API requests and capture cookies
    setTimeout(() => {
      sendContext();
    }, 1500);
  }
});

// Also listen for tab activation (when switching to already loaded tab)
browser.tabs.onActivated.addListener(async (activeInfo) => {
  const tab = await browser.tabs.get(activeInfo.tabId);

  if (isClaudeAi(tab.url)) {
    console.log('[ClaudeCompanion] Activated Claude.ai page:', tab.url);

    // Send context after a short delay
    setTimeout(() => {
      sendContext();
    }, 500);
  }
});

// Listen for messages from other parts of the extension
browser.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'testConnection') {
    // Test connection to desktop app
    fetch(`http://127.0.0.1:${currentPort}/health`)
      .then(response => response.json())
      .then(data => {
        sendResponse({ success: true, data: data });
      })
      .catch(error => {
        sendResponse({ success: false, error: error.message });
      });
    return true; // Indicates we'll send response asynchronously
  }
});

console.log('[ClaudeCompanion] Extension background script initialized');
