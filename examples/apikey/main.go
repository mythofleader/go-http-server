// This example demonstrates how to use the API Key middleware
package main

import (
	"fmt"
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkStdHTTP, "8080", false)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Configure the API Key middleware
	apiKeyConfig := &server.APIKeyConfig{
		APIKey:              "my-secret-api-key",          // The expected API key value
		UnauthorizedMessage: "Invalid or missing API key", // Custom error message
	}

	// Create a protected route group
	protected := srv.Group("/api")
	protected.Use(server.APIKeyMiddleware(apiKeyConfig))

	// Add a protected route
	protected.GET("/data", func(c server.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This is protected data",
			"status":  "success",
		})
	})

	// Add a public route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Welcome to the API. To access protected data, use the /api/data endpoint with the x-api-key header.")
	})

	// Add a route that explains how to use the API
	srv.GET("/help", func(c server.Context) {
		helpText := `
To access the protected API endpoints, you need to include the x-api-key header in your requests.

Example using curl:
curl -H "x-api-key: my-secret-api-key" http://localhost:8080/api/data

Example using JavaScript fetch:
fetch('http://localhost:8080/api/data', {
  headers: {
    'x-api-key': 'my-secret-api-key'
  }
})
.then(response => response.json())
.then(data => console.log(data));
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following commands:")
	fmt.Println("  - curl http://localhost:8080/help")
	fmt.Println("  - curl -H \"x-api-key: my-secret-api-key\" http://localhost:8080/api/data")
	fmt.Println("  - curl http://localhost:8080/api/data (this should fail with a 401 error)")
	log.Fatal(srv.Run())
}
