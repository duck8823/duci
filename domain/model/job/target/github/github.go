package github

import (
	"context"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"net/url"
	"path"
)

var instance GitHub

// GitHub describes a github client.
type GitHub interface {
	GetPullRequest(ctx context.Context, repo Repository, num int) (*go_github.PullRequest, error)
	CreateCommitStatus(ctx context.Context, src *TargetSource, state State, description string) error
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
func (c *client) CreateCommitStatus(ctx context.Context, src *TargetSource, state State, description string) error {
	// TODO: key must be const
	taskName := ctx.Value("TaskName").(string)
	if len(description) >= 50 {
		description = string([]rune(description)[:46]) + "..."
	}

	status := &go_github.RepoStatus{
		Context:     &taskName,
		Description: &description,
		State:       &state,
		TargetURL:   go_github.String(targetURL(ctx)),
	}

	ownerName, repoName, err := RepositoryName(src.GetFullName()).Split()
	if err != nil {
		return errors.WithStack(err)
	}

	if _, _, err := c.cli.Repositories.CreateStatus(
		ctx,
		ownerName,
		repoName,
		src.GetSHA().String(),
		status,
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func targetURL(ctx context.Context) string {
	// TODO: key must be const
	jobID := ctx.Value("uuid").(uuid.UUID)
	targetURL := ctx.Value("targetURL").(*url.URL)
	targetURL.Path = path.Join(targetURL.Path, "logs", jobID.String())
	return targetURL.String()
}
