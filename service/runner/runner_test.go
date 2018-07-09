package runner_test

import (
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/service/github"
	"github.com/duck8823/minimal-ci/service/github/mock_github"
	"github.com/duck8823/minimal-ci/service/runner"
	"github.com/golang/mock/gomock"
	goGithub "github.com/google/go-github/github"
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

func TestRunnerImpl_ConvertPullRequestToRef(t *testing.T) {
	r, err := runner.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	ctrl := gomock.NewController(t)
	mockGitHub := mock_github.NewMockService(ctrl)
	mockGitHub.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(&github.PullRequest{
			Head: &goGithub.PullRequestBranch{Ref: goGithub.String("master")},
		}, nil)
	r.GitHub = mockGitHub

	actual, err := r.ConvertPullRequestToRef(context.New("test/task"), &MockRepo{}, 5)
	if err != nil {
		t.Errorf("must not error. but: %+v", err)
	}

	expected := "refs/heads/master"
	if actual != expected {
		t.Errorf("wont: %+v. but got: %+v", expected, actual)
	}
}

func TestRunnerImpl_Run(t *testing.T) {
	r, err := runner.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	ctrl := gomock.NewController(t)
	mockGitHub := mock_github.NewMockService(ctrl)
	mockGitHub.EXPECT().CreateCommitStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(2).
		Return(nil)
	mockGitHub.EXPECT().Clone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		DoAndReturn(func(ctx context.Context, dir string, repo github.Repository, ref string) (plumbing.Hash, error) {
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
	r.GitHub = mockGitHub

	repo := &MockRepo{"duck8823/minimal-ci", "git@github.com:duck8823/minimal-ci.git"}
	hash, err := r.Run(context.New("test/task"), repo, "master", "Hello World.")
	if err != nil {
		t.Errorf("must not error. but: %+v", err)
	}
	var empty plumbing.Hash
	if hash == empty {
		t.Error("hash must not empty")
	}
}
