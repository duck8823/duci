package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/application/service/github/mock_github"
	"github.com/duck8823/duci/application/service/runner/mock_runner"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJobController_ServeHTTP(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("with correct payload", func(t *testing.T) {
		t.Run("when issue_comment", func(t *testing.T) {
			// given
			event := "issue_comment"

			t.Run("when github service returns no error", func(t *testing.T) {
				// given
				runner := mock_runner.NewMockRunner(ctrl)
				runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

				githubService := mock_github.NewMockService(ctrl)
				githubService.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&github.PullRequest{
						Head: &github.PullRequestBranch{},
					}, nil)

				// and
				handler := &controller.JobController{Runner: runner, GitHub: githubService}

				s := httptest.NewServer(handler)
				defer s.Close()

				// and
				payload := createIssueCommentPayload(t, "ci test")

				req := httptest.NewRequest("POST", "/", payload)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 200 {
					t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
				}
			})

			t.Run("when github service return error", func(t *testing.T) {
				// given
				runner := mock_runner.NewMockRunner(ctrl)
				runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				githubService := mock_github.NewMockService(ctrl)
				githubService.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error occur"))

				// and
				handler := &controller.JobController{Runner: runner, GitHub: githubService}

				s := httptest.NewServer(handler)
				defer s.Close()

				// and
				payload := createIssueCommentPayload(t, "ci test")

				req := httptest.NewRequest("POST", "/", payload)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 500 {
					t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
				}
			})
		})

		t.Run("when push", func(t *testing.T) {
			// given
			runner := mock_runner.NewMockRunner(ctrl)
			runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

			// and
			handler := &controller.JobController{Runner: runner}

			s := httptest.NewServer(handler)
			defer s.Close()

			// and
			payload := createPushPayload(t, "test/repo", "master")

			req := httptest.NewRequest("POST", "/", payload)
			req.Header.Set("X-GitHub-Event", "push")
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 200 {
				t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
			}
		})
	})

	t.Run("with invalid payload", func(t *testing.T) {
		// setup
		runner := mock_runner.NewMockRunner(ctrl)
		runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

		handler := &controller.JobController{Runner: runner}

		s := httptest.NewServer(handler)
		defer s.Close()

		t.Run("with invalid header", func(t *testing.T) {
			// given
			body := createIssueCommentPayload(t, "ci test")

			// and
			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Event", "hogefuga")
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 500 {
				t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
			}
		})

		t.Run("with issue_comment", func(t *testing.T) {
			// given
			event := "issue_comment"

			t.Run("without comment started ci", func(t *testing.T) {
				// given
				body := createIssueCommentPayload(t, "test")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 200 {
					t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
				}
			})

			t.Run("with invalid body", func(t *testing.T) {
				// given
				body := strings.NewReader("Invalid JSON format.")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 500 {
					t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
				}
			})
		})

		t.Run("with push", func(t *testing.T) {
			// given
			event := "push"

			t.Run("with invalid body", func(t *testing.T) {
				// given
				body := strings.NewReader("Invalid JSON format.")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 500 {
					t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
				}
			})
		})
	})
}

func createIssueCommentPayload(t *testing.T, comment string) io.Reader {
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

func createPushPayload(t *testing.T, repoName, ref string) io.Reader {
	t.Helper()

	event := github.PushEvent{
		Repo: &github.PushEventRepository{
			FullName: &repoName,
		},
		Ref: &ref,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	return bytes.NewReader(payload)
}
