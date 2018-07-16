package runner_test

import (
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/mocks/mock_git"
	"github.com/duck8823/duci/mocks/mock_github"
	"github.com/duck8823/duci/service/runner"
	"github.com/golang/mock/gomock"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"testing"
)

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

func TestRunnerImpl_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGitHub := mock_github.NewMockService(ctrl)
	mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(2).
		Return(nil)

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

	dockerClient, err := docker.New()
	if err != nil {
		t.Fatalf("error occurred: %+v", err)
	}

	r := &runner.DockerRunner{
		Name:        "test-runner",
		BaseWorkDir: path.Join(os.TempDir(), "test-runner"),
		Git:         mockGit,
		GitHub:      mockGitHub,
		Docker:      dockerClient,
	}

	repo := &MockRepo{"duck8823/duci", "git@github.com:duck8823/duci.git"}
	hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")
	if err != nil {
		t.Errorf("must not error. but: %+v", err)
	}
	var empty plumbing.Hash
	if hash == empty {
		t.Error("hash must not empty")
	}
}
