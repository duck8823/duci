package runner

import (
	"github.com/duck8823/minimal-ci/infrastructure/archive/tar"
	"github.com/duck8823/minimal-ci/infrastructure/context"
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
	Run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error)
}

type DockerRunner struct {
	GitHub      github.Service
	Docker      *docker.Client
	Name        string
	BaseWorkDir string
}

func (r *DockerRunner) Run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error) {
	commitHash := make(chan plumbing.Hash, 1)
	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	go func() {
		hash, err := r.run(ctx, repo, ref, command...)
		commitHash <- hash
		errs <- err
	}()

	select {
	case <-timeout.Done():
		hash := <-commitHash
		if timeout.Err() != nil {
			logger.Errorf(ctx.UUID(), "%+v", timeout.Err())
			r.GitHub.CreateCommitStatusWithError(ctx, repo, <-commitHash, timeout.Err())
		}
		return hash, timeout.Err()
	case err := <-errs:
		hash := <-commitHash
		if err == docker.Failure {
			logger.Error(ctx.UUID(), err.Error())
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.FAILURE)
		} else if err != nil {
			logger.Errorf(ctx.UUID(), "%+v", err)
			r.GitHub.CreateCommitStatusWithError(ctx, repo, hash, err)
		} else {
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.SUCCESS)
		}
		return hash, err
	}
}

func (r *DockerRunner) run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error) {
	workDir := path.Join(r.BaseWorkDir, strconv.FormatInt(time.Now().Unix(), 10))
	tagName := repo.GetFullName()

	head, err := r.GitHub.Clone(ctx, workDir, repo, ref)
	if err != nil {
		return plumbing.Hash{}, errors.WithStack(err)
	}

	r.GitHub.CreateCommitStatus(ctx, repo, head, github.PENDING)

	tarFilePath := path.Join(workDir, "minimal-ci.tar")
	writeFile, err := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return head, errors.WithStack(err)
	}
	defer writeFile.Close()

	if err := tar.Create(workDir, writeFile); err != nil {
		return head, errors.WithStack(err)
	}

	readFile, _ := os.Open(tarFilePath)
	defer readFile.Close()

	if err := r.Docker.Build(ctx, readFile, tagName); err != nil {
		return head, errors.WithStack(err)
	}

	if _, err = r.Docker.Run(ctx, docker.Environments{}, tagName, command...); err != nil {
		return head, errors.WithStack(err)
	}

	return head, nil
}
