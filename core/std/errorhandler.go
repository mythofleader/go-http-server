// Package std provides a standard HTTP implementation of the HTTP server abstraction.
package std

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware"
	tErrors "github.com/mythofleader/go-http-server/core/middleware/errors"
)

// ErrorHandlerMiddleware is a standard HTTP implementation of middleware.IErrorHandlerMiddleware.
type ErrorHandlerMiddleware struct{}

// Middleware returns a middleware function that handles errors for standard HTTP.
func (m *ErrorHandlerMiddleware) Middleware(config *core.ErrorHandlerConfig) core.HandlerFunc {
	if config == nil {
		config = middleware.DefaultErrorHandlerConfig()
	}

	return func(c core.Context) {
		// Get the standard HTTP context
		stdContext, ok := c.(*Context)
		if !ok {
			// Handle the case when it's not a standard HTTP context
			// Create a recovery function to catch panics
			defer func() {
				if r := recover(); r != nil {
					// Handle panic
					var err error
					switch e := r.(type) {
					case string:
						err = tErrors.NewInternalServerHttpError(fmt.Errorf("%s", e))
					case error:
						err = tErrors.NewInternalServerHttpError(e)
					default:
						err = tErrors.NewInternalServerHttpError(fmt.Errorf("unknown error: %v", e))
					}

					handleError(c, err, config)
				}
			}()

			// Continue with the next handler
			c.Next()

			// Check if there are any errors
			if errs := c.Errors(); len(errs) > 0 {
				handleError(c, errs[0], config)
			}
			return
		}

		// Create a recovery function to catch panics
		defer func() {
			if r := recover(); r != nil {
				// Handle panic
				var err error
				switch e := r.(type) {
				case string:
					err = tErrors.NewInternalServerHttpError(fmt.Errorf("%s", e))
				case error:
					err = tErrors.NewInternalServerHttpError(e)
				default:
					err = tErrors.NewInternalServerHttpError(fmt.Errorf("unknown error: %v", e))
				}

				// Handle the error based on its type
				handleError(c, err, config)
			}
		}()

		// Create a wrapper for the response writer to capture errors
		errorWriter := &errorCaptureWriter{
			ResponseWriter: stdContext.writer,
			statusCode:     http.StatusOK,
			err:            nil,
		}

		// Replace the original writer with the wrapped one
		stdContext.writer = errorWriter

		// Continue with the next middleware/handler in the chain
		c.Next()

		// Check if an error was captured
		if errorWriter.err != nil {
			// Handle the error based on its type
			handleError(c, errorWriter.err, config)
		}
	}
}

// handleError processes an error and returns an appropriate HTTP response.
func handleError(c core.Context, err error, config *core.ErrorHandlerConfig) {
	var httpErr tErrors.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode(), tErrors.NewErrorResponse(httpErr.StatusCode(), httpErr.Error()))
		return
	}
	c.JSON(config.DefaultStatusCode, tErrors.NewErrorResponse(config.DefaultStatusCode, config.DefaultErrorMessage))
}

// errorCaptureWriter is a wrapper for http.ResponseWriter that captures errors.
type errorCaptureWriter struct {
	http.ResponseWriter
	statusCode int
	err        error
}

// WriteHeader captures the status code and calls the underlying ResponseWriter's WriteHeader.
func (w *errorCaptureWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures errors based on the status code and calls the underlying ResponseWriter's Write.
func (w *errorCaptureWriter) Write(b []byte) (int, error) {
	// If the status code hasn't been set yet, set it to 200 OK
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	// If the status code indicates an error, capture it
	if w.statusCode >= 400 {
		switch w.statusCode {
		case http.StatusBadRequest:
			w.err = tErrors.NewBadRequestHttpError(fmt.Errorf("%s", string(b)))
		case http.StatusUnauthorized:
			w.err = tErrors.NewUnauthorizedHttpError(fmt.Errorf("%s", string(b)))
		case http.StatusForbidden:
			w.err = tErrors.NewForbiddenHttpError(fmt.Errorf("%s", string(b)))
		case http.StatusInternalServerError:
			w.err = tErrors.NewInternalServerHttpError(fmt.Errorf("%s", string(b)))
		default:
			w.err = fmt.Errorf("HTTP error: %d - %s", w.statusCode, string(b))
		}
	}
	return w.ResponseWriter.Write(b)
}

// SetError sets an error on the writer.
func (w *errorCaptureWriter) SetError(err error) {
	w.err = err
}

// NewErrorHandlerMiddleware creates a new ErrorHandlerMiddleware.
func NewErrorHandlerMiddleware() middleware.IErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}
