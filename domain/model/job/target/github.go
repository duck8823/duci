package target

import (
	"context"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
)

// GitHub is target with github repository
type GitHub struct {
	Repo  github.Repository
	Point github.TargetPoint
}

// Prepare working directory
func (g *GitHub) Prepare(ctx context.Context) (job.WorkDir, job.Cleanup, error) {
	tmpDir := filepath.Join(os.TempDir(), random.String(16, random.Alphanumeric, random.Numeric))
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", cleanupFunc(tmpDir), errors.WithStack(err)
	}

	git, err := git.GetInstance()
	if err != nil {
		return "", cleanupFunc(tmpDir), errors.WithStack(err)
	}

	if err := git.Clone(ctx, tmpDir, &github.TargetSource{
		Repository: g.Repo,
		Ref:        g.Point.GetRef(),
		SHA:        plumbing.NewHash(g.Point.GetHead()),
	}); err != nil {
		return "", cleanupFunc(tmpDir), errors.WithStack(err)
	}

	return job.WorkDir(tmpDir), cleanupFunc(tmpDir), nil
}
