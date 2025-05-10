// Package middleware provides common middleware functionality for HTTP servers.
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mythofleader/go-http-server/core"
)

// TimeoutConfig holds configuration for the timeout middleware.
type TimeoutConfig struct {
	// Timeout is the maximum duration to wait for a response.
	// If not set, it defaults to 2 seconds.
	Timeout time.Duration
}

// DefaultTimeoutConfig returns a default timeout configuration.
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Timeout: 2 * time.Second, // Default to 2 seconds
	}
}

// NewDefaultTimeoutMiddleware returns a middleware function with default configuration.
// This function uses the DefaultTimeoutConfig which sets a default timeout of 2 seconds.
// Example usage:
//
//	s.Use(middleware.NewDefaultTimeoutMiddleware())
//
// Or customize the configuration:
//
//	config := middleware.DefaultTimeoutConfig()
//	config.Timeout = 5 * time.Second
//	s.Use(middleware.TimeoutMiddleware(config))
func NewDefaultTimeoutMiddleware() core.HandlerFunc {
	return TimeoutMiddleware(DefaultTimeoutConfig())
}

// TimeoutMiddleware returns a middleware function that times out requests after a specified duration.
// If the handler doesn't respond within the timeout period, it returns a 503 Service Unavailable response.
func TimeoutMiddleware(config *TimeoutConfig) core.HandlerFunc {
	if config == nil {
		config = DefaultTimeoutConfig()
	}

	// Log middleware configuration
	log.Printf("[MIDDLEWARE] Timeout middleware configured:")
	log.Printf("[MIDDLEWARE]   - Timeout: %v", config.Timeout)

	return func(c core.Context) {
		// Create a channel to track if the response has been written
		responseSent := make(chan bool, 1)

		// Create a timeout channel
		timeoutCh := time.After(config.Timeout)

		// Get the original response writer
		originalWriter := c.Writer()

		// Create a goroutine to handle the timeout
		go func() {
			// Wait for the timeout
			<-timeoutCh

			// Check if a response has already been sent
			select {
			case <-responseSent:
				// Response already sent, do nothing
				return
			default:
				// No response sent yet, send timeout response
				originalWriter.WriteHeader(http.StatusServiceUnavailable)
				originalWriter.Write([]byte(fmt.Sprintf("Request timed out after %v", config.Timeout)))
				responseSent <- true
			}
		}()

		// Signal when the response is sent
		defer func() {
			select {
			case <-responseSent:
				// Response already sent by timeout handler
				return
			default:
				// Response sent by normal handler
				responseSent <- true
			}
		}()

		// Continue with the next middleware/handler in the chain
		// This will execute the actual request handler
		c.Next()
	}
}
