package github_test

import (
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSimpleTargetPoint_GetRef(t *testing.T) {
	// given
	want := "ref"

	// and
	sut := &github.SimpleTargetPoint{
		Ref: want,
	}

	// when
	got := sut.GetRef()

	// then
	if got != want {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}

func TestSimpleTargetPoint_GetHead(t *testing.T) {
	// given
	want := "sha"

	// and
	sut := &github.SimpleTargetPoint{
		SHA: want,
	}

	// when
	got := sut.GetHead()

	// then
	if got != want {
		t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
	}
}
