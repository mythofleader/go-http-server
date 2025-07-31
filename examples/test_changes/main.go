// This example demonstrates the use of the Stop method and WithDefaultLogging with console parameter
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Test with console logging enabled
	fmt.Println("=== Testing with console logging enabled ===")
	testServer(true)

	// Wait a bit to separate the logs
	time.Sleep(1 * time.Second)

	// Test with console logging disabled
	fmt.Println("\n=== Testing with console logging disabled ===")
	testServer(false)
}

func testServer(consoleLogging bool) {
	// Create a server builder with the standard HTTP implementation
	builder := server.NewServerBuilder(server.FrameworkStdHTTP, "8081")

	// Enable default logging with the console parameter
	builder.WithDefaultLogging(consoleLogging)

	// Build the server
	srv, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	// Add a simple route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting server on :8081")
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Make a request to the server
	fmt.Println("Making request to server...")
	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Printf("Response status: %s\n", resp.Status)

	// Stop the server
	fmt.Println("Stopping server...")
	if err := srv.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}
	fmt.Println("Server stopped")

	// Wait a bit to ensure the server has stopped
	time.Sleep(1 * time.Second)
}
