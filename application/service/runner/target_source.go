package runner

import (
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type TargetSource struct {
	Repo github.Repository
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
