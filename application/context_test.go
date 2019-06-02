package application_test

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestContextWithJob(t *testing.T) {
	// given
	opts := cmp.AllowUnexported(application.BuildJob{})

	// and
	want := &application.BuildJob{
		ID: job.ID(uuid.New()),
	}

	// and
	ctx := application.ContextWithJob(context.Background(), want)

	// when
	got := ctx.Value(application.GetCtxKey())

	// then
	if !cmp.Equal(got, want, opts) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want, opts))
	}
}

func TestBuildJobFromContext(t *testing.T) {
	t.Run("with value", func(t *testing.T) {
		// given
		opts := cmp.AllowUnexported(application.BuildJob{})

		// and
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
		if !cmp.Equal(got, want, opts) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want, opts))
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

func TestBuildJob_BeginAt(t *testing.T) {
	// given
	want := time.Unix(10, 10)

	// and
	sut := &application.BuildJob{}

	// when
	sut.BeginAt(want)
	got := sut.GetBeginTime()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}

func TestBuildJob_EndAt(t *testing.T) {
	// given
	want := time.Unix(10, 10)

	// and
	sut := &application.BuildJob{}

	// when
	sut.EndAt(want)
	got := sut.GetEndTime()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}

func TestBuildJob_Duration(t *testing.T) {
	tests := []struct {
		name      string
		beginTime time.Time
		endTime   time.Time
		want      string
	}{
		{
			"when duration is less than 60 seconds",
			time.Unix(10, 10),
			time.Unix(70, 0),
			"59sec",
		},
		{
			"when duration is greater than 60 seconds",
			time.Unix(10, 10),
			time.Unix(70, 11),
			"1min",
		},
		{
			"when duration is equal to 60 seconds",
			time.Unix(10, 10),
			time.Unix(70, 10),
			"1min",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			sut := &application.BuildJob{}
			sut.BeginAt(test.beginTime)
			sut.EndAt(test.endTime)

			// when
			got := sut.Duration()

			// then
			if got != test.want {
				t.Errorf("want %s, but got %s", test.want, got)
			}
		})
	}
}
