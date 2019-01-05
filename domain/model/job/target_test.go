package job_test

import (
	"github.com/duck8823/duci/domain/model/job"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestWorkDir_String(t *testing.T) {
	// given
	want := "/path/to/dir"

	// and
	sut := job.WorkDir(want)

	// when
	got := sut.String()

	// then
	if !cmp.Equal(got, want) {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}
