package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	server "github.com/mythofleader/go-http-server"
)

// User represents a user in the system
type User struct {
	ID       string
	Username string
	Role     string
}

// Common user store that both services will use
type UserStore struct {
	users map[string]User
}

// NewUserStore creates a new UserStore
func NewUserStore() *UserStore {
	// Create some sample users
	users := map[string]User{
		"user1": {
			ID:       "1",
			Username: "user1",
			Role:     "user",
		},
		"admin": {
			ID:       "2",
			Username: "admin",
			Role:     "admin",
		},
	}

	return &UserStore{
		users: users,
	}
}

// BasicAuthService implements only the BasicAuthUserLookup interface
type BasicAuthService struct {
	store *UserStore
}

// NewBasicAuthService creates a new BasicAuthService
func NewBasicAuthService(store *UserStore) *BasicAuthService {
	return &BasicAuthService{
		store: store,
	}
}

// LookupUserByBasicAuth looks up a user by username and password
func (s *BasicAuthService) LookupUserByBasicAuth(username, password string) (interface{}, error) {
	// In a real application, you would verify the password
	// For this example, we'll just check if the user exists and the password is "password"
	user, exists := s.store.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	if password != "password" {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// JWTService implements only the JWTUserLookup interface
type JWTService struct {
	store *UserStore
}

// NewJWTService creates a new JWTService
func NewJWTService(store *UserStore) *JWTService {
	return &JWTService{
		store: store,
	}
}

// LookupUserByJWT looks up a user by JWT claims
func (s *JWTService) LookupUserByJWT(claims server.MapClaims) (interface{}, error) {
	// Extract the username from the claims
	username, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing subject")
	}

	// Look up the user
	user, exists := s.store.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// LegacyUserService implements the full UserLookupInterface (both BasicAuthUserLookup and JWTUserLookup)
// This is included for backward compatibility demonstration
type LegacyUserService struct {
	store *UserStore
}

// NewLegacyUserService creates a new LegacyUserService
func NewLegacyUserService(store *UserStore) *LegacyUserService {
	return &LegacyUserService{
		store: store,
	}
}

// LookupUserByBasicAuth looks up a user by username and password
func (s *LegacyUserService) LookupUserByBasicAuth(username, password string) (interface{}, error) {
	// In a real application, you would verify the password
	// For this example, we'll just check if the user exists and the password is "password"
	user, exists := s.store.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	if password != "password" {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// LookupUserByJWT looks up a user by JWT claims
func (s *LegacyUserService) LookupUserByJWT(claims server.MapClaims) (interface{}, error) {
	// Extract the username from the claims
	username, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token: missing subject")
	}

	// Look up the user
	user, exists := s.store.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func main() {
	// Create a new server
	srv, err := server.NewServer(server.FrameworkStdHTTP, "8080")
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create a user store
	userStore := NewUserStore()

	// Create service implementations
	// Create the JWT service for our active example
	jwtService := NewJWTService(userStore)

	// Create other services for demonstration purposes
	// These are commented out in the examples below
	// Uncomment the appropriate section to use a different authentication method
	var basicAuthService = NewBasicAuthService(userStore) // For BasicAuth example
	var legacyService = NewLegacyUserService(userStore)   // For legacy example

	// Suppress unused variable warnings
	_ = basicAuthService
	_ = legacyService

	// Choose which authentication method to use
	// For this example, we'll use the specific JWTUserLookup interface
	authConfig := &server.AuthConfig{
		JWTLookup: jwtService, // Only implement JWTUserLookup
		AuthType:  server.AuthTypeJWT,
		JWTSecret: "your-secret-key",
		// Skip authentication for specific paths
		SkipPaths: []string{
			"/public",                // Exact path match
			"/api/docs/*",            // Wildcard pattern - matches all paths starting with /api/docs/
			"/api/users/:id/profile", // Parameter pattern - matches paths like /api/users/123/profile
		},
	}

	// Alternatively, you could use the specific BasicAuthUserLookup interface
	// Uncomment this section and comment out the JWT section above
	// authConfig := &server.AuthConfig{
	//     BasicAuthLookup: basicAuthService,
	//     AuthType:        server.AuthTypeBasic,
	// }

	// Or you could use the legacy UserLookupInterface (for backward compatibility)
	// Uncomment this section and comment out the JWT section above
	// authConfig := &server.AuthConfig{
	//     UserLookup: legacyService,
	//     AuthType:   server.AuthTypeJWT,
	//     JWTSecret:  "your-secret-key",
	// }

	// Create a protected route group
	protected := srv.Group("/api")
	protected.Use(server.AuthMiddleware(authConfig))

	// Add a protected route
	protected.GET("/profile", func(c server.Context) {
		// Get the authenticated user from the context
		user, ok := server.GetUserFromContext(c.Request().Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}

		// Type assert to User
		u, ok := user.(User)
		if !ok {
			c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			return
		}

		// Return the user profile
		c.JSON(http.StatusOK, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"role":     u.Role,
		})
	})

	// Add public routes to demonstrate SkipPaths functionality

	// Public route (exact path match in SkipPaths)
	srv.GET("/public", func(c server.Context) {
		c.String(http.StatusOK, "This is a public route (exact path match in SkipPaths)")
	})

	// API docs routes (wildcard pattern in SkipPaths)
	srv.GET("/api/docs/overview", func(c server.Context) {
		c.String(http.StatusOK, "API Documentation Overview - No authentication required (wildcard pattern in SkipPaths)")
	})

	srv.GET("/api/docs/endpoints", func(c server.Context) {
		c.String(http.StatusOK, "API Endpoints Documentation - No authentication required (wildcard pattern in SkipPaths)")
	})

	// User profile routes (parameter pattern in SkipPaths)
	srv.GET("/api/users/123/profile", func(c server.Context) {
		c.String(http.StatusOK, "Public profile for user 123 - No authentication required (parameter pattern in SkipPaths)")
	})

	srv.GET("/api/users/456/profile", func(c server.Context) {
		c.String(http.StatusOK, "Public profile for user 456 - No authentication required (parameter pattern in SkipPaths)")
	})

	// Add a public route
	srv.GET("/", func(c server.Context) {
		// Create a help text that explains the available routes
		helpText := `
Auth Middleware Example

This server demonstrates the auth middleware functionality with path matching.

Available endpoints:
- GET /                       - This help page (public)
- GET /public                 - Public route (not authenticated - exact path match)
- GET /api/docs/overview      - API docs overview (not authenticated - wildcard pattern)
- GET /api/docs/endpoints     - API docs endpoints (not authenticated - wildcard pattern)
- GET /api/users/123/profile  - User 123 profile (not authenticated - parameter pattern)
- GET /api/users/456/profile  - User 456 profile (not authenticated - parameter pattern)
- GET /api/profile            - User profile (authenticated)

The auth middleware is configured to skip authentication for:
1. Exact path match: "/public"
2. Wildcard pattern: "/api/docs/*" (all paths starting with /api/docs/)
3. Parameter pattern: "/api/users/:id/profile" (paths like /api/users/123/profile)

Try accessing the protected route with and without authentication.
For JWT authentication, use the Authorization header:
  Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSJ9.<signature>
`
		c.String(http.StatusOK, helpText)
	})

	// Start the server
	fmt.Println("Server running on :8080")
	log.Fatal(srv.Run())
}
