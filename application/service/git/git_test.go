package git_test

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/service/git"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	go_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"net/url"
	"os"
	"path"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("when without ssh key path", func(t *testing.T) {
		// expect
		if _, err := git.New(); err != nil {
			t.Error("error must not occur")
		}
	})

	t.Run("when missing ssh key path", func(t *testing.T) {
		// given
		application.Config.GitHub.SSHKeyPath = "/path/to/wrong"

		// expect
		if _, err := git.New(); err == nil {
			t.Error("error must occur")
		}
	})
}

func TestSshGitService_Clone(t *testing.T) {
	// setup
	application.Config.GitHub.SSHKeyPath = path.Join(os.Getenv("HOME"), ".ssh/id_rsa")

	t.Run("when failure git clone", func(t *testing.T) {
		// given
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			return nil, errors.New("test")
		})

		// and
		sut, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		// expect
		if err := sut.Clone(
			context.New("test/task", uuid.New(), &url.URL{}),
			"",
			&git.MockTargetSource{},
		); err == nil {
			t.Error("error must not nil.")
		}

		// cleanup
		git.SetPlainCloneFunc(go_git.PlainClone)
	})

	t.Run("when success git clone", func(t *testing.T) {
		// setup
		dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
		tempDir := path.Join(os.TempDir(), dirStr)
		if err := os.MkdirAll(tempDir, 0700); err != nil {
			t.Fatalf("%+v", err)
		}

		// and
		var hash plumbing.Hash
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tempDir, false)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			w, err := repo.Worktree()
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			// commit
			hash, err = w.Commit("init. commit", &go_git.CommitOptions{
				Author: &object.Signature{},
			})
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			return repo, nil
		})

		// and
		sut, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		// expect
		if err := sut.Clone(
			context.New("test/task", uuid.New(), &url.URL{}),
			"",
			&git.MockTargetSource{
				Ref: "HEAD",
				SHA: hash,
			},
		); err != nil {
			t.Errorf("error must not occur. but got %+v", err)
		}

		// cleanup
		git.SetPlainCloneFunc(go_git.PlainClone)
	})
}

func TestHttpGitService_Clone(t *testing.T) {
	// setup
	application.Config.GitHub.SSHKeyPath = ""

	t.Run("when failure git clone", func(t *testing.T) {
		// given
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			return nil, errors.New("test")
		})

		// and
		sut, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		// expect
		if err := sut.Clone(
			context.New("test/task", uuid.New(), &url.URL{}),
			"",
			&git.MockTargetSource{},
		); err == nil {
			t.Error("error must not nil.")
		}

		// cleanup
		git.SetPlainCloneFunc(go_git.PlainClone)
	})

	t.Run("when success git clone", func(t *testing.T) {
		// setup
		dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
		tempDir := path.Join(os.TempDir(), dirStr)
		if err := os.MkdirAll(tempDir, 0700); err != nil {
			t.Fatalf("%+v", err)
		}

		// and
		var hash plumbing.Hash
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tempDir, false)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			w, err := repo.Worktree()
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			// commit
			hash, err = w.Commit("init. commit", &go_git.CommitOptions{
				Author: &object.Signature{},
			})
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			return repo, nil
		})

		// and
		sut, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		// expect
		if err := sut.Clone(
			context.New("test/task", uuid.New(), &url.URL{}),
			"",
			&git.MockTargetSource{
				Ref: "HEAD",
				SHA: hash,
			},
		); err != nil {
			t.Errorf("error must not occur. but got %+v", err)
		}

		// cleanup
		git.SetPlainCloneFunc(go_git.PlainClone)
	})
}
