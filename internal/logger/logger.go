package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	logFile            *os.File
	logger             *log.Logger
	fileLoggingEnabled bool
)

// Init initializes the logger with console output only
func Init() error {
	// Set up console logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	Info("===========================================")
	Info("=== ClaudeCompanion Starting ===")
	Info("===========================================")
	Info("OS: %s, Arch: %s", runtime.GOOS, runtime.GOARCH)
	Info("Go version: %s", runtime.Version())

	return nil
}

// SetFileLogging enables or disables file logging
func SetFileLogging(enabled bool) error {
	fileLoggingEnabled = enabled

	// Close existing log file if any
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}

	if !enabled {
		// Console only
		log.SetOutput(os.Stdout)
		logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		Info("File logging disabled - console only")
		return nil
	}

	// Enable file logging
	logPath, err := getLogPath()
	if err != nil {
		Error("Failed to get log path: %v", err)
		return fmt.Errorf("failed to get log path: %w", err)
	}

	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		Error("Failed to create log directory: %v", err)
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file (append mode)
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Error("Failed to open log file: %v", err)
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Check if stdout is available (not available in GUI mode on Windows)
	var output io.Writer
	if isStdoutAvailable() {
		// Write to both file and console
		output = io.MultiWriter(os.Stdout, logFile)
	} else {
		// GUI mode - write only to file
		output = logFile
	}

	logger = log.New(output, "", log.LstdFlags|log.Lshortfile)
	log.SetOutput(output)

	Info("File logging enabled: %s", logPath)
	return nil
}

// Close closes the log file
func Close() {
	if logFile != nil {
		Info("=== ClaudeCompanion Stopped ===")
		logFile.Close()
	}
}

// getLogPath returns the log file path (next to executable)
func getLogPath() (string, error) {
	// Get executable path
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get directory containing the executable
	exeDir := filepath.Dir(exePath)

	// Add date to log filename for rotation
	dateStr := time.Now().Format("2006-01-02")
	return filepath.Join(exeDir, fmt.Sprintf("claudecompanion-%s.log", dateStr)), nil
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[INFO] "+format, v...)
	if logger != nil {
		logger.Output(2, msg)
	} else {
		log.Output(2, msg)
	}
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[ERROR] "+format, v...)
	if logger != nil {
		logger.Output(2, msg)
	} else {
		log.Output(2, msg)
	}
}

// Warning logs a warning message
func Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf("[WARNING] "+format, v...)
	if logger != nil {
		logger.Output(2, msg)
	} else {
		log.Output(2, msg)
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf("[DEBUG] "+format, v...)
	if logger != nil {
		logger.Output(2, msg)
	} else {
		log.Output(2, msg)
	}
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf("[FATAL] "+format, v...)
	if logger != nil {
		logger.Output(2, msg)
	} else {
		log.Output(2, msg)
	}
	if logFile != nil {
		logFile.Close()
	}
	os.Exit(1)
}

// GetLogPath returns the current log file path
func GetLogPath() (string, error) {
	return getLogPath()
}

// isStdoutAvailable checks if stdout is available (not in GUI mode)
func isStdoutAvailable() bool {
	// Try to get file info for stdout
	_, err := os.Stdout.Stat()
	return err == nil
}
