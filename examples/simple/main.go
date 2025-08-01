// This example demonstrates how to use the tenqube-go-http-server library
// to create a simple HTTP server.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Parse command line flags
	framework := flag.String("framework", "gin", "HTTP framework to use (gin, std)")
	lambdaMode := flag.Bool("lambda", false, "Run in AWS Lambda mode")
	port := flag.String("port", "8080", "Port to run the server on")
	env := flag.String("env", "dev", "Environment (dev, prod)")
	flag.Parse()

	// Create a new server based on the specified framework
	var s server.Server
	var err error

	switch *framework {
	case "gin":
		s, err = server.NewServer(server.FrameworkGin, *port, false)
	case "std":
		s, err = server.NewServer(server.FrameworkStdHTTP, *port, false)
	default:
		// Default to Gin
		s, err = server.NewServer(server.FrameworkGin, *port, false)
	}

	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add middleware
	// Example 1: Basic logging middleware (logs to console only)
	// loggingMiddleware := s.GetLoggingMiddleware()
	// s.Use(loggingMiddleware.Middleware(nil))

	// Example 2: Logging middleware with custom fields and console logging
	loggingConfig := &server.LoggingConfig{
		LoggingToConsole: true,  // Log to console (default)
		LoggingToRemote:  false, // Don't log to remote
		CustomFields: map[string]string{
			"environment": "development",
			"version":     server.Version,
		},
	}
	// Note: The logging middleware is added after the error handler middleware below

	// Example 3: Logging middleware with remote URL
	// Uncomment the following lines to enable remote logging
	/*
		remoteLoggingConfig := &server.LoggingConfig{
			RemoteURL: "https://your-logging-service.com/api/logs",
			LoggingToRemote: true,
			CustomFields: map[string]string{
				"environment": "development",
				"version":     server.Version,
			},
		}
		loggingMiddleware := s.GetLoggingMiddleware()
		s.Use(loggingMiddleware.Middleware(remoteLoggingConfig))
	*/

	// Example 4: Logging middleware with remote logging only
	// Logs are sent to remote URL but not printed to console
	// Uncomment the following lines to enable remote-only logging
	/*
		remoteOnlyConfig := &server.LoggingConfig{
			LoggingToConsole: false, // Don't log to console
			LoggingToRemote:  true,  // Log to remote
			RemoteURL: "https://your-logging-service.com/api/logs", // Required for remote logging
			CustomFields: map[string]string{
				"environment": "production",
				"version":     server.Version,
			},
		}
		loggingMiddleware := s.GetLoggingMiddleware()
		s.Use(loggingMiddleware.Middleware(remoteOnlyConfig))
	*/

	// The order of middleware registration is important:
	// 1. Error handler middleware (must be first)
	// 2. Timeout middleware
	// 3. CORS middleware (if used)
	// 4. Logging middleware
	// 5. Custom middleware

	// Example 7: Error handler middleware (must be first)
	// This middleware will catch errors and return appropriate HTTP responses
	errorHandlerConfig := &server.ErrorHandlerConfig{
		DefaultErrorMessage: "서버 오류가 발생했습니다", // Default error message
		DefaultStatusCode:   500,             // Default status code
	}
	errorHandlerMiddleware := s.GetErrorHandlerMiddleware()
	s.Use(errorHandlerMiddleware.Middleware(errorHandlerConfig))

	// Example 5: Timeout middleware (should be after error handler)
	// This middleware will timeout requests that take longer than the specified duration
	// If no timeout is specified, it defaults to 2 seconds
	timeoutConfig := &server.TimeoutConfig{
		Timeout: 2 * time.Second, // Set timeout to 2 seconds
	}
	s.Use(server.TimeoutMiddleware(timeoutConfig))

	// Example 6: Timeout middleware with custom timeout
	// Uncomment the following lines to set a custom timeout
	/*
		customTimeoutConfig := &server.TimeoutConfig{
			Timeout: 5 * time.Second, // Set timeout to 5 seconds
		}
		s.Use(server.TimeoutMiddleware(customTimeoutConfig))
	*/

	// Add the logging middleware after the error handler and timeout middleware
	// Use the framework-specific logging middleware implementation for accurate status code logging
	s.Use(s.GetLoggingMiddleware().Middleware(loggingConfig))

	// Register routes
	s.GET("/", helloHandler)
	s.GET("/json", jsonHandler)
	s.GET("/slow", slowHandler) // This handler will sleep for 3 seconds, triggering the timeout
	s.GET("/error/400", badRequestHandler)
	s.GET("/error/401", unauthorizedHandler)
	s.GET("/error/403", forbiddenHandler)
	s.GET("/error/500", internalServerErrorHandler)
	s.GET("/error/panic", panicHandler)

	// Create a router group
	api := s.Group("/api")
	{
		api.GET("/users", getUsersHandler)
		api.POST("/users", createUserHandler)
	}

	// Set default NoRoute handler (404 Not Found)
	s.NoRoute()

	// Set default NoMethod handler (405 Method Not Allowed)
	s.NoMethod()

	// Log the environment setting
	log.Printf("Using environment: %s (console logging %s)", *env, map[string]string{"dev": "enabled", "prod": "disabled"}[*env])

	// Start the server
	if *lambdaMode {
		// For Lambda, we need to use StartLambda instead of Run
		log.Println("Starting Lambda server with", *framework, "framework")

		// Start Lambda mode
		if err := s.StartLambda(); err != nil {
			log.Fatalf("Failed to start Lambda: %v", err)
		}
	} else {
		// For other frameworks, use Run
		log.Printf("Server starting on :%s with %s framework", *port, *framework)
		if err := s.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}

