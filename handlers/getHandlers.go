package handlers

import "net"

func WelcomeHandler(conn net.Conn, params map[string]string) {
	conn.Write([]byte("Welcome to start page!"))
}
