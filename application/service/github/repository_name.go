package github

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// RepositoryName is a github repository name.
type RepositoryName struct {
	FullName string
}

// Owner get a repository owner.
func (r *RepositoryName) Owner() (string, error) {
	ss := strings.Split(r.FullName, "/")
	if len(ss) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid repository name: %s", r.FullName))
	}
	return ss[0], nil
}

// Repo get a repository name without owner.
func (r *RepositoryName) Repo() (string, error) {
	ss := strings.Split(r.FullName, "/")
	if len(ss) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid repository name: %s", r.FullName))
	}
	return ss[1], nil
}
