package runner_test

import (
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/service/github"
	"github.com/duck8823/minimal-ci/service/runner"
	goGithub "github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"reflect"
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

type MockGitHub struct {
}

func (g *MockGitHub) GetPullRequest(ctx context.Context, repository github.Repository, num int) (*github.PullRequest, error) {
	return &github.PullRequest{
		Head: &goGithub.PullRequestBranch{
			Ref: goGithub.String("master"),
		},
	}, nil
}

func (g *MockGitHub) CreateCommitStatus(ctx context.Context, repository github.Repository, hash plumbing.Hash, status *github.Status) error {
	expected := plumbing.Hash{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(expected, hash) {
		return fmt.Errorf("hash must be equal %+v, but got %+v", expected, hash)
	}
	return nil
}

func (g *MockGitHub) Clone(ctx context.Context, dir string, repo github.Repository, ref string) (plumbing.Hash, error) {
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
}

func TestRunnerImpl_ConvertPullRequestToRef(t *testing.T) {
	r, err := runner.NewWithEnv()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	r.GitHub = &MockGitHub{}

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

	r.GitHub = &MockGitHub{}

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
