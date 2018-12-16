package git

import (
	"bufio"
	"context"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

type MockTargetSource struct {
	URL string
	Ref string
	SHA plumbing.Hash
}

func (t *MockTargetSource) GetCloneURL() string {
	return t.URL
}

func (t *MockTargetSource) GetSSHURL() string {
	return t.URL
}

func (t *MockTargetSource) GetRef() string {
	return t.Ref
}

func (t *MockTargetSource) GetSHA() plumbing.Hash {
	return t.SHA
}

func (l *ProgressLogger) SetContext(ctx context.Context) (reset func()) {
	tmp := l.ctx
	l.ctx = ctx
	return func() {
		l.ctx = tmp
	}
}
