//go:build darwin
// +build darwin

package notifier

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// NotificationState tracks which notifications have been shown
type NotificationState struct {
	mu                    sync.Mutex
	lastErrorNotification bool
	lastLowValueNotif     bool
	lastZeroNotif         bool
}

// Notifier handles system notifications
type Notifier struct {
	state *NotificationState
}

// NewNotifier creates a new notifier
// embeddedIcon parameter is accepted for API consistency but not used on macOS
func NewNotifier(embeddedIcon []byte) *Notifier {
	return &Notifier{
		state: &NotificationState{},
	}
}

// getIconPath returns the path to the app icon for notifications
func getIconPath() string {
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	exeDir := filepath.Dir(exePath)

	// Try app-icon.png in the same directory
	iconPath := filepath.Join(exeDir, "app-icon.png")
	if _, err := os.Stat(iconPath); err == nil {
		return iconPath
	}

	// Try icon96.png from extension folder (when running from source)
	iconPath = filepath.Join(exeDir, "..", "extension", "icon96.png")
	if _, err := os.Stat(iconPath); err == nil {
		absPath, _ := filepath.Abs(iconPath)
		return absPath
	}

	return ""
}

// showNotification displays a native macOS notification banner using terminal-notifier
func showNotification(title, message string) error {
	args := []string{
		"-title", "ClaudeCompanion",
		"-subtitle", title,
		"-message", message,
		"-sound", "default",
	}

	// Add icon if available (using contentImage for better visibility)
	iconPath := getIconPath()
	if iconPath != "" {
		args = append(args, "-contentImage", iconPath)
		log.Printf("[DEBUG] Using icon: %s", iconPath)
	} else {
		log.Printf("[WARNING] No icon found")
	}

	cmd := exec.Command("terminal-notifier", args...)

	// Run asynchronously so it doesn't block
	go func() {
		if err := cmd.Run(); err != nil {
			log.Printf("terminal-notifier failed: %v", err)
		}
	}()

	return nil
}

// NotifyError shows an error notification (for authorization issues)
func (n *Notifier) NotifyError(errorCount int, threshold int) {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()

	if errorCount >= threshold && !n.state.lastErrorNotification {
		title := "Проблема с авторизацией"
		message := "Сайт не принимает запросы. Возможно, сессия устарела."

		log.Printf("Attempting to show error notification")
		if err := showNotification(title, message); err != nil {
			log.Printf("Failed to show notification: %v", err)
		} else {
			log.Println("Error notification shown successfully")
		}

		n.state.lastErrorNotification = true
	}
}

// NotifyLowValue shows a notification when value is low
func (n *Notifier) NotifyLowValue(value int, phrase string) {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()

	if value > 0 && !n.state.lastLowValueNotif {
		title := "Низкая квота"
		message := phrase

		log.Printf("Attempting to show low value notification: %s", phrase)
		if err := showNotification(title, message); err != nil {
			log.Printf("Failed to show notification: %v", err)
		} else {
			log.Printf("Low value notification shown successfully: %s", phrase)
		}

		n.state.lastLowValueNotif = true
		n.state.lastZeroNotif = false // Reset zero notification state
	}
}

// NotifyZero shows a notification when value reaches zero
func (n *Notifier) NotifyZero(phrase string, resetTime string) {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()

	if !n.state.lastZeroNotif {
		title := "Квота исчерпана"
		message := phrase + " Возвращайся в " + resetTime

		log.Printf("Attempting to show zero notification: %s", message)
		if err := showNotification(title, message); err != nil {
			log.Printf("Failed to show notification: %v", err)
		} else {
			log.Printf("Zero notification shown successfully: %s", message)
		}

		n.state.lastZeroNotif = true
		n.state.lastLowValueNotif = false // Reset low value state
	}
}

// ResetErrorNotification resets the error notification state
func (n *Notifier) ResetErrorNotification() {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()
	if n.state.lastErrorNotification {
		n.state.lastErrorNotification = false
		log.Println("Error notification state reset")
	}
}

// ResetLowValueNotification resets low value notification state
func (n *Notifier) ResetLowValueNotification() {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()
	if n.state.lastLowValueNotif || n.state.lastZeroNotif {
		n.state.lastLowValueNotif = false
		n.state.lastZeroNotif = false
		log.Println("Low value notification state reset")
	}
}

// ResetAll resets all notification states
func (n *Notifier) ResetAll() {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()
	n.state.lastErrorNotification = false
	n.state.lastLowValueNotif = false
	n.state.lastZeroNotif = false
	log.Println("All notification states reset")
}

// NotifyGreeting shows a notification when greeting is sent
func (n *Notifier) NotifyGreeting() {
	title := "Утренний привет Клоду"
	message := "Сообщение отправлено успешно!"

	log.Printf("Attempting to show greeting notification")
	if err := showNotification(title, message); err != nil {
		log.Printf("Failed to show greeting notification: %v", err)
	} else {
		log.Println("Greeting notification shown successfully")
	}
}
