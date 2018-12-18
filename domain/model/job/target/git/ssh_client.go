package git

import (
	"context"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/container"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type sshGitClient struct {
	auth transport.AuthMethod
	runner.LogFunc
}

// InitializeWithSSH returns git client with ssh protocol
func InitializeWithSSH(path string, logFunc runner.LogFunc) error {
	auth, err := ssh.NewPublicKeysFromFile("git", path, "")
	if err != nil {
		return errors.WithStack(err)
	}

	git := new(Git)
	*git = &sshGitClient{auth: auth, LogFunc: logFunc}
	if err := container.Submit(git); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Clone a repository into the path with target source.
func (s *sshGitClient) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := plainClone(dir, false, &git.CloneOptions{
		URL:           src.GetSSHURL(),
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
