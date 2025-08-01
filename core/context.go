// Package core provides the core interfaces and types for the HTTP server abstraction.
package core

import (
	"context"
	"net/http"
)

// FrameworkType represents the type of HTTP framework to use.
type FrameworkType string

const (
	// FrameworkGin represents the Gin framework.
	FrameworkGin FrameworkType = "gin"
	// FrameworkStdHTTP represents the standard net/http package.
	FrameworkStdHTTP FrameworkType = "std"
)

// HttpMethod represents an HTTP method.
type HttpMethod string

const (
	// GET represents the HTTP GET method.
	GET HttpMethod = "GET"
	// POST represents the HTTP POST method.
	POST HttpMethod = "POST"
	// PUT represents the HTTP PUT method.
	PUT HttpMethod = "PUT"
	// DELETE represents the HTTP DELETE method.
	DELETE HttpMethod = "DELETE"
	// PATCH represents the HTTP PATCH method.
	PATCH HttpMethod = "PATCH"
)

// HandlerFunc is a function that handles an HTTP request.
type HandlerFunc func(c Context)

// Context represents the context of an HTTP request.
// It abstracts away the underlying framework's context.
type Context interface {
	// Request returns the underlying HTTP request.
	Request() *http.Request
	// Writer returns the underlying ResponseWriter.
	Writer() http.ResponseWriter
	// Param returns the value of the URL param.
	Param(key string) string
	// Query returns the value of the URL query parameter.
	Query(key string) string
	// DefaultQuery returns the value of the URL query parameter or the default value.
	DefaultQuery(key, defaultValue string) string
	// GetHeader returns the value of the request header.
	GetHeader(key string) string
	// SetHeader sets a response header.
	SetHeader(key, value string)
	// SetStatus sets the HTTP response status code.
	SetStatus(code int)
	// JSON serializes the given struct as JSON into the response body.
	JSON(code int, obj interface{})
	// String writes the given string into the response body.
	String(code int, format string, values ...interface{})
	// Bind binds the request body into the given struct.
	Bind(obj interface{}) error
	// BindJSON binds the JSON request body into the given struct.
	BindJSON(obj interface{}) error
	// ShouldBindJSON binds the JSON request body into the given struct.
	// If there is an error, it returns the error without aborting the request.
	ShouldBindJSON(obj interface{}) error
	// File serves a file.
	File(filepath string)
	// Redirect redirects the request to the given URL.
	Redirect(code int, location string)
	// Error adds an error to the context.
	// This is used by the error handler middleware to handle errors.
	Error(err error) error
	// Errors returns all errors added to the context.
	// This is used to retrieve all errors that occurred during request processing.
	Errors() []error
	// Next calls the next handler in the chain.
	// This is used for middleware flow control.
	Next()
	// Abort prevents pending handlers in the chain from being called.
	// This is used to stop the middleware chain execution.
	Abort()
	// Get returns the value for the given key and a boolean indicating whether the key exists.
	// This is used to retrieve values stored in the context.
	Get(key string) (interface{}, bool)
	// Set stores a value in the context for the given key.
	// This is used to store values in the context.
	Set(key string, value interface{})
}

// ILoggingMiddleware is an interface for logging middleware implementations.
// Each framework (Gin, StdHTTP) provides its own implementation of this interface:
// - Gin implementation: github.com/mythofleader/go-http-server/core/gin.LoggingMiddleware
// - Standard HTTP implementation: github.com/mythofleader/go-http-server/core/std.LoggingMiddleware
type ILoggingMiddleware interface {
	// Middleware returns a middleware function that logs API requests.
	Middleware(config *LoggingConfig) HandlerFunc
}

// IErrorHandlerMiddleware is an interface for error handler middleware implementations.
// Each framework (Gin, StdHTTP) provides its own implementation of this interface:
// - Gin implementation: github.com/mythofleader/go-http-server/core/gin.ErrorHandlerMiddleware
// - Standard HTTP implementation: github.com/mythofleader/go-http-server/core/std.ErrorHandlerMiddleware
type IErrorHandlerMiddleware interface {
	// Middleware returns a middleware function that handles errors.
	Middleware(config *ErrorHandlerConfig) HandlerFunc
}

