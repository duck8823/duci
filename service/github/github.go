package github

import (
	"context"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Repository interface {
	GetFullName() string
	GetSSHURL() string
}

type PullRequest = github.PullRequest

type Status = github.RepoStatus

type State = string

const (
	PENDING = "pending"
	SUCCESS = "success"
	ERROR   = "error"
	FAILURE = "failure"
)

type Service interface {
	GetPullRequest(ctx context.Context, repository Repository, num int) (*github.PullRequest, error)
	CreateCommitStatus(ctx context.Context, repository Repository, hash plumbing.Hash, status *Status) error
	Clone(ctx context.Context, dir string, repo Repository, ref string) (plumbing.Hash, error)
}

type serviceImpl struct {
	Client *github.Client
}

func New(ctx context.Context, token string) *serviceImpl {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &serviceImpl{github.NewClient(tc)}
}

func (s *serviceImpl) GetPullRequest(ctx context.Context, repository Repository, num int) (*github.PullRequest, error) {
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
		logger.Errorf("Failed to get pull request no. %v on %s: %+v", num, repository.GetFullName(), resp)
		return nil, errors.WithStack(err)
	}
	return pr, nil
}

func (s *serviceImpl) CreateCommitStatus(ctx context.Context, repository Repository, hash plumbing.Hash, status *Status) error {
	name := &RepositoryName{repository.GetFullName()}
	owner, err := name.Owner()
	if err != nil {
		return errors.WithStack(err)
	}
	repo, err := name.Repo()
	if err != nil {
		return errors.WithStack(err)
	}

	if _, _, err := s.Client.Repositories.CreateStatus(
		context.Background(),
		owner,
		repo,
		hash.String(),
		status,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *serviceImpl) Clone(ctx context.Context, dir string, repo Repository, ref string) (plumbing.Hash, error) {
	gitRepository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           repo.GetSSHURL(),
		Progress:      &ProgressLogger{},
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
