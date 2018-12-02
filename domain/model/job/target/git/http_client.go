package git

import (
	"context"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type httpGitClient struct {
	runner.LogFunc
}

// NewWithHTTP returns git client with http protocol
func NewWithHTTP(logFunc runner.LogFunc) *httpGitClient {
	return &httpGitClient{LogFunc: logFunc}
}

// Clone a repository into the path with target source.
func (s *httpGitClient) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           src.GetCloneURL(),
		Progress:      &ProgressLogger{ctx: ctx, LogFunc: s.LogFunc},
		ReferenceName: plumbing.ReferenceName(src.GetRef()),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkout(gitRepository, src.GetSHA()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
