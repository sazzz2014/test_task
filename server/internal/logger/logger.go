package logger

import (
	"log"
	"os"
)

type Logger struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLog.Printf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLog.Printf(format, v...)
} 