// Package middleware provides common middleware functionality for HTTP servers.
package middleware

import (
	"net/http"

	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware/errors"
)

// APIKeyConfig holds configuration for the API key middleware.
type APIKeyConfig struct {
	// APIKey is the expected API key value.
	// This value will be compared against the x-api-key header.
	APIKey string

	// Optional: custom error message
	UnauthorizedMessage string
}

// DefaultAPIKeyConfig returns a default API key configuration.
func DefaultAPIKeyConfig() *APIKeyConfig {
	return &APIKeyConfig{
		APIKey:              "", // Empty by default, must be provided
		UnauthorizedMessage: "Unauthorized: Invalid or missing API key",
	}
}

// NewDefaultAPIKeyMiddleware returns a middleware function with default configuration and the specified API key.
// This function creates a default configuration and sets the APIKey to the provided value.
// Example usage:
//
//	s.Use(middleware.NewDefaultAPIKeyMiddleware("your-api-key"))
//
// Or customize the configuration:
//
//	config := middleware.DefaultAPIKeyConfig()
//	config.APIKey = "your-api-key"
//	config.UnauthorizedMessage = "Custom unauthorized message"
//	s.Use(middleware.APIKeyMiddleware(config))
//
// Or use the APIKeyMiddleware function directly:
//
//	s.Use(middleware.APIKeyMiddleware(&middleware.APIKeyConfig{
//		APIKey: "your-api-key",
//	}))
func NewDefaultAPIKeyMiddleware(apiKey string) core.HandlerFunc {
	config := DefaultAPIKeyConfig()
	config.APIKey = apiKey
	return APIKeyMiddleware(config)
}

// APIKeyMiddleware returns a middleware function that checks for a valid API key in the x-api-key header.
// If the API key is missing or invalid, it returns a 401 Unauthorized response.
func APIKeyMiddleware(config *APIKeyConfig) core.HandlerFunc {
	if config == nil {
		config = DefaultAPIKeyConfig()
	}

	// Ensure API key is provided
	if config.APIKey == "" {
		panic("APIKeyMiddleware requires a non-empty APIKey in the configuration")
	}

	return func(c core.Context) {
		// Get the x-api-key header
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.SetStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedResponse(config.UnauthorizedMessage))
			return
		}

		// Validate the API key
		if apiKey != config.APIKey {
			c.SetStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedResponse(config.UnauthorizedMessage))
			return
		}

		// API key is valid, continue with the next middleware/handler in the chain
	}
}
