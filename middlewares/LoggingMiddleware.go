package middleware

import (
	"net"

	"github.com/defskela/httpServer/logger"
)

func LoggingMiddleware(log *logger.Logger) func(handlerFunc func(conn net.Conn, params map[string]string)) func(conn net.Conn, params map[string]string) {
	return func(next func(conn net.Conn, params map[string]string)) func(conn net.Conn, params map[string]string) {
		return func(conn net.Conn, params map[string]string) {
			log.Info("Получен запрос")
			next(conn, params)
		}
	}
}
