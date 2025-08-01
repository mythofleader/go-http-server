// This example demonstrates how to use the standardized error response structure
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

	// The order of middleware registration is important:
	// 1. Error handler middleware must be registered first to catch errors in other middleware
	// 2. Logging middleware should be registered after error handler to properly capture errors

	// Add error handler middleware using framework-specific implementation
	errorHandler := srv.GetErrorHandlerMiddleware()
	srv.Use(errorHandler.Middleware(nil))

	// Add logging middleware using framework-specific implementation
	loggingMiddleware := srv.GetLoggingMiddleware()
	srv.Use(loggingMiddleware.Middleware(nil))

	// Add routes that demonstrate different error responses
	srv.GET("/", helpHandler)
	srv.GET("/bad-request", badRequestHandler)
	srv.GET("/unauthorized", unauthorizedHandler)
	srv.GET("/forbidden", forbiddenHandler)
	srv.GET("/not-found", notFoundHandler)
	srv.GET("/conflict", conflictHandler)
	srv.GET("/internal-error", internalErrorHandler)
	srv.GET("/service-unavailable", serviceUnavailableHandler)
	srv.GET("/custom-error", customErrorHandler)
	srv.GET("/from-http-error", fromHTTPErrorHandler)
	srv.GET("/error-method", errorMethodHandler)
	srv.GET("/multiple-errors", multipleErrorsHandler)

	// Add routes that demonstrate the new error structs
	srv.GET("/new-bad-request", newBadRequestHandler)
	srv.GET("/new-unauthorized", newUnauthorizedHandler)
	srv.GET("/new-forbidden", newForbiddenHandler)
	srv.GET("/new-not-found", newNotFoundHandler)
	srv.GET("/new-internal-error", newInternalErrorHandler)
	srv.GET("/new-service-unavailable", newServiceUnavailableHandler)

	// Start the server
	fmt.Println("Server running on :8080")
	fmt.Println("Try the following endpoints:")
	fmt.Println("  - http://localhost:8080/ (help)")
	fmt.Println("  - http://localhost:8080/bad-request")
	fmt.Println("  - http://localhost:8080/unauthorized")
	fmt.Println("  - http://localhost:8080/forbidden")
	fmt.Println("  - http://localhost:8080/not-found")
	fmt.Println("  - http://localhost:8080/conflict")
	fmt.Println("  - http://localhost:8080/internal-error")
	fmt.Println("  - http://localhost:8080/service-unavailable")
	fmt.Println("  - http://localhost:8080/custom-error")
	fmt.Println("  - http://localhost:8080/from-http-error")
	fmt.Println("  - http://localhost:8080/error-method")
	fmt.Println("  - http://localhost:8080/multiple-errors")
	fmt.Println("New error structs that embed the error interface:")
	fmt.Println("  - http://localhost:8080/new-bad-request")
	fmt.Println("  - http://localhost:8080/new-unauthorized")
	fmt.Println("  - http://localhost:8080/new-forbidden")
	fmt.Println("  - http://localhost:8080/new-not-found")
	fmt.Println("  - http://localhost:8080/new-internal-error")
	fmt.Println("  - http://localhost:8080/new-service-unavailable")
	fmt.Println("Custom error struct that inherits from BadRequestHttpError:")
	fmt.Println("  - http://localhost:8080/invalid-request-param")
	log.Fatal(srv.Run())
}

// helpHandler provides instructions on how to use the example
func helpHandler(c server.Context) {
	helpText := `
Standardized Error Response Example

This server demonstrates the standardized error response structure.

Try the following endpoints:

- /bad-request - Returns a 400 Bad Request error
- /unauthorized - Returns a 401 Unauthorized error
- /forbidden - Returns a 403 Forbidden error
- /not-found - Returns a 404 Not Found error
- /conflict - Returns a 409 Conflict error
- /internal-error - Returns a 500 Internal Server Error
- /service-unavailable - Returns a 503 Service Unavailable error
- /custom-error - Returns a custom error with a specific status code
- /from-http-error - Returns an error created from an HTTPError
- /error-method - Demonstrates using the Error method of the context
- /multiple-errors - Demonstrates using the Errors method to retrieve all errors

New error structs that embed the error interface:
- /new-bad-request - Uses BadRequestHttpError
- /new-unauthorized - Uses UnauthorizedHttpError
- /new-forbidden - Uses ForbiddenHttpError
- /new-not-found - Uses NotFoundHttpError
- /new-internal-error - Uses InternalServerHttpError
- /new-service-unavailable - Uses ServiceUnavailableHttpError

Custom error struct that inherits from BadRequestHttpError:
- /invalid-request-param - Uses InvalidRequestParamError (demonstrates how to create custom error types)

All responses use the standardized error response structure:
{
  "error": {
    "code": 400,
    "message": "Bad Request"
  }
}
`
	c.String(http.StatusOK, helpText)
}

