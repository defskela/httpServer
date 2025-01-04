package models

import "net"

type HandlerFunc func(conn net.Conn, data RequestData)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type RequestData struct {
	Params   map[string]string
	Body     string
	FormData map[string]string
	Headers  map[string]string
}
