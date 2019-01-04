package application_test

import (
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/duck8823/duci/internal/container"
	"github.com/labstack/gommon/random"
	"os"
	"path"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("when singleton container is empty", func(t *testing.T) {
		// given
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		// and
		container.Clear()

		// when
		err := application.Initialize()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		git := new(git.Git)
		if err := container.Get(git); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		github := new(github.GitHub)
		if err := container.Get(github); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		jobService := new(job.Service)
		if err := container.Get(jobService); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when singleton container contains Git instance", func(t *testing.T) {
		// given
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		// and
		container.Override(new(git.Git))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when singleton container contains GitHub instance", func(t *testing.T) {
		// given
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		// and
		container.Override(new(github.GitHub))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when singleton container contains JobService instance", func(t *testing.T) {
		// given
		application.Config.Server.DatabasePath = path.Join(os.TempDir(), random.String(16, random.Alphanumeric))

		// and
		container.Override(new(job.Service))
		defer container.Clear()

		// when
		err := application.Initialize()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}
