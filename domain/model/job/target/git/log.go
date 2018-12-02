package git

import (
	"bufio"
	"context"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"strings"
	"time"
)

var now = time.Now

type cloneLogger struct {
	reader *bufio.Reader
}

// ReadLine returns LogLine.
func (l *cloneLogger) ReadLine() (*job.LogLine, error) {
	for {
		line, _, readErr := l.reader.ReadLine()
		msg := string(line)
		if readErr == io.EOF {
			return &job.LogLine{Timestamp: now(), Message: msg}, readErr
		}
		if readErr != nil {
			return nil, errors.WithStack(readErr)
		}

		if len(line) == 0 {
			continue
		}

		return &job.LogLine{Timestamp: now(), Message: msg}, readErr
	}
}

// Regexp to remove CR or later (inline progress)
var rep = regexp.MustCompile("\r.*$")

// ProgressLogger is a writer for git progress
type ProgressLogger struct {
	ctx context.Context
	runner.LogFunc
}

// Write a log without CR or later.
func (l *ProgressLogger) Write(p []byte) (n int, err error) {
	msg := rep.ReplaceAllString(string(p), "")
	log := &cloneLogger{
		reader: bufio.NewReader(strings.NewReader(msg)),
	}
	l.LogFunc(l.ctx, log)
	return 0, nil
}
