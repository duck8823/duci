package git

import (
	"github.com/duck8823/duci/application/context"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Service interface {
	Clone(ctx context.Context, dir string, sshUrl string, ref string, sha plumbing.Hash) error
}

type sshGitService struct {
	auth transport.AuthMethod
}

func New(sshKeyPath string) (Service, error) {
	auth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, "")
	if err != nil {
		return nil, err
	}
	return &sshGitService{auth: auth}, nil
}

func (s *sshGitService) Clone(ctx context.Context, dir string, sshUrl string, ref string, sha plumbing.Hash) error {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           sshUrl,
		Auth:          s.auth,
		Progress:      &ProgressLogger{ctx.UUID()},
		ReferenceName: plumbing.ReferenceName(ref),
		Depth:         1,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	wt, err := gitRepository.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:   sha,
		Branch: plumbing.ReferenceName(ref),
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
