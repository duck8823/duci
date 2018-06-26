package runner

import (
	"context"
	"fmt"
	"github.com/duck8823/minimal-ci/infrastructure/archive/tar"
	"github.com/duck8823/minimal-ci/infrastructure/docker"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
	"github.com/duck8823/minimal-ci/service/github"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path"
	"strconv"
	"time"
)

type Runner interface {
	RunWithPullRequest(ctx context.Context, repo github.Repository, num int, command ...string) error
	Run(ctx context.Context, repo github.Repository, ref string, command ...string) error
}

const NAME = "minimal-ci"

type runnerImpl struct {
	GitHub      github.Service
	Docker      *docker.Client
	Name        string
	BaseWorkDir string
}

func NewWithEnv() (*runnerImpl, error) {
	githubService := github.New(context.Background(), os.Getenv("GITHUB_API_TOKEN"))
	dockerClient, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &runnerImpl{
		GitHub:      githubService,
		Docker:      dockerClient,
		Name:        NAME,
		BaseWorkDir: path.Join(os.TempDir(), NAME),
	}, nil
}

func (r *runnerImpl) RunWithPullRequest(ctx context.Context, repo github.Repository, num int, command ...string) error {
	pr, err := r.GitHub.GetPullRequest(ctx, repo, num)
	if err != nil {
		return errors.WithStack(err)
	}
	ref := fmt.Sprintf("refs/heads/%s", pr.GetHead().GetRef())
	return r.Run(ctx, repo, ref, command...)
}

func (r *runnerImpl) Run(ctx context.Context, repo github.Repository, ref string, command ...string) error {
	workDir := path.Join(r.BaseWorkDir, strconv.FormatInt(time.Now().Unix(), 10))
	tagName := repo.GetFullName()

	head, err := r.GitHub.Clone(ctx, workDir, repo, ref)
	if err != nil {
		return errors.WithStack(err)
	}

	r.CreateCommitStatus(ctx, repo, head, github.PENDING)

	tarFilePath := path.Join(workDir, "minimal-ci.tar")
	writeFile, err := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE, 0400)
	if err != nil {
		r.CreateCommitStatusWithError(ctx, repo, head, err)
		return errors.WithStack(err)
	}
	defer writeFile.Close()

	if err := tar.Create(workDir, writeFile); err != nil {
		r.CreateCommitStatusWithError(ctx, repo, head, err)
		return errors.WithStack(err)
	}

	readFile, _ := os.Open(tarFilePath)
	defer readFile.Close()

	if err := r.Docker.Build(ctx, readFile, tagName); err != nil {
		r.CreateCommitStatusWithError(ctx, repo, head, err)
		return errors.WithStack(err)
	}

	_, err = r.Docker.Run(ctx, docker.Environments{}, tagName, command...)
	if err == docker.Failure {
		r.CreateCommitStatus(ctx, repo, head, github.FAILURE)
		return errors.WithStack(err)
	} else if err != nil {
		r.CreateCommitStatusWithError(ctx, repo, head, err)
		return errors.WithStack(err)
	}

	r.CreateCommitStatus(ctx, repo, head, github.SUCCESS)

	return nil
}

func (r *runnerImpl) CreateCommitStatus(ctx context.Context, repo github.Repository, hash plumbing.Hash, state github.State) {
	msg := fmt.Sprintf("task %s", state)
	if err := r.GitHub.CreateCommitStatus(ctx, repo, hash, &github.Status{
		Context:     &r.Name,
		Description: &msg,
		State:       &state,
	}); err != nil {
		logger.Errorf("Failed to create commit status: %+v", err)
	}
}

func (r *runnerImpl) CreateCommitStatusWithError(ctx context.Context, repo github.Repository, hash plumbing.Hash, err error) {
	msg := err.Error()
	if len(msg) >= 50 {
		msg = string([]rune(msg)[:46]) + "..."
	}
	state := github.ERROR
	if err := r.GitHub.CreateCommitStatus(ctx, repo, hash, &github.Status{
		Context:     &r.Name,
		Description: &msg,
		State:       &state,
	}); err != nil {
		logger.Errorf("Failed to create commit status: %+v", err)
	}
}
