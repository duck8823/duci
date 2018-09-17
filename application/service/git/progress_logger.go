package git

import (
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"regexp"
)

// Regexp to remove CR or later (inline progress)
var rep = regexp.MustCompile("\r.*$")

// ProgressLogger is a writer for git progress
type ProgressLogger struct {
	uuid uuid.UUID
}

// Write a log without CR or later.
func (l *ProgressLogger) Write(p []byte) (n int, err error) {
	log := rep.ReplaceAllString(string(p), "")
	if len(log) > 0 {
		logger.Info(l.uuid, log)
	}
	return 0, nil
}
