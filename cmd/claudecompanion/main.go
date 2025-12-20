package main

import (
	_ "embed"
	"math/rand"
	"time"

	"claudecompanion/internal/api"
	"claudecompanion/internal/config"
	"claudecompanion/internal/logger"
	"claudecompanion/internal/notifier"
	"claudecompanion/internal/server"
	"claudecompanion/internal/tray"

	"github.com/getlantern/systray"
	"github.com/robfig/cron/v3"
)

//go:embed icon.ico
var embeddedIcon []byte

// App represents the main application
type App struct {
	configMgr         *config.Manager
	apiClient         *api.Client
	httpServer        *server.Server
	trayMgr           *tray.TrayManager
	notifier          *notifier.Notifier
	cronScheduler     *cron.Cron
	errorCount        int
	lastValue         int
	stopChan          chan struct{}
	demoMode          bool
	demoStarted       time.Time
	demoGreetingShown bool
}

func main() {
	// Initialize logger first - THIS IS CRITICAL
	if err := logger.Init(); err != nil {
		panic("FATAL: Failed to initialize logger: " + err.Error())
	}
	defer logger.Close()

	logger.Info("===========================================")
	logger.Info("ClaudeCompanion starting...")
	logger.Info("===========================================")

	// Catch panics and log them
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("Application panicked: %v", r)
		}
	}()

	logger.Info("Initializing random seed...")
	rand.Seed(time.Now().UnixNano())

	app := &App{
		stopChan:  make(chan struct{}),
		lastValue: -1,
	}

	// Initialize configuration
	logger.Info("Loading configuration...")
	cfgMgr, err := config.NewManager()
	if err != nil {
		logger.Fatal("Failed to initialize config: %v", err)
		return
	}
	app.configMgr = cfgMgr
	logger.Info("Configuration loaded successfully from: %s", cfgMgr.GetPath())

	cfg := cfgMgr.Get()

	// Setup file logging based on config
	if err := logger.SetFileLogging(cfg.EnableFileLogging); err != nil {
		logger.Warning("Failed to setup file logging: %v", err)
	}

	logger.Debug("Config: ServerPort=%d, PollInterval=%ds, DemoMode=%v, Proxy=%s, FileLogging=%v",
		cfg.ServerPort, cfg.PollIntervalSeconds, cfg.DemoMode.Enabled, cfg.Proxy, cfg.EnableFileLogging)

	// Check demo mode
	if cfg.DemoMode.Enabled {
		logger.Warning("===========================================")
		logger.Warning("Starting in DEMO MODE")
		logger.Warning("Duration: %d seconds", cfg.DemoMode.DurationSeconds)
		logger.Warning("===========================================")
		app.demoMode = true
		app.demoStarted = time.Now()
	}

	// Initialize components
	logger.Info("Initializing components...")

	logger.Info("  - API client (proxy: %s, curl: %s)...", cfg.Proxy, cfg.CurlPath)
	app.apiClient = api.NewClient(cfg.Proxy, cfg.CurlPath)
	logger.Info("  - API client initialized")

	logger.Info("  - HTTP server on port %d...", cfg.ServerPort)
	app.httpServer = server.NewServer(cfg.ServerPort)
	logger.Info("  - HTTP server initialized")

	logger.Info("  - System tray manager...")
	app.trayMgr = tray.NewTrayManager(cfgMgr.GetPath())
	app.trayMgr.SetBrowserPath(cfg.BrowserPath)
	logger.Info("  - Tray manager initialized")

	logger.Info("  - Notifier...")
	app.notifier = notifier.NewNotifier(embeddedIcon)
	logger.Info("  - Notifier initialized")

	// Set callbacks
	logger.Info("Setting up callbacks...")
	app.trayMgr.SetExitCallback(func() {
		logger.Info("Exit callback triggered by user")
		app.Shutdown()
	})

	// Set refresh callback for manual statistics update
	app.trayMgr.SetRefreshCallback(func() {
		logger.Info("Manual refresh requested by user")
		app.pollManual()
	})

	// Don't set OpenSettings callback - use default implementation from tray.go
	// which opens the file in notepad.exe on Windows

	// Set server callback for context updates
	app.httpServer.SetContextCallback(func(cookies, targetURL, organizationID string, headers map[string]string) {
		logger.Info(">>> Context received from browser extension")
		logger.Info("    URL: %s", targetURL)
		logger.Info("    Organization ID: %s", organizationID)
		logger.Info("    Cookies length: %d characters", len(cookies))
		logger.Info("    Headers count: %d", len(headers))
		if ua, ok := headers["User-Agent"]; ok {
			logger.Info("    User-Agent: %s", ua)
		}
		app.apiClient.SetContext(cookies, targetURL, organizationID, headers)
		app.trayMgr.UpdateTargetURL(targetURL)
		// Reset error count when new cookies arrive
		app.errorCount = 0
		app.notifier.ResetAll()
		logger.Info("    Context updated successfully, error count reset")

		// Setup greeting scheduler when context is received
		app.setupGreetingScheduler()
	})

	// Watch for config changes
	cfgMgr.OnChange(func(newCfg *config.Config) {
		logger.Info(">>> Configuration file changed")

		// Update file logging setting
		if err := logger.SetFileLogging(newCfg.EnableFileLogging); err != nil {
			logger.Warning("    Failed to update file logging: %v", err)
		}

		// Update browser path
		app.trayMgr.SetBrowserPath(newCfg.BrowserPath)
		logger.Info("    Browser path updated: %s", newCfg.BrowserPath)

		// Update API client settings (preserves cookies and context)
		logger.Info("    Updating API client settings (cookies preserved)...")
		app.apiClient.UpdateSettings(newCfg.Proxy, newCfg.CurlPath)
		logger.Info("    API client settings updated successfully")
	})

	// Start HTTP server (unless in demo mode)
	if !app.demoMode {
		logger.Info("Starting HTTP server on http://127.0.0.1:%d ...", cfg.ServerPort)
		if err := app.httpServer.Start(); err != nil {
			logger.Fatal("Failed to start HTTP server: %v", err)
			return
		}
		logger.Info("HTTP server started successfully")
		logger.Info("  - Endpoint: POST http://127.0.0.1:%d/set-context", cfg.ServerPort)
		logger.Info("  - Health: GET http://127.0.0.1:%d/health", cfg.ServerPort)
	} else {
		logger.Info("Demo mode: HTTP server NOT started")
	}

	// Run systray
	logger.Info("===========================================")
	logger.Info("Initializing system tray...")
	logger.Info("===========================================")
	systray.Run(func() {
		logger.Info(">>> Systray ready, initializing tray manager...")
		app.trayMgr.Initialize()
		logger.Info(">>> Tray manager initialized successfully")

		// Wait a bit for systray to fully initialize
		logger.Info("Waiting for systray to stabilize...")
		time.Sleep(500 * time.Millisecond)

		// Start polling AFTER systray is ready
		if app.demoMode {
			logger.Info("Starting DEMO mode loop (interval: 2 seconds)...")
			go app.demoLoop()
		} else {
			logger.Info("Starting API poll loop (interval: %d seconds)...", cfg.PollIntervalSeconds)
			go app.pollLoop()
		}

		logger.Info("===========================================")
		logger.Info("Application startup complete!")
		logger.Info("===========================================")
	}, func() {
		logger.Info("Systray exit callback triggered")
		app.Shutdown()
	})

	logger.Info("Main function completed, application should now be running in tray")
}

