package controller_test

import (
	ctx "context"
	"github.com/duck8823/duci/application/service/log/mock_log"
	"github.com/duck8823/duci/domain/model"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http/httptest"
	"testing"
)

func TestLogsController_ServeHTTP(t *testing.T) {
	t.Run("with valid uuid", func(t *testing.T) {
		// setup
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_log.NewMockStoreService(ctrl)
		handler := &controller.LogController{LogService: mockService}

		// given
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		t.Run("when service returns error", func(t *testing.T) {
			// and
			mockService.EXPECT().
				Get(gomock.Eq(id)).
				Return(nil, errors.New("hello error"))

			// and
			s := httptest.NewServer(handler)
			defer s.Close()

			// and
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("uuid", id.String())

			req := httptest.NewRequest("GET", "/", nil).
				WithContext(ctx.WithValue(ctx.Background(), chi.RouteCtxKey, chiCtx))
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 500 {
				t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
			}
		})

		t.Run("when service returns correct job", func(t *testing.T) {
			// and
			job := &model.Job{
				Finished: true,
				Stream: []model.Message{{
					Time: clock.Now(),
					Text: "Hello World",
				}},
			}

			mockService.EXPECT().
				Get(gomock.Eq(id)).
				Return(job, nil)

			// and
			s := httptest.NewServer(handler)
			defer s.Close()

			// and
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("uuid", id.String())

			req := httptest.NewRequest("GET", "/", nil).
				WithContext(ctx.WithValue(ctx.Background(), chi.RouteCtxKey, chiCtx))
			rec := httptest.NewRecorder()

			// when
			handler.ServeHTTP(rec, req)

			// then
			if rec.Code != 200 {
				t.Errorf("status must equal %+v, but got %+v", 200, rec.Code)
			}
		})
	})

	t.Run("with invalid uuid", func(t *testing.T) {
		// setup
		handler := &controller.LogController{}

		// given
		s := httptest.NewServer(handler)
		defer s.Close()

		// and
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("uuid", "invalid_uuid")

		req := httptest.NewRequest("GET", "/", nil).
			WithContext(ctx.WithValue(ctx.Background(), chi.RouteCtxKey, chiCtx))
		rec := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rec, req)

		// then
		if rec.Code != 500 {
			t.Errorf("status must equal %+v, but got %+v", 500, rec.Code)
		}
	})
}
