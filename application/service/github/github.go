package github

import (
	ctx "context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"
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
	GetPullRequest(ctx context.Context, repository Repository, num int) (*PullRequest, error)
	CreateCommitStatus(ctx context.Context, repo Repository, hash plumbing.Hash, state State, description string) error
}

type serviceImpl struct {
	cli *github.Client
}

func New() (Service, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: string(application.Config.GitHub.APIToken)},
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
	pr, resp, err := s.cli.PullRequests.Get(
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
	targetUrl := *ctx.Url()
	targetUrl.Path = path.Join(targetUrl.Path, "logs", ctx.UUID().String())
	targetUrlStr := targetUrl.String()
	status := &Status{
		Context:     &taskName,
		Description: &description,
		State:       &state,
		TargetURL:   &targetUrlStr,
	}

	if _, _, err := s.cli.Repositories.CreateStatus(
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
