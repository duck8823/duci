package git

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var instance Git

var plainClone = git.PlainClone

// TargetSource is a interface returns clone URLs, Ref and SHA for target
type TargetSource interface {
	GetSSHURL() string
	GetCloneURL() string
	GetRef() string
	GetSHA() plumbing.Hash
}

// Git describes a git service.
type Git interface {
	Clone(ctx context.Context, dir string, src TargetSource) error
}

// GetInstance returns a git client
func GetInstance() (Git, error) {
	if instance == nil {
		return nil, errors.New("instance still not initialized.")
	}

	return instance, nil
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
