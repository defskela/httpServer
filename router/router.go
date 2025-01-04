package router

import (
	"net"
	"regexp"

	"github.com/defskela/httpServer/http"
	"github.com/defskela/httpServer/logger"
	"github.com/defskela/httpServer/models"
)

type Router struct {
	routes      map[string]map[*regexp.Regexp]models.HandlerFunc
	middlewares []models.MiddlewareFunc
}

func NewRouter() *Router {
	return &Router{
		routes:      make(map[string]map[*regexp.Regexp]models.HandlerFunc),
		middlewares: []models.MiddlewareFunc{},
	}
}

func (r *Router) addRoute(method, path string, handler models.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[*regexp.Regexp]models.HandlerFunc)
	}
	regexPath := regexp.MustCompile(path)
	r.routes[method][regexPath] = handler
}

func (r *Router) Get(path string, handler models.HandlerFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Router) Post(path string, handler models.HandlerFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Router) Put(path string, handler models.HandlerFunc) {
	r.addRoute("PUT", path, handler)
}

func (r *Router) Patch(path string, handler models.HandlerFunc) {
	r.addRoute("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler models.HandlerFunc) {
	r.addRoute("DELETE", path, handler)
}

func (r *Router) HandleRequest(conn net.Conn, request *http.HTTPRequest) {
	if handlers, ok := r.routes[request.Method]; ok {
		for regex, handler := range handlers {
			if matches := regex.FindStringSubmatch(request.Path); matches != nil {
				params := make(map[string]string)
				for i, name := range regex.SubexpNames() {
					if i != 0 && name != "" {
						params[name] = matches[i]
					}
				}

				reqData := models.RequestData{
					Params:   params,
					Body:     request.Body,
					FormData: request.FormData,
					Headers:  request.Headers,
				}

				handler = r.applyMiddlewares(handler)
				handler(conn, reqData)
				return
			} else {
				logger.Debug("Неверно указан url")
			}
		}
	} else {
		logger.Debug("Неверно указан метод")
	}

	r.handleNotFound(conn)
}

func (r *Router) handleNotFound(conn net.Conn) {
	response := "HTTP/1.1 404 Not Found\r\n\r\nRoute not found"
	conn.Write([]byte(response))
}

func (r *Router) Use(middleware models.MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) applyMiddlewares(handler models.HandlerFunc) models.HandlerFunc {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}
