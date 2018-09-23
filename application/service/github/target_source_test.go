package github_test

import (
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"reflect"
	"testing"
)

func TestTargetSource_ToGitTargetSource(t *testing.T) {
	// given
	expected := git.TargetSource{
		URL: "url",
		Ref: "ref",
		SHA: plumbing.Hash{},
	}

	// and
	src := github.TargetSource{
		Repo: &MockRepository{FullName: "fullName", SSHURL: expected.URL},
		Ref:  expected.Ref,
		SHA:  expected.SHA,
	}

	// expect
	if !reflect.DeepEqual(expected, src.ToGitTargetSource()) {
		t.Errorf("must be equals, but not. diff: %+v", cmp.Diff(expected, src.ToGitTargetSource()))
	}
}

type MockRepository struct {
	FullName string
	SSHURL   string
	CloneURL string
}

func (r *MockRepository) GetFullName() string {
	return r.FullName
}

func (r *MockRepository) GetSSHURL() string {
	return r.SSHURL
}

func (r *MockRepository) GetCloneURL() string {
	return r.CloneURL
}
