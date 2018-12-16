package target

import "github.com/duck8823/duci/domain/model/job/target/git"

type GithubPush = githubPush

func (g *GithubPush) SetGit(git git.Git) (reset func()) {
	tmp := g.git
	g.git = git
	return func() {
		g.git = tmp
	}
}

type MockRepository struct {
	FullName string
	URL      string
}

func (r *MockRepository) GetFullName() string {
	return r.FullName
}

func (r *MockRepository) GetSSHURL() string {
	return r.URL
}

func (r *MockRepository) GetCloneURL() string {
	return r.URL
}