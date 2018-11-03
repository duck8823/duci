package git_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
	"path/filepath"
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
	application.Config.GitHub.SSHKeyPath = createTemporaryKey(t)

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
		tempDir := filepath.Join(os.TempDir(), dirStr)
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

	t.Run("when failure git checkout", func(t *testing.T) {
		// setup
		dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
		tempDir := filepath.Join(os.TempDir(), dirStr)
		if err := os.MkdirAll(tempDir, 0700); err != nil {
			t.Fatalf("%+v", err)
		}

		// and
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tempDir, false)
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
			},
		); err == nil {
			t.Error("error must occur. but got nil")
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
		tempDir := filepath.Join(os.TempDir(), dirStr)
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

	t.Run("when failure git checkout", func(t *testing.T) {
		// setup
		dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
		tempDir := filepath.Join(os.TempDir(), dirStr)
		if err := os.MkdirAll(tempDir, 0700); err != nil {
			t.Fatalf("%+v", err)
		}

		// and
		git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tempDir, false)
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
			},
		); err == nil {
			t.Error("error must occur. but got nil")
		}

		// cleanup
		git.SetPlainCloneFunc(go_git.PlainClone)
	})
}

func createTemporaryKey(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 256)
	if err != nil {
		t.Fatalf("error occur: %+v", err)
	}
	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))

	tempDir := filepath.Join(os.TempDir(), random.String(16, random.Alphanumeric))
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		t.Fatalf("error occur: %+v", err)
	}
	keyPath := filepath.Join(tempDir, "id_rsa")
	file, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		t.Fatalf("error occur: %+v", err)
	}

	if _, err := file.WriteString(privateKeyPem); err != nil {
		t.Fatalf("error occur: %+v", err)
	}

	return keyPath
}
