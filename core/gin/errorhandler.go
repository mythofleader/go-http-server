// Package gin provides a Gin implementation of the HTTP server abstraction.
package gin

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware"
	tErrors "github.com/mythofleader/go-http-server/core/middleware/errors"
)

// ErrorHandlerMiddleware is a Gin-specific implementation of middleware.IErrorHandlerMiddleware.
type ErrorHandlerMiddleware struct {
	// This is just to make the linter happy about the gin import
	_ gin.HandlerFunc
}

// Middleware returns a middleware function that handles errors for Gin.
func (m *ErrorHandlerMiddleware) Middleware(config *core.ErrorHandlerConfig) core.HandlerFunc {
	if config == nil {
		config = middleware.DefaultErrorHandlerConfig()
	}

	return func(c core.Context) {
		// Get the Gin context
		ginContext, ok := c.(*Context)
		if !ok {
			// Handle the case when it's not a Gin context
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

		// Get the underlying gin.Context
		gc := ginContext.ginContext

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

				// Abort the request
				gc.Abort()
			}
		}()

		// Use Gin's built-in error handling
		gc.Next()

		// Check if there are any errors
		if len(gc.Errors) > 0 {
			// Get the first error
			handleError(c, gc.Errors[0].Err, config)
			// Abort the request
			gc.Abort()
		}
	}
}

func handleError(c core.Context, err error, config *core.ErrorHandlerConfig) {
	var httpErr tErrors.HTTPError
	if errors.As(err, &httpErr) {
		c.JSON(httpErr.StatusCode(), tErrors.NewErrorResponse(httpErr.StatusCode(), httpErr.Error()))
		return
	}
	c.JSON(config.DefaultStatusCode, tErrors.NewErrorResponse(config.DefaultStatusCode, config.DefaultErrorMessage))
}

// NewErrorHandlerMiddleware creates a new ErrorHandlerMiddleware.
func NewErrorHandlerMiddleware() middleware.IErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}
