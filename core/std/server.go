// Package std provides a standard HTTP implementation of the HTTP server abstraction.
package std

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync"

	"github.com/mythofleader/go-http-server/core"
	httperrors "github.com/mythofleader/go-http-server/core/middleware/errors"
)

// Context is an implementation of core.Context using the standard net/http package.
type Context struct {
	req        *http.Request
	writer     http.ResponseWriter
	params     map[string]string
	queryCache map[string]string
	errs       []error                // Errors that occurred during request processing
	keys       map[string]interface{} // Key-value store for context data
	mu         sync.RWMutex           // Mutex to protect concurrent access to keys

	// Fields for middleware flow control
	handlers     []core.HandlerFunc // All handlers (middleware + route handlers)
	index        int                // Current handler index
	handlerCount int                // Total number of handlers
}

// Request implements core.Context.Request
func (c *Context) Request() *http.Request {
	return c.req
}

// Writer implements core.Context.Writer
func (c *Context) Writer() http.ResponseWriter {
	return c.writer
}

// Param implements core.Context.Param
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Query implements core.Context.Query
func (c *Context) Query(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.queryCache == nil {
		c.queryCache = make(map[string]string)
	}
	if val, ok := c.queryCache[key]; ok {
		return val
	}
	val := c.req.URL.Query().Get(key)
	c.queryCache[key] = val
	return val
}

// DefaultQuery implements core.Context.DefaultQuery
func (c *Context) DefaultQuery(key, defaultValue string) string {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// GetHeader implements core.Context.GetHeader
func (c *Context) GetHeader(key string) string {
	return c.req.Header.Get(key)
}

// SetHeader implements core.Context.SetHeader
func (c *Context) SetHeader(key, value string) {
	c.writer.Header().Set(key, value)
}

// SetStatus implements core.Context.SetStatus
func (c *Context) SetStatus(code int) {
	c.writer.WriteHeader(code)
}

// JSON implements core.Context.JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	// Use a JSON encoder to write the response
	if err := json.NewEncoder(c.writer).Encode(obj); err != nil {
		http.Error(c.writer, err.Error(), http.StatusInternalServerError)
	}
}

// String implements core.Context.String
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	fmt.Fprintf(c.writer, format, values...)
}

// Bind implements core.Context.Bind
func (c *Context) Bind(obj interface{}) error {
	// This is a simplified implementation
	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" {
		return c.BindJSON(obj)
	}
	return errors.New("unsupported content type")
}

// BindJSON implements core.Context.BindJSON
func (c *Context) BindJSON(obj interface{}) error {
	return json.NewDecoder(c.req.Body).Decode(obj)
}

// ShouldBindJSON implements core.Context.ShouldBindJSON
func (c *Context) ShouldBindJSON(obj interface{}) error {
	return json.NewDecoder(c.req.Body).Decode(obj)
}

// File implements core.Context.File
func (c *Context) File(filepath string) {
	http.ServeFile(c.writer, c.req, filepath)
}

// Redirect implements core.Context.Redirect
func (c *Context) Redirect(code int, location string) {
	http.Redirect(c.writer, c.req, location, code)
}

// Error implements core.Context.Error
// Since the standard HTTP package doesn't have a built-in error handling mechanism,
// this implementation stores the error in the context and returns it.
func (c *Context) Error(err error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Initialize the errs slice if it's nil
	if c.errs == nil {
		c.errs = make([]error, 0)
	}

	// Add the error to the errs slice
	c.errs = append(c.errs, err)

	// Return the error
	return err
}

// Errors implements core.Context.Errors
// It returns all errors added to the context.
func (c *Context) Errors() []error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.errs
}

// Next implements core.Context.Next
// It calls the next handler in the chain.
func (c *Context) Next() {
	c.index++
	for c.index < c.handlerCount {
		c.handlers[c.index](c)
		c.index++
	}
}

// Get implements core.Context.Get
// It returns the value for the given key and a boolean indicating whether the key exists.
func (c *Context) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.keys == nil {
		return nil, false
	}
	value, exists := c.keys[key]
	return value, exists
}

