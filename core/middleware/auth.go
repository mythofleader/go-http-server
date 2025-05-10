// Package middleware provides common middleware functionality for HTTP servers.
package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mythofleader/go-http-server/core"
	httperrors "github.com/mythofleader/go-http-server/core/middleware/errors"
)

// MapClaims represents JWT claims as a map
type MapClaims map[string]interface{}

// BasicAuthUserLookup defines the interface for looking up users based on Basic Auth credentials
type BasicAuthUserLookup interface {
	// LookupUserByBasicAuth looks up a user by username and password
	LookupUserByBasicAuth(username, password string) (interface{}, error)
}

// JWTUserLookup defines the interface for looking up users based on JWT claims
type JWTUserLookup interface {
	// LookupUserByJWT looks up a user by JWT claims
	LookupUserByJWT(claims MapClaims) (interface{}, error)
}

// UserLookupInterface defines the interface for looking up users based on credentials
// This is kept for backward compatibility
type UserLookupInterface interface {
	BasicAuthUserLookup
	JWTUserLookup
}

// AuthType represents the type of authentication to use
type AuthType string

const (
	// AuthTypeBasic represents HTTP Basic authentication
	AuthTypeBasic AuthType = "basic"
	// AuthTypeJWT represents JWT Bearer token authentication
	AuthTypeJWT AuthType = "jwt"
)

// AuthConfig holds configuration for the authorization middleware
type AuthConfig struct {
	// UserLookup is the implementation of UserLookupInterface
	// This is kept for backward compatibility
	UserLookup UserLookupInterface

	// BasicAuthLookup is the implementation of BasicAuthUserLookup
	// Used when AuthType is AuthTypeBasic
	BasicAuthLookup BasicAuthUserLookup

	// JWTLookup is the implementation of JWTUserLookup
	// Used when AuthType is AuthTypeJWT
	JWTLookup JWTUserLookup

	// AuthType specifies which authentication method to use (basic or jwt)
	// If not specified, it defaults to jwt
	AuthType AuthType

	// JWTSecret is the secret key used to validate JWT tokens
	// Required when AuthType is AuthTypeJWT
	JWTSecret string

	// Optional: custom error messages
	UnauthorizedMessage string
	ForbiddenMessage    string

	// SkipPaths is a list of paths to ignore for authentication
	SkipPaths []string
}

// DefaultAuthConfig returns a default auth configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		AuthType:            AuthTypeJWT, // Default to JWT authentication
		UnauthorizedMessage: "Unauthorized",
		ForbiddenMessage:    "Forbidden",
		SkipPaths:           []string{},
		// UserLookup, BasicAuthLookup, and JWTLookup are nil by default
		// and must be provided by the user
	}
}

// NewDefaultJWTAuthMiddleware returns a middleware function with default JWT authentication configuration.
// This function creates a default configuration with AuthType set to AuthTypeJWT and sets the provided parameters.
// Example usage:
//
//	s.Use(middleware.NewDefaultJWTAuthMiddleware(myJWTLookup, "your-jwt-secret"))
//
// Or customize the configuration:
//
//	config := middleware.DefaultAuthConfig()
//	config.AuthType = middleware.AuthTypeJWT
//	config.JWTLookup = myJWTLookup
//	config.JWTSecret = "your-jwt-secret"
//	config.UnauthorizedMessage = "Custom unauthorized message"
//	s.Use(middleware.AuthMiddleware(config))
func NewDefaultJWTAuthMiddleware(jwtLookup JWTUserLookup, jwtSecret string) core.HandlerFunc {
	config := DefaultAuthConfig()
	config.AuthType = AuthTypeJWT
	config.JWTLookup = jwtLookup
	config.JWTSecret = jwtSecret
	return AuthMiddleware(config)
}

// NewDefaultBasicAuthMiddleware returns a middleware function with default Basic authentication configuration.
// This function creates a default configuration with AuthType set to AuthTypeBasic and sets the provided parameter.
// Example usage:
//
//	s.Use(middleware.NewDefaultBasicAuthMiddleware(myBasicAuthLookup))
//
// Or customize the configuration:
//
//	config := middleware.DefaultAuthConfig()
//	config.AuthType = middleware.AuthTypeBasic
//	config.BasicAuthLookup = myBasicAuthLookup
//	config.UnauthorizedMessage = "Custom unauthorized message"
//	s.Use(middleware.AuthMiddleware(config))
func NewDefaultBasicAuthMiddleware(basicAuthLookup BasicAuthUserLookup) core.HandlerFunc {
	config := DefaultAuthConfig()
	config.AuthType = AuthTypeBasic
	config.BasicAuthLookup = basicAuthLookup
	return AuthMiddleware(config)
}

