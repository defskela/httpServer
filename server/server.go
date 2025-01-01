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

func StartServ(router *router.Router, levelLogger int) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(levelLogger)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("Сервер работает на порту 8080...")

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
					log.Error(err)
					continue
				}
			}
			go connection(conn, log, router)
		}
	}()

	<-stop
	log.Info("Получен сигнал завершения, завершаем работу сервера...")
	gracefulShutdown(listener, log, done)
}

func gracefulShutdown(listener net.Listener, log *logger.Logger, done chan struct{}) {
	listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()
	close(done)
	log.Info("Сервер успешно завершил работу")
}

func connection(conn net.Conn, log *logger.Logger, router *router.Router) {
	log.Info(fmt.Sprintf("Соединение установлено %s %s", conn.LocalAddr().Network(), conn.LocalAddr().String()))

	defer func() {
		log.Info(fmt.Sprintf("Соединение закрыто %s %s", conn.LocalAddr().Network(), conn.LocalAddr().String()))
		conn.Close()
	}()

	request, err := utils.ReadHTTPRequest(conn)
	if err != nil {
		log.Warn(fmt.Sprintf("Ошибка чтения HTTP-запроса: %s", err))
		return
	}

	router.HandleRequest(conn, log, request.Method, request.Path)
}
