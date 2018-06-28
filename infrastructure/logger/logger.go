package logger

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/op/go-logging"
	"io"
)

var logger = logging.MustGetLogger("minimal-ci")
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02 15:04:05.000} [%{level}]%{color:reset} %{message}`,
)

func Init(writer io.Writer, level logging.Level) {
	backend := logging.NewLogBackend(writer, "", 0)

	formatter := logging.NewBackendFormatter(backend, format)

	leveled := logging.AddModuleLevel(formatter)
	leveled.SetLevel(level, "")

	logging.SetBackend(leveled)
}

func Debug(uuid uuid.UUID, message string) {
	logger.Debug(fmt.Sprintf("[%s] %s", uuid, message))
}

func Debugf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Debug(uuid, message)
}

func Info(uuid uuid.UUID, message string) {
	logger.Info(fmt.Sprintf("[%s] %s", uuid, message))
}

func Infof(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Info(uuid, message)
}

func Error(uuid uuid.UUID, message string) {
	logger.Error(fmt.Sprintf("[%s] %s", uuid, message))
}

func Errorf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Error(uuid, message)
}
