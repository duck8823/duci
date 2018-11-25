package git

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func SetPlainCloneFunc(f func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)) (reset func()) {
	tmp := plainClone
	plainClone = f
	return func() {
		plainClone = tmp
	}
}

type MockTargetSource struct {
	URL string
	Ref string
	SHA plumbing.Hash
}

func (t *MockTargetSource) GetURL() string {
	return t.URL
}

func (t *MockTargetSource) GetRef() string {
	return t.Ref
}

func (t *MockTargetSource) GetSHA() plumbing.Hash {
	return t.SHA
}
