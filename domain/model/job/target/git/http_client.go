package git

import (
	"context"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/container"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type httpGitClient struct {
	auth transport.AuthMethod
	runner.LogFunc
}

// InitializeWithHTTP initialize git client with http protocol
func InitializeWithHTTP(token string, logFunc runner.LogFunc) error {
	git := new(Git)
	*git = &httpGitClient{
		auth: &http.BasicAuth{
			Username: "",
			Password: token,
		},
		LogFunc: logFunc,
	}
	if err := container.Submit(git); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Clone a repository into the path with target source.
func (s *httpGitClient) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := plainClone(dir, false, &git.CloneOptions{
		URL:           src.GetCloneURL(),
		Auth:          s.auth,
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
