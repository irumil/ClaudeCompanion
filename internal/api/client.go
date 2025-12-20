package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

// UsageResponse represents the API response
type UsageResponse struct {
	FiveHour struct {
		Utilization float64    `json:"utilization"`
		ResetsAt    *time.Time `json:"resets_at"`
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64    `json:"utilization"`
		ResetsAt    *time.Time `json:"resets_at"`
	} `json:"seven_day"`
}

// Client handles API requests
type Client struct {
	cookies        string
	targetURL      string
	organizationID string
	headers        map[string]string // Includes User-Agent
	proxy          string
	curlPath       string
}

// NewClient creates a new API client
func NewClient(proxy, curlPath string) *Client {
	return &Client{
		proxy:    proxy,
		curlPath: curlPath,
	}
}

// SetContext updates cookies, target URL, organization ID and headers (includes User-Agent)
func (c *Client) SetContext(cookies, targetURL, organizationID string, headers map[string]string) {
	c.cookies = cookies
	c.targetURL = targetURL
	c.organizationID = organizationID
	c.headers = headers
	log.Printf("Context updated: URL=%s, OrgID=%s, Cookies length=%d, Headers count=%d",
		targetURL, organizationID, len(cookies), len(headers))
}

// HasContext returns true if cookies are set
func (c *Client) HasContext() bool {
	return c.cookies != "" && c.targetURL != ""
}

// GetUsage fetches the current usage from the API using curl
func (c *Client) GetUsage() (*UsageResponse, error) {
	if !c.HasContext() {
		return nil, fmt.Errorf("no context set (cookies not received from extension)")
	}

	// Use curl directly (HTTP client removed - curl works better with proxy)
	return c.fetchWithCurl()
}

// fetchWithCurl performs the request using system curl
func (c *Client) fetchWithCurl() (*UsageResponse, error) {
	curlPath := c.getCurlPath()

	args := []string{
		"-X", "GET",
		c.targetURL,
		"-H", fmt.Sprintf("Cookie: %s", c.cookies), // All cookies from browser
	}

	// Add all browser headers to emulate real browser request
	for key, value := range c.headers {
		// Skip Accept-Encoding because curl doesn't handle gzip automatically
		if key == "Accept-Encoding" {
			continue
		}
		// Add all headers including User-Agent from browser
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value))
	}

	// Add proxy if configured
	if c.proxy != "" {
		args = append([]string{"-x", c.proxy}, args...)
	}

	// Log curl command
	log.Printf("========================================")
	log.Printf("CURL Request:")
	log.Printf("  Command: %s", curlPath)
	log.Printf("  URL: %s", c.targetURL)
	log.Printf("  Proxy: %s", c.proxy)
	log.Printf("  Cookie preview: %.100s...", c.cookies)
	log.Printf("  Full command: %s %v", curlPath, args)
	log.Printf("========================================")

	cmd := exec.Command(curlPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Hide console window on Windows
	hideWindow(cmd)

	if err := cmd.Run(); err != nil {
		log.Printf("CURL execution failed: %v", err)
		log.Printf("CURL stderr: %s", stderr.String())
		log.Printf("CURL stdout: %s", stdout.String())
		log.Printf("========================================")
		return nil, fmt.Errorf("curl execution failed: %w, stderr: %s", err, stderr.String())
	}

	// With -v flag, headers go to stderr, body goes to stdout
	log.Printf("CURL Response Headers (stderr):")
	log.Printf("%s", stderr.String())
	log.Printf("CURL Response Body (stdout):")
	log.Printf("%s", stdout.String())
	log.Printf("========================================")

	// Parse JSON from stdout (body only)
	var usage UsageResponse
	if err := json.Unmarshal(stdout.Bytes(), &usage); err != nil {
		log.Printf("Failed to parse CURL JSON: %v", err)
		log.Printf("JSON body: %s", stdout.String())
		return nil, fmt.Errorf("failed to parse curl output: %w, output: %s", err, stdout.String())
	}

	log.Printf("CURL Success! Parsed usage data")
	log.Printf("========================================")

	return &usage, nil
}

// getCurlPath returns the platform-specific curl path
func (c *Client) getCurlPath() string {
	// Use custom curl path if configured
	if c.curlPath != "" {
		return c.curlPath
	}

	// Otherwise use platform defaults
	switch runtime.GOOS {
	case "windows":
		return "curl.exe"
	case "darwin":
		// macOS: use Homebrew curl to avoid Cloudflare issues
		return "/opt/homebrew/opt/curl/bin/curl"
	default:
		// Linux and others: use system curl
		return "curl"
	}
}

