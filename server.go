// Package server provides an abstraction layer for HTTP servers.
// It wraps popular frameworks like Gin to provide a consistent API.
package server

import (
	"fmt"

	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/gin"
	"github.com/mythofleader/go-http-server/core/middleware"
	"github.com/mythofleader/go-http-server/core/middleware/errors"
	"github.com/mythofleader/go-http-server/core/std"
)

// Re-export types from core package
type (
	// Server is an interface for HTTP servers.
	Server = core.Server
	// Context represents the context of an HTTP request.
	Context = core.Context
	// FrameworkType represents the type of HTTP framework to use.
	FrameworkType = core.FrameworkType
	// HandlerFunc is a function that handles an HTTP request.
	HandlerFunc = core.HandlerFunc
	// RouterGroup is a group of routes.
	RouterGroup = core.RouterGroup
	// LoggingConfig holds configuration for the logging middleware.
	LoggingConfig = core.LoggingConfig
	// ErrorHandlerConfig holds configuration for the error handler middleware.
	ErrorHandlerConfig = core.ErrorHandlerConfig
	// HttpMethod represents an HTTP method.
	HttpMethod = core.HttpMethod
)

// Re-export types from middleware package
type (
	// TimeoutConfig holds configuration for the timeout middleware.
	TimeoutConfig = middleware.TimeoutConfig
	// AuthConfig holds configuration for the authorization middleware.
	AuthConfig = middleware.AuthConfig
	// APIKeyConfig holds configuration for the API key middleware.
	APIKeyConfig = middleware.APIKeyConfig
	// CORSConfig holds configuration for the CORS middleware.
	CORSConfig = middleware.CORSConfig
	// DuplicateRequestConfig holds configuration for the duplicate request prevention middleware.
	DuplicateRequestConfig = middleware.DuplicateRequestConfig
	// RequestIDGenerator defines the interface for generating request IDs.
	RequestIDGenerator = middleware.RequestIDGenerator
	// RequestIDStorage defines the interface for checking and storing request IDs.
	RequestIDStorage = middleware.RequestIDStorage
	// BasicAuthUserLookup defines the interface for looking up users based on Basic Auth credentials.
	BasicAuthUserLookup = middleware.BasicAuthUserLookup
	// JWTUserLookup defines the interface for looking up users based on JWT claims.
	JWTUserLookup = middleware.JWTUserLookup
	// MapClaims represents JWT claims as a map.
	MapClaims = middleware.MapClaims
	// AuthType represents the type of authentication to use.
	AuthType = middleware.AuthType
)

// Re-export types from middleware/errors package
type (
	// ErrorDetail represents the structure of an error detail in the response.
	ErrorDetail = errors.ErrorDetail
	// ErrorResponse represents the structure of an error response.
	ErrorResponse = errors.ErrorResponse

	// Error structs that embed the error interface
	// BadRequestHttpError represents a 400 Bad Request error.
	BadRequestHttpError = errors.BadRequestHttpError
	// UnauthorizedHttpError represents a 401 Unauthorized error.
	UnauthorizedHttpError = errors.UnauthorizedHttpError
	// ForbiddenHttpError represents a 403 Forbidden error.
	ForbiddenHttpError = errors.ForbiddenHttpError
	// NotFoundHttpError represents a 404 Not Found error.
	NotFoundHttpError = errors.NotFoundHttpError
	// MethodNotAllowedHttpError represents a 405 Method Not Allowed error.
	MethodNotAllowedHttpError = errors.MethodNotAllowedHttpError
	// InternalServerHttpError represents a 500 Internal Server Error.
	InternalServerHttpError = errors.InternalServerHttpError
	// ServiceUnavailableHttpError represents a 503 Service Unavailable error.
	ServiceUnavailableHttpError = errors.ServiceUnavailableHttpError
)

// Re-export constants from core package
const (
	// FrameworkGin represents the Gin framework.
	FrameworkGin = core.FrameworkGin
	// FrameworkStdHTTP represents the standard net/http package.
	FrameworkStdHTTP = core.FrameworkStdHTTP

	// HTTP methods
	// GET represents the HTTP GET method.
	GET = core.GET
	// POST represents the HTTP POST method.
	POST = core.POST
	// PUT represents the HTTP PUT method.
	PUT = core.PUT
	// DELETE represents the HTTP DELETE method.
	DELETE = core.DELETE
	// PATCH represents the HTTP PATCH method.
	PATCH = core.PATCH
)

// Re-export constants from middleware package
const (
	// AuthTypeBasic represents HTTP Basic authentication.
	AuthTypeBasic = middleware.AuthTypeBasic
	// AuthTypeJWT represents JWT Bearer token authentication.
	AuthTypeJWT = middleware.AuthTypeJWT
)

// Re-export types from gin package
type (
	// GinServer is an implementation of Server using the Gin framework.
	GinServer = gin.Server
)

// Re-export types from std package
type (
	// StdServer is an implementation of Server using the standard net/http package.
	StdServer = std.Server
)

