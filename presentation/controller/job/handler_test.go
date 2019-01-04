package job_test

import (
	"context"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/application/service/job/mock_job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/internal/container"
	jobController "github.com/duck8823/duci/presentation/controller/job"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	t.Run("when there is service in container", func(t *testing.T) {
		// given
		service := new(jobService.Service)

		container.Override(service)
		defer container.Clear()

		// and
		want := &jobController.Handler{}
		defer want.SetService(*service)()

		// when
		got, err := jobController.NewHandler()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		opts := cmp.Options{
			cmp.AllowUnexported(jobController.Handler{}),
		}
		if !cmp.Equal(got, want, opts) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want, opts))
		}
	})

	t.Run("when there are no service in container", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := jobController.NewHandler()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}

	})
}

func TestHandler_ServeHTTP(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		id := job.ID(uuid.New())

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("uuid", uuid.UUID(id).String())
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(&job.Job{
				ID:       id,
				Finished: true,
				Stream: []job.LogLine{
					{Timestamp: time.Now(), Message: "Hello Test"},
				},
			}, nil)

		// and
		sut := &jobController.Handler{}
		defer sut.SetService(service)()

		// when
		sut.ServeHTTP(rec, req.WithContext(ctx))

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("with invalid path param", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("uuid", "")
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)
		req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

		// and
		sut := &jobController.Handler{}

		// when
		sut.ServeHTTP(rec, req)

		// then
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("must be %d, but got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("when service returns error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		id := job.ID(uuid.New())

		routeCtx := chi.NewRouteContext()
		routeCtx.URLParams.Add("uuid", uuid.UUID(id).String())
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, routeCtx)

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := mock_job_service.NewMockService(ctrl)
		service.EXPECT().
			FindBy(gomock.Eq(id)).
			Times(1).
			Return(nil, errors.New("test error"))

		// and
		sut := &jobController.Handler{}
		defer sut.SetService(service)()

		// when
		sut.ServeHTTP(rec, req.WithContext(ctx))

		// then
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("response code must be %d, but got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}
