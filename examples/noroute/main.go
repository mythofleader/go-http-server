// This example demonstrates how to use NoRoute and NoMethod handlers
package main

import (
	"fmt"
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkGin, "8080")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add error handler middleware
	errorHandler := srv.GetErrorHandlerMiddleware()
	srv.Use(errorHandler.Middleware(nil))

	// Add logging middleware
	loggingMiddleware := srv.GetLoggingMiddleware()
	srv.Use(loggingMiddleware.Middleware(nil))

	// Register a route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Register a NoRoute handler for 404 Not Found errors
	srv.NoRoute(func(c server.Context) {
		path := c.Request().URL.Path
		log.Printf("404 Not Found: %s", path)

		// Option 1: Return a custom JSON response
		c.JSON(http.StatusNotFound, server.NewErrorResponse(http.StatusNotFound,
			fmt.Sprintf("Page not found: %s", path)))
	})

	// Register a NoMethod handler for 405 Method Not Allowed errors
	srv.NoMethod(func(c server.Context) {
		method := c.Request().Method
		path := c.Request().URL.Path
		log.Printf("405 Method Not Allowed: %s %s", method, path)

		// Option 2: Use the error handler middleware
		err := fmt.Errorf("Method %s not allowed for path %s", method, path)
		c.Error(server.NewMethodNotAllowedHttpError(err))
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following:")
	fmt.Println("  - http://localhost:8080/ (valid route)")
	fmt.Println("  - http://localhost:8080/nonexistent (404 Not Found)")
	fmt.Println("  - POST http://localhost:8080/ (405 Method Not Allowed)")
	log.Fatal(srv.Run())
}
