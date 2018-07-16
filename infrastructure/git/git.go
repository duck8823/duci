package git

import (
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os"
	"path"
)

type Client interface {
	Clone(ctx context.Context, dir string, sshUrl string, ref string) (plumbing.Hash, error)
}

type sshGitClient struct {
	auth transport.AuthMethod
}

func New() (Client, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", path.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "")
	if err != nil {
		return nil, err
	}
	return &sshGitClient{auth: auth}, nil
}

func (s *sshGitClient) Clone(ctx context.Context, dir string, sshUrl string, ref string) (plumbing.Hash, error) {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           sshUrl,
		Auth:          s.auth,
		Progress:      &ProgressLogger{ctx.UUID()},
		ReferenceName: plumbing.ReferenceName(ref),
	})
	if err != nil {
		return plumbing.Hash{}, errors.WithStack(err)
	}

	reference, err := gitRepository.Head()
	if err != nil {
		return plumbing.Hash{}, errors.WithStack(err)
	}
	return reference.Hash(), nil
}
