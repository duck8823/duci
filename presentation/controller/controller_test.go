package controller

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/minimal-ci/service/github/mock_github"
	"github.com/duck8823/minimal-ci/service/runner/mock_runner"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJobController_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("with correct payload", func(t *testing.T) {

		t.Run("when issue_comment", func(t *testing.T) {
			runner := mock_runner.NewMockRunner(ctrl)
			runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

			githubService := mock_github.NewMockService(ctrl)
			githubService.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).Return(&github.PullRequest{
				Head: &github.PullRequestBranch{},
			}, nil)

			handler := &jobController{runner: runner, github: githubService}

			s := httptest.NewServer(handler)
			defer s.Close()

			body := CreateBody(t, "ci test")

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "issue_comment")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 200
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})

		t.Run("when push", func(t *testing.T) {
			mock := mock_runner.NewMockRunner(ctrl)
			mock.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

			handler := &jobController{runner: mock}

			s := httptest.NewServer(handler)
			defer s.Close()

			repoName := "test/repo"
			ref := "ref"
			body, err := json.Marshal(&github.PushEvent{
				Repo: &github.PushEventRepository{
					FullName: &repoName,
				},
				Ref: &ref,
			})
			if err != nil {
				t.Fatalf("error occured: %+v", err)
			}

			req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Set("X-GitHub-Event", "push")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 200
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})

	})

	t.Run("with invalid payload", func(t *testing.T) {
		mock := mock_runner.NewMockRunner(ctrl)
		mock.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		handler := &jobController{runner: mock}

		s := httptest.NewServer(handler)
		defer s.Close()

		t.Run("with invalid header", func(t *testing.T) {
			body := CreateBody(t, "ci test")

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "hogefuga")
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			actual := rec.Code
			expected := 500
			if actual != expected {
				t.Errorf("status must equal %+v, but got %+v", expected, actual)
			}
		})

		t.Run("with issue_comment", func(t *testing.T) {
			event := "issue_comment"

			t.Run("without comment started ci", func(t *testing.T) {
				body := CreateBody(t, "test")

				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				actual := rec.Code
				expected := 200
				if actual != expected {
					t.Errorf("status must equal %+v, but got %+v", expected, actual)
				}
			})

			t.Run("with invalid body", func(t *testing.T) {
				body := strings.NewReader("Invalid JSON format.")

				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				actual := rec.Code
				expected := 500
				if actual != expected {
					t.Errorf("status must equal %+v, but got %+v", expected, actual)
				}
			})
		})

		t.Run("with push", func(t *testing.T) {
			event := "push"

			t.Run("with invalid body", func(t *testing.T) {
				body := strings.NewReader("Invalid JSON format.")

				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				handler.ServeHTTP(rec, req)

				actual := rec.Code
				expected := 500
				if actual != expected {
					t.Errorf("status must equal %+v, but got %+v", expected, actual)
				}
			})
		})
	})
}

func CreateBody(t *testing.T, comment string) io.Reader {
	t.Helper()

	number := 1
	event := &github.IssueCommentEvent{
		Repo: &github.Repository{},
		Issue: &github.Issue{
			Number: &number,
		},
		Comment: &github.IssueComment{
			Body: &comment,
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}
	return bytes.NewReader(payload)
}
