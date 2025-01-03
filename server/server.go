package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/defskela/httpServer/logger"
	"github.com/defskela/httpServer/router"
	"github.com/defskela/httpServer/utils"
)

// Need to transfer the router and logger level
func StartServ(router *router.Router, levelLogger int) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	logger.Info("Сервер работает на порту 8080...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-done:
					return
				default:
					if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
						return
					}
					logger.Error(err)
					continue
				}
			}
			go connection(conn, router)
		}
	}()

	<-stop
	logger.Info("Получен сигнал завершения, завершаем работу сервера...")
	gracefulShutdown(listener, done)
}

func gracefulShutdown(listener net.Listener, done chan struct{}) {
	listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()
	close(done)
	logger.Info("Сервер успешно завершил работу")
}

func connection(conn net.Conn, router *router.Router) {
	logger.Info(fmt.Sprintf("Соединение установлено %s %s", conn.LocalAddr().Network(), conn.LocalAddr().String()))

	defer func() {
		logger.Info("Соединение закрыто", conn.LocalAddr().Network(), conn.LocalAddr().String())
		conn.Close()
	}()

	request, err := utils.ReadHTTPRequest(conn)
	if err != nil {
		logger.Warn(fmt.Sprintf("Ошибка чтения HTTP-запроса: %s", err))
		return
	}

	logger.Info("Получен HTTP-запрос:", request.Method, request.Path, request.Version, request.FormData)

	router.HandleRequest(conn, request)
}
