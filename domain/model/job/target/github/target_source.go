package github

import (
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// TargetSource stores Repo, Ref and SHA for target
type TargetSource struct {
	Repository
	Ref string
	SHA plumbing.Hash
}

// GetRef returns a ref
func (s *TargetSource) GetRef() string {
	return s.Ref
}

// GetSHA returns a hash
func (s *TargetSource) GetSHA() plumbing.Hash {
	return s.SHA
}