// demoLoop runs demo mode with 2 second interval
func (a *App) demoLoop() {
	logger.Info("Demo loop started")
	cfg := a.configMgr.Get()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial update
	a.handleDemoMode(cfg)

	for {
		select {
		case <-ticker.C:
			a.handleDemoMode(cfg)
		case <-a.stopChan:
			logger.Info("Demo loop stopped")
			return
		}
	}
}

// pollLoop polls the API periodically
func (a *App) pollLoop() {
	logger.Info("Poll loop started")
	cfg := a.configMgr.Get()
	ticker := time.NewTicker(time.Duration(cfg.PollIntervalSeconds) * time.Second)
	defer ticker.Stop()

	// Poll immediately on start
	logger.Info("Performing initial poll...")
	a.poll()

	for {
		select {
		case <-ticker.C:
			logger.Debug("Poll tick (interval: %ds)", cfg.PollIntervalSeconds)
			a.poll()
			// Update ticker interval if config changed
			newCfg := a.configMgr.Get()
			if newCfg.PollIntervalSeconds != cfg.PollIntervalSeconds {
				logger.Info("Poll interval changed: %ds -> %ds", cfg.PollIntervalSeconds, newCfg.PollIntervalSeconds)
				cfg = newCfg
				ticker.Reset(time.Duration(cfg.PollIntervalSeconds) * time.Second)
			}
		case <-a.stopChan:
			logger.Info("Poll loop stopped")
			return
		}
	}
}

