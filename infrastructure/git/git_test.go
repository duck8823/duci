package git_test

import (
	"fmt"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/git"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("when missing ssh key", func(t *testing.T) {
		// setup
		home := os.Getenv("HOME")
		os.Setenv("HOME", "/path/to/wrong/")

		// expect
		if _, err := git.New(); err == nil {
			t.Error("error must occur")
		}

		// cleanup
		os.Setenv("HOME", home)
	})
}

func TestSshGitClient_Clone(t *testing.T) {
	t.Run("with correct key", func(t *testing.T) {
		// setup
		client, err := git.New()
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		t.Run("when target directory exists", func(t *testing.T) {
			// setup
			tempDir := path.Join(os.TempDir(), fmt.Sprintf("duci_test_%v", time.Now().Unix()))
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// given
			var empty plumbing.Hash

			// when
			hash, err := client.Clone(
				context.New("test/task"),
				tempDir,
				"git@github.com:duck8823/duci.git",
				"refs/heads/master",
			)

			// then
			if err != nil {
				t.Errorf("must not error. %+v", err)
			}

			if _, err := os.Stat(path.Join(tempDir, ".git")); err != nil {
				t.Errorf("must create dir: %s", path.Join(tempDir, ".git"))
			}

			if hash == empty {
				t.Errorf("commit hash must not be empty")
			}
		})

		t.Run("when target directory not exists", func(t *testing.T) {
			// given
			wrongPath := "/path/to/not/exists"

			// when
			_, err := client.Clone(
				context.New("test/task"),
				wrongPath,
				"git@github.com:duck8823/duci.git",
				"refs/heads/master",
			)

			// then
			if err == nil {
				t.Error("erro must occur")
			}
		})
	})
}
