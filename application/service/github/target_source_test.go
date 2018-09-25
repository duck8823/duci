package github_test

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/github"
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

}

func TestTargetSource_GetSHA(t *testing.T) {

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
