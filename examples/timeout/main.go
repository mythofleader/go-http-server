// This example demonstrates how to use the timeout middleware
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkGin, "8080", false)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add error handler middleware first using framework-specific implementation
	errorHandler := srv.GetErrorHandlerMiddleware()
	srv.Use(errorHandler.Middleware(nil))

	// Example 1: Default timeout middleware
	// This will timeout requests after 2 seconds (default)
	// Uncomment the following line to use the default configuration
	// srv.Use(server.NewDefaultTimeoutMiddleware())

	// Example 2: Timeout middleware with custom timeout
	// This will timeout requests after 3 seconds
	timeoutConfig := &server.TimeoutConfig{
		Timeout: 3 * time.Second, // Set timeout to 3 seconds
	}
	srv.Use(server.TimeoutMiddleware(timeoutConfig))

	// Add logging middleware for better visibility using framework-specific implementation
	loggingMiddleware := srv.GetLoggingMiddleware()
	srv.Use(loggingMiddleware.Middleware(nil))

	// Add routes to demonstrate timeout
	// Normal route - responds immediately
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World! This request responds immediately.")
	})

	// Route that responds after a delay, but before timeout
	srv.GET("/delay/1", func(c server.Context) {
		// Sleep for 1 second (less than the timeout)
		time.Sleep(1 * time.Second)
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This response was delayed by 1 second, but completed before the timeout",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Route that responds after a delay, but before timeout
	srv.GET("/delay/2", func(c server.Context) {
		// Sleep for 2 seconds (less than the timeout)
		time.Sleep(2 * time.Second)
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This response was delayed by 2 seconds, but completed before the timeout",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Route that times out
	srv.GET("/timeout", func(c server.Context) {
		// Sleep for 5 seconds (more than the timeout)
		time.Sleep(5 * time.Second)
		// This response will never be sent because the timeout middleware
		// will respond with a 503 Service Unavailable after 3 seconds
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This response will never be sent due to timeout",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Add a route that explains how to test the timeout middleware
	srv.GET("/help", func(c server.Context) {
		helpText := `
Timeout Middleware Example

This server demonstrates the timeout middleware functionality.

Available endpoints:
- GET /          - Returns a simple text response immediately
- GET /delay/1   - Returns a response after a 1-second delay (before timeout)
- GET /delay/2   - Returns a response after a 2-second delay (before timeout)
- GET /timeout   - Attempts to return a response after a 5-second delay (will timeout)
- GET /help      - This help page

The timeout middleware is configured to:
1. Timeout requests after 3 seconds
2. Return a 503 Service Unavailable response when a timeout occurs

Try accessing the /timeout endpoint and observe that it returns a 503 response
after 3 seconds, instead of waiting for the full 5 seconds.
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Open http://localhost:8080/help in your browser for instructions")
	log.Fatal(srv.Run())
}
