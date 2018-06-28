package logger

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"time"
)

var (
	timeFormat           = "2006-01-02 15:04:05.000"
	Writer     io.Writer = os.Stdout
)

func Debug(uuid uuid.UUID, message string) {
	if message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[36;1m[DEBUG]\033[0m %s", uuid, time.Now().Format(timeFormat), message)))
}

func Debugf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Debug(uuid, message)
}

func Info(uuid uuid.UUID, message string) {
	if message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[1m[INFO]\033[0m %s", uuid, time.Now().Format(timeFormat), message)))
}

func Infof(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Info(uuid, message)
}

func Error(uuid uuid.UUID, message string) {
	if message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[41;1m[ERROR]\033[0m %s", uuid, time.Now().Format(timeFormat), message)))
}

func Errorf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Error(uuid, message)
}
