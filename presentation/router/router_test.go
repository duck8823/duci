package router_test

import (
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/duck8823/duci/presentation/router"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		container.Override(new(job_service.Service))
		container.Override(new(github.GitHub))
		defer container.Clear()

		// when
		_, err := router.New()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when component not enough", func(t *testing.T) {
		// given
		container.Clear()

		// when
		_, err := router.New()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("with invalid environment variable for docker client", func(t *testing.T) {
		// given
		container.Override(new(job_service.Service))
		container.Override(new(github.GitHub))
		defer container.Clear()

		// and
		DOCKER_HOST := os.Getenv("DOCKER_HOST")
		if err := os.Setenv("DOCKER_HOST", "invalid_host"); err != nil {
			t.Fatalf("error occur: %+v", err)
		}
		defer func() {
			_ = os.Setenv("DOCKER_HOST", DOCKER_HOST)
		}()

		// when
		_, err := router.New()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}
