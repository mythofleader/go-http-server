// This example demonstrates how to use the ServerBuilder to create a server with controllers and middleware.
package main

import (
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

// UserController is a simple controller for user-related routes.
type UserController struct{}

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

// SkipLogging returns whether to skip logging for this controller.
func (c *UserController) SkipLogging() bool {
	return true
}

// SkipAuthCheck returns whether to skip authentication checks for this controller.
func (c *UserController) SkipAuthCheck() bool {
	return true
}

// HealthController is a simple controller for health check routes.
type HealthController struct{}

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

// SkipLogging returns whether to skip logging for this controller.
func (c *HealthController) SkipLogging() bool {
	return true
}

// SkipAuthCheck returns whether to skip authentication checks for this controller.
func (c *HealthController) SkipAuthCheck() bool {
	return true
}

func main() {
	// Create a server builder with Gin framework and port 8080
	builder := server.NewGinServerBuilder()

	// Add controllers
	builder.AddControllers(
		&UserController{},
		&HealthController{},
	)

	// Configure logging
	builder.WithLogging(map[string]string{
		"environment": "development",
		"version":     "1.0.0",
	})

	// Configure timeout
	builder.WithTimeout(server.TimeoutConfig{
		Timeout: 5 * time.Second,
	})

	// Configure CORS
	builder.WithCORS(server.CORSConfig{
		AllowedDomains: []string{"*"},
		AllowedMethods: "GET, POST, PUT, DELETE, PATCH",
	})

	// Configure error handler
	builder.WithErrorHandler(server.ErrorHandlerConfig{
		DefaultErrorMessage: "Internal Server Error",
		DefaultStatusCode:   500,
	})

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
