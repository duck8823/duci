package github_test

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/google/go-cmp/cmp"
	go_github "github.com/google/go-github/github"
	"github.com/labstack/gommon/random"
	"gopkg.in/h2non/gock.v1"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/url"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		container.Clear()

		// when
		err := github.Initialize("github_api_token")

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when instance is not nil", func(t *testing.T) {
		// given
		container.Override(&github.StubClient{})
		defer container.Clear()

		// when
		err := github.Initialize("github_api_token")

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestGetInstance(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := github.GetInstance()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})

	t.Run("when instance is not nil", func(t *testing.T) {
		// given
		want := &github.StubClient{}

		// and
		container.Override(want)
		defer container.Clear()

		// when
		got, err := github.GetInstance()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}

	})
}

func TestClient_GetPullRequest(t *testing.T) {
	// given
	_ = github.Initialize("github_api_token")
	sut, err := github.GetInstance()
	if err != nil {
		t.Fatalf("error occurred. %+v", err)
	}

	t.Run("when github server returns status ok", func(t *testing.T) {
		// given
		repo := &github.MockRepository{
			FullName: "duck8823/duci",
		}
		num := 19

		// and
		want := "hello world"

		// and
		gock.New("https://api.github.com").
			Get(fmt.Sprintf("/repos/%s/pulls/%d", repo.FullName, num)).
			Reply(200).
			JSON(&go_github.PullRequest{
				Title: go_github.String(want),
			})
		defer gock.Clean()

		// when
		got, err := sut.GetPullRequest(context.Background(), repo, num)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if got.GetTitle() != want {
			t.Errorf("must be equal, but diff %+v", cmp.Diff(got.GetTitle(), want))
		}
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		// given
		repo := &github.MockRepository{
			FullName: "duck8823/duci",
		}
		num := 19

		// and
		gock.New("https://api.github.com").
			Get(fmt.Sprintf("/repos/%s/pulls/%d", repo.FullName, num)).
			Reply(404)
		defer gock.Clean()

		// when
		pr, err := sut.GetPullRequest(context.Background(), repo, num)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if pr != nil {
			t.Errorf("must be nil, but got %+v", pr)
		}
	})

	t.Run("with invalid repository", func(t *testing.T) {
		// given
		repo := &github.MockRepository{
			FullName: "",
		}

		// expect
		if _, err := sut.GetPullRequest(context.Background(), repo, 19); err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestClient_CreateCommitStatus(t *testing.T) {
	// given
	_ = github.Initialize("github_api_token")
	sut, err := github.GetInstance()
	if err != nil {
		t.Fatalf("error occurred. %+v", err)
	}

	t.Run("when github server returns status ok", func(t *testing.T) {
		// given
		status := github.CommitStatus{
			TargetSource: &github.TargetSource{
				Repository: &github.MockRepository{
					FullName: "duck8823/duci",
				},
				SHA: plumbing.ComputeHash(plumbing.AnyObject, []byte(random.String(16, random.Alphanumeric))),
			},
			State:       github.SUCCESS,
			Description: "hello world",
			Context:     "duci test",
			TargetURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
		}

		// and
		gock.New("https://api.github.com").
			Post(fmt.Sprintf("/repos/%s/statuses/%s", status.TargetSource.Repository.GetFullName(), status.TargetSource.SHA)).
			Reply(200)
		defer gock.Clean()

		// expect
		if err := sut.CreateCommitStatus(
			context.Background(),
			status,
		); err != nil {
			t.Errorf("error must be nil: but got %+v", err)
		}
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		// given
		status := github.CommitStatus{
			TargetSource: &github.TargetSource{
				Repository: &github.MockRepository{
					FullName: "duck8823/duci",
				},
				SHA: plumbing.ComputeHash(plumbing.AnyObject, []byte(random.String(16, random.Alphanumeric))),
			},
			State:       github.SUCCESS,
			Description: "hello world",
			Context:     "duci test",
			TargetURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
		}

		// and
		gock.New("https://api.github.com").
			Post(fmt.Sprintf("/repos/%s/statuses/%s", status.TargetSource.Repository.GetFullName(), status.TargetSource.SHA)).
			Reply(404)
		defer gock.Clean()

		// expect
		if err := sut.CreateCommitStatus(
			context.Background(),
			status,
		); err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("with invalid repository", func(t *testing.T) {
		// given
		status := github.CommitStatus{
			TargetSource: &github.TargetSource{
				Repository: &github.MockRepository{
					FullName: "",
				},
			},
			TargetURL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
		}

		// expect
		if err := sut.CreateCommitStatus(
			context.Background(),
			status,
		); err == nil {
			t.Error("error must not be nil")
		}
	})
}
