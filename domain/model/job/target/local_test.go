package target_test

import (
	"github.com/duck8823/duci/domain/model/job/target"
	"testing"
)

func TestLocal_Prepare(t *testing.T) {
	// given
	sut := &target.Local{Path: "."}

	// when
	dir, cleanup, err := sut.Prepare()
	defer cleanup()

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}

	// and
	if len(dir) == 0 {
		t.Errorf("must not be empty")
	}
}