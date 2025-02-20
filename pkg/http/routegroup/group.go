package routegroup

import (
	"net/http"
	"regexp"
)

// Group represents a group of routes with associated middleware.
type Group struct {
	mux         *http.ServeMux                    // the underlying mux to register the routes to
	basePath    string                            // base path for the group
	middlewares []func(http.Handler) http.Handler // middlewares stack
}

// New creates a new Group.
func New(mux *http.ServeMux) *Group {
	return Mount(mux, "")
}

// Mount creates a new group with a specified base path.
func Mount(mux *http.ServeMux, basePath string) *Group {
	return &Group{
		mux:         mux,
		basePath:    basePath,
		middlewares: nil,
	}
}

// ServeHTTP implements the http.Handler interface.
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Mount creates a new group with a specified base path on top of the existing group.
func (g *Group) Mount(basePath string) *Group {
	newGroup := g.clone() // copy the middlewares to avoid modifying the original
	newGroup.basePath += basePath

	return newGroup
}

// WithGroup allows for configuring the Group inside the configure function.
func (g *Group) WithGroup(basePath string, configure func(*Group)) {
	newGroup := g.Mount(basePath)
	configure(newGroup)
}

// Use adds middleware(s) to the Group.
func (g *Group) Use(middlewares ...func(http.Handler) http.Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// With adds new middleware(s) to the Group and returns a new Group with the updated middleware stack.
// The With method is similar to Use, but instead of modifying the current Group,
// it returns a new Group instance with the added middleware(s).
// This allows for creating chain of middleware without affecting the original Group.
func (g *Group) With(middlewares ...func(http.Handler) http.Handler) *Group {
	newMiddlewares := make([]func(http.Handler) http.Handler, len(g.middlewares), len(g.middlewares)+len(middlewares))
	copy(newMiddlewares, g.middlewares)
	newMiddlewares = append(newMiddlewares, middlewares...)

	return &Group{
		mux:         g.mux,
		basePath:    g.basePath,
		middlewares: newMiddlewares,
	}
}

// Handle adds a new route to the Group's mux, applying all middlewares to the handler.
func (g *Group) Handle(pattern string, handler http.Handler) {
	g.register(pattern, handler.ServeHTTP)
}

// HandleFunc registers the handler function for the given pattern to the Group's mux.
// The handler is wrapped with the Group's middlewares.
func (g *Group) HandleFunc(pattern string, handler http.HandlerFunc) {
	g.register(pattern, handler)
}

// Handler returns the handler and the pattern that matches the request.
// It always returns a non-nil handler, see http.ServeMux.Handler documentation for details.
func (g *Group) Handler(r *http.Request) (h http.Handler, pattern string) { //nolint:nonamedreturns // std lib pattern
	return g.mux.Handler(r)
}

// Matches non-space characters, spaces, then anything, i.e. "GET /path/to/resource".
var _reGo122 = regexp.MustCompile(`^(\S+)\s+(.+)$`)

func (g *Group) register(pattern string, handler http.HandlerFunc) {
	var path, method string

	matches := _reGo122.FindStringSubmatch(pattern)
	if len(matches) > 2 { //nolint:mnd // path in the form "GET /path/to/resource"
		method = matches[1]
		path = matches[2]
		pattern = method + " " + g.basePath + path
	} else { // path is just "/path/to/resource"
		// method is not set intentionally here, the request pattern had no method part
		path = pattern
		pattern = g.basePath + path
	}

	g.mux.HandleFunc(pattern, g.wrapMiddleware(handler).ServeHTTP)
}

// wrapMiddleware applies the registered middlewares to a handler.
func (g *Group) wrapMiddleware(handler http.Handler) http.Handler {
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		handler = g.middlewares[i](handler)
	}

	return handler
}

func (g *Group) clone() *Group {
	middlewares := make([]func(http.Handler) http.Handler, len(g.middlewares))
	copy(middlewares, g.middlewares)

	return &Group{
		mux:         g.mux,
		basePath:    g.basePath,
		middlewares: middlewares,
	}
}
