package job_test

import (
	"github.com/duck8823/duci/domain/model/job"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
)

func TestJob_AppendLog(t *testing.T) {
	// given
	want := []job.LogLine{{
		Message: "Hello World",
	}}

	// and
	sut := job.Job{}

	// when
	sut.AppendLog(want[0])

	// then
	got := sut.Stream
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}

func TestJob_Finish(t *testing.T) {
	// given
	sut := job.Job{}

	// when
	sut.Finish()

	// then
	if !sut.Finished {
		t.Errorf("must be true, but false")
	}
}

func TestJob_ToBytes(t *testing.T) {
	t.Run("when success marshal", func(t *testing.T) {
		// given
		want := []byte("{\"ID\":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],\"finished\":false,\"stream\":[]}")

		// and
		sut := job.Job{
			ID: job.ID(uuid.Nil),
			Finished: false,
			Stream: []job.LogLine{},
		}

		// when
		got, err := sut.ToBytes()

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

func TestID_ToSlice(t *testing.T) {
	// given
	want := []byte(uuid.New().String())

	// and
	id, _ := uuid.ParseBytes(want)
	sut := job.ID(id)

	// when
	got := sut.ToSlice()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(string(got), string(want)))
	}
}