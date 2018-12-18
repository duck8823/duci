package git

import (
	"context"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/container"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type httpGitClient struct {
	runner.LogFunc
}

// InitializeWithHTTP initialize git client with http protocol
func InitializeWithHTTP(logFunc runner.LogFunc) error {
	git := new(Git)
	*git = &httpGitClient{LogFunc: logFunc}
	if err := container.Submit(git); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Clone a repository into the path with target source.
func (s *httpGitClient) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := plainClone(dir, false, &git.CloneOptions{
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
