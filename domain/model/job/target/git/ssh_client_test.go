package git_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/git/mock_git"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/duck8823/duci/internal/container"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	go_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeWithSSH(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		t.Run("with correct key path", func(t *testing.T) {
			// given
			container.Clear()

			// and
			path, reset := createTemporaryKey(t)
			defer reset()

			// when
			err := git.InitializeWithSSH(path, func(_ context.Context, _ job.Log) {})

			// then
			if err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}
		})

		t.Run("with wrong key path", func(t *testing.T) {
			// given
			container.Clear()

			// when
			err := git.InitializeWithSSH("/path/to/nothing", func(_ context.Context, _ job.Log) {})

			// then
			if err == nil {
				t.Error("error must not be nil")
			}
		})
	})

	t.Run("when instance is not nil", func(t *testing.T) {
		// given
		container.Override(new(git.Git))
		defer container.Clear()

		// and
		path, reset := createTemporaryKey(t)
		defer reset()

		// when
		err := git.InitializeWithSSH(path, func(_ context.Context, _ job.Log) {})

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestSshGitClient_Clone(t *testing.T) {
	t.Run("when failure git clone", func(t *testing.T) {
		// given
		reset := git.SetPlainCloneFunc(func(_ string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			return nil, errors.New("test")
		})
		defer reset()

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetSrc := mock_git.NewMockTargetSource(ctrl)
		targetSrc.EXPECT().
			GetSSHURL().
			Times(1).
			Return("git@github.com:duck8823/duci.git")
		targetSrc.EXPECT().
			GetRef().
			Times(1).
			Return("HEAD")

		// and
		sut := &git.SshGitClient{LogFunc: runner.NothingToDo}

		// expect
		if err := sut.Clone(
			context.Background(),
			"/path/to/dummy",
			targetSrc,
		); err == nil {
			t.Error("error must not nil.")
		}
	})

	t.Run("when success git clone and checkout", func(t *testing.T) {
		// given
		var hash plumbing.Hash
		defer git.SetPlainCloneFunc(func(tmpDir string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tmpDir, false)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}

			w, err := repo.Worktree()
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}

			// initial commit ( for success checkout )
			hash, err = w.Commit("init. commit", &go_git.CommitOptions{Author: &object.Signature{}})
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}

			return repo, nil
		})()

		tmpDir, reset := createTmpDir(t)
		defer reset()

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetSrc := mock_git.NewMockTargetSource(ctrl)
		targetSrc.EXPECT().
			GetSSHURL().
			Times(1)
		targetSrc.EXPECT().
			GetRef().
			Times(1).
			Return("HEAD")
		targetSrc.EXPECT().
			GetSHA().
			Times(1).
			Return(hash)

		// and
		sut := &git.SshGitClient{LogFunc: runner.NothingToDo}

		// expect
		if err := sut.Clone(context.Background(), tmpDir, targetSrc); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when success git clone but failure checkout", func(t *testing.T) {
		// given
		defer git.SetPlainCloneFunc(func(tmpDir string, _ bool, _ *go_git.CloneOptions) (*go_git.Repository, error) {
			// git init
			repo, err := go_git.PlainInit(tmpDir, false)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}

			return repo, nil
		})()

		tmpDir, reset := createTmpDir(t)
		defer reset()

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		targetSrc := mock_git.NewMockTargetSource(ctrl)
		targetSrc.EXPECT().
			GetSSHURL().
			Times(1).
			Return("git@github.com:duck8823/duci.git")
		targetSrc.EXPECT().
			GetRef().
			Times(1).
			Return("HEAD")
		targetSrc.EXPECT().
			GetSHA().
			Times(1).
			Return(plumbing.ZeroHash)

		// and
		sut := &git.SshGitClient{LogFunc: runner.NothingToDo}

		// expect
		if err := sut.Clone(context.Background(), tmpDir, targetSrc); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func createTemporaryKey(t *testing.T) (path string, reset func()) {
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

	return keyPath, func() {
		_ = os.RemoveAll(tempDir)
	}
}
