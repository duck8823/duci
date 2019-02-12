package webhook_test

import (
	"context"
	"fmt"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor/mock_executor"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/domain/model/job/target/github/mock_github"
	"github.com/duck8823/duci/internal/container"
	"github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	t.Run("when there are job service and github in container", func(t *testing.T) {
		// given
		container.Override(new(jobService.Service))
		container.Override(new(github.GitHub))
		defer container.Clear()

		// when
		_, err := webhook.NewHandler()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when there are not enough instance in container", func(t *testing.T) {
		// where
		for _, tt := range []struct {
			name  string
			given func()
		}{
			{
				name: "without job service",
				given: func() {
					container.Override(new(github.GitHub))
				},
			},
			{
				name: "without github",
				given: func() {
					container.Override(new(jobService.Service))
				},
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				// given
				container.Clear()
				tt.given()

				// when
				_, err := webhook.NewHandler()

				// then
				if err == nil {
					t.Error("error must not be nil")
				}

				// cleanup
				container.Clear()
			})
		}
	})
}

func TestHandler_ServeHTTP(t *testing.T) {
	for _, tt := range []struct {
		event   string
		payload string
	}{
		{
			event:   "push",
			payload: "testdata/push.correct.json",
		},
		{
			event:   "pull_request",
			payload: "testdata/pr.synchronize.json",
		},
	} {
		t.Run(fmt.Sprintf("when %s event", tt.event), func(t *testing.T) {
			// given
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			// and
			req.Header.Set("X-GitHub-Event", tt.event)
			req.Header.Set("X-GitHub-Delivery", "72d3162e-cc78-11e3-81ab-4c9367dc0958")

			// and
			f, err := os.Open(tt.payload)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			req.Body = f

			// and
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			executor := mock_executor.NewMockExecutor(ctrl)
			executor.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Times(1).
				Return(nil)

			// and
			sut := &webhook.Handler{}
			reset := sut.SetExecutor(executor)
			defer func() {
				time.Sleep(10 * time.Millisecond) // for goroutine
				reset()
			}()

			// when
			sut.ServeHTTP(rec, req)

			// then
			if rec.Code != http.StatusOK {
				t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
			}
		})
	}

	t.Run("when pull request comment event", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header.Set("X-GitHub-Event", "issue_comment")
		req.Header.Set("X-GitHub-Delivery", "72d3162e-cc78-11e3-81ab-4c9367dc0958")

		// and
		f, err := os.Open("testdata/issue_comment.correct.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gh := mock_github.NewMockGitHub(ctrl)
		gh.EXPECT().
			GetPullRequest(gomock.Any(), gomock.Any(), gomock.Eq(2)).
			Times(1).
			Return(&go_github.PullRequest{
				Head: &go_github.PullRequestBranch{
					Ref: go_github.String("refs/test/dummy"),
					SHA: go_github.String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				},
			}, nil)
		container.Override(gh)
		defer container.Clear()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.ServeHTTP(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("when other event", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		sut := &webhook.Handler{}

		// when
		sut.ServeHTTP(rec, req)

		// then
		if rec.Code != http.StatusBadRequest {
			t.Errorf("response code must be %d, but got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestHandler_PushEvent(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/push.correct.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(1).
			Do(func(ctx context.Context, target job.Target) {
				got, err := application.BuildJobFromContext(ctx)
				if err != nil {
					t.Errorf("must not be nil, but got %+v", err)
				}

				want := &application.BuildJob{
					ID: job.ID(uuid.Must(uuid.Parse("72d3162e-cc78-11e3-81ab-4c9367dc0958"))),
					TargetSource: &github.TargetSource{
						Repository: &go_github.PushEventRepository{
							ID:       go_github.Int64(135493233),
							FullName: go_github.String("Codertocat/Hello-World"),
							SSHURL:   go_github.String("git@github.com:Codertocat/Hello-World.git"),
							CloneURL: go_github.String("https://github.com/Codertocat/Hello-World.git"),
						},
						Ref: "refs/tags/simple-tag",
						SHA: plumbing.ZeroHash,
					},
					TaskName:  "duci/push",
					TargetURL: webhook.URLMust(url.Parse("http://example.com/logs/72d3162e-cc78-11e3-81ab-4c9367dc0958")),
				}

				opt := webhook.CmpOptsAllowFields(go_github.PushEventRepository{}, "ID", "FullName", "SSHURL", "CloneURL")
				if !cmp.Equal(got, want, opt) {
					t.Errorf("must be equal but: %+v", cmp.Diff(got, want, opt))
				}

				typ := reflect.TypeOf(target).String()
				if typ != "*target.GitHub" {
					t.Errorf("type must be *target.GitHub, but got %s", typ)
				}
			}).
			Return(nil)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.PushEvent(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("when url param is invalid format uuid", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"invalid format"},
		}

		// and
		f, err := os.Open("testdata/push.correct.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.PushEvent(rec, req)

		// then
		if rec.Code != http.StatusBadRequest {
			t.Errorf("response code must be %d, but got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("with invalid payload", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		req.Body = ioutils.NewReadCloserWrapper(strings.NewReader("invalid payload"), func() error {
			return nil
		})

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.PushEvent(rec, req)

		// then
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("response code must be %d, but got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestHandler_IssueCommentEvent_Normal(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.correct.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gh := mock_github.NewMockGitHub(ctrl)
		gh.EXPECT().
			GetPullRequest(gomock.Any(), gomock.Any(), gomock.Eq(2)).
			Times(1).
			Return(&go_github.PullRequest{
				Head: &go_github.PullRequestBranch{
					Ref: go_github.String("dummy"),
					SHA: go_github.String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				},
			}, nil)
		container.Override(gh)
		defer container.Clear()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Do(func(ctx context.Context, target job.Target, cmd ...string) {
				got, err := application.BuildJobFromContext(ctx)
				if err != nil {
					t.Errorf("must not be nil, but got %+v", err)
				}

				want := &application.BuildJob{
					ID: job.ID(uuid.Must(uuid.Parse("72d3162e-cc78-11e3-81ab-4c9367dc0958"))),
					TargetSource: &github.TargetSource{
						Repository: &go_github.Repository{
							ID:       go_github.Int64(135493233),
							FullName: go_github.String("Codertocat/Hello-World"),
							SSHURL:   go_github.String("git@github.com:Codertocat/Hello-World.git"),
							CloneURL: go_github.String("https://github.com/Codertocat/Hello-World.git"),
						},
						Ref: "refs/heads/dummy",
						SHA: plumbing.NewHash("aa218f56b14c9653891f9e74264a383fa43fefbd"),
					},
					TaskName:  "duci/pr/build",
					TargetURL: webhook.URLMust(url.Parse("http://example.com/logs/72d3162e-cc78-11e3-81ab-4c9367dc0958")),
				}

				opt := webhook.CmpOptsAllowFields(go_github.Repository{}, "ID", "FullName", "SSHURL", "CloneURL")
				if !cmp.Equal(got, want, opt) {
					t.Errorf("must be equal but: %+v", cmp.Diff(got, want, opt))
				}

				typ := reflect.TypeOf(target).String()
				if typ != "*target.GitHub" {
					t.Errorf("type must be *target.GitHub, but got %s", typ)
				}
			}).
			Return(nil)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("when no match comment", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.skip_comment.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gh := mock_github.NewMockGitHub(ctrl)
		gh.EXPECT().
			GetPullRequest(gomock.Any(), gomock.Any(), gomock.Eq(2)).
			Times(1).
			Return(&go_github.PullRequest{
				Head: &go_github.PullRequestBranch{
					Ref: go_github.String("refs/test/dummy"),
					SHA: go_github.String("aa218f56b14c9653891f9e74264a383fa43fefbd"),
				},
			}, nil)
		container.Override(gh)
		defer container.Clear()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}

		// and
		got := rec.Body.String()
		if got != `{"message":"skip build"}` {
			t.Errorf("must be equal. want %s, but got %s", `{"message":"skip build"}`, got)
		}
	})

	t.Run("when action is deleted", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.deleted.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}

		// and
		got := rec.Body.String()
		if got != `{"message":"skip build"}` {
			t.Errorf("must be equal. want %s, but got %s", `{"message":"skip build"}`, got)
		}
	})
}

func TestHandler_IssueCommentEvent_UnNormal(t *testing.T) {
	t.Run("with invalid payload body", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		req.Body = ioutils.NewReadCloserWrapper(strings.NewReader("invalid payload"), func() error {
			return nil
		})

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("response code must be %d, but got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("when url param is invalid format uuid", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"invalid format"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.correct.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusBadRequest {
			t.Errorf("response code must be %d, but got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("when fail to get pull request", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.skip_comment.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gh := mock_github.NewMockGitHub(ctrl)
		gh.EXPECT().
			GetPullRequest(gomock.Any(), gomock.Any(), gomock.Eq(2)).
			Times(1).
			Return(nil, errors.New("test error"))
		container.Override(gh)
		defer container.Clear()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusBadRequest {
			t.Errorf("response code must be %d, but got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("when fail to get github instance", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/issue_comment.skip_comment.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		container.Clear()

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		defer sut.SetExecutor(executor)()

		// when
		sut.IssueCommentEvent(rec, req)

		// then
		if rec.Code != http.StatusBadRequest {
			t.Errorf("response code must be %d, but got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestHandler_PullRequestEvent(t *testing.T) {
	for _, tt := range []struct {
		action  string
		payload string
	}{
		{
			action:  "opened",
			payload: "testdata/pr.opened.json",
		},
		{
			action:  "synchronize",
			payload: "testdata/pr.synchronize.json",
		},
	} {
		t.Run(fmt.Sprintf("when action is %s", tt.action), func(t *testing.T) {
			// given
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			// and
			req.Header = http.Header{
				"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
			}

			// and
			f, err := os.Open(tt.payload)
			if err != nil {
				t.Fatalf("error occur: %+v", err)
			}
			req.Body = f

			// and
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			executor := mock_executor.NewMockExecutor(ctrl)
			executor.EXPECT().
				Execute(gomock.Any(), gomock.Any()).
				Times(1).
				Do(func(ctx context.Context, target job.Target) {
					got, err := application.BuildJobFromContext(ctx)
					if err != nil {
						t.Errorf("must not be nil, but got %+v", err)
					}

					want := &application.BuildJob{
						ID: job.ID(uuid.Must(uuid.Parse("72d3162e-cc78-11e3-81ab-4c9367dc0958"))),
						TargetSource: &github.TargetSource{
							Repository: &go_github.Repository{
								ID:       go_github.Int64(135493233),
								FullName: go_github.String("Codertocat/Hello-World"),
								SSHURL:   go_github.String("git@github.com:Codertocat/Hello-World.git"),
								CloneURL: go_github.String("https://github.com/Codertocat/Hello-World.git"),
							},
							Ref: "refs/heads/changes",
							SHA: plumbing.NewHash("34c5c7793cb3b279e22454cb6750c80560547b3a"),
						},
						TaskName:  "duci/pr",
						TargetURL: webhook.URLMust(url.Parse("http://example.com/logs/72d3162e-cc78-11e3-81ab-4c9367dc0958")),
					}

					opt := webhook.CmpOptsAllowFields(go_github.Repository{}, "ID", "FullName", "SSHURL", "CloneURL")
					if !cmp.Equal(got, want, opt) {
						t.Errorf("must be equal but: %+v", cmp.Diff(got, want, opt))
					}

					typ := reflect.TypeOf(target).String()
					if typ != "*target.GitHub" {
						t.Errorf("type must be *target.GitHub, but got %s", typ)
					}
				}).
				Return(nil)

			// and
			sut := &webhook.Handler{}
			reset := sut.SetExecutor(executor)
			defer func() {
				time.Sleep(10 * time.Millisecond) // for goroutine
				reset()
			}()

			// when
			sut.PullRequestEvent(rec, req)

			// then
			if rec.Code != http.StatusOK {
				t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
			}
		})
	}

	t.Run("when pull request closed", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		req.Header = http.Header{
			"X-Github-Delivery": []string{"72d3162e-cc78-11e3-81ab-4c9367dc0958"},
		}

		// and
		f, err := os.Open("testdata/pr.closed.json")
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		req.Body = f

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		executor := mock_executor.NewMockExecutor(ctrl)
		executor.EXPECT().
			Execute(gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &webhook.Handler{}
		reset := sut.SetExecutor(executor)
		defer func() {
			time.Sleep(10 * time.Millisecond) // for goroutine
			reset()
		}()

		// when
		sut.PullRequestEvent(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})
}