// badRequestHandler returns a 400 Bad Request error
func badRequestHandler(c server.Context) {
	c.JSON(http.StatusBadRequest, server.NewBadRequestResponse("Invalid request parameters"))
}

// unauthorizedHandler returns a 401 Unauthorized error
func unauthorizedHandler(c server.Context) {
	c.JSON(http.StatusUnauthorized, server.NewUnauthorizedResponse("Authentication required"))
}

// forbiddenHandler returns a 403 Forbidden error
func forbiddenHandler(c server.Context) {
	c.JSON(http.StatusForbidden, server.NewForbiddenResponse("Insufficient permissions"))
}

// notFoundHandler returns a 404 Not Found error
func notFoundHandler(c server.Context) {
	c.JSON(http.StatusNotFound, server.NewNotFoundResponse("Resource not found"))
}

// conflictHandler returns a 409 Conflict error
func conflictHandler(c server.Context) {
	c.JSON(http.StatusConflict, server.NewConflictResponse("Resource already exists"))
}

// internalErrorHandler returns a 500 Internal Server Error
func internalErrorHandler(c server.Context) {
	c.JSON(http.StatusInternalServerError, server.NewInternalServerErrorResponse("An unexpected error occurred"))
}

// serviceUnavailableHandler returns a 503 Service Unavailable error
func serviceUnavailableHandler(c server.Context) {
	c.JSON(http.StatusServiceUnavailable, server.NewServiceUnavailableResponse("Service is currently unavailable"))
}

// customErrorHandler returns a custom error with a specific status code
func customErrorHandler(c server.Context) {
	c.JSON(http.StatusTeapot, server.NewErrorResponse(http.StatusTeapot, "I'm a teapot"))
}

// fromHTTPErrorHandler returns an error created from a BadRequestHttpError
func fromHTTPErrorHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Invalid request")

	// Create an ErrorResponse and return it
	c.JSON(http.StatusBadRequest, server.NewBadRequestResponse(err.Error()))
}

// errorMethodHandler demonstrates using the Error method of the context
func errorMethodHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Invalid request using Error method")

	// Wrap it in a BadRequestHttpError
	httpErr := server.NewBadRequestHttpError(err)

	// Use the Error method of the context to set the error
	c.Error(httpErr)
}

// multipleErrorsHandler demonstrates using the Errors method to retrieve all errors
func multipleErrorsHandler(c server.Context) {
	// Add multiple errors to the context
	c.Error(fmt.Errorf("First error: Something went wrong"))

	// Create and add a BadRequestHttpError
	err2 := fmt.Errorf("Second error: Invalid parameter")
	c.Error(server.NewBadRequestHttpError(err2))

	// Create and add an UnauthorizedHttpError
	err3 := fmt.Errorf("Third error: Authentication required")
	c.Error(server.NewUnauthorizedHttpError(err3))

	// Get all errors from the context
	errors := c.Errors()

	// Create a response with all errors
	response := map[string]interface{}{
		"message":     "Multiple errors example",
		"error_count": len(errors),
		"errors":      make([]string, len(errors)),
	}

	// Add each error message to the response
	for i, err := range errors {
		response["errors"].([]string)[i] = err.Error()
	}

	// Return the response
	c.JSON(http.StatusOK, response)
}

// newBadRequestHandler demonstrates using the new BadRequestHttpError struct
func newBadRequestHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Invalid request parameters")

	// Wrap it in a BadRequestHttpError
	httpErr := server.NewBadRequestHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}

// newUnauthorizedHandler demonstrates using the new UnauthorizedHttpError struct
func newUnauthorizedHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Authentication required")

	// Wrap it in an UnauthorizedHttpError
	httpErr := server.NewUnauthorizedHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}

// newForbiddenHandler demonstrates using the new ForbiddenHttpError struct
func newForbiddenHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Insufficient permissions")

	// Wrap it in a ForbiddenHttpError
	httpErr := server.NewForbiddenHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}

// newNotFoundHandler demonstrates using the new NotFoundHttpError struct
func newNotFoundHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Resource not found")

	// Wrap it in a NotFoundHttpError
	httpErr := server.NewNotFoundHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}

// newInternalErrorHandler demonstrates using the new InternalServerHttpError struct
func newInternalErrorHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("An unexpected error occurred")

	// Wrap it in an InternalServerHttpError
	httpErr := server.NewInternalServerHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}

// newServiceUnavailableHandler demonstrates using the new ServiceUnavailableHttpError struct
func newServiceUnavailableHandler(c server.Context) {
	// Create a standard error
	err := fmt.Errorf("Service is currently unavailable")

	// Wrap it in a ServiceUnavailableHttpError
	httpErr := server.NewServiceUnavailableHttpError(err)

	// Panic with the error to trigger the error handler middleware
	panic(httpErr)
}
