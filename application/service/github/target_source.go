package github

import (
	"github.com/duck8823/duci/application/service/git"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// TargetSource stores Repo, Ref and SHA for target
type TargetSource struct {
	Repo Repository
	Ref  string
	SHA  plumbing.Hash
}

// ToGitTargetSource returns a git.TargetSource.
func (s TargetSource) ToGitTargetSource() git.TargetSource {
	return git.TargetSource{
		URL: s.Repo.GetSSHURL(),
		Ref: s.Ref,
		SHA: s.SHA,
	}
}
