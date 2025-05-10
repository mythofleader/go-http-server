// Package middleware provides common middleware functionality for HTTP servers.
// This package contains default implementations and interfaces for middleware components.
// Framework-specific implementations of these middleware components can be found in their
// respective packages:
// - Gin implementation: github.com/tenqube/tenqube-go-http-server/core/gin
// - Standard HTTP implementation: github.com/tenqube/tenqube-go-http-server/core/std
package middleware

import (
	"errors"
	"net/http"

	"github.com/mythofleader/go-http-server/core"
	tErrors "github.com/mythofleader/go-http-server/core/middleware/errors"
)

// DefaultErrorHandlerConfig returns a default error handler configuration.
func DefaultErrorHandlerConfig() *core.ErrorHandlerConfig {
	return &core.ErrorHandlerConfig{
		DefaultErrorMessage: "Internal Server Error",
		DefaultStatusCode:   http.StatusInternalServerError,
	}
}

// errorHandlerContext is a wrapper for core.Context that catches errors.
type errorHandlerContext struct {
	core.Context
	err error
}

// Error implements core.Context.Error
func (c *errorHandlerContext) Error(err error) error {
	c.err = err
	return c.Context.Error(err)
}

// Errors implements core.Context.Errors
func (c *errorHandlerContext) Errors() []error {
	return c.Context.Errors()
}

// handleError processes an error and returns an appropriate HTTP response.
func handleError(c core.Context, err error, config *core.ErrorHandlerConfig) {
	var httpErr tErrors.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode(), tErrors.NewErrorResponse(httpErr.StatusCode(), httpErr.Error()))
		return
	}
	c.JSON(config.DefaultStatusCode, tErrors.NewErrorResponse(config.DefaultStatusCode, config.DefaultErrorMessage))
}

// IErrorHandlerMiddleware is an interface for error handler middleware implementations.
// Each framework (Gin, StdHTTP) provides its own implementation of this interface:
// - Gin implementation: github.com/tenqube/tenqube-go-http-server/core/gin.ErrorHandlerMiddleware
// - Standard HTTP implementation: github.com/tenqube/tenqube-go-http-server/core/std.ErrorHandlerMiddleware
type IErrorHandlerMiddleware interface {
	// Middleware returns a middleware function that handles errors.
	Middleware(config *core.ErrorHandlerConfig) core.HandlerFunc
}
