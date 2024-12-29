package server

import (
	"fmt"
	"httpServer/server/handlers"
	"net"
	"os"
)

func StartServ() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server is running on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		} else {
			go handlers.HandleConnection(conn)
		}

	}
}
