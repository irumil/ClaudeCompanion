package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	ServerPort            int                   `yaml:"server_port"`
	PollIntervalSeconds   int                   `yaml:"poll_interval_seconds"`
	GrayModeThreshold     int                   `yaml:"gray_mode_threshold"`
	NotificationThreshold int                   `yaml:"notification_threshold"`
	Proxy                 string                `yaml:"proxy"`
	EnableFileLogging     bool                  `yaml:"enable_file_logging"`
	BrowserPath           string                `yaml:"browser_path"`
	CurlPath              string                `yaml:"curl_path"` // Custom path to curl binary
	LowValueNotifications LowValueNotifications `yaml:"low_value_notifications"`
	DemoMode              DemoMode              `yaml:"demo_mode"`
	Greeting              Greeting              `yaml:"greeting"`
	WorkHours             WorkHours             `yaml:"work_hours"`
}

type LowValueNotifications struct {
	Enabled     bool     `yaml:"enabled"`
	Threshold   int      `yaml:"threshold"`
	Phrases     []string `yaml:"phrases"`
	ZeroPhrases []string `yaml:"zero_phrases"`
}

type DemoMode struct {
	Enabled         bool `yaml:"enabled"`
	DurationSeconds int  `yaml:"duration_seconds"`
}

type Greeting struct {
	Cron   string `yaml:"greeting_cron"`
	Text   string `yaml:"greeting_text"`
	ChatID string `yaml:"greeting_chat_id"`
}

type WorkHours struct {
	Enabled bool   `yaml:"enabled"`
	Start   string `yaml:"start"` // Format: "08:00"
	End     string `yaml:"end"`   // Format: "20:00"
}

// Manager handles configuration loading and hot-reloading
type Manager struct {
	mu          sync.RWMutex
	config      *Config
	configPath  string
	lastModTime time.Time
	onChange    []func(*Config)
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	m := &Manager{
		configPath: configPath,
		onChange:   make([]func(*Config), 0),
	}

	// Create config file with defaults if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := m.createDefaultConfig(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	// Load initial configuration
	if err := m.reload(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Start watching for changes
	go m.watchChanges()

	return m, nil
}

// Get returns a copy of the current configuration
func (m *Manager) Get() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return *m.config
}

// GetPath returns the configuration file path
func (m *Manager) GetPath() string {
	return m.configPath
}

// OnChange registers a callback to be called when configuration changes
func (m *Manager) OnChange(callback func(*Config)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = append(m.onChange, callback)
}

// reload reads and parses the configuration file
func (m *Manager) reload() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	// Apply defaults
	if config.ServerPort == 0 {
		config.ServerPort = 8383
	}
	if config.PollIntervalSeconds == 0 {
		config.PollIntervalSeconds = 60 // Changed from 30 to 60 for safety
	}
	if config.GrayModeThreshold == 0 {
		config.GrayModeThreshold = 5
	}
	if config.NotificationThreshold == 0 {
		config.NotificationThreshold = 10
	}

	m.mu.Lock()
	m.config = &config
	callbacks := make([]func(*Config), len(m.onChange))
	copy(callbacks, m.onChange)
	m.mu.Unlock()

	// Call callbacks
	for _, callback := range callbacks {
		callback(&config)
	}

	log.Println("Configuration reloaded successfully")
	return nil
}

// watchChanges monitors the config file for changes
func (m *Manager) watchChanges() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		info, err := os.Stat(m.configPath)
		if err != nil {
			continue
		}

		if info.ModTime().After(m.lastModTime) {
			m.lastModTime = info.ModTime()
			if err := m.reload(); err != nil {
				log.Printf("Failed to reload config: %v", err)
			}
		}
	}
}

// createDefaultConfig creates a default configuration file
func (m *Manager) createDefaultConfig() error {
	configDir := filepath.Dir(m.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	defaultConfig := &Config{
		ServerPort:            8383,
		PollIntervalSeconds:   60, // Changed from 30 to 60 for safety
		GrayModeThreshold:     5,
		NotificationThreshold: 10,
		Proxy:                 "",
		EnableFileLogging:     true,
		BrowserPath:           "",
		LowValueNotifications: LowValueNotifications{
			Enabled:   true,
			Threshold: 20,
			Phrases: []string{
				"Пора идти домой! 🏡",
				"Система устала. Вы — тоже. 😴",
				"Время отдохнуть! ☕",
				"API говорит: хватит на сегодня! 🛑",
			},
			ZeroPhrases: []string{
				"Всё, капут! 💥",
				"0 — это не число, это приговор. 🛌",
				"Game over! 🎮",
				"Лимит исчерпан! 🚫",
			},
		},
		DemoMode: DemoMode{
			Enabled:         false,
			DurationSeconds: 60,
		},
		Greeting: Greeting{
			Cron:   "0 8 * * *", // 8 AM every day
			Text:   "Ok",
			ChatID: "", // User must specify chat UUID
		},
		WorkHours: WorkHours{
			Enabled: true,    // Enabled by default
			Start:   "08:00", // 8 AM
			End:     "20:00", // 8 PM
		},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// getConfigPath returns the config file path
func getConfigPath() (string, error) {
	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get directory containing the executable
	exeDir := filepath.Dir(exePath)

	// On macOS, use Application Support directory
	// Check if running from .app bundle
	if filepath.Base(exeDir) == "MacOS" && filepath.Base(filepath.Dir(exeDir)) == "Contents" {
		// Running from .app bundle, use ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir := filepath.Join(homeDir, "Library", "Application Support", "ClaudeCompanion")
		return filepath.Join(configDir, "config.yaml"), nil
	}

	// For Windows/Linux or when running from source, use executable directory
	return filepath.Join(exeDir, "config.yaml"), nil
}

// GetRandomPhrase returns a random phrase from the list
func GetRandomPhrase(phrases []string) string {
	if len(phrases) == 0 {
		return ""
	}
	return phrases[rand.Intn(len(phrases))]
}

// IsWithinWorkHours checks if current time is within configured work hours
func (wh *WorkHours) IsWithinWorkHours() bool {
	if !wh.Enabled {
		return true // Always allow if work hours not enabled
	}

	now := time.Now()
	currentTime := now.Format("15:04")

	// Parse start and end times
	start := wh.Start
	end := wh.End

	// Simple string comparison works for HH:MM format
	if start <= end {
		// Normal case: 08:00 - 20:00
		return currentTime >= start && currentTime < end
	} else {
		// Overnight case: 20:00 - 08:00 (next day)
		return currentTime >= start || currentTime < end
	}
}
