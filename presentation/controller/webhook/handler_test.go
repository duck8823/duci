package webhook_test

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor/mock_executor"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	go_github "github.com/google/go-github/github"
	"github.com/google/uuid"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("when there are job service and github in container", func(t *testing.T) {
		// given
		container.Override(new(job_service.Service))
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
					container.Override(new(job_service.Service))
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

func TestHandler_PushEvent(t *testing.T) {
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
				TargetURL: webhook.URLMust(url.Parse("http://example.com/")),
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
	defer sut.SetExecutor(executor)()

	// when
	sut.PushEvent(rec, req)

	// then
	if rec.Code != http.StatusOK {
		t.Errorf("response code must be %d, but got %d", http.StatusOK, rec.Code)
	}
}
