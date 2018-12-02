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

// GitHubPR is target with github pull request
type GitHubPR struct {
	Repo github.Repository
	Num  int
}

// Prepare working directory
func (g *GitHubPR) Prepare() (job.WorkDir, job.Cleanup, error) {
	tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric, random.Numeric))
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", nil, errors.WithStack(err)
	}

	githubCli, err := github.GetInstance()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	pr, err := githubCli.GetPullRequest(context.Background(), g.Repo, g.Num)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	gitCli, err := git.GetInstance()
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if err := gitCli.Clone(context.Background(), tmpDir, &github.TargetSource{
		Repository: g.Repo,
		Ref:        fmt.Sprintf("refs/heads/%s", pr.GetHead().GetRef()),
		SHA:        plumbing.NewHash(pr.GetHead().GetSHA()),
	}); err != nil {
		return "", nil, errors.WithStack(err)
	}

	return job.WorkDir(tmpDir), cleanupFunc(tmpDir), nil
}
