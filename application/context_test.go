package application_test

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
)

func TestContextWithJob(t *testing.T) {
	// given
	want := &application.BuildJob{
		ID: job.ID(uuid.New()),
	}

	// and
	ctx := application.ContextWithJob(context.Background(), want)

	// when
	got := ctx.Value(application.GetCtxKey())

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}

func TestBuildJobFromContext(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// given
		want := &application.BuildJob{
			ID: job.ID(uuid.New()),
		}

		sut := context.WithValue(context.Background(), application.GetCtxKey(), want)

		// when
		got, err := application.BuildJobFromContext(sut)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})

	t.Run("without value", func(t *testing.T) {
		sut := context.Background()

		// when
		got, err := application.BuildJobFromContext(sut)

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
