package webhook_test

import (
	"context"
	"github.com/duck8823/duci/application/service/executor/mock_executor"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/duck8823/duci/presentation/controller/webhook"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"os"
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
			// TODO: check argument
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
