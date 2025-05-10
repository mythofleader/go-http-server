// Package server provides an abstraction layer for HTTP servers.
// It wraps popular frameworks like Gin to provide a consistent API.
package server

import (
	"github.com/mythofleader/go-http-server/core"
)

// ServerBuilder is a builder for creating a server with controllers and middleware.
type ServerBuilder struct {
	frameworkType    core.FrameworkType
	port             string
	controllers      []core.Controller
	middleware       []core.HandlerFunc
	loggingConfig    *core.LoggingConfig
	timeoutConfig    *TimeoutConfig
	corsConfig       *CORSConfig
	errorConfig      *core.ErrorHandlerConfig
	noRouteHandlers  []core.HandlerFunc // Handlers for 404 Not Found errors
	noMethodHandlers []core.HandlerFunc // Handlers for 405 Method Not Allowed errors

	// Flags for default middleware
	useDefaultLogging      bool
	useDefaultTimeout      bool
	useDefaultCORS         bool
	useDefaultErrorHandler bool
}

// NewServerBuilder creates a new ServerBuilder with the specified framework type and port.
func NewServerBuilder(frameworkType core.FrameworkType, port string) *ServerBuilder {
	return &ServerBuilder{
		frameworkType:    frameworkType,
		port:             port,
		controllers:      make([]core.Controller, 0),
		middleware:       make([]core.HandlerFunc, 0),
		noRouteHandlers:  make([]core.HandlerFunc, 0),
		noMethodHandlers: make([]core.HandlerFunc, 0),
	}
}

// NewGinServerBuilder creates a new ServerBuilder with the Gin framework and port 8080.
// This is a convenience function that doesn't require any arguments.
func NewGinServerBuilder() *ServerBuilder {
	return NewServerBuilder(core.FrameworkGin, "8080")
}

// AddController adds a controller to the builder.
func (b *ServerBuilder) AddController(controller core.Controller) *ServerBuilder {
	b.controllers = append(b.controllers, controller)
	return b
}

// AddControllers adds multiple controllers to the builder.
func (b *ServerBuilder) AddControllers(controllers ...core.Controller) *ServerBuilder {
	b.controllers = append(b.controllers, controllers...)
	return b
}

// AddMiddleware adds a middleware to the builder.
func (b *ServerBuilder) AddMiddleware(middleware core.HandlerFunc) *ServerBuilder {
	b.middleware = append(b.middleware, middleware)
	return b
}

// AddMiddlewares adds multiple middleware to the builder.
func (b *ServerBuilder) AddMiddlewares(middleware ...core.HandlerFunc) *ServerBuilder {
	b.middleware = append(b.middleware, middleware...)
	return b
}

// WithLogging configures the logging middleware with the specified custom fields.
func (b *ServerBuilder) WithLogging(customFields map[string]string) *ServerBuilder {
	b.loggingConfig = &core.LoggingConfig{
		RemoteURL:        "",
		CustomFields:     customFields,
		LoggingToConsole: true,
		LoggingToRemote:  false,
		SkipPaths:        []string{},
	}
	return b
}

// WithRemoteLogging configures the logging middleware with the specified remote URL and custom fields.
func (b *ServerBuilder) WithRemoteLogging(remoteURL string, customFields map[string]string) *ServerBuilder {
	b.loggingConfig = &core.LoggingConfig{
		RemoteURL:        remoteURL,
		CustomFields:     customFields,
		LoggingToConsole: true,
		LoggingToRemote:  true,
		SkipPaths:        []string{},
	}
	return b
}

// WithTimeout configures the timeout middleware with the specified timeout.
func (b *ServerBuilder) WithTimeout(timeout TimeoutConfig) *ServerBuilder {
	b.timeoutConfig = &timeout
	return b
}

// WithCORS configures the CORS middleware with the specified configuration.
func (b *ServerBuilder) WithCORS(cors CORSConfig) *ServerBuilder {
	b.corsConfig = &cors
	return b
}

// WithErrorHandler configures the error handler middleware with the specified configuration.
func (b *ServerBuilder) WithErrorHandler(errorConfig core.ErrorHandlerConfig) *ServerBuilder {
	b.errorConfig = &errorConfig
	return b
}

// WithDefaultLogging enables the default logging middleware.
func (b *ServerBuilder) WithDefaultLogging() *ServerBuilder {
	b.useDefaultLogging = true
	return b
}

// WithDefaultTimeout enables the default timeout middleware.
func (b *ServerBuilder) WithDefaultTimeout() *ServerBuilder {
	b.useDefaultTimeout = true
	return b
}

