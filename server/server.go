package server

import (
	"fmt"
	"httpServer/logger"
	"httpServer/server/handlers"
	"log"
	"net"
	"os"
	"strconv"

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

	log.Info("Сервер работает на порту 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
		} else {
			go handlers.HandleConnection(conn, log)
		}
	}
}
