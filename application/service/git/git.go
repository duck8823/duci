package git

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// TargetSource stores clone URL, Ref and SHA for target
type TargetSource struct {
	URL string
	Ref string
	SHA plumbing.Hash
}

// Service describes a git service.
type Service interface {
	Clone(ctx context.Context, dir string, src TargetSource) error
}

type sshGitService struct {
	auth transport.AuthMethod
}

type httpGitService struct{}

// New returns the Service.
func New() (Service, error) {
	if application.Config.GitHub.SSHKeyPath == "" {
		return &httpGitService{}, nil
	}
	auth, err := ssh.NewPublicKeysFromFile("git", application.Config.GitHub.SSHKeyPath, "")
	if err != nil {
		return nil, err
	}
	return &sshGitService{auth: auth}, nil
}

// Clone a repository into the path with target source.
func (s *sshGitService) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           src.URL,
		Auth:          s.auth,
		Progress:      &ProgressLogger{ctx.UUID()},
		ReferenceName: plumbing.ReferenceName(src.Ref),
		Depth:         1,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkout(gitRepository, src.SHA); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Clone a repository into the path with target source.
func (s *httpGitService) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           src.URL,
		Progress:      &ProgressLogger{ctx.UUID()},
		ReferenceName: plumbing.ReferenceName(src.Ref),
		Depth:         1,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkout(gitRepository, src.SHA); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func checkout(repo *git.Repository, sha plumbing.Hash) error {
	wt, err := repo.Worktree()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:   sha,
		Branch: plumbing.ReferenceName(sha.String()),
		Create: true,
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
