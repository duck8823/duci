package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/application/service/github/mock_github"
	"github.com/duck8823/duci/application/service/runner/mock_runner"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/github"
	"github.com/google/uuid"
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
		// given
		requestId, _ := uuid.NewRandom()

		t.Run("when issue_comment", func(t *testing.T) {
			// given
			event := "issue_comment"

			t.Run("when github service returns no error", func(t *testing.T) {
				// given
				runner := mock_runner.NewMockRunner(ctrl)
				runner.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

				githubService := mock_github.NewMockService(ctrl)
				githubService.EXPECT().GetPullRequest(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(&github.PullRequest{
						Head: &github.PullRequestBranch{},
					}, nil)

				// and
				handler := &controller.JobController{Runner: runner, GitHub: githubService}

				s := httptest.NewServer(handler)
				defer s.Close()

				t.Run("with valid action", func(t *testing.T) {
					actions := []string{"created", "edited"}
					for _, action := range actions {
						// and
						payload := createIssueCommentPayload(t, action, "ci test")

						req := httptest.NewRequest("POST", "/", payload)
						req.Header.Set("X-GitHub-Delivery", requestId.String())
						req.Header.Set("X-GitHub-Event", event)
						rec := httptest.NewRecorder()

						// when
						handler.ServeHTTP(rec, req)

						// then
						if rec.Code != 200 {
							t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
						}
					}
				})

				t.Run("with invalid action", func(t *testing.T) {
					actions := []string{"deleted", "foo", ""}
					for _, action := range actions {
						t.Run(action, func(t *testing.T) {
							// given
							body := createIssueCommentPayload(t, action, "ci test")

							// and
							req := httptest.NewRequest("POST", "/", body)
							req.Header.Set("X-GitHub-Delivery", requestId.String())
							req.Header.Set("X-GitHub-Event", event)
							rec := httptest.NewRecorder()

							// when
							handler.ServeHTTP(rec, req)

							// then
							if rec.Code != 200 {
								t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
							}

							if rec.Body.String() != "build skip" {
								t.Errorf("body must equal %+v, but got %+v", "build skip", rec.Body.String())
							}
						})
					}
				})
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
				payload := createIssueCommentPayload(t, "created", "ci test")

				req := httptest.NewRequest("POST", "/", payload)
				req.Header.Set("X-GitHub-Delivery", requestId.String())
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
			req.Header.Set("X-GitHub-Delivery", requestId.String())
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

		t.Run("with invalid `X-GitHub-Event` header", func(t *testing.T) {
			// given
			body := createIssueCommentPayload(t, "created", "ci test")

			// and
			requestId, _ := uuid.NewRandom()

			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Delivery", requestId.String())
			req.Header.Set("X-GitHub-Event", "hogefuga")
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 500 {
				t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
			}
		})

		t.Run("with invalid `X-GitHub-Delivery` header", func(t *testing.T) {
			// given
			body := createIssueCommentPayload(t, "created", "ci test")

			// and
			req := httptest.NewRequest("POST", "/", body)
			req.Header.Set("X-GitHub-Delivery", "hogefuga")
			req.Header.Set("X-GitHub-Event", "push")
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 400 {
				t.Errorf("status must equal %+v, but got %+v", 400, rec.Code)
			}
		})

		t.Run("with issue_comment", func(t *testing.T) {
			// given
			event := "issue_comment"
			requestId, _ := uuid.NewRandom()

			t.Run("without comment started ci", func(t *testing.T) {
				// given
				body := createIssueCommentPayload(t, "created", "test")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Delivery", requestId.String())
				req.Header.Set("X-GitHub-Event", event)
				rec := httptest.NewRecorder()

				// when
				handler.ServeHTTP(rec, req)

				// then
				if rec.Code != 200 {
					t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
				}

				if rec.Body.String() != "build skip" {
					t.Errorf("body must equal %+v, but got %+v", "build skip", rec.Body.String())
				}
			})

			t.Run("with invalid body", func(t *testing.T) {
				// given
				body := strings.NewReader("Invalid JSON format.")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Delivery", requestId.String())
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
			requestId, _ := uuid.NewRandom()

			t.Run("with invalid body", func(t *testing.T) {
				// given
				body := strings.NewReader("Invalid JSON format.")

				// and
				req := httptest.NewRequest("POST", "/", body)
				req.Header.Set("X-GitHub-Delivery", requestId.String())
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

func createIssueCommentPayload(t *testing.T, action, comment string) io.Reader {
	t.Helper()

	number := 1
	event := &github.IssueCommentEvent{
		Repo:   &github.Repository{},
		Action: &action,
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
