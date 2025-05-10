// Package errors provides error classes and response structures for HTTP status codes.
package errors

import (
	"net/http"
)

// ErrorDetail represents the structure of an error detail in the response.
type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse with the given status code and message.
func NewErrorResponse(statusCode int, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:    statusCode,
			Message: message,
		},
	}
}

// NewBadRequestResponse creates a new ErrorResponse for a 400 Bad Request error.
func NewBadRequestResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Bad Request"
	}
	return NewErrorResponse(http.StatusBadRequest, message)
}

// NewUnauthorizedResponse creates a new ErrorResponse for a 401 Unauthorized error.
func NewUnauthorizedResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Unauthorized"
	}
	return NewErrorResponse(http.StatusUnauthorized, message)
}

// NewForbiddenResponse creates a new ErrorResponse for a 403 Forbidden error.
func NewForbiddenResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Forbidden"
	}
	return NewErrorResponse(http.StatusForbidden, message)
}

// NewNotFoundResponse creates a new ErrorResponse for a 404 Not Found error.
func NewNotFoundResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Not Found"
	}
	return NewErrorResponse(http.StatusNotFound, message)
}

// NewConflictResponse creates a new ErrorResponse for a 409 Conflict error.
func NewConflictResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Conflict"
	}
	return NewErrorResponse(http.StatusConflict, message)
}

// NewInternalServerErrorResponse creates a new ErrorResponse for a 500 Internal Server Error.
func NewInternalServerErrorResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Internal Server Error"
	}
	return NewErrorResponse(http.StatusInternalServerError, message)
}

// NewServiceUnavailableResponse creates a new ErrorResponse for a 503 Service Unavailable error.
func NewServiceUnavailableResponse(message string) *ErrorResponse {
	if message == "" {
		message = "Service Unavailable"
	}
	return NewErrorResponse(http.StatusServiceUnavailable, message)
}
