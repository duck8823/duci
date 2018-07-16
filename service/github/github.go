package github

import (
	ctx "context"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
)

type State = string

const (
	PENDING State = "pending"
	SUCCESS State = "success"
	ERROR   State = "error"
	FAILURE State = "failure"
)

type Service interface {
	GetPullRequest(ctx context.Context, repository Repository, num int) (*PullRequest, error)
	CreateCommitStatus(ctx context.Context, repo Repository, hash plumbing.Hash, state State, description string) error
}

type serviceImpl struct {
	Client *github.Client
}

// TODO change return type to interface using mock server ( gock? )
func NewWithEnv() (*serviceImpl, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx.Background(), ts)

	return &serviceImpl{github.NewClient(tc)}, nil
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
