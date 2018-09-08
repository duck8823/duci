package docker

import (
	"bufio"
	"time"
)

func SetNowFunc(f func() time.Time) {
	now = f
}

type BuildLogger = buildLogger

func (l *buildLogger) SetReader(r *bufio.Reader) {
	l.reader = r
}

type RunLogger = runLogger

func (l *runLogger) SetReader(r *bufio.Reader) {
	l.reader = r
}
