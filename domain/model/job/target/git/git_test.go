package git_test

import (
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetInstance(t *testing.T) {
	t.Run("when instance is nil", func(t *testing.T) {
		// given
		defer git.SetInstance(nil)()

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
		want := &git.HttpGitClient{}

		// and
		defer git.SetInstance(want)()

		// when
		got, err := git.GetInstance()

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
