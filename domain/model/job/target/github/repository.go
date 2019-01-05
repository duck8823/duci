package github

import (
	"fmt"
	"strings"
)

// Repository is a github repository
type Repository interface {
	GetFullName() string
	GetSSHURL() string
	GetCloneURL() string
}

// RepositoryName is a github repository name.
type RepositoryName string

// Owner get a repository owner.
func (r RepositoryName) Owner() (string, error) {
	ss := strings.Split(string(r), "/")
	if len(ss) != 2 {
		return "", fmt.Errorf("Invalid repository name: %s ", r)
	}
	return ss[0], nil
}

// Repo get a repository name without owner.
func (r RepositoryName) Repo() (string, error) {
	ss := strings.Split(string(r), "/")
	if len(ss) != 2 {
		return "", fmt.Errorf("Invalid repository name: %s ", r)
	}
	return ss[1], nil
}

// Split repository name to owner and repo
func (r RepositoryName) Split() (owner string, repo string, err error) {
	ss := strings.Split(string(r), "/")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("Invalid repository name: %s ", r)
	}
	return ss[0], ss[1], nil
}
