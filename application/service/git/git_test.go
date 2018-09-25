package git_test

import (
	"fmt"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/service/git"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"net/url"
	"os"
	"path"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("when without ssh key path", func(t *testing.T) {
		// expect
		if _, err := git.New(); err != nil {
			t.Error("error must not occur")
		}
	})

	t.Run("when missing ssh key path", func(t *testing.T) {
		// given
		application.Config.GitHub.SSHKeyPath = "/path/to/wrong"

		// expect
		if _, err := git.New(); err == nil {
			t.Error("error must occur")
		}
	})
}

func TestSshGitService_Clone(t *testing.T) {
	// setup
	application.Config.GitHub.SSHKeyPath = path.Join(os.Getenv("HOME"), ".ssh/id_rsa")

	t.Run("with correct key", func(t *testing.T) {
		// setup
		client, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		t.Run("when target directory exists", func(t *testing.T) {
			// setup
			dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
			tempDir := path.Join(os.TempDir(), dirStr)
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				tempDir,
				&MockTargetSource{
					URL: "git@github.com:duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.ZeroHash,
				},
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
				&MockTargetSource{
					URL: "git@github.com:duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.ZeroHash,
				},
			)

			// then
			if err == nil {
				t.Error("erro must occur")
			}
		})

		t.Run("with wrong sha", func(t *testing.T) {
			// setup
			dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
			tempDir := path.Join(os.TempDir(), dirStr)
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				tempDir,
				&MockTargetSource{
					URL: "git@github.com:duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.NewHash(uuid.New().String()),
				},
			)

			// then
			if err == nil {
				t.Error("error must occur")
			}
		})
	})
}

func TestHttpGitService_Clone(t *testing.T) {
	// setup
	application.Config.GitHub.SSHKeyPath = ""

	t.Run("with correct key", func(t *testing.T) {
		// setup
		client, err := git.New()
		if err != nil {
			t.Fatalf("error occurred. %+v", err)
		}

		t.Run("when target directory exists", func(t *testing.T) {
			// setup
			dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
			tempDir := path.Join(os.TempDir(), dirStr)
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				tempDir,
				&MockTargetSource{
					URL: "https://github.com/duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.ZeroHash,
				},
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
				&MockTargetSource{
					URL: "https://github.com/duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.ZeroHash,
				},
			)

			// then
			if err == nil {
				t.Error("erro must occur")
			}
		})

		t.Run("with wrong sha", func(t *testing.T) {
			// setup
			dirStr := fmt.Sprintf("duci_test_%s", random.String(16, random.Alphanumeric))
			tempDir := path.Join(os.TempDir(), dirStr)
			if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
				t.Fatalf("%+v", err)
			}

			// when
			err := client.Clone(
				context.New("test/task", uuid.New(), &url.URL{}),
				tempDir,
				&MockTargetSource{
					URL: "git@github.com:duck8823/duci.git",
					Ref: "refs/heads/master",
					SHA: plumbing.NewHash(uuid.New().String()),
				},
			)

			// then
			if err == nil {
				t.Error("error must occur")
			}
		})
	})
}

type MockTargetSource struct {
	URL string
	Ref string
	SHA plumbing.Hash
}

func (t *MockTargetSource) GetURL() string {
	return t.URL
}

func (t *MockTargetSource) GetRef() string {
	return t.Ref
}

func (t *MockTargetSource) GetSHA() plumbing.Hash {
	return t.SHA
}