// AuthMiddleware returns a middleware function that checks authorization
// It supports either Basic HTTP authentication or Bearer JWT tokens based on the configuration
func AuthMiddleware(config *AuthConfig) core.HandlerFunc {
	if config == nil {
		config = DefaultAuthConfig()
	}

	// Validate the configuration based on the authentication type
	switch config.AuthType {
	case AuthTypeBasic:
		// For Basic authentication, we need either UserLookup or BasicAuthLookup
		if config.UserLookup == nil && config.BasicAuthLookup == nil {
			panic("AuthMiddleware with AuthTypeBasic requires either UserLookup or BasicAuthLookup implementation")
		}
	case AuthTypeJWT:
		// For JWT authentication, we need either UserLookup or JWTLookup
		if config.UserLookup == nil && config.JWTLookup == nil {
			panic("AuthMiddleware with AuthTypeJWT requires either UserLookup or JWTLookup implementation")
		}
		// Also check for JWTSecret
		if config.JWTSecret == "" {
			panic("JWTSecret is required when using JWT authentication")
		}
	default:
		panic("Invalid AuthType specified")
	}

	return func(c core.Context) {
		// Get request path
		path := c.Request().URL.Path

		// Check if the path is in the skip paths list
		for _, ignorePath := range config.SkipPaths {
			if path == ignorePath {
				// Skip authentication for this path
				return
			}
		}

		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.SetStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, httperrors.NewUnauthorizedResponse(config.UnauthorizedMessage))
			return
		}

		// Split the header into type and credentials
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.SetStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, httperrors.NewUnauthorizedResponse("Invalid authorization format"))
			return
		}

		authType := parts[0]
		credentials := parts[1]

		var user interface{}
		var err error

		// Handle the authentication based on the configured type
		switch config.AuthType {
		case AuthTypeBasic:
			// Only accept Basic authentication
			if authType != "Basic" {
				c.SetStatus(http.StatusUnauthorized)
				c.JSON(http.StatusUnauthorized, httperrors.NewUnauthorizedResponse("Basic authentication required"))
				return
			}

			// Use the appropriate lookup interface
			var basicLookup BasicAuthUserLookup
			if config.BasicAuthLookup != nil {
				basicLookup = config.BasicAuthLookup
			} else {
				// Fall back to UserLookup for backward compatibility
				basicLookup = config.UserLookup
			}

			user, err = handleBasicAuth(credentials, basicLookup)
		case AuthTypeJWT:
			// Only accept Bearer token authentication
			if authType != "Bearer" {
				c.SetStatus(http.StatusUnauthorized)
				c.JSON(http.StatusUnauthorized, httperrors.NewUnauthorizedResponse("Bearer token required"))
				return
			}

			// Use the appropriate lookup interface
			var jwtLookup JWTUserLookup
			if config.JWTLookup != nil {
				jwtLookup = config.JWTLookup
			} else {
				// Fall back to UserLookup for backward compatibility
				jwtLookup = config.UserLookup
			}

			user, err = handleBearerToken(credentials, config.JWTSecret, jwtLookup)
		default:
			c.SetStatus(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, httperrors.NewInternalServerErrorResponse("Invalid authentication configuration"))
			return
		}

		if err != nil {
			statusCode := http.StatusUnauthorized
			message := config.UnauthorizedMessage

			if errors.Is(err, ErrForbidden) {
				statusCode = http.StatusForbidden
				message = config.ForbiddenMessage
			}

			c.SetStatus(statusCode)
			if statusCode == http.StatusUnauthorized {
				c.JSON(statusCode, httperrors.NewUnauthorizedResponse(message))
			} else {
				c.JSON(statusCode, httperrors.NewForbiddenResponse(message))
			}
			return
		}

		// Store the user in the request context for later use
		req := c.Request()
		newCtx := context.WithValue(req.Context(), UserContextKey, user)

		// Create a new request with the updated context
		newReq := req.WithContext(newCtx)

		// Update the request in the context
		*req = *newReq
	}
}

// UserContextKey is the key used to store the user in the request context
type contextKey string

// Define the context key for the user
const UserContextKey contextKey = "user"

// ErrForbidden is returned when the user is authenticated but not authorized
var ErrForbidden = errors.New("forbidden")

// handleBasicAuth processes HTTP Basic Authentication
func handleBasicAuth(credentials string, lookup BasicAuthUserLookup) (interface{}, error) {
	decoded, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return nil, errors.New("invalid basic auth format")
	}

	userPass := string(decoded)
	parts := strings.SplitN(userPass, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid basic auth format")
	}

	username := parts[0]
	password := parts[1]

	user, err := lookup.LookupUserByBasicAuth(username, password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return user, nil
}

// handleBearerToken processes JWT Bearer tokens
func handleBearerToken(tokenString string, secret string, lookup JWTUserLookup) (interface{}, error) {
	// Parse and validate the JWT token
	claims, err := parseJWT(tokenString, secret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Look up the user based on the JWT claims
	user, err := lookup.LookupUserByJWT(claims)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return user, nil
}

// parseJWT parses and validates a JWT token
func parseJWT(tokenString string, secret string) (MapClaims, error) {
	// Split the token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Decode the header
	headerJSON, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid token header: %w", err)
	}

	// Parse the header
	var header map[string]interface{}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("invalid token header: %w", err)
	}

	// Check the algorithm
	alg, ok := header["alg"].(string)
	if !ok || alg != "HS256" {
		return nil, errors.New("unsupported signing method")
	}

	// Decode the payload
	payloadJSON, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token payload: %w", err)
	}

	// Parse the payload
	var claims MapClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid token payload: %w", err)
	}

	// Verify the signature
	signatureBytes, err := base64URLDecode(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid token signature: %w", err)
	}

	// Create the signature
	signatureString := parts[0] + "." + parts[1]
	expectedSignature := createHmacSignature(signatureString, secret)

	// Compare the signatures
	if !hmac.Equal(signatureBytes, expectedSignature) {
		return nil, errors.New("invalid token signature")
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token expired")
		}
	}

	return claims, nil
}

// base64URLDecode decodes a base64url encoded string
func base64URLDecode(s string) ([]byte, error) {
	// Add padding if needed
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}
	// Replace URL encoding with standard base64 encoding
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	return base64.StdEncoding.DecodeString(s)
}

// createHmacSignature creates an HMAC signature for a JWT token
func createHmacSignature(data, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return h.Sum(nil)
}

// GetUserFromContext retrieves the authenticated user from the context
func GetUserFromContext(ctx context.Context) (interface{}, bool) {
	user := ctx.Value(UserContextKey)
	if user == nil {
		return nil, false
	}
	return user, true
}
