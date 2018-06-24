package logger

import (
	"github.com/op/go-logging"
	"os"
)

var logger = logging.MustGetLogger("minimal-ci")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} [%{level:-8s}]%{color:reset} %{message}`,
)

func Init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(formatter)
}

func Debug(message string) {
	logger.Debug(message)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

func Info(message string) {
	logger.Info(message)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

func Error(message string) {
	logger.Error(message)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
}

func Fatal(message string) {
	logger.Fatal(message)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
}
