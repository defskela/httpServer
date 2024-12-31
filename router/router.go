package router

import (
	"httpServer/logger"
	"net"
	"regexp"
)

type HandlerFunc func(conn net.Conn, params map[string]string)

type Router struct {
	routes map[string]map[*regexp.Regexp]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[*regexp.Regexp]HandlerFunc),
	}
}

func (r *Router) addRoute(method, path string, handler HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[*regexp.Regexp]HandlerFunc)
	}
	regexPath := regexp.MustCompile(path)
	r.routes[method][regexPath] = handler
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.addRoute("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) HandleRequest(conn net.Conn, log *logger.Logger, method, path string) {
	if handlers, ok := r.routes[method]; ok {
		for regex, handler := range handlers {
			if matches := regex.FindStringSubmatch(path); matches != nil {
				params := make(map[string]string)
				for i, name := range regex.SubexpNames() {
					if i != 0 && name != "" {
						params[name] = matches[i]
					}
				}
				handler(conn, params)
				return
			} else {
				log.Debug("Неверно указан url")
			}
		}
	} else {
		log.Debug("Неверно указан метод")
	}

	r.handleNotFound(conn)
}

func (r *Router) handleNotFound(conn net.Conn) {
	response := "HTTP/1.1 404 Not Found\r\n\r\nRoute not found"
	conn.Write([]byte(response))
}
