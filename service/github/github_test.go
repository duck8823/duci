package github_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duck8823/minimal-ci/service/github"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"
	"time"
)

type MockHandler struct {
	Body   interface{}
	Status int
}

func (h *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(h.Body)
	w.Write(payload)
	w.WriteHeader(h.Status)
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

func TestService_GetPullRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/repos/duck8823/minimal-ci/pulls/5", &MockHandler{
		Body: struct {
			Id int64 `json:"id"`
		}{Id: 19},
		Status: 200,
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	baseUrl, err := url.Parse(ts.URL + "/")
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	s := github.New(context.Background(), "")
	s.Client.BaseURL = baseUrl

	repo := &MockRepo{
		FullName: "duck8823/minimal-ci",
		SSHURL:   "git@github.com:duck8823/minimal-ci.git",
	}
	pr, err := s.GetPullRequest(context.Background(), repo, 5)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	actual := pr.GetID()
	expected := 19
	t.Logf("%+v", pr)
	if pr.GetID() != 19 {
		t.Errorf("id must be equal %+v, but got %+v. \npr=%+v", expected, actual, pr)
	}
}

func TestService_CreateCommitStatus(t *testing.T) {
	t.Run("when github server returns status ok", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/repos/duck8823/minimal-ci/statuses/0000000000000000000000000000000000000000", &MockHandler{
			Status: 200,
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		baseUrl, err := url.Parse(ts.URL + "/")
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		s := github.New(context.Background(), "")
		s.Client.BaseURL = baseUrl

		repo := &MockRepo{
			FullName: "duck8823/minimal-ci",
			SSHURL:   "git@github.com:duck8823/minimal-ci.git",
		}
		if err := s.CreateCommitStatus(context.Background(), repo, plumbing.Hash{}, &github.Status{}); err != nil {
			t.Errorf("failed to create commit status: %+v", err)
		}
	})

	t.Run("when github server returns status not found", func(t *testing.T) {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()

		baseUrl, err := url.Parse(ts.URL + "/")
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		s := github.New(context.Background(), "")
		s.Client.BaseURL = baseUrl

		repo := &MockRepo{
			FullName: "duck8823/minimal-ci",
			SSHURL:   "git@github.com:duck8823/minimal-ci.git",
		}
		if err := s.CreateCommitStatus(context.Background(), repo, plumbing.Hash{}, &github.Status{}); err == nil {
			t.Error("must error")
		}
	})
}

func TestService_Clone(t *testing.T) {
	tempDir := path.Join(os.TempDir(), fmt.Sprintf("minimal-ci_test_%v", time.Now().Unix()))
	if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	s := github.New(context.Background(), "")

	repo := &MockRepo{
		FullName: "duck8823/minimal-ci",
		SSHURL:   "git@github.com:duck8823/minimal-ci.git",
	}

	_, err := s.Clone(context.Background(), tempDir, repo, "refs/heads/master")
	if err != nil {
		t.Errorf("must not error. %+v", err)
	}

	if _, err := os.Stat(path.Join(tempDir, ".git")); err != nil {
		t.Errorf("must be created dir: %s", path.Join(tempDir, ".git"))
	}

}
