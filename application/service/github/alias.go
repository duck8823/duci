package github

import "github.com/google/go-github/github"

// Repository is a interface to get information of git repository.
type Repository interface {
	GetFullName() string
	GetSSHURL() string
}

// PullRequest is a type alias of github.PullRequest
type PullRequest = github.PullRequest

// Status is a type alias of github.RepoStatus
type Status = github.RepoStatus
