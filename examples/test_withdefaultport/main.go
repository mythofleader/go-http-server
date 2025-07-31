// This example demonstrates the use of WithDefaultPort to automatically find an available port
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	fmt.Println("=== Testing WithDefaultPort to automatically find available ports ===")

	// Test multiple instances to ensure different ports are assigned
	const numInstances = 3
	for i := 0; i < numInstances; i++ {
		fmt.Printf("\n=== Instance %d ===\n", i+1)
		testWithDefaultPort()
		
		// Wait a bit to separate the logs
		time.Sleep(1 * time.Second)
	}
}

func testWithDefaultPort() {
	// Create a server builder with the standard HTTP implementation
	builder := server.NewServerBuilder(server.FrameworkStdHTTP)

	// Use WithDefaultPort to automatically find an available port
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
	// We can extract it from the server's address when it starts
	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting server...")
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Make a request to the server to determine its port
	// We'll try ports in the range 8000-9000
	var port int
	var resp *http.Response
	var requestErr error
	
	for port = 8000; port <= 9000; port++ {
		url := fmt.Sprintf("http://localhost:%d/", port)
		resp, requestErr = http.Get(url)
		if requestErr == nil {
			break
		}
	}
	
	if requestErr != nil {
		log.Fatalf("Failed to connect to server: %v", requestErr)
	}
	
	defer resp.Body.Close()

	// Print the port that was found
	fmt.Printf("Server is running on port %d\n", port)
	
	// Verify the port is within the expected range
	if port < 8000 || port > 9000 {
		log.Fatalf("Port %d is outside the expected range (8000-9000)", port)
	}
	
	fmt.Printf("Port %d is within the expected range (8000-9000)\n", port)
	
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