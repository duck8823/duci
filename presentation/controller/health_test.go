package controller_test

import (
	"github.com/duck8823/duci/application/service/docker/mock_docker"
	"github.com/duck8823/duci/presentation/controller"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckController_ServeHTTP(t *testing.T) {
	t.Run("without error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_docker.NewMockService(ctrl)
		mockService.EXPECT().
			Status().
			Return(nil)

		handler := &controller.HealthController{Docker: mockService}

		// and
		s := httptest.NewServer(handler)
		defer s.Close()

		// and
		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rec, req)

		// then
		if rec.Code != 200 {
			t.Errorf("status code must be 200, but got %+v", rec.Code)
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_docker.NewMockService(ctrl)
		mockService.EXPECT().
			Status().
			Return(errors.New("test"))

		handler := &controller.HealthController{Docker: mockService}

		// and
		s := httptest.NewServer(handler)
		defer s.Close()

		// and
		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rec, req)

		// then
		if rec.Code != 500 {
			t.Errorf("status code must be 500, but got %+v", rec.Code)
		}
	})

}
