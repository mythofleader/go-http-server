// This example demonstrates how custom NoRoute and NoMethod handlers work
package main

import (
	"fmt"
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkGin, "8080", false)
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

	// Register a custom NoRoute handler
	srv.NoRoute(func(c server.Context) {
		path := c.Request().URL.Path
		log.Printf("Custom 404 handler: %s", path)
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Custom 404 handler",
			"path":    path,
			"message": "The page you're looking for doesn't exist",
		})
	})

	// Register a custom NoMethod handler
	srv.NoMethod(func(c server.Context) {
		method := c.Request().Method
		path := c.Request().URL.Path
		log.Printf("Custom 405 handler: %s %s", method, path)
		c.JSON(http.StatusMethodNotAllowed, map[string]interface{}{
			"error":   "Custom 405 handler",
			"method":  method,
			"path":    path,
			"message": "The method you're using is not allowed for this path",
		})
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following:")
	fmt.Println("  - http://localhost:8080/ (valid route)")
	fmt.Println("  - http://localhost:8080/nonexistent (404 Not Found - custom NoRoute handler)")
	fmt.Println("  - POST http://localhost:8080/ (405 Method Not Allowed - custom NoMethod handler)")
	log.Fatal(srv.Run())
}
