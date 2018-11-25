package docker

import (
	"bufio"
	"time"
)

func SetNowFunc(f func() time.Time) (reset func()) {
	tmp := now
	now = f
	return func() {
		now = tmp
	}
}

type BuildLogger = buildLogger

func (l *buildLogger) SetReader(r *bufio.Reader) (reset func()) {
	tmp := l.reader
	l.reader = r
	return func() {
		l.reader = tmp
	}
}

type RunLogger = runLogger

func (l *runLogger) SetReader(r *bufio.Reader) (reset func()) {
	tmp := l.reader
	l.reader = r
	return func() {
		l.reader = tmp
	}
}

func (c *clientImpl) SetMoby(m Moby) (reset func()) {
	tmp := c.moby
	c.moby = m
	return func() {
		c.moby = tmp
	}
}
