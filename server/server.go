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
	"github.com/defskela/httpServer/models"
	"github.com/defskela/httpServer/router"
	"github.com/defskela/httpServer/utils"
)

type server struct {
	listener net.Listener
	router   *router.Router
	done     chan struct{}
	stop     chan os.Signal
}

func NewServer(router *router.Router) *server {
	return &server{
		router: router,
		done:   make(chan struct{}),
		stop:   make(chan os.Signal, 1),
	}
}

// Need to transfer the router and logger level
func (s *server) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	s.listener = listener

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Сервер работает на порту 8080...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	s.router.Get("/ping", pingHandler)

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
			go connection(conn, s.router)
		}
	}()

	<-stop
	logger.Info("Получен сигнал завершения, завершаем работу сервера...")
	s.Shutdown()

	return nil

}

func (s *server) Shutdown() {
	logger.Info("Завершаем работу сервера...")
	s.listener.Close()

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	close(s.done)

	<-ctx.Done()
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

func pingHandler(conn net.Conn, reqData models.RequestData) {
	conn.Write([]byte(utils.CreateHTTPResponse(200, "pong")))
}
