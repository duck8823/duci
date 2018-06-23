package github

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type RepositoryName struct {
	FullName string
}

func (r *RepositoryName) Owner() (string, error) {
	ss := strings.Split(r.FullName, "/")
	if len(ss) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid repository name: %s", r.FullName))
	}
	return ss[0], nil
}

func (r *RepositoryName) Repo() (string, error) {
	ss := strings.Split(r.FullName, "/")
	if len(ss) != 2 {
		return "", errors.New(fmt.Sprintf("Invalid repository name: %s", r.FullName))
	}
	return ss[1], nil
}
