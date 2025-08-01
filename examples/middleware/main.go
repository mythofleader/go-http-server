package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

// CustomMiddleware demonstrates how to use the Next() method for flow control in middleware.
// It logs the request before and after the handler is executed.
func CustomMiddleware() server.HandlerFunc {
	return func(c server.Context) {
		// Log before handler execution
		log.Printf("Before handler: %s %s", c.Request().Method, c.Request().URL.Path)

		// Record start time
		start := time.Now()

		// Call the next handler in the chain
		c.Next()

		// Log after handler execution
		duration := time.Since(start)
		log.Printf("After handler: %s %s (took %v)", c.Request().Method, c.Request().URL.Path, duration)
	}
}

// AnotherMiddleware demonstrates how to use the Next() method for conditional flow control.
// It only calls the next handler if the request path doesn't start with "/skip".
func AnotherMiddleware() server.HandlerFunc {
	return func(c server.Context) {
		path := c.Request().URL.Path

		if path == "/skip" {
			log.Printf("Skipping handler for path: %s", path)
			c.String(http.StatusOK, "Handler skipped!")
			return
		}

		log.Printf("Continuing to next handler for path: %s", path)
		c.Next()
	}
}

func main() {
	// Create a new server
	s, err := server.NewServer(server.FrameworkGin, "8080", false)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// The order of middleware registration is important:
	// 1. Error handler middleware must be registered first to catch errors in other middleware
	// 2. Timeout middleware (if used)
	// 3. CORS middleware (if used)
	// 4. Logging middleware should be registered after error handler to properly capture errors
	// 5. Custom middleware

	// Add error handler middleware first
	errorHandler := s.GetErrorHandlerMiddleware()
	s.Use(errorHandler.Middleware(nil))

	// Add logging middleware after error handler
	loggingConfig := &server.LoggingConfig{
		LoggingToConsole: true,
		CustomFields: map[string]string{
			"example": "middleware-flow-control",
		},
	}
	s.Use(s.GetLoggingMiddleware().Middleware(loggingConfig))

	// Add our custom middleware last
	s.Use(CustomMiddleware())
	s.Use(AnotherMiddleware())

	// Add a route
	s.GET("/", func(c server.Context) {
		log.Println("Handler executing...")
		time.Sleep(100 * time.Millisecond) // Simulate work
		c.String(http.StatusOK, "Hello, World!")
	})

	// Add a route that will be skipped by AnotherMiddleware
	s.GET("/skip", func(c server.Context) {
		log.Println("This handler should not be executed!")
		c.String(http.StatusOK, "This should not be seen!")
	})

	// Add a route that demonstrates the order of middleware execution
	s.GET("/order", func(c server.Context) {
		log.Println("Main handler for /order")
		c.String(http.StatusOK, "Check the logs to see the order of middleware execution!")
	})

	// Add a help route
	s.GET("/help", func(c server.Context) {
		helpText := `
Middleware Flow Control Example

This server demonstrates how to use the Next() method for flow control in middleware.

Available endpoints:
- GET /         - Returns a simple text response (both middleware execute)
- GET /skip     - AnotherMiddleware skips the handler for this path
- GET /order    - Demonstrates the order of middleware execution
- GET /help     - This help page

Check your console to see the logs for each request.
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Open http://localhost:8080/help in your browser for instructions")
	log.Fatal(s.Run())
}
