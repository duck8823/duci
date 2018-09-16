package github

import (
	"github.com/duck8823/duci/application/service/git"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type TargetSource struct {
	Repo Repository
	Ref  string
	SHA  plumbing.Hash
}

func (s TargetSource) ToGitTargetSource() git.TargetSource {
	return git.TargetSource{
		URL: s.Repo.GetSSHURL(),
		Ref: s.Ref,
		SHA: s.SHA,
	}
}
