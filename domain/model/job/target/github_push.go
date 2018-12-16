package target

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
)

// githubPush is target with github repository
type githubPush struct {
	git   git.Git
	Repo  github.Repository
	Point github.TargetPoint
}

// NewGitHubPush returns target for github push event
func NewGithubPush(repo github.Repository, point github.TargetPoint) (*githubPush, error) {
	cli, err := git.GetInstance()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &githubPush{
		git:   cli,
		Repo:  repo,
		Point: point,
	}, nil
}

// Prepare working directory
func (g *githubPush) Prepare() (job.WorkDir, job.Cleanup, error) {
	tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric, random.Numeric))
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", cleanupFunc(tmpDir), errors.WithStack(err)
	}

	if err := g.git.Clone(context.Background(), tmpDir, &github.TargetSource{
		Repository: g.Repo,
		Ref:        fmt.Sprintf("refs/heads/%s", g.Point.GetRef()),
		SHA:        plumbing.NewHash(g.Point.GetHead()),
	}); err != nil {
		return "", cleanupFunc(tmpDir), errors.WithStack(err)
	}

	return job.WorkDir(tmpDir), cleanupFunc(tmpDir), nil
}
