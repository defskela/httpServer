package handlers

import (
	"fmt"
	"httpServer/utils"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	request, err := utils.ReadHTTPRequest(conn)
	if err != nil {
		fmt.Println("Error reading HTTP request:", err)
		return
	}

	var response string
	if request.Path == "/" {
		response = utils.CreateHTTPResponse(200, "Welcome to my HTTP Server!")
	} else {
		response = utils.CreateHTTPResponse(404, "Page not found")
	}

	conn.Write([]byte(response))
}
