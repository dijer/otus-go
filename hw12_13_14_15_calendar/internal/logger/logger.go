package logger

import (
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

type Logger struct {
	level LogLevel
}

func New(level LogLevel) *Logger {
	return &Logger{level: level}
}

func (l Logger) Info(msg ...string) {
	l.log(INFO, msg...)
}

func (l Logger) Error(msg ...string) {
	l.log(ERROR, msg...)
}

func (l Logger) Warn(msg ...string) {
	l.log(WARN, msg...)
}

func (l Logger) Debug(msg ...string) {
	l.log(DEBUG, msg...)
}

func (l Logger) log(prefix LogLevel, msg ...string) {
	log.SetOutput(os.Stdout)
	msgs := strings.Join(msg, " ")
	log.Printf("%s %s: %s", time.Now().Format("2000-01-01 10:31:05"), prefix, msgs)
}
