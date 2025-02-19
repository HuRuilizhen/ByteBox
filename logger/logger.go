package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	prefix   string
	logLevel LogLevel
}

var (
	instance *Logger
	once     sync.Once
)

func newLogger(prefix string, logLevel LogLevel) *Logger {
	return &Logger{
		prefix:   prefix,
		logLevel: logLevel,
	}
}

func GetLoggerInstance() *Logger {
	once.Do(func() { instance = newLogger("", INFO) })
	return instance
}

func convertString(logLevel LogLevel) string {
	switch logLevel {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEFAUT"
	}
}

func (logger *Logger) log(logLevel LogLevel, msg string) {
	if logLevel >= logger.logLevel {
		fmt.Printf("[%-6s] [%s]: %s\n", convertString(logLevel), time.Now().Local().Format("2006-01-02 15:04:05"), msg)
	}
}

func (logger *Logger) Debug(msg string) {
	logger.log(DEBUG, msg)
}

func (logger *Logger) Info(msg string) {
	logger.log(INFO, msg)
}

func (logger *Logger) Warn(msg string) {
	logger.log(WARN, msg)
}

func (logger *Logger) Error(msg string) {
	logger.log(ERROR, msg)
}

func (logger *Logger) Fatal(msg string) {
	logger.log(FATAL, msg)
	os.Exit(1)
}

func (logger *Logger) SetLogLevel(logLevel LogLevel) {
	logger.logLevel = logLevel
}

func (logger *Logger) GetLogLevel() LogLevel {
	return logger.logLevel
}

func (logger *Logger) GetLogLevelString() string {
	return convertString(logger.logLevel)
}
