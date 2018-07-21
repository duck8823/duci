package github_test

import (
	"github.com/duck8823/duci/application/service/github"
	"testing"
)

func TestRepositoryName_Owner(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		// given
		name := &github.RepositoryName{FullName: "hoge/fuga"}

		// and
		expected := "hoge"

		// when
		owner, err := name.Owner()

		// then
		if err != nil {
			t.Fatalf("must not error. %+v", err)
		}
		if owner != expected {
			t.Errorf("owner must be equal %+v, but got %+v", expected, owner)
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		// given
		name := &github.RepositoryName{FullName: "hoge"}

		// when
		_, err := name.Owner()

		// then
		if err == nil {
			t.Fatalf("must error.")
		}
	})
}

func TestRepositoryName_Repo(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		// given
		name := &github.RepositoryName{FullName: "hoge/fuga"}

		// and
		expected := "fuga"

		// when
		owner, err := name.Repo()

		// then
		if err != nil {
			t.Fatalf("must not error. %+v", err)
		}
		if owner != expected {
			t.Errorf("owner must be equal %+v, but got %+v", expected, owner)
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		// given
		name := &github.RepositoryName{FullName: "hoge"}

		// when
		_, err := name.Repo()

		// then
		if err == nil {
			t.Fatalf("must error.")
		}
	})
}
