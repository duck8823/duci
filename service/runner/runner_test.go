package runner_test

import (
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/service/github"
	"github.com/duck8823/minimal-ci/service/github/mock_github"
	"github.com/duck8823/minimal-ci/service/runner"
	"github.com/golang/mock/gomock"
	goGithub "github.com/google/go-github/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

	actual, err := r.ConvertPullRequestToRef(context.New(), &MockRepo{}, 5)
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
		Times(1).
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
	hash, err := r.Run(context.New(), repo, "master", "Hello World.")
	if err != nil {
		t.Errorf("must not error. but: %+v", err)
	}
	var empty plumbing.Hash
	if hash == empty {
		t.Error("hash must not empty")
	}
}

func TestRunnerImpl_CreateCommitStatus(t *testing.T) {
	r, err := runner.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	t.Run("with no error", func(t *testing.T) {
		reader, writer, _ := os.Pipe()
		logger.Writer = writer

		ctx := context.New()
		repo := &MockRepo{"full/name", ""}
		hash := plumbing.Hash{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		state := github.SUCCESS
		statusContext := "minimal-ci"
		description := "task success"
		status := goGithub.RepoStatus{
			State:       &state,
			Description: &description,
			Context:     &statusContext,
		}

		ctrl := gomock.NewController(t)
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(ctx, repo, hash, &status).Times(1).Return(nil)

		r.GitHub = mockGitHub

		r.CreateCommitStatus(ctx, repo, hash, github.SUCCESS)

		writer.Close()
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("error occured: %+v", err)
		}

		expected := ""
		actual := string(bytes)
		if actual != expected {
			t.Errorf("log wont: %+v, but got: %+v", expected, actual)
		}
		logger.Writer = os.Stdout
	})

	t.Run("with error", func(t *testing.T) {
		reader, writer, _ := os.Pipe()
		logger.Writer = writer

		ctx := context.New()
		repo := &MockRepo{"full/name", ""}
		hash := plumbing.Hash{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		state := github.SUCCESS
		statusContext := "minimal-ci"
		description := "task success"
		status := goGithub.RepoStatus{
			State:       &state,
			Description: &description,
			Context:     &statusContext,
		}

		ctrl := gomock.NewController(t)
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(ctx, repo, hash, &status).Times(1).Return(errors.New("hello error"))

		r.GitHub = mockGitHub

		r.CreateCommitStatus(ctx, repo, hash, github.SUCCESS)

		writer.Close()
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("error occured: %+v", err)
		}

		expected := "Failed to create commit status: hello error"
		log := string(bytes)
		if !strings.Contains(log, expected) {
			t.Errorf("log must contains %+v, but got: %+v", expected, log)
		}
		logger.Writer = os.Stdout
	})
}

func TestRunnerImpl_CreateCommitStatusWithError(t *testing.T) {
	r, err := runner.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	t.Run("with no error", func(t *testing.T) {
		reader, writer, _ := os.Pipe()
		logger.Writer = writer

		ctx := context.New()
		repo := &MockRepo{"full/name", ""}
		hash := plumbing.Hash{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		state := github.ERROR
		statusContext := "minimal-ci"
		description := "long error description error1 / error2 / error..."
		status := goGithub.RepoStatus{
			State:       &state,
			Description: &description,
			Context:     &statusContext,
		}

		ctrl := gomock.NewController(t)
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(ctx, repo, hash, &status).Times(1).Return(nil)

		r.GitHub = mockGitHub

		r.CreateCommitStatusWithError(ctx, repo, hash, errors.New("long error description error1 / error2 / error3 / error4 / error5 / error6 / error7"))

		writer.Close()
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("error occured: %+v", err)
		}

		expected := ""
		actual := string(bytes)
		if actual != expected {
			t.Errorf("log wont: %+v, but got: %+v", expected, actual)
		}
		logger.Writer = os.Stdout
	})

	t.Run("with error", func(t *testing.T) {
		reader, writer, _ := os.Pipe()
		logger.Writer = writer

		ctx := context.New()
		repo := &MockRepo{"full/name", ""}
		hash := plumbing.Hash{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		state := github.ERROR
		statusContext := "minimal-ci"
		description := "err"
		status := goGithub.RepoStatus{
			State:       &state,
			Description: &description,
			Context:     &statusContext,
		}

		ctrl := gomock.NewController(t)
		mockGitHub := mock_github.NewMockService(ctrl)
		mockGitHub.EXPECT().CreateCommitStatus(ctx, repo, hash, &status).Times(1).Return(errors.New("hello error"))

		r.GitHub = mockGitHub

		r.CreateCommitStatusWithError(ctx, repo, hash, errors.New("err"))

		writer.Close()
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			t.Errorf("error occured: %+v", err)
		}

		expected := "Failed to create commit status: hello error"
		log := string(bytes)
		if !strings.Contains(log, expected) {
			t.Errorf("log must contains %+v, but got: %+v", expected, log)
		}
		logger.Writer = os.Stdout
	})
}