// Note: Use the framework-specific logging middleware implementation for accurate status code logging

// Handler examples
func helloHandler(c server.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func jsonHandler(c server.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Hello, JSON!",
		"status":  "success",
	})
}

func getUsersHandler(c server.Context) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	c.JSON(http.StatusOK, users)
}

func createUserHandler(c server.Context) {
	var user struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// In a real application, you would save the user to a database
	c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      3, // Just an example ID
		"name":    user.Name,
		"message": "User created successfully",
	})
}

// slowHandler is a handler that sleeps for 3 seconds to demonstrate the timeout middleware
func slowHandler(c server.Context) {
	// Sleep for 3 seconds
	time.Sleep(3 * time.Second)

	// This response will not be sent if the timeout middleware times out first
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "This response was delayed by 3 seconds",
	})
}

// badRequestHandler demonstrates a 400 Bad Request error
func badRequestHandler(c server.Context) {
	// Directly set the status code instead of using panic
	c.JSON(http.StatusBadRequest, server.NewBadRequestResponse("잘못된 요청 파라미터"))
}

// unauthorizedHandler demonstrates a 401 Unauthorized error
func unauthorizedHandler(c server.Context) {
	// Create a 401 Unauthorized error
	err := server.NewUnauthorizedHttpError(fmt.Errorf("인증이 필요합니다"))

	// Panic with the error to trigger the error handler middleware
	panic(err)
}

// forbiddenHandler demonstrates a 403 Forbidden error
func forbiddenHandler(c server.Context) {
	// Create a 403 Forbidden error
	err := server.NewForbiddenHttpError(fmt.Errorf("접근 권한이 없습니다"))

	// Panic with the error to trigger the error handler middleware
	panic(err)
}

// internalServerErrorHandler demonstrates a 500 Internal Server Error
func internalServerErrorHandler(c server.Context) {
	// Create a 500 Internal Server Error
	err := server.NewInternalServerHttpError(fmt.Errorf("서버 오류가 발생했습니다"))

	// Panic with the error to trigger the error handler middleware
	panic(err)
}

// panicHandler demonstrates a panic that will be caught by the error handler middleware
func panicHandler(c server.Context) {
	// Panic with a string
	panic("This is a panic that will be caught by the error handler middleware")
}
