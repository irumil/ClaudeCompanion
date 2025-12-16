package tray

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"claudecompanion/internal/icon"

	"github.com/getlantern/systray"
)

// TrayManager manages the system tray icon
type TrayManager struct {
	iconGen        *icon.Generator
	currentValue   int
	currentMode    icon.ColorMode
	hasError       bool
	tooltip        string
	targetURL      string
	configPath     string
	browserPath    string
	onExit         func()
	onOpenSettings func()
	onClick        func()
	onRefresh      func()
}

// NewTrayManager creates a new tray manager
func NewTrayManager(configPath string) *TrayManager {
	return &TrayManager{
		iconGen:      icon.NewGenerator(),
		currentValue: -1,
		currentMode:  icon.ColorGray,
		tooltip:      "Ожидаю куки от расширения",
		targetURL:    "https://claude.ai",
		configPath:   configPath,
	}
}

// SetExitCallback sets the callback for exit action
func (t *TrayManager) SetExitCallback(callback func()) {
	t.onExit = callback
}

// SetOpenSettingsCallback sets the callback for opening settings
func (t *TrayManager) SetOpenSettingsCallback(callback func()) {
	t.onOpenSettings = callback
}

// SetClickCallback sets the callback for icon click
func (t *TrayManager) SetClickCallback(callback func()) {
	t.onClick = callback
}

// SetRefreshCallback sets the callback for manual refresh
func (t *TrayManager) SetRefreshCallback(callback func()) {
	t.onRefresh = callback
}

// SetBrowserPath sets the browser path to use for opening URLs
func (t *TrayManager) SetBrowserPath(path string) {
	t.browserPath = path
}

// Initialize sets up the system tray
func (t *TrayManager) Initialize() {
	// Set initial tooltip and icon
	systray.SetTooltip(t.tooltip)

	// Set default gray icon
	iconData := icon.GetDefaultIcon()
	systray.SetIcon(iconData)

	// Create menu items
	mOpenClaude := systray.AddMenuItem("Открыть Claude.ai", "Открыть сайт Claude.ai в браузере")
	mRefresh := systray.AddMenuItem("Получить статистику", "Обновить статистику сейчас")
	systray.AddSeparator()
	mOpenSettings := systray.AddMenuItem("Открыть настройки", "Открыть конфигурационный файл")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Выход", "Выйти из приложения")

	// Handle menu clicks in goroutine
	go func() {
		for {
			select {
			case <-mOpenClaude.ClickedCh:
				t.openURL("https://claude.ai")
			case <-mRefresh.ClickedCh:
				if t.onRefresh != nil {
					t.onRefresh()
				}
			case <-mOpenSettings.ClickedCh:
				t.handleOpenSettings()
			case <-mQuit.ClickedCh:
				if t.onExit != nil {
					t.onExit()
				}
				systray.Quit()
				return
			}
		}
	}()
}

// UpdateIcon updates the tray icon with new value
func (t *TrayManager) UpdateIcon(value int, hasError bool, tooltip string) {
	t.currentValue = value
	t.hasError = hasError
	t.tooltip = tooltip
	t.currentMode = icon.GetColorMode(value, hasError)

	// Set tooltip without prefix
	systray.SetTooltip(tooltip)

	// Generate icon with text showing percentage
	var text string
	if value < 0 {
		text = "--"
	} else {
		text = fmt.Sprintf("%d", value)
	}

	iconData, err := t.iconGen.Generate(text, t.currentMode)
	if err != nil {
		log.Printf("Failed to generate icon: %v", err)
		// Fallback to static icon
		iconData = icon.GetIconByMode(t.currentMode)
	}

	systray.SetIcon(iconData)
	log.Printf("Tray updated: value=%d, text=%s, hasError=%v", value, text, hasError)
}

// UpdateTargetURL updates the target URL for click action
func (t *TrayManager) UpdateTargetURL(url string) {
	if url != "" {
		t.targetURL = url
	}
}

// updateIcon generates and sets the tray icon
func (t *TrayManager) updateIcon() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in updateIcon: %v", r)
		}
	}()

	var text string
	if t.currentValue < 0 {
		text = "--"
	} else {
		text = fmt.Sprintf("%d", t.currentValue)
	}

	iconData, err := t.iconGen.Generate(text, t.currentMode)
	if err != nil {
		log.Printf("Failed to generate icon: %v", err)
		return
	}

	// Try to set icon, but don't fail if it doesn't work
	// (Windows systray can be finicky)
	log.Printf("Setting icon: text=%s, size=%d bytes", text, len(iconData))
	systray.SetIcon(iconData)
}

// handleOpenSettings opens the configuration file
func (t *TrayManager) handleOpenSettings() {
	if t.onOpenSettings != nil {
		t.onOpenSettings()
		return
	}

	// Default implementation
	if err := t.openFileInEditor(t.configPath); err != nil {
		log.Printf("Failed to open settings: %v", err)
	}
}

// openFileInEditor opens a file in the system's default editor
func (t *TrayManager) openFileInEditor(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad.exe", path)
	case "darwin":
		cmd = exec.Command("open", "-t", path)
	default: // Linux and others
		// Try common editors
		editors := []string{"xdg-open", "gedit", "nano", "vim"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				cmd = exec.Command(editor, path)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable editor found")
		}
	}

	return cmd.Start()
}

// HandleClick handles icon click events
func (t *TrayManager) HandleClick() {
	if t.onClick != nil {
		t.onClick()
		return
	}

	// Default implementation: open target URL in browser
	t.openURL(t.targetURL)
}

// openURL opens a URL in the default browser or specified browser
func (t *TrayManager) openURL(url string) error {
	var cmd *exec.Cmd

	// Use custom browser if specified
	if t.browserPath != "" {
		cmd = exec.Command(t.browserPath, url)
	} else {
		// Use default browser
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		default: // Linux and others
			cmd = exec.Command("xdg-open", url)
		}
	}

	return cmd.Start()
}

// Quit exits the tray application
func (t *TrayManager) Quit() {
	systray.Quit()
}

// OnReady is called when systray is ready
func OnReady(initFunc func()) {
	initFunc()
}

// OnExit is called when systray is exiting
func OnExit() {
	log.Println("Systray exiting...")
}
