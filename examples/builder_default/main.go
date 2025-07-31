// This example demonstrates how to use the ServerBuilder with default middleware.
package main

import (
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

// UserController is a simple controller for user-related routes.
type UserController struct{}

func (c *UserController) SkipLogging() bool {
	return true
}

func (c *UserController) SkipAuthCheck() bool {
	return true
}

// GetHttpMethod returns the HTTP method for the route.
func (c *UserController) GetHttpMethod() server.HttpMethod {
	return server.GET
}

// GetPath returns the path for the route.
func (c *UserController) GetPath() string {
	return "/users"
}

// Handler returns handler functions for the route.
func (c *UserController) Handler() []server.HandlerFunc {
	return []server.HandlerFunc{
		func(ctx server.Context) {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"users": []map[string]interface{}{
					{"id": 1, "name": "Alice"},
					{"id": 2, "name": "Bob"},
				},
			})
		},
	}
}

// GetLogIgnorePath returns a path to ignore for logging.
func (c *UserController) GetLogIgnorePath() string {
	return "/users/health"
}

// GetAuthCheckIgnorePath returns a path to ignore for authentication checks.
func (c *UserController) GetAuthCheckIgnorePath() string {
	return "/users/public"
}

// HealthController is a simple controller for health check routes.
type HealthController struct{}

func (c *HealthController) SkipLogging() bool {
	return true
}

func (c *HealthController) SkipAuthCheck() bool {
	return true
}

// GetHttpMethod returns the HTTP method for the route.
func (c *HealthController) GetHttpMethod() server.HttpMethod {
	return server.GET
}

// GetPath returns the path for the route.
func (c *HealthController) GetPath() string {
	return "/health"
}

// Handler returns handler functions for the route.
func (c *HealthController) Handler() []server.HandlerFunc {
	return []server.HandlerFunc{
		func(ctx server.Context) {
			ctx.JSON(http.StatusOK, map[string]string{
				"status": "ok",
			})
		},
	}
}

// GetLogIgnorePath returns a path to ignore for logging.
func (c *HealthController) GetLogIgnorePath() string {
	return "/health"
}

// GetAuthCheckIgnorePath returns a path to ignore for authentication checks.
func (c *HealthController) GetAuthCheckIgnorePath() string {
	return "/health"
}

func main() {
	// Create a server builder with Gin framework and port 8080
	builder := server.NewGinServerBuilder()

	// Add controllers
	builder.AddControllers(
		&UserController{},
		&HealthController{},
	)

	// Enable default middleware
	builder.WithDefaultLogging(true)   // Enable default logging middleware with console logging enabled
	builder.WithDefaultTimeout()       // Enable default timeout middleware
	builder.WithDefaultCORS()          // Enable default CORS middleware
	builder.WithDefaultErrorHandling() // Enable default error handler middleware

	// Add custom middleware
	builder.AddMiddleware(func(c server.Context) {
		log.Printf("Request: %s %s", c.Request().Method, c.Request().URL.Path)
	})

	// Build the server
	s, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	// Start the server
	log.Println("Server starting on :8080")
	if err := s.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
