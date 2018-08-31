package git_test

import (
	"fmt"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/service/git"
	"github.com/google/uuid"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/url"
	"os"
	"path"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("when missing ssh key", func(t *testing.T) {
		// expect
		if _, err := git.New("/path/to/wrong/"); err == nil {
			t.Error("error must occur")
		}
	})
}

func TestSshGitService_Clone(t *testing.T) {
	t.Run("with correct key", func(t *testing.T) {
		// setup
		client, err := git.New(path.Join(os.Getenv("HOME"), ".ssh/id_rsa"))
		if err != nil {
			t.Fatalf("error occured. %+v", err)
		}

		t.Run("when target directory exists", func(t *testing.T) {
			// setup
			tempDir := path.Join(os.TempDir(), fmt.Sprintf("duci_test_%v", time.Now().Unix()))
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				tempDir,
				"git@github.com:duck8823/duci.git",
				"refs/heads/master",
				plumbing.ZeroHash,
			)

			// then
			if err != nil {
				t.Errorf("must not error. %+v", err)
			}

			if _, err := os.Stat(path.Join(tempDir, ".git")); err != nil {
				t.Errorf("must create dir: %s", path.Join(tempDir, ".git"))
			}
		})

		t.Run("when target directory not exists", func(t *testing.T) {
			if os.Getuid() == 0 {
				t.Skip("skip if root user")
			}

			// given
			wrongPath := "/path/to/not/exists"

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				wrongPath,
				"git@github.com:duck8823/duci.git",
				"refs/heads/master",
				plumbing.ZeroHash,
			)

			// then
			if err == nil {
				t.Error("erro must occur")
			}
		})
	})
}
