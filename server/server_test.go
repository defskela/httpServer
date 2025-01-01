package server

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	envContent := `LEVEL_LOGGER="0"`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания временного .env файла: %v", err)
	}
	defer os.Remove(".env")

	go func() {
		StartServ()
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
