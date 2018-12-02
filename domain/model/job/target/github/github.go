package github

import (
	"context"
	"github.com/duck8823/duci/application/service/github"
	go_github "github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// GitHub describes a github client.
type GitHub interface {
	GetPullRequest(ctx context.Context, repo github.Repository, num int) (*go_github.PullRequest, error)
}

type client struct {
	cli *go_github.Client
}

var instance GitHub

// Initialize create a github client.
func Initialize(token string) (GitHub, error) {
	if instance != nil {
		return nil, errors.New("instance already initialized.")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)

	return &client{go_github.NewClient(tc)}, nil
}

// GetInstance returns a github client
func GetInstance() (GitHub, error) {
	if instance == nil {
		return nil, errors.New("instance still not initialized.")
	}

	return instance, nil
}

// GetPullRequest returns a pull request with specific repository and number.
func (c *client) GetPullRequest(ctx context.Context, repo github.Repository, num int) (*go_github.PullRequest, error) {
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
