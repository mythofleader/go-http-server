// Package middleware provides common middleware functionality for HTTP servers.
package middleware

import (
	"context"
	"net/http"

	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware/errors"
)

// RequestIDGenerator defines the interface for generating request IDs
type RequestIDGenerator interface {
	// GenerateRequestID generates a unique request ID from the context
	GenerateRequestID(ctx context.Context) (string, error)
}

// RequestIDStorage defines the interface for checking and storing request IDs
type RequestIDStorage interface {
	// CheckRequestID checks if a request ID exists in the storage
	CheckRequestID(requestID string) (bool, error)

	// SaveRequestID saves a request ID to the storage
	SaveRequestID(requestID string) error
}

// DuplicateRequestConfig holds configuration for the duplicate request prevention middleware
type DuplicateRequestConfig struct {
	// RequestIDGenerator is the implementation of RequestIDGenerator
	RequestIDGenerator RequestIDGenerator

	// RequestIDStorage is the implementation of RequestIDStorage
	RequestIDStorage RequestIDStorage

	// Optional: custom error message
	ConflictMessage string
}

// DefaultDuplicateRequestConfig returns a default duplicate request configuration
func DefaultDuplicateRequestConfig() *DuplicateRequestConfig {
	return &DuplicateRequestConfig{
		ConflictMessage: "Duplicate request detected",
		// RequestIDGenerator and RequestIDStorage are nil by default
		// and must be provided by the user
	}
}

// NewDefaultDuplicateRequestMiddleware returns a middleware function with default configuration.
// Note: This function panics because DuplicateRequestMiddleware requires additional configuration:
// - RequestIDGenerator must be provided
// - RequestIDStorage must be provided
// You must set these fields in the configuration before using this middleware.
// Example usage:
//
//	config := middleware.DefaultDuplicateRequestConfig()
//	config.RequestIDGenerator = myRequestIDGenerator
//	config.RequestIDStorage = myRequestIDStorage
//	s.Use(middleware.DuplicateRequestMiddleware(config))
//
// Or use the DuplicateRequestMiddleware function directly:
//
//	s.Use(middleware.DuplicateRequestMiddleware(&middleware.DuplicateRequestConfig{
//		RequestIDGenerator: myRequestIDGenerator,
//		RequestIDStorage:   myRequestIDStorage,
//		ConflictMessage:    "Custom conflict message",
//	}))
func NewDefaultDuplicateRequestMiddleware() core.HandlerFunc {
	return DuplicateRequestMiddleware(DefaultDuplicateRequestConfig())
}

// DuplicateRequestMiddleware returns a middleware function that prevents duplicate requests
// It generates a request ID using the provided generator, checks if it exists in the storage,
// and if it does, returns a 409 Conflict response. Otherwise, it saves the ID and continues.
func DuplicateRequestMiddleware(config *DuplicateRequestConfig) core.HandlerFunc {
	if config == nil {
		config = DefaultDuplicateRequestConfig()
	}

	// Validate the configuration
	if config.RequestIDGenerator == nil {
		panic("DuplicateRequestMiddleware requires a RequestIDGenerator implementation")
	}

	if config.RequestIDStorage == nil {
		panic("DuplicateRequestMiddleware requires a RequestIDStorage implementation")
	}

	return func(c core.Context) {
		// Get the request context
		ctx := c.Request().Context()

		// Generate a request ID
		requestID, err := config.RequestIDGenerator.GenerateRequestID(ctx)
		if err != nil {
			// If we can't generate a request ID, return an internal server error
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerErrorResponse("Failed to generate request ID"))
			return
		}

		// Check if the request ID exists in the storage
		exists, err := config.RequestIDStorage.CheckRequestID(requestID)
		if err != nil {
			// If we can't check the request ID, return an internal server error
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerErrorResponse("Failed to check request ID"))
			return
		}

		// If the request ID exists, return a conflict error
		if exists {
			c.JSON(http.StatusConflict, errors.NewConflictResponse(config.ConflictMessage))
			return
		}

		// Save the request ID to the storage
		if err := config.RequestIDStorage.SaveRequestID(requestID); err != nil {
			// If we can't save the request ID, return an internal server error
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerErrorResponse("Failed to save request ID"))
			return
		}

		// Continue with the next middleware/handler in the chain
		c.Next()
	}
}
