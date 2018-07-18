package runner_test

import (
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/mocks/mock_docker"
	"github.com/duck8823/duci/mocks/mock_git"
	"github.com/duck8823/duci/mocks/mock_github"
	"github.com/duck8823/duci/service/runner"
	"github.com/golang/mock/gomock"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"testing"
)

func TestRunnerImpl_Run(t *testing.T) {
	t.Run("with correct return values", func(t *testing.T) {
		// setup
		ctrl := gomock.NewController(t)

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
			Build(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)
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