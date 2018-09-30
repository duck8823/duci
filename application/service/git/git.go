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

var plainClone = git.PlainClone

// TargetSource is a interface returns clone URL, Ref and SHA for target
type TargetSource interface {
	GetURL() string
	GetRef() string
	GetSHA() plumbing.Hash
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
	gitRepository, err := plainClone(dir, false, &git.CloneOptions{
		URL:           src.GetURL(),
		Auth:          s.auth,
		Progress:      &ProgressLogger{ctx.UUID()},
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

// Clone a repository into the path with target source.
func (s *httpGitService) Clone(ctx context.Context, dir string, src TargetSource) error {
	gitRepository, err := plainClone(dir, false, &git.CloneOptions{
		URL:           src.GetURL(),
		Progress:      &ProgressLogger{ctx.UUID()},
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
