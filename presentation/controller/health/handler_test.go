package health_test

import (
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/docker/mock_docker"
	"github.com/duck8823/duci/presentation/controller/health"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewHandler(t *testing.T) {
	t.Run("with default", func(t *testing.T) {
		// given
		defaultDocker, err := docker.New()
		if err != nil {
			t.Fatalf("error occur: %+v", err)
		}

		// and
		want := &health.Handler{}
		defer want.SetDocker(defaultDocker)()

		// when
		got, err := health.NewHandler()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		opts := cmp.Options{
			cmp.AllowUnexported(health.Handler{}),
			cmpopts.IgnoreInterfaces(struct{ docker.Moby }{}),
		}
		if !cmp.Equal(got, want, opts) {
			t.Errorf("must be equal, but: %+v", cmp.Diff(got, want, opts))
		}
	})

	t.Run("with invalid environment variable", func(t *testing.T) {
		// given
		DOCKER_HOST := os.Getenv("DOCKER_HOST")
		if err := os.Setenv("DOCKER_HOST", "invalid host"); err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		defer os.Setenv("DOCKER_HOST", DOCKER_HOST)

		// when
		got, err := health.NewHandler()

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
	t.Run("when docker status returns no error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		docker := mock_docker.NewMockDocker(ctrl)
		docker.EXPECT().
			Status().
			Return(nil)

		// and
		sut := &health.Handler{}
		defer sut.SetDocker(docker)()

		// when
		sut.ServeHTTP(rec, req)

		// then
		if rec.Code != http.StatusOK {
			t.Errorf("must be %d, but got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("when docker status returns error", func(t *testing.T) {
		// given
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		docker := mock_docker.NewMockDocker(ctrl)
		docker.EXPECT().
			Status().
			Return(errors.New("test error"))

		// and
		sut := &health.Handler{}
		defer sut.SetDocker(docker)()

		// when
		sut.ServeHTTP(rec, req)

		// then
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("response code must be %d, but got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}
