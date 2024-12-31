package router

import (
	"net"
)

type HandlerFunc func(conn net.Conn, params map[string]string)

type Router struct {
	routes map[string]map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]HandlerFunc),
	}
}

func (r *Router) AddRoute(method, path string, handler HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.AddRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.AddRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.AddRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.AddRoute("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.AddRoute("DELETE", path, handler)
}

func (r *Router) HandleRequest(conn net.Conn, method, path string) {
	if handlers, ok := r.routes[method]; ok {
		if handler, ok := handlers[path]; ok {
			handler(conn, nil)
			return
		}
	}

	r.handleNotFound(conn)
}

func (r *Router) handleNotFound(conn net.Conn) {
	response := "HTTP/1.1 404 Not Found\r\n\r\nRoute not found"
	conn.Write([]byte(response))
}
