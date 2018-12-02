package git

import (
	"context"
	"github.com/duck8823/duci/domain/model/runner"
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

// NewWithSSH returns git client with ssh protocol
func NewWithSSH(path string, logFunc runner.LogFunc) (*sshGitClient, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", path, "")
	if err != nil {
		return nil, err
	}
	return &sshGitClient{auth: auth, LogFunc: logFunc}, nil
}

// Clone a repository into the path with target source.
func (s *sshGitClient) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           src.GetSSHURL(),
		Auth:          s.auth,
		Progress:      &ProgressLogger{ctx: ctx, LogFunc: s.LogFunc},
		ReferenceName: plumbing.ReferenceName(src.GetRef()),
		Depth:         1,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkout(gitRepository, src.GetSHA()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
