package logger

import (
	"fmt"
	"log"
	"os"
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

func NewLogger(level int) (*Logger, error) {
	if level == 0 || level == 1 || level == 2 || level == 3 {
		return &Logger{
			infoLogger:  log.New(os.Stdout, "\033[1;34mINFO\033[0m: ", log.Ldate|log.Ltime),
			errorLogger: log.New(os.Stderr, "\033[1;31mERROR\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile),
			warnLogger:  log.New(os.Stdout, "\033[31mWARN\033[0m: ", log.Ldate|log.Ltime),
			debugLogger: log.New(os.Stdout, "\033[1;32mDEBUG\033[0m: ", log.Ldate|log.Ltime),
			level:       level,
		}, nil
	} else {
		return nil, fmt.Errorf("Logger level must be 0 or 1 or 2 or 3 or 4")
	}

}

func (l *Logger) Info(message string) {
	if l.level < 1 {
		l.infoLogger.Println(message)
	}
}

func (l *Logger) Debug(message string) {
	if l.level < 2 {
		l.debugLogger.Println(message)
	}
}

func (l *Logger) Warn(message string) {
	if l.level < 3 {
		l.warnLogger.Println(message)
	}
}

func (l *Logger) Error(err error) {
	if l.level < 4 {
		l.errorLogger.Println(err)
	}
}