// poll performs a single API poll (automatic)
func (a *App) poll() {
	a.doPoll(false)
}

// pollManual performs a manual API poll (user requested)
func (a *App) pollManual() {
	a.doPoll(true)
}

// doPoll performs a single API poll
func (a *App) doPoll(isManual bool) {
	cfg := a.configMgr.Get()

	// Check if we have context
	if !a.apiClient.HasContext() {
		logger.Debug("No cookies received yet from extension")
		a.updateTrayNoCookies()
		return
	}

	// Check work hours (skip check if manual request)
	if !isManual && !cfg.WorkHours.IsWithinWorkHours() {
		logger.Debug("Outside work hours, skipping automatic poll")
		return
	}

	// Fetch usage
	if isManual {
		logger.Info("Manual poll: Fetching usage from API...")
	} else {
		logger.Debug("Automatic poll: Fetching usage from API...")
	}

	usage, err := a.apiClient.GetUsage()
	if err != nil {
		logger.Error("API request failed (error #%d): %v", a.errorCount+1, err)
		a.errorCount++
		a.handleError(cfg)
		return
	}

	// Success - reset error count
	if a.errorCount > 0 {
		logger.Info("API request succeeded after %d errors", a.errorCount)
	}
	a.errorCount = 0
	a.notifier.ResetErrorNotification()

	// Get inverted value (remaining quota)
	value := usage.GetInvertedValue()
	tooltip := usage.FormatTooltip()

	logger.Debug("API response: remaining=%d%%, tooltip=%s", value, tooltip)

	// Update tray
	a.trayMgr.UpdateIcon(value, false, tooltip)
	a.lastValue = value

	// Check for low value notifications
	a.checkLowValueNotifications(value, cfg, usage)
}

// handleError handles API errors
func (a *App) handleError(cfg config.Config) {
	// Show gray icon after threshold
	if a.errorCount >= cfg.GrayModeThreshold {
		tooltip := "Ошибка подключения к API"
		logger.Warning("Error count (%d) reached gray mode threshold (%d)", a.errorCount, cfg.GrayModeThreshold)
		a.trayMgr.UpdateIcon(a.lastValue, true, tooltip)
	}

	// Show notification after threshold
	if a.errorCount >= cfg.NotificationThreshold {
		logger.Warning("Error count (%d) reached notification threshold (%d)", a.errorCount, cfg.NotificationThreshold)
		a.notifier.NotifyError(a.errorCount, cfg.NotificationThreshold)
	}
}

// checkLowValueNotifications checks if we should show low value notifications
func (a *App) checkLowValueNotifications(value int, cfg config.Config, usage *api.UsageResponse) {
	if !cfg.LowValueNotifications.Enabled {
		return
	}

	// Check for zero value (show notification once)
	if value == 0 {
		phrase := config.GetRandomPhrase(cfg.LowValueNotifications.ZeroPhrases)
		var resetTime string
		if usage.FiveHour.ResetsAt != nil {
			resetTime = usage.FiveHour.ResetsAt.Local().Format("15:04:05")
		} else {
			resetTime = "—"
		}
		logger.Warning("Quota reached ZERO, showing notification")
		a.notifier.NotifyZero(phrase, resetTime)
		return
	}

	// Check for low value (show notification once when dropping below threshold)
	if value <= cfg.LowValueNotifications.Threshold {
		phrase := config.GetRandomPhrase(cfg.LowValueNotifications.Phrases)
		logger.Warning("Quota low (%d%% <= %d%%), showing notification", value, cfg.LowValueNotifications.Threshold)
		a.notifier.NotifyLowValue(value, phrase)
	} else {
		// Reset notification state when value goes above threshold
		// This also resets zero notification state
		a.notifier.ResetLowValueNotification()
	}
}

// updateTrayNoCookies updates the tray when no cookies are available
func (a *App) updateTrayNoCookies() {
	a.trayMgr.UpdateIcon(-1, false, "Ожидаю куки от расширения")
}

