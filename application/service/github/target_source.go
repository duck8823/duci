package github

import (
	"github.com/duck8823/duci/application"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// TargetSource stores Repo, Ref and SHA for target
type TargetSource struct {
	Repo Repository
	Ref  string
	SHA  plumbing.Hash
}

// GetURL returns a clone URL
func (s *TargetSource) GetURL() string {
	if application.Config.GitHub.SSHKeyPath != "" {
		return s.Repo.GetSSHURL()
	}
	return s.Repo.GetCloneURL()
}

// GetRef returns a ref
func (s *TargetSource) GetRef() string {
	return s.Ref
}

// GetSHA returns a hash
func (s *TargetSource) GetSHA() plumbing.Hash {
	return s.SHA
}
