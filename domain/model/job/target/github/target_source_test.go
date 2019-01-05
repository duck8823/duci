package github_test

import (
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/gommon/random"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"testing"
)

func TestTargetSource_GetRef(t *testing.T) {
	// given
	want := "ref"

	// and
	sut := &github.TargetSource{
		Ref: want,
	}

	// when
	got := sut.GetRef()

	// then
	if got != want {
		t.Errorf("must be euqal, but %+v", cmp.Diff(got, want))
	}

}

func TestTargetSource_GetSHA(t *testing.T) {
	// given
	want := plumbing.ComputeHash(plumbing.AnyObject, []byte(random.String(16, random.Alphanumeric)))

	// and
	sut := &github.TargetSource{
		SHA: want,
	}

	// when
	got := sut.GetSHA()

	// then
	if got != want {
		t.Errorf("must be euqal, but %+v", cmp.Diff(got, want))
	}
}
