package git

import (
	"bufio"
	"context"
	"gopkg.in/src-d/go-git.v4"
	"time"
)

type HttpGitClient = httpGitClient

type SshGitClient = sshGitClient

type CloneLogger = cloneLogger

func (l *CloneLogger) SetReader(r *bufio.Reader) (reset func()) {
	tmp := l.reader
	l.reader = r
	return func() {
		l.reader = tmp
	}
}

func SetPlainCloneFunc(f func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)) (reset func()) {
	tmp := plainClone
	plainClone = f
	return func() {
		plainClone = tmp
	}
}

func SetNowFunc(f func() time.Time) (reset func()) {
	tmp := now
	now = f
	return func() {
		now = tmp
	}
}

func (l *ProgressLogger) SetContext(ctx context.Context) (reset func()) {
	tmp := l.ctx
	l.ctx = ctx
	return func() {
		l.ctx = tmp
	}
}
