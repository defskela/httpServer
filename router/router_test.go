package router

import (
	"httpServer/logger"
	"net"
	"testing"
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
	router.HandleRequest(conn, nil, "GET", "/")

	expected := "Welcome to start page!"
	if conn.written != expected {
		t.Errorf("expected %q, got %q", expected, conn.written)
	}
}

func TestRouter_NotFound(t *testing.T) {
	router := NewRouter()

	conn := &mockConn{}
	log, _ := logger.NewLogger(0)
	router.HandleRequest(conn, log, "GET", "/notfound")

	expected := "HTTP/1.1 404 Not Found\r\n\r\nRoute not found"
	if conn.written != expected {
		t.Errorf("expected %q, got %q", expected, conn.written)
	}
}
