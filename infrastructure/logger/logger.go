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
	now                  = time.Now
)

// Debug logs with the Debug severity.
func Debug(uuid uuid.UUID, message string) {
	if len(message) < 1 || message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[36;1m[DEBUG]\033[0m %s", uuid, now().Format(timeFormat), message)))
}

// Debugf logs with the Debug severity.
func Debugf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Debug(uuid, message)
}

// Info logs with the Info severity.
func Info(uuid uuid.UUID, message string) {
	if len(message) < 1 || message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[1m[INFO]\033[0m %s", uuid, now().Format(timeFormat), message)))
}

// Infof logs with the Info severity.
func Infof(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Info(uuid, message)
}

// Error logs with the Error severity.
func Error(uuid uuid.UUID, message string) {
	if len(message) < 1 || message[len(message)-1] != '\n' {
		message += "\n"
	}
	Writer.Write([]byte(fmt.Sprintf("[%s] %s \033[41;1m[ERROR]\033[0m %s", uuid, now().Format(timeFormat), message)))
}

// Errorf logs with the Error severity.
func Errorf(uuid uuid.UUID, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Error(uuid, message)
}
