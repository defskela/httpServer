package server

import (
	"context"
	"fmt"
	"httpServer/handlers"
	"httpServer/logger"
	"httpServer/router"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func StartServ() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		os.Exit(1)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
	levelLogger := os.Getenv("LEVEL_LOGGER")
	level, err := strconv.Atoi(levelLogger)
	if err != nil {
		fmt.Println("Ошибка преобразования уровня логгера:", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(level)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	router := router.NewRouter()
	handlers.InitHandlers(router)

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
			go handlers.HandleConnection(conn, log, router)
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