// Set implements core.Context.Set
// It stores a value in the context for the given key.
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// Server is an implementation of core.Server using the standard net/http package.
type Server struct {
	mux              *http.ServeMux
	server           *http.Server
	routes           map[string]map[string][]core.HandlerFunc // method -> path -> handlers
	middleware       []core.HandlerFunc
	port             string
	middlewareLog    []string           // Track middleware names for logging
	noRouteHandlers  []core.HandlerFunc // Handlers for 404 Not Found errors
	noMethodHandlers []core.HandlerFunc // Handlers for 405 Method Not Allowed errors
}

// GetLoggingMiddleware returns a standard HTTP-specific logging middleware.
func (s *Server) GetLoggingMiddleware() core.ILoggingMiddleware {
	return NewLoggingMiddleware()
}

// GetErrorHandlerMiddleware returns a standard HTTP-specific error handler middleware.
func (s *Server) GetErrorHandlerMiddleware() core.IErrorHandlerMiddleware {
	return NewErrorHandlerMiddleware()
}

// GET implements core.Server.GET for Server
func (s *Server) GET(path string, handlers ...core.HandlerFunc) {
	if s.routes == nil {
		s.routes = make(map[string]map[string][]core.HandlerFunc)
	}
	if s.routes["GET"] == nil {
		s.routes["GET"] = make(map[string][]core.HandlerFunc)
	}
	s.routes["GET"][path] = handlers
	s.mux.HandleFunc(path, s.handleHTTP("GET", path))
}

// POST implements core.Server.POST for Server
func (s *Server) POST(path string, handlers ...core.HandlerFunc) {
	if s.routes == nil {
		s.routes = make(map[string]map[string][]core.HandlerFunc)
	}
	if s.routes["POST"] == nil {
		s.routes["POST"] = make(map[string][]core.HandlerFunc)
	}
	s.routes["POST"][path] = handlers
	s.mux.HandleFunc(path, s.handleHTTP("POST", path))
}

// PUT implements core.Server.PUT for Server
func (s *Server) PUT(path string, handlers ...core.HandlerFunc) {
	if s.routes == nil {
		s.routes = make(map[string]map[string][]core.HandlerFunc)
	}
	if s.routes["PUT"] == nil {
		s.routes["PUT"] = make(map[string][]core.HandlerFunc)
	}
	s.routes["PUT"][path] = handlers
	s.mux.HandleFunc(path, s.handleHTTP("PUT", path))
}

// DELETE implements core.Server.DELETE for Server
func (s *Server) DELETE(path string, handlers ...core.HandlerFunc) {
	if s.routes == nil {
		s.routes = make(map[string]map[string][]core.HandlerFunc)
	}
	if s.routes["DELETE"] == nil {
		s.routes["DELETE"] = make(map[string][]core.HandlerFunc)
	}
	s.routes["DELETE"][path] = handlers
	s.mux.HandleFunc(path, s.handleHTTP("DELETE", path))
}

// PATCH implements core.Server.PATCH for Server
func (s *Server) PATCH(path string, handlers ...core.HandlerFunc) {
	if s.routes == nil {
		s.routes = make(map[string]map[string][]core.HandlerFunc)
	}
	if s.routes["PATCH"] == nil {
		s.routes["PATCH"] = make(map[string][]core.HandlerFunc)
	}
	s.routes["PATCH"][path] = handlers
	s.mux.HandleFunc(path, s.handleHTTP("PATCH", path))
}

// Group implements core.Server.Group for Server
func (s *Server) Group(path string) core.RouterGroup {
	return &RouterGroup{
		server: s,
		prefix: path,
	}
}

// Use implements core.Server.Use for Server
func (s *Server) Use(middleware ...core.HandlerFunc) {
	for _, m := range middleware {
		// Get the function name for logging
		funcValue := reflect.ValueOf(m)
		middlewareName := runtime.FuncForPC(funcValue.Pointer()).Name()
		s.middlewareLog = append(s.middlewareLog, middlewareName)

		// Log middleware addition
		log.Printf("[STD] Adding middleware: %s", middlewareName)
	}

	s.middleware = append(s.middleware, middleware...)
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

		// Log controller registration
		log.Printf("[STD] Registered controller with method: %s, path: %s, skip logging: %t, skip auth check: %t",
			method, path, controller.SkipLogging(), controller.SkipAuthCheck())
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
				_ = c.Error(httperrors.NewNotFoundHttpError(err))
			},
		}
		log.Printf("[STD] Using default NoRoute handler")
	}

	s.noRouteHandlers = handlers
	log.Printf("[STD] Registered NoRoute handler")
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
				_ = c.Error(httperrors.NewMethodNotAllowedHttpError(err))
			},
		}
		log.Printf("[STD] Using default NoMethod handler")
	}

	s.noMethodHandlers = handlers
	log.Printf("[STD] Registered NoMethod handler")
}

