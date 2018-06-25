package logger

import (
	"github.com/op/go-logging"
	"io"
)

var logger = logging.MustGetLogger("minimal-ci")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} [%{level:-8s}]%{color:reset} %{message}`,
)

func Init(writer io.Writer, level logging.Level) {
	backend := logging.NewLogBackend(writer, "", 0)

	formatter := logging.NewBackendFormatter(backend, format)

	leveled := logging.AddModuleLevel(formatter)
	leveled.SetLevel(level, "")

	logging.SetBackend(leveled)
}

func Debug(message string) {
	logger.Debug(message)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(message string) {
	logger.Info(message)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Error(message string) {
	logger.Error(message)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
