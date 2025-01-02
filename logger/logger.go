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
	defaultLogger, _ = NewLogger(0)
}

func SetLevel(level int) {
	defaultLogger.level = level
}

func NewLogger(level int) (*Logger, error) {
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

func (l *Logger) Info(message ...any) {
	if l.level < 1 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.infoLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) Debug(message ...any) {
	if l.level < 2 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.debugLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) Warn(message ...any) {
	if l.level < 3 {
		var strMessages []string
		for _, m := range message {
			strMessages = append(strMessages, fmt.Sprint(m))
		}
		l.warnLogger.Println(strings.Join(strMessages, " "))
	}
}

func (l *Logger) Error(err error) {
	if l.level < 4 {
		l.errorLogger.Println(err)
	}
}

// Глобальные функции
func Info(message ...any) {
	defaultLogger.Info(message...)
}

func Debug(message ...any) {
	defaultLogger.Debug(message...)
}

func Warn(message ...any) {
	defaultLogger.Warn(message...)
}

func Error(err error) {
	defaultLogger.Error(err)
}