// Run implements core.Server.Run for Server
func (s *Server) Run() error {
	addr := ":" + s.port

	// Log server information
	log.Printf("[STD] Server starting on %s", addr)
	log.Printf("[STD] Using standard net/http package")

	// Log middleware information
	if len(s.middlewareLog) > 0 {
		log.Println("[STD] Middleware registered:")
		for i, middleware := range s.middlewareLog {
			log.Printf("[STD]   %d. %s", i+1, middleware)
		}
	} else {
		log.Println("[STD] No middleware registered")
	}

	// Log routes information
	routeCount := 0
	for _, paths := range s.routes {
		routeCount += len(paths)
	}

	if routeCount > 0 {
		log.Println("[STD] Routes registered:")
		routeIndex := 1
		for method, paths := range s.routes {
			for path := range paths {
				log.Printf("[STD]   %d. %s %s", routeIndex, method, path)
				routeIndex++
			}
		}
	} else {
		log.Println("[STD] No routes registered")
	}

	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	log.Printf("[STD] Server is ready to handle requests")
	return s.server.ListenAndServe()
}

// RunTLS implements core.Server.RunTLS for Server
func (s *Server) RunTLS(addr, certFile, keyFile string) error {
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	return s.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown implements core.Server.Shutdown for Server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// StartLambda starts the server in AWS Lambda mode.
// This method should be called instead of Run or RunTLS when running in AWS Lambda.
// This method uses the httpadapter library to convert the standard HTTP handler to a Lambda handler.
//
// Example usage:
//
//	import (
//	    "github.com/mythofleader/go-http-server"
//	)
//
//	func main() {
//	    s, _ := server.NewServer(server.FrameworkStdHTTP, "8080")
//	    // ... configure your server ...
//	    if err := s.StartLambda(); err != nil {
//	        // Handle error
//	    }
//	}
func (s *Server) StartLambda() error {
	// Lambda is only supported with the Gin framework
	return errors.New("Lambda is only supported with the Gin framework")
}

// handleHTTP creates an http.HandlerFunc that handles the request based on the method and path
func (s *Server) handleHTTP(method, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			// Method not allowed
			if len(s.noMethodHandlers) > 0 {
				// Use custom NoMethod handlers
				allHandlers := make([]core.HandlerFunc, 0, len(s.middleware)+len(s.noMethodHandlers))
				allHandlers = append(allHandlers, s.middleware...)
				allHandlers = append(allHandlers, s.noMethodHandlers...)

				ctx := &Context{
					req:          r,
					writer:       w,
					params:       make(map[string]string),
					keys:         make(map[string]interface{}),
					handlers:     allHandlers,
					index:        -1,
					handlerCount: len(allHandlers),
				}

				// Add a MethodNotAllowedHttpError to the context
				ctx.Error(fmt.Errorf("Method %s not allowed for path %s", r.Method, path))

				// Start the middleware chain
				ctx.Next()
			} else {
				// Use default error response
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		handlers, ok := s.routes[method][path]
		if !ok {
			// Route not found
			if len(s.noRouteHandlers) > 0 {
				// Use custom NoRoute handlers
				allHandlers := make([]core.HandlerFunc, 0, len(s.middleware)+len(s.noRouteHandlers))
				allHandlers = append(allHandlers, s.middleware...)
				allHandlers = append(allHandlers, s.noRouteHandlers...)

				ctx := &Context{
					req:          r,
					writer:       w,
					params:       make(map[string]string),
					keys:         make(map[string]interface{}),
					handlers:     allHandlers,
					index:        -1,
					handlerCount: len(allHandlers),
				}

				// Add a NotFoundHttpError to the context
				ctx.Error(fmt.Errorf("Route %s not found", path))

				// Start the middleware chain
				ctx.Next()
			} else {
				// Use default error response
				http.NotFound(w, r)
			}
			return
		}

		// Combine middleware and route handlers into a single slice
		allHandlers := make([]core.HandlerFunc, 0, len(s.middleware)+len(handlers))
		allHandlers = append(allHandlers, s.middleware...)
		allHandlers = append(allHandlers, handlers...)

		ctx := &Context{
			req:          r,
			writer:       w,
			params:       make(map[string]string),
			keys:         make(map[string]interface{}),
			handlers:     allHandlers,
			index:        -1,
			handlerCount: len(allHandlers),
		}

		// Log middleware execution
		for i := range s.middleware {
			if i < len(s.middlewareLog) {
				log.Printf("[STD] Middleware registered: %s for %s %s", s.middlewareLog[i], method, path)
			}
		}

		// Start the middleware chain
		ctx.Next()
	}
}

// RouterGroup is an implementation of core.RouterGroup using the standard net/http package.
type RouterGroup struct {
	server     *Server
	prefix     string
	middleware []core.HandlerFunc
}

// GET implements core.RouterGroup.GET for RouterGroup
func (g *RouterGroup) GET(path string, handlers ...core.HandlerFunc) {
	fullPath := g.prefix + path
	wrappedHandlers := make([]core.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = g.wrapHandler(handler)
	}
	g.server.GET(fullPath, wrappedHandlers...)
}

// POST implements core.RouterGroup.POST for RouterGroup
func (g *RouterGroup) POST(path string, handlers ...core.HandlerFunc) {
	fullPath := g.prefix + path
	wrappedHandlers := make([]core.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = g.wrapHandler(handler)
	}
	g.server.POST(fullPath, wrappedHandlers...)
}

// PUT implements core.RouterGroup.PUT for RouterGroup
func (g *RouterGroup) PUT(path string, handlers ...core.HandlerFunc) {
	fullPath := g.prefix + path
	wrappedHandlers := make([]core.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = g.wrapHandler(handler)
	}
	g.server.PUT(fullPath, wrappedHandlers...)
}

// DELETE implements core.RouterGroup.DELETE for RouterGroup
func (g *RouterGroup) DELETE(path string, handlers ...core.HandlerFunc) {
	fullPath := g.prefix + path
	wrappedHandlers := make([]core.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = g.wrapHandler(handler)
	}
	g.server.DELETE(fullPath, wrappedHandlers...)
}

// PATCH implements core.RouterGroup.PATCH for RouterGroup
func (g *RouterGroup) PATCH(path string, handlers ...core.HandlerFunc) {
	fullPath := g.prefix + path
	wrappedHandlers := make([]core.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		wrappedHandlers[i] = g.wrapHandler(handler)
	}
	g.server.PATCH(fullPath, wrappedHandlers...)
}

// Group implements core.RouterGroup.Group for RouterGroup
func (g *RouterGroup) Group(path string) core.RouterGroup {
	return &RouterGroup{
		server:     g.server,
		prefix:     g.prefix + path,
		middleware: g.middleware,
	}
}

// Use implements core.RouterGroup.Use for RouterGroup
func (g *RouterGroup) Use(middleware ...core.HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
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
		log.Printf("[STD] Registered controller with method: %s, path: %s, skip logging: %t, skip auth check: %t",
			method, path, controller.SkipLogging(), controller.SkipAuthCheck())
	}
}

// wrapHandler wraps a core.HandlerFunc to apply middleware
func (g *RouterGroup) wrapHandler(handler core.HandlerFunc) core.HandlerFunc {
	return func(c core.Context) {
		// Apply middleware
		for _, m := range g.middleware {
			m(c)
		}
		handler(c)
	}
}

// NewServer creates a new Server instance using the standard HTTP package.
func NewServer(port string) *Server {
	log.Printf("[STD] Creating new standard HTTP server on port %s", port)
	return &Server{
		mux:              http.NewServeMux(),
		port:             port,
		middlewareLog:    make([]string, 0),
		noRouteHandlers:  make([]core.HandlerFunc, 0),
		noMethodHandlers: make([]core.HandlerFunc, 0),
	}
}