// GetInvertedValue returns the inverted utilization value (100 - utilization)
func (ur *UsageResponse) GetInvertedValue() int {
	utilization := ur.FiveHour.Utilization
	remaining := 100 - utilization
	if remaining < 0 {
		return 0
	}
	if remaining > 100 {
		return 100
	}
	return int(remaining)
}

// GetResetTime returns the time when the quota resets
func (ur *UsageResponse) GetResetTime() *time.Time {
	return ur.FiveHour.ResetsAt
}

// GetSevenDayValue returns the inverted 7-day utilization value
func (ur *UsageResponse) GetSevenDayInvertedValue() int {
	utilization := ur.SevenDay.Utilization
	remaining := 100 - utilization
	if remaining < 0 {
		return 0
	}
	if remaining > 100 {
		return 100
	}
	return int(remaining)
}

// FormatTooltip creates a formatted tooltip string
func (ur *UsageResponse) FormatTooltip() string {
	fiveHourUtilization := ur.FiveHour.Utilization
	sevenDayUtilization := ur.SevenDay.Utilization

	var fiveHourResetTime string
	var sevenDayResetTime string

	if ur.FiveHour.ResetsAt != nil {
		fiveHourResetTime = ur.FiveHour.ResetsAt.Local().Format("15:04:05")
	} else {
		fiveHourResetTime = "—"
	}

	if ur.SevenDay.ResetsAt != nil {
		sevenDayResetTime = ur.SevenDay.ResetsAt.Local().Format("02.01.2006 15:04:05")
	} else {
		sevenDayResetTime = "—"
	}

	// Use \r\n for Windows multiline tooltips
	return fmt.Sprintf("5 часов: %.0f%%, Сброс: %s\r\nНеделя: %.0f%%, Сброс: %s",
		fiveHourUtilization, fiveHourResetTime,
		sevenDayUtilization, sevenDayResetTime)
}

// SendGreeting sends a greeting message to specified chat
func (c *Client) SendGreeting(chatID, text string) error {
	if !c.HasContext() {
		return fmt.Errorf("no context set (cookies not received from extension)")
	}

	if c.organizationID == "" {
		return fmt.Errorf("organization ID not set")
	}

	if chatID == "" {
		return fmt.Errorf("chat ID not specified")
	}

	curlPath := c.getCurlPath()

	// Build URL
	url := fmt.Sprintf("https://claude.ai/api/organizations/%s/chat_conversations/%s/completion",
		c.organizationID, chatID)

	// Build JSON body
	payload := map[string]string{"prompt": text}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	args := []string{
		"-X", "POST",
		url,
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Cookie: %s", c.cookies), // All cookies from browser
		"-d", string(payloadBytes),
	}

	// Add all browser headers to emulate real browser request
	for key, value := range c.headers {
		// Skip Content-Type as it's already added above
		// Skip Accept-Encoding because curl doesn't handle gzip automatically
		if key == "Content-Type" || key == "Accept-Encoding" {
			continue
		}
		// Add all headers including User-Agent from browser
		args = append(args, "-H", fmt.Sprintf("%s: %s", key, value))
	}

	// Add proxy if configured
	if c.proxy != "" {
		args = append([]string{"-x", c.proxy}, args...)
	}

	// Log greeting request
	log.Printf("========================================")
	log.Printf("GREETING Request:")
	log.Printf("  URL: %s", url)
	log.Printf("  Text: %s", text)
	log.Printf("  Chat ID: %s", chatID)
	log.Printf("========================================")

	cmd := exec.Command(curlPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Hide console window on Windows
	hideWindow(cmd)

	if err := cmd.Run(); err != nil {
		log.Printf("GREETING: CURL execution failed: %v", err)
		log.Printf("GREETING: CURL stderr: %s", stderr.String())
		log.Printf("GREETING: CURL stdout: %s", stdout.String())
		log.Printf("========================================")
		return fmt.Errorf("greeting request failed: %w, stderr: %s", err, stderr.String())
	}

	log.Printf("GREETING: Success!")
	log.Printf("GREETING: Response: %s", stdout.String())
	log.Printf("========================================")

	return nil
}
