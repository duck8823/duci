package git_test

import (
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/internal/container"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetInstance(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		container.Clear()

		// when
		got, err := git.GetInstance()

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
		want := &git.HTTPGitClient{}

		// and
		container.Override(want)
		defer container.Clear()

		// when
		got, err := git.GetInstance()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		opt := cmp.AllowUnexported(git.HTTPGitClient{})
		if !cmp.Equal(got, want, opt) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want, opt))
		}
	})
}
