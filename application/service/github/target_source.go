package github

import (
	"github.com/duck8823/duci/application"
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
	var url string
	if application.Config.GitHub.SSHKeyPath == "" {
		url = s.Repo.GetSSHURL()
	} else {
		url = s.Repo.GetCloneURL()
	}
	return git.TargetSource{
		URL: url,
		Ref: s.Ref,
		SHA: s.SHA,
	}
}
