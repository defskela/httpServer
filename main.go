package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
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
			go handleConnection(conn)
		}

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection established")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Received:", line)

		if line == "ping" {
			conn.Write([]byte("pong\n"))
		} else {
			conn.Write([]byte("You said: " + line + "\n"))
		}
	}
	fmt.Println("Connection closed")
}
