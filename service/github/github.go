package github

import (
	ctx "context"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"os"
	"path"
)

type State = string

const (
	PENDING State = "pending"
	SUCCESS State = "success"
	ERROR   State = "error"
	FAILURE State = "failure"
)

type Service interface {
	GetPullRequest(ctx context.Context, repository Repository, num int) (*github.PullRequest, error)
	CreateCommitStatus(ctx context.Context, repo Repository, hash plumbing.Hash, state State, description string) error
	Clone(ctx context.Context, dir string, repo Repository, ref string) (plumbing.Hash, error)
}

type serviceImpl struct {
	Client *github.Client
	auth   transport.AuthMethod
}

func New(token string) (*serviceImpl, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx.Background(), ts)

	auth, err := ssh.NewPublicKeysFromFile("git", path.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "")
	if err != nil {
		return nil, err
	}

	return &serviceImpl{Client: github.NewClient(tc), auth: auth}, nil
}

func (s *serviceImpl) GetPullRequest(ctx context.Context, repository Repository, num int) (*PullRequest, error) {
	name := &RepositoryName{repository.GetFullName()}
	owner, err := name.Owner()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	repo, err := name.Repo()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	pr, resp, err := s.Client.PullRequests.Get(
		ctx,
		owner,
		repo,
		num,
	)
	if err != nil {
		logger.Errorf(ctx.UUID(), "Failed to get pull request no. %v on %s: %+v", num, repository.GetFullName(), resp)
		return nil, errors.WithStack(err)
	}
	return pr, nil
}

func (s *serviceImpl) CreateCommitStatus(ctx context.Context, repository Repository, hash plumbing.Hash, state State, description string) error {
	name := &RepositoryName{repository.GetFullName()}
	owner, err := name.Owner()
	if err != nil {
		return errors.WithStack(err)
	}
	repo, err := name.Repo()
	if err != nil {
		return errors.WithStack(err)
	}

	taskName := ctx.TaskName()
	if len(description) >= 50 {
		description = string([]rune(description)[:46]) + "..."
	}
	status := &Status{
		Context:     &taskName,
		Description: &description,
		State:       &state,
	}

	if _, _, err := s.Client.Repositories.CreateStatus(
		ctx,
		owner,
		repo,
		hash.String(),
		status,
	); err != nil {
		logger.Errorf(ctx.UUID(), "Failed to create commit status: %+v", err)
		return errors.WithStack(err)
	}
	return nil
}

func (s *serviceImpl) Clone(ctx context.Context, dir string, repo Repository, ref string) (plumbing.Hash, error) {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           repo.GetSSHURL(),
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
