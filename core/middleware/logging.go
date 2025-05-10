// Package middleware provides common middleware functionality for HTTP servers.
// This package contains default implementations and interfaces for middleware components.
// Framework-specific implementations of these middleware components can be found in their
// respective packages:
// - Gin implementation: github.com/tenqube/tenqube-go-http-server/core/gin
// - Standard HTTP implementation: github.com/tenqube/tenqube-go-http-server/core/std
package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mythofleader/go-http-server/core"
)

// ApiLog represents the structure of a log entry for API requests.
type ApiLog struct {
	ClientIp      string            `json:"client_ip"`
	Timestamp     string            `json:"timestamp"`
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Protocol      string            `json:"protocol"`
	StatusCode    int               `json:"status_code"`
	Latency       int64             `json:"latency"`
	UserAgent     string            `json:"user_agent"`
	Error         string            `json:"error"`
	RequestId     string            `json:"request_id"`
	Authorization string            `json:"authorization"`
	CustomFields  map[string]string `json:"custom_fields,omitempty"`
}

// DefaultLoggingConfig returns a default logging configuration.
func DefaultLoggingConfig() *core.LoggingConfig {
	return &core.LoggingConfig{
		RemoteURL:        "",
		CustomFields:     make(map[string]string),
		LoggingToConsole: true,  // Default to logging to console
		LoggingToRemote:  false, // Default to not logging to remote
		SkipPaths:        []string{},
	}
}

// NewDefaultConsoleLogging returns a logging configuration for console-only logging
// with the specified ignore path list and custom fields.
//
// Example usage:
//
//	config := middleware.NewDefaultConsoleLogging(
//		[]string{"/health", "/metrics"},
//		map[string]string{"version": "1.0.0", "environment": "production"}
//	)
//	s.Use(middleware.LoggingMiddleware(config))
func NewDefaultConsoleLogging(skipPaths []string, customFields map[string]string) *core.LoggingConfig {
	return &core.LoggingConfig{
		RemoteURL:        "",
		CustomFields:     customFields,
		LoggingToConsole: true,  // Enable console logging
		LoggingToRemote:  false, // Disable remote logging
		SkipPaths:        skipPaths,
	}
}

// BaseLoggingMiddleware provides common functionality for logging middleware implementations.
// This struct is embedded by framework-specific logging middleware implementations:
// - Gin implementation: github.com/tenqube/tenqube-go-http-server/core/gin.LoggingMiddleware
// - Standard HTTP implementation: github.com/tenqube/tenqube-go-http-server/core/std.LoggingMiddleware
// It provides methods for creating and processing log entries that are used by all implementations.
type BaseLoggingMiddleware struct{}

// CreateLogEntry creates a log entry from the request details.
func (m *BaseLoggingMiddleware) CreateLogEntry(req *http.Request, statusCode int, latency int64, requestID string, config *core.LoggingConfig) *ApiLog {
	clientIP := getClientIP(req)
	method := req.Method
	path := req.URL.Path
	protocol := req.Proto
	userAgent := req.UserAgent()
	authorization := req.Header.Get("Authorization")

	// Determine whether to mask authorization based on LoggingToConsole
	// If logging to console, we don't mask for easier debugging
	maskAuth := !config.LoggingToConsole

	return &ApiLog{
		ClientIp:      clientIP,
		Timestamp:     time.Now().Format(time.RFC3339),
		Method:        method,
		Path:          path,
		Protocol:      protocol,
		StatusCode:    statusCode,
		Latency:       latency,
		UserAgent:     userAgent,
		Error:         "none", // Default value, can be overridden by framework-specific implementation
		RequestId:     requestID,
		Authorization: maskAuthorizationBool(authorization, maskAuth),
		CustomFields:  config.CustomFields,
	}
}

// ProcessLog logs the entry to the console and sends it to the remote URL if configured.
func (m *BaseLoggingMiddleware) ProcessLog(logEntry *ApiLog, config *core.LoggingConfig) {
	// Log to console if LoggingToConsole is true
	if config.LoggingToConsole {
		logToConsole(logEntry)
	}

	// Send to remote URL if LoggingToRemote is true and RemoteURL is configured
	if config.LoggingToRemote && config.RemoteURL != "" {
		go sendLogToRemote(config.RemoteURL, logEntry)
	}
}

// ResponseWriterWrapper is a wrapper for http.ResponseWriter that captures the status code.
type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code and calls the underlying ResponseWriter's WriteHeader.
func (w *ResponseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the status code (if not already set) and calls the underlying ResponseWriter's Write.
func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	if !w.written {
		w.statusCode = http.StatusOK
		w.written = true
	}
	return w.ResponseWriter.Write(b)
}

// Status returns the captured status code.
func (w *ResponseWriterWrapper) Status() int {
	if !w.written {
		return http.StatusOK
	}
	return w.statusCode
}

// Hijack implements the http.Hijacker interface to pass through to the underlying ResponseWriter.
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
}

// Flush implements the http.Flusher interface to pass through to the underlying ResponseWriter.
func (w *ResponseWriterWrapper) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// CloseNotify implements the http.CloseNotifier interface to pass through to the underlying ResponseWriter.
func (w *ResponseWriterWrapper) CloseNotify() <-chan bool {
	if closeNotifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return closeNotifier.CloseNotify()
	}
	return nil
}

// Push implements the http.Pusher interface to pass through to the underlying ResponseWriter.
func (w *ResponseWriterWrapper) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("underlying ResponseWriter does not implement http.Pusher")
}

// getClientIP extracts the client IP address from the request.
func getClientIP(req *http.Request) string {
	// Try X-Forwarded-For header first
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try X-Real-IP header
	if xrip := req.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fall back to RemoteAddr
	ip := req.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// maskAuthorizationBool masks the authorization token for security.
// If maskAuth is false, the token is not masked.
func maskAuthorizationBool(auth string, maskAuth bool) string {
	if auth == "" {
		return ""
	}

	// If maskAuth is false, don't mask the token
	if !maskAuth {
		return auth
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "[MASKED]"
	}

	// Keep the auth type (e.g., "Bearer") but mask the token
	return parts[0] + " [MASKED]"
}

// logToConsole logs the API request to the console.
func logToConsole(logEntry *ApiLog) {
	jsonData, err := json.MarshalIndent(logEntry, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling log entry: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// sendLogToRemote sends the log entry to a remote URL.
func sendLogToRemote(url string, logEntry *ApiLog) {
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Printf("Error marshaling log entry: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending log to remote URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("Remote logging server returned error status: %d\n", resp.StatusCode)
	}
}
