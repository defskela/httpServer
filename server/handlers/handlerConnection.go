package handlers

import (
	"fmt"
	"httpServer/logger"
	"httpServer/utils"
	"net"
)

func HandleConnection(conn net.Conn, log *logger.Logger) {
	log.Info(fmt.Sprintf("Соединение установлено %s %s", conn.LocalAddr().Network(), conn.LocalAddr().String()))

	defer conn.Close()

	request, err := utils.ReadHTTPRequest(conn)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка чтения HTTP-запроса: %s", err))
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
