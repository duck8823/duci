package github_test

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"testing"
)

func TestTargetSource_GetURL(t *testing.T) {
	t.Run("without ssh key path", func(t *testing.T) {
		// given
		application.Config.GitHub.SSHKeyPath = ""

		// and
		expected := "clone_url"

		sut := github.TargetSource{
			Repo: &MockRepository{SSHURL: "ssh_url", CloneURL: expected},
		}

		// expect
		if sut.GetURL() != expected {
			t.Errorf("url must equal. wont %#v, but got %#v", expected, sut.GetURL())
		}
	})

	t.Run("without ssh key path", func(t *testing.T) {
		// given
		application.Config.GitHub.SSHKeyPath = "path/to/ssh_key"

		// and
		expected := "ssh_url"

		sut := github.TargetSource{
			Repo: &MockRepository{SSHURL: expected, CloneURL: "clone_url"},
		}

		// expect
		if sut.GetURL() != expected {
			t.Errorf("url must equal. wont %#v, but got %#v", expected, sut.GetURL())
		}
	})
}

func TestTargetSource_GetRef(t *testing.T) {
	// given
	expected := "ref"

	// and
	sut := github.TargetSource{Ref: expected}

	// when
	actual := sut.GetRef()

	// expect
	if actual != expected {
		t.Errorf("must equal. wont %#v, but got %#v", expected, actual)
	}
}

func TestTargetSource_GetSHA(t *testing.T) {
	// given
	expected := plumbing.NewHash("hello world.")

	// and
	sut := github.TargetSource{SHA: expected}

	// when
	actual := sut.GetSHA()

	// expect
	if actual != expected {
		t.Errorf("must equal. wont %#v, but got %#v", expected, actual)
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
