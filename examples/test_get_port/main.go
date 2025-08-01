// This example demonstrates the use of the GetPort method
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	fmt.Println("=== Testing GetPort method ===")

	// Test with Gin framework
	fmt.Println("\n=== Testing with Gin framework ===")
	testGetPort(server.FrameworkGin)

	// Wait a bit to separate the logs
	time.Sleep(1 * time.Second)

	// Test with standard HTTP framework
	fmt.Println("\n=== Testing with standard HTTP framework ===")
	testGetPort(server.FrameworkStdHTTP)
}

func testGetPort(frameworkType server.FrameworkType) {
	// Test with explicit port
	fmt.Println("\n--- Testing with explicit port ---")
	testWithExplicitPort(frameworkType)

	// Wait a bit to separate the logs
	time.Sleep(1 * time.Second)

	// Test with default port
	fmt.Println("\n--- Testing with default port ---")
	testWithDefaultPort(frameworkType)

	// Wait a bit to separate the logs
	time.Sleep(1 * time.Second)

	// Test with random port
	fmt.Println("\n--- Testing with random port ---")
	testWithRandomPort(frameworkType)
}

func testWithExplicitPort(frameworkType server.FrameworkType) {
	// Create a server with an explicit port
	srv, err := server.NewServer(frameworkType, "8081", true)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add a simple route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Get the port from the server
	port := srv.GetPort()
	fmt.Printf("Server port: %s\n", port)

	// Verify the port is correct
	if port != "8081" {
		log.Fatalf("Expected port 8081, got %s", port)
	}
	fmt.Println("Port is correct: 8081")
}

func testWithDefaultPort(frameworkType server.FrameworkType) {
	// Create a server builder
	builder := server.NewServerBuilder(frameworkType)

	// Set the default port (8080)
	builder.WithDefaultPort()

	// Build the server
	srv, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	// Add a simple route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Get the port from the server
	port := srv.GetPort()
	fmt.Printf("Server port: %s\n", port)

	// Verify the port is correct
	if port != "8080" {
		log.Fatalf("Expected port 8080, got %s", port)
	}
	fmt.Println("Port is correct: 8080")
}

func testWithRandomPort(frameworkType server.FrameworkType) {
	// Create a server builder
	builder := server.NewServerBuilder(frameworkType)

	// Set a random port
	builder.WithDefaultRandomPort()

	// Build the server
	srv, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build server: %v", err)
	}

	// Add a simple route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Get the port from the server
	port := srv.GetPort()
	fmt.Printf("Server port: %s\n", port)

	// Verify the port is not empty
	if port == "" {
		log.Fatalf("Expected non-empty port, got empty string")
	}
	fmt.Printf("Port is not empty: %s\n", port)

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Starting server on port %s...\n", port)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Verify the port is valid by connecting to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/", port))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Successfully connected to server on port %s\n", port)

	// Stop the server
	fmt.Println("Stopping server...")
	if err := srv.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}
	fmt.Println("Server stopped")
}
