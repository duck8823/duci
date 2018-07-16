package github

import "github.com/google/go-github/github"

type Repository interface {
	GetFullName() string
	GetSSHURL() string
}

type PullRequest = github.PullRequest

type Status = github.RepoStatus
