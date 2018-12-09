package github

import (
	"context"
	go_github "github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var instance GitHub

// GitHub describes a github client.
type GitHub interface {
	GetPullRequest(ctx context.Context, repo Repository, num int) (*go_github.PullRequest, error)
	CreateCommitStatus(ctx context.Context, status CommitStatus) error
}

type client struct {
	cli *go_github.Client
}

// Initialize create a github client.
func Initialize(token string) error {
	if instance != nil {
		return errors.New("instance already initialized.")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	instance = &client{go_github.NewClient(tc)}
	return nil
}

// GetInstance returns a github client
func GetInstance() (GitHub, error) {
	if instance == nil {
		return nil, errors.New("instance still not initialized.")
	}

	return instance, nil
}

// GetPullRequest returns a pull request with specific repository and number.
func (c *client) GetPullRequest(ctx context.Context, repo Repository, num int) (*go_github.PullRequest, error) {
	ownerName, repoName, err := RepositoryName(repo.GetFullName()).Split()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pr, _, err := c.cli.PullRequests.Get(
		ctx,
		ownerName,
		repoName,
		num,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return pr, nil
}

// CreateCommitStatus create commit status to github.
func (c *client) CreateCommitStatus(ctx context.Context, status CommitStatus) error {
	repoStatus := &go_github.RepoStatus{
		Context:     go_github.String(status.Context),
		Description: go_github.String(status.Description.TrimmedString()),
		State:       go_github.String(status.State.String()),
		TargetURL:   go_github.String(status.TargetURL.String()),
	}

	ownerName, repoName, err := RepositoryName(status.TargetSource.GetFullName()).Split()
	if err != nil {
		return errors.WithStack(err)
	}

	if _, _, err := c.cli.Repositories.CreateStatus(
		ctx,
		ownerName,
		repoName,
		status.TargetSource.GetSHA().String(),
		repoStatus,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
