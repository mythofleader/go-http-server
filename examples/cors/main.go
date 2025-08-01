// This example demonstrates how to use the CORS middleware
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

	// Configure the CORS middleware
	corsConfig := &server.CORSConfig{
		// Specify allowed domains (uncomment to restrict to specific domains)
		// AllowedDomains: []string{
		//     "http://localhost:3000",
		//     "https://example.com",
		// },

		// These are the default values, you can customize them if needed
		AllowedMethods:   "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowedHeaders:   "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Requested-With",
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}

	// Add the CORS middleware to the server
	srv.Use(server.CORSMiddleware(corsConfig))

	// Add a route that returns JSON data
	srv.GET("/api/data", func(c server.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "This is CORS-enabled data",
			"status":  "success",
		})
	})

	// Add a route that explains how to test CORS
	srv.GET("/", func(c server.Context) {
		helpText := `
CORS Middleware Example

This server demonstrates the CORS middleware functionality.

To test CORS:

1. With all domains allowed (default):
   - The server is configured to allow requests from any domain.
   - Try accessing the API from any origin, and it should work.

2. With specific domains only:
   - Uncomment the AllowedDomains section in the code.
   - Restart the server.
   - Try accessing the API from a non-allowed domain, and it should be blocked.
   - Try accessing the API from an allowed domain, and it should work.

API Endpoint:
- GET /api/data - Returns JSON data

Example using JavaScript fetch:
fetch('http://localhost:8080/api/data', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  },
  credentials: 'include'
})
.then(response => response.json())
.then(data => console.log(data));
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Open http://localhost:8080 in your browser for instructions")
	log.Fatal(srv.Run())
}