// ErrorHandlerConfig holds configuration for the error handler middleware.
type ErrorHandlerConfig struct {
	// DefaultErrorMessage is the message to use for non-HTTP errors.
	DefaultErrorMessage string
	// DefaultStatusCode is the status code to use for non-HTTP errors.
	DefaultStatusCode int
}

// LoggingConfig holds configuration for the logging middleware.
type LoggingConfig struct {
	RemoteURL        string
	CustomFields     map[string]string
	LoggingToConsole bool     // Whether to log to console
	LoggingToRemote  bool     // Whether to log to remote
	SkipPaths        []string // List of paths to ignore for logging
}

// Controller is an interface for defining routes.
type Controller interface {
	// GetHttpMethod returns the HTTP method for the route
	GetHttpMethod() HttpMethod
	// GetPath returns the path for the route
	GetPath() string
	// Handler returns handler functions for the route
	Handler() []HandlerFunc
	// SkipLogging returns whether to skip logging for this controller
	SkipLogging() bool
	// SkipAuthCheck returns whether to skip authentication checks for this controller
	SkipAuthCheck() bool
}

// Server is an interface for HTTP servers.
// It abstracts away the underlying framework.
type Server interface {
	// GET registers a route for GET requests
	GET(path string, handlers ...HandlerFunc)
	// POST registers a route for POST requests
	POST(path string, handlers ...HandlerFunc)
	// PUT registers a route for PUT requests
	PUT(path string, handlers ...HandlerFunc)
	// DELETE registers a route for DELETE requests
	DELETE(path string, handlers ...HandlerFunc)
	// PATCH registers a route for PATCH requests
	PATCH(path string, handlers ...HandlerFunc)
	// Group creates a new router group
	Group(path string) RouterGroup
	// Use adds middleware to the server
	Use(middleware ...HandlerFunc)
	// RegisterRouter registers routes from Controller objects
	RegisterRouter(controllers ...Controller)
	// NoRoute registers handlers for 404 Not Found errors
	NoRoute(handlers ...HandlerFunc)
	// NoMethod registers handlers for 405 Method Not Allowed errors
	NoMethod(handlers ...HandlerFunc)
	// Run starts the server
	Run() error
	// Stop stops the server immediately
	Stop() error
	// RunTLS starts the server with TLS
	RunTLS(addr, certFile, keyFile string) error
	// Shutdown gracefully shuts down the server
	Shutdown(ctx context.Context) error
	// GetLoggingMiddleware returns a framework-specific logging middleware
	GetLoggingMiddleware() ILoggingMiddleware
	// GetErrorHandlerMiddleware returns a framework-specific error handler middleware
	GetErrorHandlerMiddleware() IErrorHandlerMiddleware
	// StartLambda starts the server in AWS Lambda mode.
	// This method should be called instead of Run or RunTLS when running in AWS Lambda.
	// It returns an error if the framework does not support Lambda.
	StartLambda() error
	// GetPort returns the port the server is configured to run on.
	// This is useful when using random ports.
	GetPort() string
}

// RouterGroup is a group of routes.
type RouterGroup interface {
	// GET registers a route for GET requests
	GET(path string, handlers ...HandlerFunc)
	// POST registers a route for POST requests
	POST(path string, handlers ...HandlerFunc)
	// PUT registers a route for PUT requests
	PUT(path string, handlers ...HandlerFunc)
	// DELETE registers a route for DELETE requests
	DELETE(path string, handlers ...HandlerFunc)
	// PATCH registers a route for PATCH requests
	PATCH(path string, handlers ...HandlerFunc)
	// Group creates a new router group
	Group(path string) RouterGroup
	// Use adds middleware to the group
	Use(middleware ...HandlerFunc)
	// RegisterRouter registers routes from Controller objects
	RegisterRouter(controllers ...Controller)
}
