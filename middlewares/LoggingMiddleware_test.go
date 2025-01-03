package middleware

import (
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

func TestLoggingMiddleware(t *testing.T) {
	middleware := LoggingMiddleware()

	handler := func(conn net.Conn, params map[string]string) {
		conn.Write([]byte("Handler executed"))
	}

	conn := &mockConn{}
	middleware(handler)(conn, nil)

	expected := "Handler executed"
	if conn.written != expected {
		t.Errorf("expected %q, got %q", expected, conn.written)
	}
}
