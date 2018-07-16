package github_test

import (
	"github.com/duck8823/duci/service/github"
	"testing"
)

func TestRepositoryName_Owner(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		name := &github.RepositoryName{FullName: "hoge/fuga"}
		owner, err := name.Owner()
		if err != nil {
			t.Fatalf("must not error. %+v", err)
		}

		expected := "hoge"
		if owner != expected {
			t.Errorf("owner must be equal %+v, but got %+v", expected, owner)
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		name := &github.RepositoryName{FullName: "hoge"}
		_, err := name.Owner()
		if err == nil {
			t.Fatalf("must error.")
		}
	})
}

func TestRepositoryName_Repo(t *testing.T) {
	t.Run("with correct name", func(t *testing.T) {
		name := &github.RepositoryName{FullName: "hoge/fuga"}
		owner, err := name.Repo()
		if err != nil {
			t.Fatalf("must not error. %+v", err)
		}

		expected := "fuga"
		if owner != expected {
			t.Errorf("owner must be equal %+v, but got %+v", expected, owner)
		}
	})

	t.Run("with invalid name", func(t *testing.T) {
		name := &github.RepositoryName{FullName: "hoge"}
		_, err := name.Repo()
		if err == nil {
			t.Fatalf("must error.")
		}
	})
}