// WithDefaultCORS enables the default CORS middleware.
func (b *ServerBuilder) WithDefaultCORS() *ServerBuilder {
	b.useDefaultCORS = true
	return b
}

// WithDefaultErrorHandling enables the default error handler middleware.
func (b *ServerBuilder) WithDefaultErrorHandling() *ServerBuilder {
	b.useDefaultErrorHandler = true
	return b
}

// WithNoRoute configures custom handlers for 404 Not Found errors.
func (b *ServerBuilder) WithNoRoute(handlers ...core.HandlerFunc) *ServerBuilder {
	b.noRouteHandlers = handlers
	return b
}

// WithNoMethod configures custom handlers for 405 Method Not Allowed errors.
func (b *ServerBuilder) WithNoMethod(handlers ...core.HandlerFunc) *ServerBuilder {
	b.noMethodHandlers = handlers
	return b
}

// Build creates a server with the configured controllers and middleware.
func (b *ServerBuilder) Build() (core.Server, error) {
	// Create a new server
	server, err := NewServer(b.frameworkType, b.port)
	if err != nil {
		return nil, err
	}

	// Collect controllers that should be skipped for logging and auth checks
	var skipLogPaths []string
	var skipAuthCheckPaths []string
	for _, controller := range b.controllers {
		if controller.SkipLogging() {
			path := controller.GetPath()
			if path != "" {
				skipLogPaths = append(skipLogPaths, path)
			}
		}
		if controller.SkipAuthCheck() {
			path := controller.GetPath()
			if path != "" {
				skipAuthCheckPaths = append(skipAuthCheckPaths, path)
			}
		}
	}

	// Add middleware in the correct order
	// The order of middleware registration is important:
	//
	// 1. Error handler middleware (must be first)
	//    - This middleware catches errors and panics from all subsequent middleware
	//    - It must be registered first to properly handle errors in other middleware
	//
	// 2. Timeout middleware
	//    - Controls request timeout and prevents long-running requests
	//
	// 3. CORS middleware
	//    - Handles Cross-Origin Resource Sharing headers
	//
	// 4. Logging middleware (must be after error handler)
	//    - This middleware logs request details including status codes and errors
	//    - It must be registered after the error handler to properly capture errors
	//
	// 5. Custom middleware
	//    - Any additional middleware provided by the application

	// 1. Error handler middleware (must be first)
	if b.errorConfig != nil {
		// Use framework-specific error handler middleware
		errorHandler := server.GetErrorHandlerMiddleware()
		server.Use(errorHandler.Middleware(b.errorConfig))
	} else if b.useDefaultErrorHandler {
		// Use framework-specific error handler middleware with default config
		errorHandler := server.GetErrorHandlerMiddleware()
		server.Use(errorHandler.Middleware(nil))
	}

	// 2. Timeout middleware
	if b.timeoutConfig != nil {
		server.Use(TimeoutMiddleware(b.timeoutConfig))
	} else if b.useDefaultTimeout {
		server.Use(NewDefaultTimeoutMiddleware())
	}

	// 3. CORS middleware
	if b.corsConfig != nil {
		server.Use(CORSMiddleware(b.corsConfig))
	} else if b.useDefaultCORS {
		server.Use(NewDefaultCORSMiddleware())
	}

	// 4. Logging middleware (must be after error handler)
	if b.loggingConfig != nil {
		// Add skip paths from controllers
		b.loggingConfig.SkipPaths = append(b.loggingConfig.SkipPaths, skipLogPaths...)
		// Use framework-specific logging middleware
		loggingMiddleware := server.GetLoggingMiddleware()
		server.Use(loggingMiddleware.Middleware(b.loggingConfig))
	} else if b.useDefaultLogging {
		// Create a default logging config with skip paths from controllers
		loggingConfig := &core.LoggingConfig{
			RemoteURL:        "",
			CustomFields:     make(map[string]string),
			LoggingToConsole: true,
			LoggingToRemote:  false,
			SkipPaths:        skipLogPaths,
		}
		// Use framework-specific logging middleware with default config
		loggingMiddleware := server.GetLoggingMiddleware()
		server.Use(loggingMiddleware.Middleware(loggingConfig))
	}

	// 5. Custom middleware
	for _, middleware := range b.middleware {
		server.Use(middleware)
	}

	// Register controllers
	if len(b.controllers) > 0 {
		server.RegisterRouter(b.controllers...)
	}

	// Set NoRoute handlers if provided, otherwise use default handlers
	server.NoRoute(b.noRouteHandlers...)

	// Set NoMethod handlers if provided, otherwise use default handlers
	server.NoMethod(b.noMethodHandlers...)

	return server, nil
}
