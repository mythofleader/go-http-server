// Package middleware provides common middleware functionality for HTTP servers.
package middleware

import (
	"net/http"
	"strconv"

	"github.com/mythofleader/go-http-server/core"
)

// CORSConfig holds configuration for the CORS middleware.
type CORSConfig struct {
	// AllowedDomains is a list of domains that are allowed to access the API.
	// If empty, all domains are allowed.
	AllowedDomains []string

	// AllowedMethods is a list of HTTP methods that are allowed.
	// Default: "GET, POST, PUT, DELETE, OPTIONS, PATCH"
	AllowedMethods string

	// AllowedHeaders is a list of HTTP headers that are allowed.
	// Default: "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Requested-With"
	AllowedHeaders string

	// AllowCredentials indicates whether the request can include user credentials.
	// Default: true
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	// Default: 86400 (24 hours)
	MaxAge int
}

// DefaultCORSConfig returns a default CORS configuration.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedDomains:   []string{}, // Empty by default, which means all domains are allowed
		AllowedMethods:   "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowedHeaders:   "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Requested-With",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// NewDefaultCORSMiddleware returns a middleware function with default configuration.
// This function uses the DefaultCORSConfig which allows all domains.
// Example usage:
//
//	s.Use(middleware.NewDefaultCORSMiddleware())
//
// Or customize the configuration:
//
//	config := middleware.DefaultCORSConfig()
//	config.AllowedDomains = []string{"https://example.com"}
//	s.Use(middleware.CORSMiddleware(config))
func NewDefaultCORSMiddleware() core.HandlerFunc {
	return CORSMiddleware(DefaultCORSConfig())
}

// CORSMiddleware returns a middleware function that handles CORS (Cross-Origin Resource Sharing).
// If AllowedDomains is empty, all domains are allowed.
// If AllowedDomains contains specific domains, only those domains are allowed.
func CORSMiddleware(config *CORSConfig) core.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(c core.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			// Not a CORS request, continue with the next middleware/handler in the chain
			return
		}

		// Check if the origin is allowed
		allowOrigin := "*" // Default to allow all
		if len(config.AllowedDomains) > 0 {
			// Check if the origin is in the allowed domains list
			allowed := false
			for _, domain := range config.AllowedDomains {
				if domain == origin {
					allowed = true
					allowOrigin = origin // Set the specific origin
					break
				}
			}

			if !allowed {
				// Origin not allowed, continue without setting CORS headers
				return
			}
		}

		// Set CORS headers
		c.SetHeader("Access-Control-Allow-Origin", allowOrigin)
		c.SetHeader("Access-Control-Allow-Methods", config.AllowedMethods)
		c.SetHeader("Access-Control-Allow-Headers", config.AllowedHeaders)

		if config.AllowCredentials {
			c.SetHeader("Access-Control-Allow-Credentials", "true")
		}

		c.SetHeader("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

		// Handle preflight requests
		if c.Request().Method == "OPTIONS" {
			c.SetStatus(http.StatusOK)
			c.Abort()
			return
		}

		// Continue with the next middleware/handler in the chain
	}
}
