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

	// Add a public route
	srv.GET("/", func(c server.Context) {
		c.String(http.StatusOK, "Welcome to the API")
	})

	// Start the server
	fmt.Println("Server running on :8080")
	log.Fatal(srv.Run())
}
