package logger

import (
	"fmt"
	"os"
)

func Error(err error) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("%+v", err))
}
