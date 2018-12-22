package job_service_test

import (
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/internal/container"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/gommon/random"
	"os"
	"path"
	"testing"
)

func TestInitialize(t *testing.T) {
	t.Run("with temporary directory", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// when
		err := job_service.Initialize(tmpDir)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("with invalid directory", func(t *testing.T) {
		// given
		tmpDir := path.Join("/path/to/invalid/dir")
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// when
		err := job_service.Initialize(tmpDir)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func TestGetInstance(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := job_service.GetInstance()

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", err)
		}
	})

	t.Run("when instance is not nil", func(t *testing.T) {
		// given
		want := &job_service.StubService{
			ID: random.String(16, random.Alphanumeric),
		}

		// and
		container.Override(want)
		defer container.Clear()

		// when
		got, err := job_service.GetInstance()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})
}
