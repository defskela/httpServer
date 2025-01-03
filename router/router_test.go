package router

import (
	"net"
	"testing"

	"github.com/defskela/httpServer/http"
)

type mockConn struct {
	net.Conn
	written string
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	m.written = string(b)
	return len(b), nil
}

func TestRouter_Get(t *testing.T) {
	router := NewRouter()
	router.Get("/", func(conn net.Conn, params map[string]string) {
		conn.Write([]byte("Welcome to start page!"))
	})

	conn := &mockConn{}
	request := &http.HTTPRequest{Method: "GET", Path: "/"}
	router.HandleRequest(conn, request)

	expected := "Welcome to start page!"
	if conn.written != expected {
		t.Errorf("expected %q, got %q", expected, conn.written)
	}
}

func TestRouter_NotFound(t *testing.T) {
	router := NewRouter()

	conn := &mockConn{}
	request := &http.HTTPRequest{Method: "GET", Path: "/"}
	router.HandleRequest(conn, request)

	expected := "HTTP/1.1 404 Not Found\r\n\r\nRoute not found"
	if conn.written != expected {
		t.Errorf("expected %q, got %q", expected, conn.written)
	}
}
