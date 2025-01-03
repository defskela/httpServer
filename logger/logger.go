package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	/*
		level 0 - info, debug, warn, error
		level 1 - debug, warn, error
		level 2 - warn, error
		level 3 - error
		level 4 - nothing
	*/
	level int
}

var defaultLogger *Logger

func init() {
	defaultLogger, _ = newLogger(0)
}

/*
level 1 убирает info-логи,
level 2 убирает info и debug-логи,
level 3 убирает info, debug и warn-логи,
level 4 убирает info, debug, warn и error-логи,
*/
func SetLevel(level int) {
	defaultLogger.level = level
}

func newLogger(level int) (*Logger, error) {
	if level >= 0 && level <= 4 {
		return &Logger{
			infoLogger:  log.New(os.Stdout, "\033[1;34mINFO\033[0m: ", log.Ldate|log.Ltime),
			errorLogger: log.New(os.Stderr, "\033[1;31mERROR\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLogger:  log.New(os.Stdout, "\033[31mWARN\033[0m: ", log.Ldate|log.Ltime),
			debugLogger: log.New(os.Stdout, "\033[1;32mDEBUG\033[0m: ", log.Ldate|log.Ltime),
			level:       level,
		}, nil
	} else {
		return nil, fmt.Errorf("Logger level must be between 0 and 4")
	}
}

func (l *Logger) info(message ...any) {
	if l.level < 1 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.infoLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) debug(message ...any) {
	if l.level < 2 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.debugLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) warn(message ...any) {
	if l.level < 3 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.warnLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) error(err error) {
	if l.level < 4 {
		l.errorLogger.Println(err)
	}
}

// Для логирования ожидаемого поведения системы
func Info(message ...any) {
	defaultLogger.info(message...)
}

// Для разработки и тестирования, чтобы отслеживать шаги выполнения кода
func Debug(message ...any) {
	defaultLogger.debug(message...)
}

// Для указания на нежелательное поведение, которое может потребовать внимания, но не критично для продолжения работы
func Warn(message ...any) {
	defaultLogger.warn(message...)
}

// Для логирования ошибок, которые требуют вмешательства
func Error(err error) {
	defaultLogger.error(err)
}
