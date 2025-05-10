// Package gin provides a Gin implementation of the HTTP server abstraction.
package gin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware"
)

// LoggingMiddleware is a Gin-specific implementation of core.ILoggingMiddleware.
// It works with the Gin framework (github.com/gin-gonic/gin).
type LoggingMiddleware struct {
	middleware.BaseLoggingMiddleware
	// This is just to make the linter happy about the gin import
	_ gin.HandlerFunc
}

// Middleware returns a middleware function that logs API requests for Gin.
// This implementation can capture the actual status code set by the handler.
func (m *LoggingMiddleware) Middleware(config *core.LoggingConfig) core.HandlerFunc {
	if config == nil {
		config = middleware.DefaultLoggingConfig()
	}

	return func(c core.Context) {
		// Get the Gin context
		ginContext, ok := c.(*Context)
		if !ok {
			// Handle the case when it's not a Gin context
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

		// Get the underlying gin.Context
		gc := ginContext.ginContext

		// Use Gin's built-in middleware to capture the status code
		gc.Next()

		// Calculate latency
		latency := time.Since(start).Milliseconds()

		// Get the status code from the Gin context
		statusCode := gc.Writer.Status()

		// Get error information if available
		var errorMsg string
		if len(gc.Errors) > 0 {
			errorMsg = gc.Errors.String()
		}

		// Create log entry with the actual status code
		logEntry := m.BaseLoggingMiddleware.CreateLogEntry(req, statusCode, latency, requestID, config)
		logEntry.Error = errorMsg

		// Process the log
		m.BaseLoggingMiddleware.ProcessLog(logEntry, config)
	}
}

// NewLoggingMiddleware creates a new LoggingMiddleware.
func NewLoggingMiddleware() core.ILoggingMiddleware {
	return &LoggingMiddleware{}
}
