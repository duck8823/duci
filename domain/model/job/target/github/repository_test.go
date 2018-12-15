package github_test

import (
	"fmt"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestRepositoryName_Owner(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		// given
		want := "duck8823"

		// and
		sut := github.RepositoryName(fmt.Sprintf("%s/duci", want))

		// when
		got, err := sut.Owner()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if got != want {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		// where
		for _, tt := range []struct {
			name string
		}{
			{
				name: "",
			},
			{
				name: "duck8823",
			},
			{
				name: "duck8823/duci/domain",
			},
		} {
			t.Run(fmt.Sprintf("when name is %s", tt.name), func(t *testing.T) {
				// given
				sut := github.RepositoryName(tt.name)

				// when
				got, err := sut.Owner()

				// then
				if err == nil {
					t.Error("error must not be nil")
				}

				// and
				if got != "" {
					t.Errorf("must be empty, but got %+v", got)
				}
			})
		}
	})
}

func TestRepositoryName_Repo(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		// given
		want := "duci"

		// and
		sut := github.RepositoryName(fmt.Sprintf("duck8823/%s", want))

		// when
		got, err := sut.Repo()

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if got != want {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		// where
		for _, tt := range []struct {
			name string
		}{
			{
				name: "",
			},
			{
				name: "duci",
			},
			{
				name: "duck8823/duci/domain",
			},
		} {
			t.Run(fmt.Sprintf("when name is %s", tt.name), func(t *testing.T) {
				// given
				sut := github.RepositoryName(tt.name)

				// when
				got, err := sut.Repo()

				// then
				if err == nil {
					t.Error("error must not be nil")
				}

				// and
				if got != "" {
					t.Errorf("must be empty, but got %+v", got)
				}
			})
		}
	})
}
