package handlers

import (
	"fmt"
	"httpServer/logger"
	"httpServer/router"
	"httpServer/utils"
	"net"
)

func HandleConnection(conn net.Conn, log *logger.Logger, router *router.Router) {
	log.Info(fmt.Sprintf("Соединение установлено %s %s", conn.LocalAddr().Network(), conn.LocalAddr().String()))

	defer conn.Close()

	request, err := utils.ReadHTTPRequest(conn)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка чтения HTTP-запроса: %s", err))
		return
	}

	router.HandleRequest(conn, request.Method, request.Path)

}
