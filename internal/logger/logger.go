package logger

import (
	"fmt"
	"os"
)

// Error print stack to stderr
func Error(err error) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("%+v", err))
}
