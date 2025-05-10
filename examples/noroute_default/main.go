// This example demonstrates how the default NoRoute and NoMethod handlers are applied automatically
package main

import (
	"fmt"
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server without explicitly setting NoRoute or NoMethod handlers
	srv, err := server.NewServer(server.FrameworkGin, "8080")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add error handler middleware to handle the errors from the default handlers
	errorHandler := srv.GetErrorHandlerMiddleware()
	srv.Use(errorHandler.Middleware(nil))

	// Add logging middleware
	loggingMiddleware := srv.GetLoggingMiddleware()
	srv.Use(loggingMiddleware.Middleware(nil))

	// Register a route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following:")
	fmt.Println("  - http://localhost:8080/ (valid route)")
	fmt.Println("  - http://localhost:8080/nonexistent (404 Not Found - default NoRoute handler)")
	fmt.Println("  - POST http://localhost:8080/ (405 Method Not Allowed - default NoMethod handler)")
	log.Fatal(srv.Run())
}
