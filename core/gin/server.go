// Package gin provides a Gin implementation of the HTTP server abstraction.
package gin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/mythofleader/go-http-server/core"
	"github.com/mythofleader/go-http-server/core/middleware/errors"
)

// Context is an implementation of core.Context using the Gin framework.
type Context struct {
	ginContext *gin.Context
}

// Request implements core.Context.Request
func (c *Context) Request() *http.Request {
	return c.ginContext.Request
}

// Writer implements core.Context.Writer
func (c *Context) Writer() http.ResponseWriter {
	return c.ginContext.Writer
}

// Param implements core.Context.Param
func (c *Context) Param(key string) string {
	return c.ginContext.Param(key)
}

// Query implements core.Context.Query
func (c *Context) Query(key string) string {
	return c.ginContext.Query(key)
}

// DefaultQuery implements core.Context.DefaultQuery
func (c *Context) DefaultQuery(key, defaultValue string) string {
	return c.ginContext.DefaultQuery(key, defaultValue)
}

// GetHeader implements core.Context.GetHeader
func (c *Context) GetHeader(key string) string {
	return c.ginContext.GetHeader(key)
}

// SetHeader implements core.Context.SetHeader
func (c *Context) SetHeader(key, value string) {
	c.ginContext.Header(key, value)
}

// SetStatus implements core.Context.SetStatus
func (c *Context) SetStatus(code int) {
	c.ginContext.Status(code)
}

// JSON implements core.Context.JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.ginContext.JSON(code, obj)
}

// String implements core.Context.String
func (c *Context) String(code int, format string, values ...interface{}) {
	c.ginContext.String(code, format, values...)
}

// Bind implements core.Context.Bind
func (c *Context) Bind(obj interface{}) error {
	return c.ginContext.Bind(obj)
}

// BindJSON implements core.Context.BindJSON
func (c *Context) BindJSON(obj interface{}) error {
	return c.ginContext.BindJSON(obj)
}

// ShouldBindJSON implements core.Context.ShouldBindJSON
func (c *Context) ShouldBindJSON(obj interface{}) error {
	return c.ginContext.ShouldBindJSON(obj)
}

// File implements core.Context.File
func (c *Context) File(filepath string) {
	c.ginContext.File(filepath)
}

// Redirect implements core.Context.Redirect
func (c *Context) Redirect(code int, location string) {
	c.ginContext.Redirect(code, location)
}

// Error implements core.Context.Error
func (c *Context) Error(err error) error {
	return c.ginContext.Error(err)
}

// Errors implements core.Context.Errors
func (c *Context) Errors() []error {
	if len(c.ginContext.Errors) == 0 {
		return nil
	}

	errors := make([]error, len(c.ginContext.Errors))
	for i, err := range c.ginContext.Errors {
		errors[i] = err.Err
	}
	return errors
}

// Next implements core.Context.Next
func (c *Context) Next() {
	c.ginContext.Next()
}

// Abort implements core.Context.Abort
func (c *Context) Abort() {
	c.ginContext.Abort()
}

// Get implements core.Context.Get
func (c *Context) Get(key string) (interface{}, bool) {
	value, exists := c.ginContext.Get(key)
	return value, exists
}

// Set implements core.Context.Set
func (c *Context) Set(key string, value interface{}) {
	c.ginContext.Set(key, value)
}

// Server is an implementation of core.Server using the Gin framework.
type Server struct {
	engine      *gin.Engine
	server      *http.Server
	port        string
	middlewares []string // Track middleware names
	showLogs    bool     // Controls whether framework logs are shown
}

// GetLoggingMiddleware returns a Gin-specific logging middleware.
func (s *Server) GetLoggingMiddleware() core.ILoggingMiddleware {
	return NewLoggingMiddleware()
}

// GetErrorHandlerMiddleware returns a Gin-specific error handler middleware.
func (s *Server) GetErrorHandlerMiddleware() core.IErrorHandlerMiddleware {
	return NewErrorHandlerMiddleware()
}

// RouterGroup is an implementation of core.RouterGroup using the Gin framework.
type RouterGroup struct {
	group *gin.RouterGroup
}

// GET implements core.Server.GET
func (s *Server) GET(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.GET(path, ginHandlers...)
}

// POST implements core.Server.POST
func (s *Server) POST(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.POST(path, ginHandlers...)
}

