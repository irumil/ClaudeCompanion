package notifier

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-toast/toast"
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
func NewNotifier() *Notifier {
	return &Notifier{
		state: &NotificationState{},
	}
}

// getIconPath returns the path to the application icon
func getIconPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	exeDir := filepath.Dir(exePath)
	iconPath := filepath.Join(exeDir, "icon.ico")

	// If icon.ico doesn't exist, use the exe itself (Windows can extract icon from exe)
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		return exePath
	}
	return iconPath
}

// NotifyError shows an error notification (for authorization issues)
func (n *Notifier) NotifyError(errorCount int, threshold int) {
	n.state.mu.Lock()
	defer n.state.mu.Unlock()

	if errorCount >= threshold && !n.state.lastErrorNotification {
		title := "Проблема с авторизацией"
		message := "Сайт не принимает запросы. Возможно, сессия устарела. Пожалуйста, зайдите на сайт и обновите авторизацию. 🔐"

		log.Printf("Attempting to show error notification")
		notification := toast.Notification{
			AppID:   "ClaudeCompanion",
			Title:   title,
			Message: message,
			Icon:    getIconPath(),
		}
		if err := notification.Push(); err != nil {
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
		notification := toast.Notification{
			AppID:   "ClaudeCompanion",
			Title:   title,
			Message: message,
			Icon:    getIconPath(),
		}
		if err := notification.Push(); err != nil {
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
		message := phrase + "\nВозвращайся в " + resetTime

		log.Printf("Attempting to show zero notification: %s", message)
		notification := toast.Notification{
			AppID:   "ClaudeCompanion",
			Title:   title,
			Message: message,
			Icon:    getIconPath(),
		}
		if err := notification.Push(); err != nil {
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
	title := "Утренний привет Клоду ☀️"
	//message := "Сообщение отправлено успешно! ☀️"

	log.Printf("Attempting to show greeting notification")
	notification := toast.Notification{
		AppID: "ClaudeCompanion",
		Title: title,
		//Message: message,
		Icon: getIconPath(),
	}
	if err := notification.Push(); err != nil {
		log.Printf("Failed to show greeting notification: %v", err)
	} else {
		log.Println("Greeting notification shown successfully")
	}
}
