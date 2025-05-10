// Package std provides a standard HTTP implementation of the HTTP server abstraction.
package std

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware"
)

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

// LoggingMiddleware is a standard HTTP implementation of core.ILoggingMiddleware.
type LoggingMiddleware struct {
	middleware.BaseLoggingMiddleware
}

// Middleware returns a middleware function that logs API requests for standard HTTP.
// This implementation can capture the actual status code set by the handler.
func (m *LoggingMiddleware) Middleware(config *core.LoggingConfig) core.HandlerFunc {
	if config == nil {
		config = middleware.DefaultLoggingConfig()
	}

	return func(c core.Context) {
		// Get the standard HTTP context
		stdContext, ok := c.(*Context)
		if !ok {
			// Handle the case when it's not a standard HTTP context
			// Get request path
			path := c.Request().URL.Path

			// Check if the path is in the skip paths list
			for _, skipPath := range config.SkipPaths {
				if path == skipPath {
					// Skip logging for this path
					c.Next()
					return
				}
			}

			// Start timer
			start := time.Now()

			// Get request details before processing
			req := c.Request()
			requestID := req.Header.Get("X-Request-ID")

			// Ensure request ID is in response
			if requestID == "" {
				requestID = fmt.Sprintf("%d", time.Now().UnixNano())
				c.SetHeader("X-Request-ID", requestID)
			} else {
				c.SetHeader("X-Request-ID", requestID)
			}

			// Continue with the next handler
			c.Next()

			// Calculate latency
			latency := time.Since(start).Milliseconds()

			// Create log entry
			logEntry := m.BaseLoggingMiddleware.CreateLogEntry(req, 200, latency, requestID, config)

			// Process the log
			m.BaseLoggingMiddleware.ProcessLog(logEntry, config)
			return
		}

		// Start timer
		start := time.Now()

		// Get request details before processing
		req := c.Request()
		requestID := req.Header.Get("X-Request-ID")

		// Ensure request ID is in response
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			c.SetHeader("X-Request-ID", requestID)
		} else {
			c.SetHeader("X-Request-ID", requestID)
		}

		// Store the original writer to restore it later
		originalWriter := stdContext.writer

		// Wrap the response writer to capture the status code
		wrappedWriter := &ResponseWriterWrapper{
			ResponseWriter: originalWriter,
			statusCode:     http.StatusOK,
		}

		// Replace the original writer with the wrapped one
		stdContext.writer = wrappedWriter

		// Continue with the next middleware/handler in the chain
		c.Next()

		// Calculate latency
		latency := time.Since(start).Milliseconds()

		// Get the status code from the wrapped writer
		statusCode := wrappedWriter.Status()

		// Create log entry with the actual status code
		logEntry := m.BaseLoggingMiddleware.CreateLogEntry(req, statusCode, latency, requestID, config)

		// Set error message based on status code
		if statusCode >= 400 {
			// For 4xx and 5xx status codes, set an error message
			logEntry.Error = fmt.Sprintf("HTTP error: %d", statusCode)
		}

		// Process the log
		m.BaseLoggingMiddleware.ProcessLog(logEntry, config)

		// Restore the original writer
		stdContext.writer = originalWriter
	}
}

// NewLoggingMiddleware creates a new LoggingMiddleware.
func NewLoggingMiddleware() core.ILoggingMiddleware {
	return &LoggingMiddleware{}
}
