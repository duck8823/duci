package github

import (
	"context"
	"github.com/google/go-github/github"
)

type StubClient struct {
}

func (*StubClient) GetPullRequest(ctx context.Context, repo Repository, num int) (*github.PullRequest, error) {
	return nil, nil
}

func (*StubClient) CreateCommitStatus(ctx context.Context, status CommitStatus) error {
	return nil
}

type MockRepository struct {
	FullName string
	URL      string
}

func (r *MockRepository) GetFullName() string {
	return r.FullName
}

func (r *MockRepository) GetSSHURL() string {
	return r.URL
}

func (r *MockRepository) GetCloneURL() string {
	return r.URL
}

func SetInstance(github GitHub) (reset func()) {
	tmp := instance
	instance = github
	return func() {
		instance = tmp
	}
}