// Re-export functions from middleware package
var (
	// TimeoutMiddleware returns a middleware function that times out requests after a specified duration.
	TimeoutMiddleware = middleware.TimeoutMiddleware
	// AuthMiddleware returns a middleware function that checks authorization.
	AuthMiddleware = middleware.AuthMiddleware
	// APIKeyMiddleware returns a middleware function that checks for a valid API key.
	APIKeyMiddleware = middleware.APIKeyMiddleware
	// CORSMiddleware returns a middleware function that handles CORS (Cross-Origin Resource Sharing).
	CORSMiddleware = middleware.CORSMiddleware
	// DuplicateRequestMiddleware returns a middleware function that prevents duplicate requests.
	DuplicateRequestMiddleware = middleware.DuplicateRequestMiddleware
	// GetUserFromContext retrieves the authenticated user from the context.
	GetUserFromContext = middleware.GetUserFromContext

	// NewDefaultAPIKeyMiddleware returns a middleware function with default configuration and the specified API key.
	NewDefaultAPIKeyMiddleware = middleware.NewDefaultAPIKeyMiddleware
	// NewDefaultJWTAuthMiddleware returns a middleware function with default JWT authentication configuration.
	NewDefaultJWTAuthMiddleware = middleware.NewDefaultJWTAuthMiddleware
	// NewDefaultBasicAuthMiddleware returns a middleware function with default Basic authentication configuration.
	NewDefaultBasicAuthMiddleware = middleware.NewDefaultBasicAuthMiddleware
	// NewDefaultCORSMiddleware returns a middleware function with default configuration.
	NewDefaultCORSMiddleware = middleware.NewDefaultCORSMiddleware
	// NewDefaultDuplicateRequestMiddleware returns a middleware function with default configuration.
	NewDefaultDuplicateRequestMiddleware = middleware.NewDefaultDuplicateRequestMiddleware
	// NewDefaultConsoleLogging returns a logging configuration for console-only logging with the specified ignore path list and custom fields.
	NewDefaultConsoleLogging = middleware.NewDefaultConsoleLogging
	// NewDefaultTimeoutMiddleware returns a middleware function with default configuration.
	NewDefaultTimeoutMiddleware = middleware.NewDefaultTimeoutMiddleware
)

// Re-export functions from middleware/errors package
var (
	// NewErrorResponse creates a new ErrorResponse with the given status code and message.
	NewErrorResponse = errors.NewErrorResponse
	// NewBadRequestResponse creates a new ErrorResponse for a 400 Bad Request error.
	NewBadRequestResponse = errors.NewBadRequestResponse
	// NewUnauthorizedResponse creates a new ErrorResponse for a 401 Unauthorized error.
	NewUnauthorizedResponse = errors.NewUnauthorizedResponse
	// NewForbiddenResponse creates a new ErrorResponse for a 403 Forbidden error.
	NewForbiddenResponse = errors.NewForbiddenResponse
	// NewNotFoundResponse creates a new ErrorResponse for a 404 Not Found error.
	NewNotFoundResponse = errors.NewNotFoundResponse
	// NewConflictResponse creates a new ErrorResponse for a 409 Conflict error.
	NewConflictResponse = errors.NewConflictResponse
	// NewInternalServerErrorResponse creates a new ErrorResponse for a 500 Internal Server Error.
	NewInternalServerErrorResponse = errors.NewInternalServerErrorResponse
	// NewServiceUnavailableResponse creates a new ErrorResponse for a 503 Service Unavailable error.
	NewServiceUnavailableResponse = errors.NewServiceUnavailableResponse

	// Constructor functions for the error structs
	// NewBadRequestHttpError creates a new BadRequestHttpError.
	NewBadRequestHttpError = errors.NewBadRequestHttpError
	// NewUnauthorizedHttpError creates a new UnauthorizedHttpError.
	NewUnauthorizedHttpError = errors.NewUnauthorizedHttpError
	// NewForbiddenHttpError creates a new ForbiddenHttpError.
	NewForbiddenHttpError = errors.NewForbiddenHttpError
	// NewNotFoundHttpError creates a new NotFoundHttpError.
	NewNotFoundHttpError = errors.NewNotFoundHttpError
	// NewMethodNotAllowedHttpError creates a new MethodNotAllowedHttpError.
	NewMethodNotAllowedHttpError = errors.NewMethodNotAllowedHttpError
	// NewInternalServerHttpError creates a new InternalServerHttpError.
	NewInternalServerHttpError = errors.NewInternalServerHttpError
	// NewServiceUnavailableHttpError creates a new ServiceUnavailableHttpError.
	NewServiceUnavailableHttpError = errors.NewServiceUnavailableHttpError
)

// NewServer creates a new Server instance.
// By default, it uses the Gin framework if no framework type is specified.
// If port is not provided, it defaults to "8080".
// If showFrameworkLogs is true, logs about the framework, middleware, and routes will be printed to the console.
// If showFrameworkLogs is false, these logs will be suppressed.
func NewServer(frameworkType core.FrameworkType, port string, showFrameworkLogs bool) (core.Server, error) {
	// Default port to "8080" if not provided
	if port == "" {
		port = "8080"
	}

	// Use the specified framework
	switch frameworkType {
	case core.FrameworkGin:
		return gin.NewServer(port, showFrameworkLogs), nil
	case core.FrameworkStdHTTP:
		return std.NewServer(port, showFrameworkLogs), nil
	default:
		return nil, fmt.Errorf("unsupported framework type: %s", frameworkType)
	}
}