// handleDemoMode simulates declining values in demo mode (infinite loop)
func (a *App) handleDemoMode(cfg config.Config) {
	elapsed := time.Since(a.demoStarted).Seconds()
	duration := float64(cfg.DemoMode.DurationSeconds)

	// Loop infinitely: calculate position in current cycle
	cyclePosition := elapsed
	for cyclePosition >= duration {
		cyclePosition -= duration
		// Reset demo start for next cycle
		a.demoStarted = a.demoStarted.Add(time.Duration(duration) * time.Second)
	}

	// Calculate current value (100 -> 0)
	progress := cyclePosition / duration
	value := int(100 - (progress * 100))

	// Force zero in the last 3 seconds to ensure it's reached
	if cyclePosition >= duration-3 {
		value = 0
	}

	if value < 0 {
		value = 0
	}

	logger.Debug("Demo: cycle_pos=%.1fs, progress=%.2f, value=%d", cyclePosition, progress, value)

	// Create fake usage response for demo mode
	utilization := float64(100 - value)
	fakeUsage := &api.UsageResponse{}
	fakeUsage.FiveHour.Utilization = utilization
	fiveHourReset := time.Now().Add(time.Hour * 2)
	fakeUsage.FiveHour.ResetsAt = &fiveHourReset
	fakeUsage.SevenDay.Utilization = utilization / 2 // Half for weekly
	sevenDayReset := time.Now().Add(time.Hour * 24 * 7)
	fakeUsage.SevenDay.ResetsAt = &sevenDayReset

	tooltip := fakeUsage.FormatTooltip()
	a.trayMgr.UpdateIcon(value, false, tooltip)
	a.lastValue = value

	// Reset greeting notification flag at the start of each cycle
	if value >= 95 && value <= 100 {
		logger.Debug("Demo: New cycle starting")
		a.demoGreetingShown = false // Reset greeting notification flag for new cycle
	}

	// Show greeting notification at the start of each cycle (once per cycle)
	if value >= 95 && value <= 100 && !a.demoGreetingShown {
		logger.Info("Demo: Showing greeting notification")
		a.notifier.NotifyGreeting()
		a.demoGreetingShown = true
	}

	// Trigger notifications in demo mode
	// checkLowValueNotifications handles reset when value goes above threshold
	a.checkLowValueNotifications(value, cfg, fakeUsage)

	// No error simulation in demo mode - let it run clean
	a.errorCount = 0
}

// setupGreetingScheduler sets up cron scheduler for greeting messages
func (a *App) setupGreetingScheduler() {
	cfg := a.configMgr.Get()

	// Stop existing scheduler if any
	if a.cronScheduler != nil {
		a.cronScheduler.Stop()
		logger.Info("Stopped existing greeting scheduler")
	}

	// Check if greeting is configured
	if cfg.Greeting.Cron == "" || cfg.Greeting.ChatID == "" {
		logger.Info("Greeting not configured (cron or chat_id missing), scheduler not started")
		return
	}

	// Create new cron scheduler
	a.cronScheduler = cron.New()

	// Add greeting job
	_, err := a.cronScheduler.AddFunc(cfg.Greeting.Cron, func() {
		logger.Info(">>> Cron triggered: Sending greeting message")
		a.sendGreeting()
	})

	if err != nil {
		logger.Error("Failed to setup greeting scheduler: %v", err)
		return
	}

	// Start scheduler
	a.cronScheduler.Start()
	logger.Info("Greeting scheduler started with cron: %s", cfg.Greeting.Cron)
}

// sendGreeting sends a greeting message to configured chat
func (a *App) sendGreeting() {
	cfg := a.configMgr.Get()

	if cfg.Greeting.ChatID == "" {
		logger.Warning("Greeting chat ID not configured")
		return
	}

	text := cfg.Greeting.Text
	if text == "" {
		text = "Ok"
	}

	logger.Info("Sending greeting: '%s' to chat %s", text, cfg.Greeting.ChatID)

	if err := a.apiClient.SendGreeting(cfg.Greeting.ChatID, text); err != nil {
		logger.Error("Failed to send greeting: %v", err)
		return
	}

	logger.Info("Greeting sent successfully!")
	a.notifier.NotifyGreeting()
}

// Shutdown performs cleanup before exit
func (a *App) Shutdown() {
	logger.Info("===========================================")
	logger.Info("Shutting down ClaudeCompanion...")
	logger.Info("===========================================")
	close(a.stopChan)
	if a.cronScheduler != nil {
		logger.Info("Stopping cron scheduler...")
		a.cronScheduler.Stop()
		logger.Info("Cron scheduler stopped")
	}
	if a.httpServer != nil {
		logger.Info("Stopping HTTP server...")
		a.httpServer.Stop()
		logger.Info("HTTP server stopped")
	}
	logger.Info("Shutdown complete. Goodbye!")
	logger.Info("===========================================")
}
