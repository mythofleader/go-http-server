// This example demonstrates how to use the ServerBuilder with NoRoute and NoMethod handlers.
package main

import (
	"fmt"
	"log"
	"net/http"

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
	return false
}

// SkipAuthCheck returns whether to skip authentication checks for this controller.
func (c *UserController) SkipAuthCheck() bool {
	return true
}

func main() {
	// Example 1: Using default NoRoute and NoMethod handlers
	fmt.Println("Example 1: Using default NoRoute and NoMethod handlers")
	builder1 := server.NewGinServerBuilder()

	// Add controllers
	builder1.AddController(&UserController{})

	// Enable default middleware
	builder1.WithDefaultLogging()
	builder1.WithDefaultErrorHandling()

	// Default NoRoute and NoMethod handlers are now applied automatically
	// No need to explicitly enable them

	// Build the server
	_, err := builder1.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	fmt.Println("Server 1 built successfully with default NoRoute and NoMethod handlers")
	fmt.Println("If this were a real server, it would return:")
	fmt.Println("- 404 Not Found with a message for non-existent routes")
	fmt.Println("- 405 Method Not Allowed with a message for invalid methods")

	// Example 2: Using custom NoRoute and NoMethod handlers
	fmt.Println("\nExample 2: Using custom NoRoute and NoMethod handlers")
	builder2 := server.NewGinServerBuilder()

	// Add controllers
	builder2.AddController(&UserController{})

	// Enable default middleware
	builder2.WithDefaultLogging()
	builder2.WithDefaultErrorHandling()

	// Set custom NoRoute handler
	builder2.WithNoRoute(func(c server.Context) {
		path := c.Request().URL.Path
		log.Printf("Custom 404 handler: %s", path)
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "Custom 404 handler",
			"path":  path,
		})
	})

	// Set custom NoMethod handler
	builder2.WithNoMethod(func(c server.Context) {
		method := c.Request().Method
		path := c.Request().URL.Path
		log.Printf("Custom 405 handler: %s %s", method, path)
		c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{
			"error":  "Custom 405 handler",
			"method": method,
			"path":   path,
		})
	})

	// Build the server
	_, err = builder2.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	fmt.Println("Server 2 built successfully with custom NoRoute and NoMethod handlers")
	fmt.Println("If this were a real server, it would return:")
	fmt.Println("- Custom 404 response for non-existent routes")
	fmt.Println("- Custom 405 response for invalid methods")

	// Note: We're not actually starting the servers in this example
	// This is just to demonstrate how to configure the ServerBuilder

	// If you want to start one of the servers, uncomment one of these lines:
	// log.Fatal(s1.Run())
	// log.Fatal(s2.Run())
}
