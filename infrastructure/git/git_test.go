package git_test

import (
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/context"
	"github.com/duck8823/minimal-ci/infrastructure/git"
	"os"
	"path"
	"testing"
	"time"
)

func TestSshGitClient_Clone(t *testing.T) {
	tempDir := path.Join(os.TempDir(), fmt.Sprintf("minimal-ci_test_%v", time.Now().Unix()))
	if err := os.MkdirAll(path.Join(tempDir, "dir"), 0700); err != nil {
		t.Fatalf("%+v", err)
	}

	s, err := git.New()
	if err != nil {
		t.Fatalf("error occured. %+v", err)
	}

	if _, err := s.Clone(
		context.New("test/task"),
		tempDir,
		"git@github.com:duck8823/minimal-ci.git",
		"refs/heads/master",
	); err != nil {
		t.Errorf("must not error. %+v", err)
	}

	if _, err := os.Stat(path.Join(tempDir, ".git")); err != nil {
		t.Errorf("must be created dir: %s", path.Join(tempDir, ".git"))
	}

}
