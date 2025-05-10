// This example demonstrates how to use the logging middleware
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server using the standard HTTP implementation
	srv, err := server.NewServer(server.FrameworkStdHTTP, "8081")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// The order of middleware registration is important:
	// 1. Error handler middleware must be registered first to catch errors in other middleware
	// 2. Logging middleware should be registered after error handler to properly capture errors
	//    and log the correct status codes

	// Add error handler middleware first
	// Use framework-specific error handler middleware
	errorHandler := srv.GetErrorHandlerMiddleware()
	srv.Use(errorHandler.Middleware(nil))

	// Example 1: Basic logging middleware with default configuration
	// This will log requests to the console with default settings
	// Uncomment the following line to use the default configuration
	// loggingMiddleware := srv.GetLoggingMiddleware()
	// srv.Use(loggingMiddleware.Middleware(nil))

	// Example 2: Logging middleware with custom fields
	// This will log requests to the console with custom fields
	loggingConfig := &server.LoggingConfig{
		LoggingToConsole: true,  // Log to console (default)
		LoggingToRemote:  false, // Don't log to remote
		CustomFields: map[string]string{
			"environment": "development",
			"version":     server.Version,
			"app_name":    "logging-example",
		},
		// Paths to skip for logging
		SkipPaths: []string{
			"/health",
			"/favicon.ico",
		},
	}
	srv.Use(srv.GetLoggingMiddleware().Middleware(loggingConfig))

	// Example 3: Console-only logging with skip paths
	// Uncomment the following lines to use the console-only logging with skip paths
	/*
		skipPaths := []string{"/health", "/favicon.ico"}
		customFields := map[string]string{
			"environment": "development",
			"version":     server.Version,
		}
		consoleConfig := server.NewDefaultConsoleLogging(skipPaths, customFields)
		loggingMiddleware := srv.GetLoggingMiddleware()
		srv.Use(loggingMiddleware.Middleware(consoleConfig))
	*/

	// Example 4: Remote logging configuration
	// Uncomment the following lines to enable remote logging
	/*
		remoteLoggingConfig := &server.LoggingConfig{
			RemoteURL: "https://your-logging-service.com/api/logs",
			LoggingToConsole: true,  // Also log to console
			LoggingToRemote:  true,  // Enable remote logging
			CustomFields: map[string]string{
				"environment": "development",
				"version":     server.Version,
			},
		}
		loggingMiddleware := srv.GetLoggingMiddleware()
		srv.Use(loggingMiddleware.Middleware(remoteLoggingConfig))
	*/

	// Add routes to demonstrate logging
	// Normal route - will be logged
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World! This request will be logged.")
	})

	// Route that returns JSON - will be logged
	srv.GET("/json", func(c server.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This JSON response will be logged",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Route that returns an error - will be logged with error status
	srv.GET("/error", func(c server.Context) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Bad Request",
			"message": "This error response will be logged with status code 400",
		})
	})

	// Route that sleeps for a specific amount of time to simulate a slow handler
	srv.GET("/slow", func(c server.Context) {
		// Sleep for 1 second to simulate a slow handler
		time.Sleep(1 * time.Second)
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This response was delayed by 1 second",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Health check route - will be skipped by the logging middleware
	srv.GET("/health", func(c server.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Add a route that explains how to test the logging middleware
	srv.GET("/help", func(c server.Context) {
		helpText := `
Logging Middleware Example

This server demonstrates the logging middleware functionality.

Available endpoints:
- GET /         - Returns a simple text response (logged)
- GET /json     - Returns a JSON response (logged)
- GET /error    - Returns an error response (logged with error status)
- GET /slow     - Returns a response after a 1-second delay (logged with latency)
- GET /health   - Health check endpoint (not logged due to SkipPaths)
- GET /help     - This help page (logged)

The logging middleware is configured to:
1. Log requests to the console
2. Include custom fields (environment, version, app_name)
3. Skip requests to /health and /favicon.ico

Check your console to see the logs for each request.
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8081")
	fmt.Println("Open http://localhost:8081/help in your browser for instructions")
	log.Fatal(srv.Run())
}
