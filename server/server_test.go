package server

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	envContent := "LEVEL_LOGGER=0\n"
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания временного .env файла: %v", err)
	}
	defer os.Remove(".env")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		StartServ()
		close(done)
	}()

	time.Sleep(1 * time.Second)
	stop <- syscall.SIGINT

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("Сервер не завершил работу в течение 15 секунд")
	}
}
