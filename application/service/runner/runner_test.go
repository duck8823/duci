package runner_test

import (
	"github.com/duck8823/duci/application/service/github/mock_github"
	"github.com/duck8823/duci/application/service/runner"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/docker/mock_docker"
	"github.com/duck8823/duci/infrastructure/git/mock_git"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"testing"
)

func TestRunnerImpl_Run(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)

	t.Run("with correct return values", func(t *testing.T) {
		t.Run("when Dockerfile in proj root", func(t *testing.T) {
			// given
			mockGitHub := mock_github.NewMockService(ctrl)
			mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(2).
				Return(nil)

			// and
			mockGit := mock_git.NewMockClient(ctrl)
			mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				DoAndReturn(func(_ context.Context, dir string, _ string, _ string) (plumbing.Hash, error) {
					if err := os.MkdirAll(dir, 0700); err != nil {
						return plumbing.Hash{}, err
					}

					dockerfile, err := os.OpenFile(path.Join(dir, "Dockerfile"), os.O_RDWR|os.O_CREATE, 0600)
					if err != nil {
						return plumbing.Hash{}, err
					}
					defer dockerfile.Close()

					dockerfile.WriteString("FROM alpine\nENTRYPOINT [\"echo\"]")

					return plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil
				})

			// and
			mockDocker := mock_docker.NewMockClient(ctrl)
			mockDocker.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq("./Dockerfile")).
				Times(1).
				Return(nil)
			mockDocker.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not("./Dockerfile")).
				Return(errors.New("must not call this"))
			mockDocker.EXPECT().
				Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return("", nil)

			r := &runner.DockerRunner{
				Name:        "test-runner",
				BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
				Git:         mockGit,
				GitHub:      mockGitHub,
				Docker:      mockDocker,
			}

			// and
			repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

			// and
			var empty plumbing.Hash

			// when
			hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

			// then
			if err != nil {
				t.Errorf("must not error. but: %+v", err)
			}

			if hash == empty {
				t.Error("hash must not empty")
			}
		})

		t.Run("when Dockerfile in sub directory", func(t *testing.T) {
			// given
			mockGitHub := mock_github.NewMockService(ctrl)
			mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(2).
				Return(nil)

			// and
			mockGit := mock_git.NewMockClient(ctrl)
			mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				DoAndReturn(func(_ context.Context, dir string, _ string, _ string) (plumbing.Hash, error) {
					if err := os.MkdirAll(path.Join(dir, ".duci"), 0700); err != nil {
						return plumbing.Hash{}, err
					}

					dockerfile, err := os.OpenFile(path.Join(dir, ".duci/Dockerfile"), os.O_RDWR|os.O_CREATE, 0600)
					if err != nil {
						return plumbing.Hash{}, err
					}
					defer dockerfile.Close()

					dockerfile.WriteString("FROM alpine\nENTRYPOINT [\"echo\"]")

					return plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil
				})

			// and
			mockDocker := mock_docker.NewMockClient(ctrl)
			mockDocker.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(".duci/Dockerfile")).
				Return(nil)
			mockDocker.EXPECT().
				Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(".duci/Dockerfile")).
				Return(errors.New("must not call this"))
			mockDocker.EXPECT().
				Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Times(1).
				Return("", nil)

			r := &runner.DockerRunner{
				Name:        "test-runner",
				BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
				Git:         mockGit,
				GitHub:      mockGitHub,
				Docker:      mockDocker,
			}

			// and
			repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

			// and
			var empty plumbing.Hash

			// when
			hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

			// then
			if err != nil {
				t.Errorf("must not error. but: %+v", err)
			}

			if hash == empty {
				t.Error("hash must not empty")
			}
		})
	})

	t.Run("with config file", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			DoAndReturn(func(_ context.Context, dir string, _ string, _ string) (plumbing.Hash, error) {
				if err := os.MkdirAll(path.Join(dir, ".duci"), 0700); err != nil {
					return plumbing.Hash{}, err
				}

				dockerfile, err := os.OpenFile(path.Join(dir, ".duci/config.yml"), os.O_RDWR|os.O_CREATE, 0600)
				if err != nil {
					return plumbing.Hash{}, err
				}
				defer dockerfile.Close()

				dockerfile.WriteString("---\nvolumes:\n  - /hello:/hello")

				return plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil
			})

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Eq(docker.RuntimeOptions{Volumes: []string{"/hello:/hello"}}), gomock.Any(), gomock.Any()).
			Times(1).
			Return("", nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Not(docker.RuntimeOptions{Volumes: []string{"/hello:/hello"}}), gomock.Any(), gomock.Any()).
			Return("", errors.New("must not call this"))

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// and
		var empty plumbing.Hash

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err != nil {
			t.Errorf("must not error. but: %+v", err)
		}

		if hash == empty {
			t.Error("hash must not empty")
		}
	})

	t.Run("when failed to git clone", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, errors.New("error"))

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		var empty plumbing.Hash

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err == nil {
			t.Error("must occur error")
		}

		if hash != empty {
			t.Errorf("commit hash must be equal empty, but got %+v", hash)
		}
	})

	t.Run("when workdir not exists", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: "/path/to/not/exists/dir",
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		var empty plumbing.Hash

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err == nil {
			t.Error("must occur error")
		}

		if hash == empty {
			t.Error("hash must not empty")
		}
	})

	t.Run("when docker build failure", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("test"))
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		var empty plumbing.Hash

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err == nil {
			t.Error("must occur error")
		}

		if hash == empty {
			t.Error("hash must not empty")
		}
	})

	t.Run("when docker run error", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return("", errors.New("test"))

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		var empty plumbing.Hash

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err == nil {
			t.Error("must occur error")
		}

		if hash == empty {
			t.Error("hash must not empty")
		}
	})

	t.Run("when docker run failure", func(t *testing.T) {
		// given
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(2).
			Return(nil)

		// and
		mockGit := mock_git.NewMockClient(ctrl)
		mockGit.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)

		// and
		mockDocker := mock_docker.NewMockClient(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return("", docker.Failure)

		r := &runner.DockerRunner{
			Name:        "test-runner",
			BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
			Git:         mockGit,
			GitHub:      mockGitHub,
			Docker:      mockDocker,
		}

		// and
		var empty plumbing.Hash

		// and
		repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}

		// when
		hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")

		// then
		if err != docker.Failure {
			t.Errorf("error must be docker.Failure, but got %+v", err)
		}

		if hash == empty {
			t.Error("hash must not empty")
		}
	})
}

type MockRepo struct {
	FullName string
	SSHURL   string
}

func (r *MockRepo) GetFullName() string {
	return r.FullName
}

func (r *MockRepo) GetSSHURL() string {
	return r.SSHURL
}
