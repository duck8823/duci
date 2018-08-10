package runner

import (
	"bytes"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/git"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Runner interface {
	Run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error)
}

type DockerRunner struct {
	Git         git.Client
	GitHub      github.Service
	Docker      docker.Client
	Name        string
	BaseWorkDir string
}

func (r *DockerRunner) Run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error) {
	commitHash := make(chan plumbing.Hash, 1)
	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		hash, err := r.run(ctx, repo, ref, command...)
		commitHash <- hash
		errs <- err
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		hash := <-commitHash
		if timeout.Err() != nil {
			logger.Errorf(ctx.UUID(), "%+v", timeout.Err())
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.ERROR, timeout.Err().Error())
		}
		return hash, timeout.Err()
	case err := <-errs:
		hash := <-commitHash
		if err == docker.Failure {
			logger.Error(ctx.UUID(), err.Error())
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.FAILURE, "failure job")
		} else if err != nil {
			logger.Errorf(ctx.UUID(), "%+v", err)
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.ERROR, err.Error())
		} else {
			r.GitHub.CreateCommitStatus(ctx, repo, hash, github.SUCCESS, "success")
		}
		return hash, err
	}
}

func (r *DockerRunner) run(ctx context.Context, repo github.Repository, ref string, command ...string) (plumbing.Hash, error) {
	workDir := path.Join(r.BaseWorkDir, strconv.FormatInt(clock.Now().Unix(), 10))
	tagName := repo.GetFullName()

	head, err := r.Git.Clone(ctx, workDir, repo.GetSSHURL(), ref)
	if err != nil {
		return plumbing.Hash{}, errors.WithStack(err)
	}

	r.GitHub.CreateCommitStatus(ctx, repo, head, github.PENDING, "started job")

	tarFilePath := path.Join(workDir, "duci.tar")
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

	dockerfile := "./Dockerfile"
	if exists(path.Join(workDir, ".duci/Dockerfile")) {
		dockerfile = ".duci/Dockerfile"
	}
	if _, err := r.Docker.Build(ctx, readFile, tagName, dockerfile); err != nil {
		return head, errors.WithStack(err)
	}
	var opts docker.RuntimeOptions
	if exists(path.Join(workDir, ".duci/config.yml")) {
		content, err := ioutil.ReadFile(path.Join(workDir, ".duci/config.yml"))
		if err != nil {
			return head, errors.WithStack(err)
		}
		content = []byte(os.ExpandEnv(string(content)))
		if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&opts); err != nil {
			return head, errors.WithStack(err)
		}
	}

	_, _, err = r.Docker.Run(ctx, opts, tagName, command...)

	return head, err
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