// PUT implements core.Server.PUT
func (s *Server) PUT(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.PUT(path, ginHandlers...)
}

// DELETE implements core.Server.DELETE
func (s *Server) DELETE(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.DELETE(path, ginHandlers...)
}

// PATCH implements core.Server.PATCH
func (s *Server) PATCH(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.PATCH(path, ginHandlers...)
}

// Group implements core.Server.Group
func (s *Server) Group(path string) core.RouterGroup {
	return &RouterGroup{
		group: s.engine.Group(path),
	}
}

// Use implements core.Server.Use
func (s *Server) Use(middleware ...core.HandlerFunc) {
	for _, m := range middleware {
		// Get the function name for logging
		funcValue := reflect.ValueOf(m)
		middlewareName := runtime.FuncForPC(funcValue.Pointer()).Name()
		s.middlewares = append(s.middlewares, middlewareName)

		// Log middleware addition if showLogs is true
		if s.showLogs {
			log.Printf("[GIN] Adding middleware: %s", middlewareName)
		}

		s.engine.Use(wrapHandler(m))
	}
}

// RegisterRouter implements core.Server.RegisterRouter
func (s *Server) RegisterRouter(controllers ...core.Controller) {
	for _, controller := range controllers {
		// Get HTTP method, path, and handlers from the controller
		method := controller.GetHttpMethod()
		path := controller.GetPath()
		handlers := controller.Handler()

		// Register the route based on the HTTP method
		switch method {
		case core.GET:
			s.GET(path, handlers...)
		case core.POST:
			s.POST(path, handlers...)
		case core.PUT:
			s.PUT(path, handlers...)
		case core.DELETE:
			s.DELETE(path, handlers...)
		case core.PATCH:
			s.PATCH(path, handlers...)
		}

		// Log controller registration if showLogs is true
		if s.showLogs {
			log.Printf("[GIN] Registered controller with method: %s, path: %s, skip logging: %t, skip auth check: %t",
				method, path, controller.SkipLogging(), controller.SkipAuthCheck())
		}
	}
}

// NoRoute implements core.Server.NoRoute
func (s *Server) NoRoute(handlers ...core.HandlerFunc) {
	// If no handlers are provided, use default handler
	if len(handlers) == 0 {
		// Default handler returns a 404 Not Found error
		handlers = []core.HandlerFunc{
			func(c core.Context) {
				path := c.Request().URL.Path
				err := fmt.Errorf("route not found: %s", path)
				_ = c.Error(errors.NewNotFoundHttpError(err))
			},
		}
		if s.showLogs {
			log.Printf("[GIN] Using default NoRoute handler")
		}
	}

	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.NoRoute(ginHandlers...)
	if s.showLogs {
		log.Printf("[GIN] Registered NoRoute handler")
	}
}

// NoMethod implements core.Server.NoMethod
func (s *Server) NoMethod(handlers ...core.HandlerFunc) {
	// If no handlers are provided, use default handler
	if len(handlers) == 0 {
		// Default handler returns a 405 Method Not Allowed error
		handlers = []core.HandlerFunc{
			func(c core.Context) {
				method := c.Request().Method
				path := c.Request().URL.Path
				err := fmt.Errorf("method %s not allowed for path %s", method, path)
				_ = c.Error(errors.NewMethodNotAllowedHttpError(err))
			},
		}
		if s.showLogs {
			log.Printf("[GIN] Using default NoMethod handler")
		}
	}

	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	s.engine.NoMethod(ginHandlers...)
	if s.showLogs {
		log.Printf("[GIN] Registered NoMethod handler")
	}
}

// Run implements core.Server.Run
func (s *Server) Run() error {
	addr := ":" + s.port

	// Log server information if showLogs is true
	if s.showLogs {
		log.Printf("[GIN] Server starting on %s", addr)
		log.Printf("[GIN] Using Gin framework version: %s", gin.Version)

		// Log middleware information
		if len(s.middlewares) > 0 {
			log.Println("[GIN] Middleware registered:")
			for i, middleware := range s.middlewares {
				log.Printf("[GIN]   %d. %s", i+1, middleware)
			}
		} else {
			log.Println("[GIN] No middleware registered")
		}
	}

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	// Log routes information if showLogs is true
	if s.showLogs {
		routes := s.engine.Routes()
		if len(routes) > 0 {
			log.Println("[GIN] Routes registered:")
			for i, route := range routes {
				log.Printf("[GIN]   %d. %s %s", i+1, route.Method, route.Path)
			}
		} else {
			log.Println("[GIN] No routes registered")
		}

		log.Printf("[GIN] Server is ready to handle requests")
	}

	return s.engine.Run(addr)
}

