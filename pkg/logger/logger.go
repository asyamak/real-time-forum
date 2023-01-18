package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log("\033[34m[INFO]\033[0m", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log("\033[31m[ERROR]\033[0m", format, v...)
	os.Exit(1)
}

func (l *Logger) log(prefix, format string, v ...interface{}) {
	fmt.Printf("\033[32m%s\033[0m %s: %s\n", l.prefix, prefix, fmt.Sprintf(format, v...))
}
