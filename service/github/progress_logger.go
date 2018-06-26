package github

import (
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"regexp"
)

// Remove CR or later (inline progress)
var rep = regexp.MustCompile("\r.*$")

type ProgressLogger struct {
}

func (l *ProgressLogger) Write(p []byte) (n int, err error) {
	log := rep.ReplaceAllString(string(p), "")
	if len(log) > 0 {
		logger.Info(log)
	}
	return 0, nil
}