// RunTLS implements core.Server.RunTLS
func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	return s.server.ListenAndServeTLS(certFile, keyFile)
}

// Stop implements core.Server.Stop
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}
	return s.server.Close()
}

// Shutdown implements core.Server.Shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// GetPort implements core.Server.GetPort
func (s *Server) GetPort() string {
	return s.port
}

// StartLambda starts the server in AWS Lambda mode.
// This method should be called instead of Run or RunTLS when running in AWS Lambda.
// This method uses the ginadapter library to convert the Gin engine to a Lambda handler.
//
// Example usage:
//
//	import (
//	    "github.com/mythofleader/go-http-server"
//	)
//
//	func main() {
//	    s, _ := server.NewServer(server.FrameworkGin, "8080")
//	    // ... configure your server ...
//	    if err := s.StartLambda(); err != nil {
//	        // Handle error
//	    }
//	}
func (s *Server) StartLambda() error {
	// Create a new ALB adapter for the Gin engine
	ginLambda := ginadapter.NewALB(s.engine)

	// Start the Lambda handler
	lambda.Start(func(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
		// Process the request
		return ginLambda.ProxyWithContext(ctx, req)
	})

	// This line is never reached because lambda.Start() doesn't return
	return nil
}

// GET implements core.RouterGroup.GET
func (g *RouterGroup) GET(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	g.group.GET(path, ginHandlers...)
}

// POST implements core.RouterGroup.POST
func (g *RouterGroup) POST(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	g.group.POST(path, ginHandlers...)
}

// PUT implements core.RouterGroup.PUT
func (g *RouterGroup) PUT(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	g.group.PUT(path, ginHandlers...)
}

// DELETE implements core.RouterGroup.DELETE
func (g *RouterGroup) DELETE(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	g.group.DELETE(path, ginHandlers...)
}

// PATCH implements core.RouterGroup.PATCH
func (g *RouterGroup) PATCH(path string, handlers ...core.HandlerFunc) {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = wrapHandler(handler)
	}
	g.group.PATCH(path, ginHandlers...)
}

// Group implements core.RouterGroup.Group
func (g *RouterGroup) Group(path string) core.RouterGroup {
	return &RouterGroup{
		group: g.group.Group(path),
	}
}

// Use implements core.RouterGroup.Use
func (g *RouterGroup) Use(middleware ...core.HandlerFunc) {
	for _, m := range middleware {
		g.group.Use(wrapHandler(m))
	}
}

// RegisterRouter implements core.RouterGroup.RegisterRouter
func (g *RouterGroup) RegisterRouter(controllers ...core.Controller) {
	for _, controller := range controllers {
		// Get HTTP method, path, and handlers from the controller
		method := controller.GetHttpMethod()
		path := controller.GetPath()
		handlers := controller.Handler()

		// Register the route based on the HTTP method
		switch method {
		case core.GET:
			g.GET(path, handlers...)
		case core.POST:
			g.POST(path, handlers...)
		case core.PUT:
			g.PUT(path, handlers...)
		case core.DELETE:
			g.DELETE(path, handlers...)
		case core.PATCH:
			g.PATCH(path, handlers...)
		}

		// Log controller registration
		log.Printf("[GIN] Registered controller with method: %s, path: %s, skip logging: %t, skip auth check: %t",
			method, path, controller.SkipLogging(), controller.SkipAuthCheck())
	}
}

// wrapHandler wraps a core.HandlerFunc to a gin.HandlerFunc
func wrapHandler(handler core.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&Context{ginContext: c})
	}
}

// NewServer creates a new Server instance using the Gin framework.
// If showLogs is true, logs about the framework, middleware, and routes will be printed to the console.
// If showLogs is false, these logs will be suppressed.
func NewServer(port string, showLogs bool) *Server {
	gin.SetMode(gin.ReleaseMode)

	// Only log if showLogs is true
	if showLogs {
		log.Printf("[GIN] Creating new Gin server on port %s", port)
	}

	return &Server{
		engine:      gin.New(),
		port:        port,
		middlewares: make([]string, 0),
		showLogs:    showLogs,
	}
}
