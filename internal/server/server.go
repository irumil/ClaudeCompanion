package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// ContextData represents the data received from browser extension
type ContextData struct {
	Cookies        string            `json:"cookies"`
	TargetURL      string            `json:"targetUrl"`
	OrganizationID string            `json:"organizationId"`
	Headers        map[string]string `json:"headers"` // Includes User-Agent
}

// Server handles HTTP requests from browser extension
type Server struct {
	port         int
	mu           sync.RWMutex
	contextData  *ContextData
	onContextSet func(cookies, targetURL, organizationID string, headers map[string]string)
	httpServer   *http.Server
}

// NewServer creates a new HTTP server
func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

// SetContextCallback sets the callback to be called when context is updated
func (s *Server) SetContextCallback(callback func(cookies, targetURL, organizationID string, headers map[string]string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onContextSet = callback
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/set-context", s.handleSetContext)
	mux.HandleFunc("/health", s.handleHealth)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", s.port),
		Handler: s.corsMiddleware(mux),
	}

	log.Printf("Starting HTTP server on port %d", s.port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}

// UpdatePort updates the server port (requires restart)
func (s *Server) UpdatePort(newPort int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.port != newPort {
		s.port = newPort
		log.Printf("Server port updated to %d (restart required)", newPort)
	}
}

// handleSetContext handles the /set-context endpoint
func (s *Server) handleSetContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data ContextData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("Failed to decode context data: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if data.Cookies == "" || data.TargetURL == "" {
		http.Error(w, "Missing cookies or targetUrl", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.contextData = &data
	callback := s.onContextSet
	s.mu.Unlock()

	log.Printf("Context received from extension: URL=%s, OrgID=%s, Cookies length=%d, Headers count=%d",
		data.TargetURL, data.OrganizationID, len(data.Cookies), len(data.Headers))

	// Call callback if set
	if callback != nil {
		callback(data.Cookies, data.TargetURL, data.OrganizationID, data.Headers)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Context updated successfully",
	})
}

// handleHealth handles the /health endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"version": "1.0.0",
	})
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from browser extensions
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetContext returns the current context data
func (s *Server) GetContext() *ContextData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.contextData
}

// HasContext returns true if context is set
func (s *Server) HasContext() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.contextData != nil && s.contextData.Cookies != ""
}
