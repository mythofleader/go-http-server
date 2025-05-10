// This example demonstrates how to use the duplicate request prevention middleware
package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	server "github.com/mythofleader/go-http-server"
)

// SimpleRequestIDGenerator is a basic implementation of the RequestIDGenerator interface
// It generates a request ID based on the request path and body
type SimpleRequestIDGenerator struct{}

// GenerateRequestID generates a unique request ID from the context
func (g *SimpleRequestIDGenerator) GenerateRequestID(ctx context.Context) (string, error) {
	// Get the request from the context
	req, ok := ctx.Value(requestKey{}).(*http.Request)
	if !ok {
		// If we can't get the request from the context, use the current time as a fallback
		return fmt.Sprintf("time-%d", time.Now().UnixNano()), nil
	}

	// Use the request path as part of the ID
	path := req.URL.Path

	// Read the request body
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}

	// Important: Restore the request body so it can be read again by handlers
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Create a hash of the body
	hash := sha256.Sum256(bodyBytes)
	bodyHash := hex.EncodeToString(hash[:])

	// Combine path and body hash to create a unique ID
	requestID := fmt.Sprintf("%s-%s", path, bodyHash)
	return requestID, nil
}

// requestKey is used to store and retrieve the request from the context
type requestKey struct{}

// InMemoryRequestIDStorage is a simple in-memory implementation of the RequestIDStorage interface
type InMemoryRequestIDStorage struct {
	requestIDs map[string]bool
	mutex      sync.RWMutex
	expiry     time.Duration
	cleanup    *time.Ticker
}

// NewInMemoryRequestIDStorage creates a new InMemoryRequestIDStorage
func NewInMemoryRequestIDStorage(expiry time.Duration) *InMemoryRequestIDStorage {
	storage := &InMemoryRequestIDStorage{
		requestIDs: make(map[string]bool),
		mutex:      sync.RWMutex{},
		expiry:     expiry,
		cleanup:    time.NewTicker(expiry / 2), // Run cleanup at half the expiry time
	}

	// Start a goroutine to clean up expired request IDs
	go func() {
		for range storage.cleanup.C {
			storage.cleanupExpiredIDs()
		}
	}()

	return storage
}

// CheckRequestID checks if a request ID exists in the storage
func (s *InMemoryRequestIDStorage) CheckRequestID(requestID string) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.requestIDs[requestID]
	return exists, nil
}

// SaveRequestID saves a request ID to the storage
func (s *InMemoryRequestIDStorage) SaveRequestID(requestID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.requestIDs[requestID] = true
	return nil
}

// cleanupExpiredIDs removes expired request IDs from the storage
// This is a simplified implementation that just clears all IDs
// In a real application, you would want to track when each ID was added
func (s *InMemoryRequestIDStorage) cleanupExpiredIDs() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// In a real application, you would only remove expired IDs
	// For simplicity, we're just clearing all IDs in this example
	s.requestIDs = make(map[string]bool)
	log.Println("Cleaned up expired request IDs")
}

// Close stops the cleanup ticker
func (s *InMemoryRequestIDStorage) Close() {
	s.cleanup.Stop()
}

// requestMiddleware is a middleware that stores the request in the context
func requestMiddleware() server.HandlerFunc {
	return func(c server.Context) {
		// Store the request in the context
		req := c.Request()
		ctx := context.WithValue(req.Context(), requestKey{}, req)

		// Create a new request with the updated context
		newReq := req.WithContext(ctx)

		// Replace the original request with the new one
		// Note: This is a simplified approach and may not work in all frameworks
		// In a real application, you might need to use a different approach
		*req = *newReq
	}
}

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkStdHTTP, "8080")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create the request ID generator and storage
	idGenerator := &SimpleRequestIDGenerator{}
	idStorage := NewInMemoryRequestIDStorage(5 * time.Minute) // IDs expire after 5 minutes
	defer idStorage.Close()

	// Configure the duplicate request middleware
	dupReqConfig := &server.DuplicateRequestConfig{
		RequestIDGenerator: idGenerator,
		RequestIDStorage:   idStorage,
		ConflictMessage:    "Duplicate request detected",
	}

	// Add the request middleware to store the request in the context
	srv.Use(requestMiddleware())

	// Create a protected API group
	api := srv.Group("/api")

	// Add the duplicate request middleware to the API group
	api.Use(server.DuplicateRequestMiddleware(dupReqConfig))

	// Add a route that will be protected from duplicate requests
	api.POST("/orders", createOrderHandler)

	// Add a route to explain how to test the middleware
	srv.GET("/", helpHandler)

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following commands:")
	fmt.Println("  - curl http://localhost:8080/ (to see instructions)")
	fmt.Println("  - curl -X POST -d '{\"id\":\"123\",\"product\":\"example\"}' http://localhost:8080/api/orders")
	fmt.Println("  - Run the same command again to see the duplicate request detection")
	log.Fatal(srv.Run())
}

// createOrderHandler handles the creation of a new order
func createOrderHandler(c server.Context) {
	// In a real application, you would parse the request body and create an order
	c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Order created successfully",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// helpHandler provides instructions on how to test the middleware
func helpHandler(c server.Context) {
	helpText := `
Duplicate Request Prevention Middleware Example

This server demonstrates the duplicate request prevention middleware functionality.

To test the middleware:

1. Send a POST request to create an order:
   curl -X POST -d '{"id":"123","product":"example"}' http://localhost:8080/api/orders

2. Send the exact same request again:
   curl -X POST -d '{"id":"123","product":"example"}' http://localhost:8080/api/orders

   The second request should fail with a 409 Conflict response because it's a duplicate.

3. Change the request body slightly and try again:
   curl -X POST -d '{"id":"124","product":"example"}' http://localhost:8080/api/orders

   This should succeed because it's not a duplicate (different request body).

4. Wait for 5 minutes and try the original request again:
   curl -X POST -d '{"id":"123","product":"example"}' http://localhost:8080/api/orders

   This should succeed because the original request ID has expired.
`
	c.String(http.StatusOK, helpText)
}
