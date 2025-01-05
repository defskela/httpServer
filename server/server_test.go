package server

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/defskela/httpServer/router"
)

func TestGracefulShutdown(t *testing.T) {
	server := NewServer(router.NewRouter())
	go func() {
		server.Start("8080")
	}()

	time.Sleep(1 * time.Second)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Не удалось найти текущий процесс: %v", err)
	}

	err = p.Signal(syscall.SIGINT)
	if err != nil {
		t.Fatalf("Не удалось отправить сигнал SIGINT: %v", err)
	}

	time.Sleep(6 * time.Second)

}
