package runner

import (
	"bytes"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/logstore"
	"github.com/duck8823/duci/data/model"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var Failure = errors.New("Task Failure")

type Runner interface {
	Run(ctx context.Context, repo github.Repository, ref string, sha plumbing.Hash, command ...string) error
}

type DockerRunner struct {
	Git         git.Service
	GitHub      github.Service
	Docker      docker.Client
	LogStore    logstore.Service
	Name        string
	BaseWorkDir string
}

func (r *DockerRunner) Run(ctx context.Context, repo github.Repository, ref string, sha plumbing.Hash, command ...string) error {
	if err := r.LogStore.Start(ctx.UUID()); err != nil {
		r.GitHub.CreateCommitStatus(ctx, repo, sha, github.ERROR, err.Error())
		return errors.WithStack(err)
	}

	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		err := r.run(ctx, repo, ref, sha, command...)
		errs <- err
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		if timeout.Err() != nil {
			logger.Errorf(ctx.UUID(), "%+v", timeout.Err())
			r.GitHub.CreateCommitStatus(ctx, repo, sha, github.ERROR, timeout.Err().Error())
		}
		r.LogStore.Finish(ctx.UUID())
		return timeout.Err()
	case err := <-errs:
		if err == Failure {
			logger.Error(ctx.UUID(), err.Error())
			r.GitHub.CreateCommitStatus(ctx, repo, sha, github.FAILURE, "failure job")
		} else if err != nil {
			logger.Errorf(ctx.UUID(), "%+v", err)
			r.GitHub.CreateCommitStatus(ctx, repo, sha, github.ERROR, err.Error())
		} else {
			r.GitHub.CreateCommitStatus(ctx, repo, sha, github.SUCCESS, "success")
		}
		r.LogStore.Finish(ctx.UUID())
		return err
	}
}

func (r *DockerRunner) run(ctx context.Context, repo github.Repository, ref string, sha plumbing.Hash, command ...string) error {
	workDir := path.Join(r.BaseWorkDir, strconv.FormatInt(clock.Now().Unix(), 10))
	tagName := repo.GetFullName()

	_, err := r.Git.Clone(ctx, workDir, repo.GetSSHURL(), ref)
	if err != nil {
		return errors.WithStack(err)
	}

	r.GitHub.CreateCommitStatus(ctx, repo, sha, github.PENDING, "started job")

	tarFilePath := path.Join(workDir, "duci.tar")
	writeFile, err := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errors.WithStack(err)
	}
	defer writeFile.Close()

	if err := tar.Create(workDir, writeFile); err != nil {
		return errors.WithStack(err)
	}

	readFile, _ := os.Open(tarFilePath)
	defer readFile.Close()

	dockerfile := "./Dockerfile"
	if exists(path.Join(workDir, ".duci/Dockerfile")) {
		dockerfile = ".duci/Dockerfile"
	}
	buildLog, err := r.Docker.Build(ctx, readFile, tagName, dockerfile)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := r.logAppend(ctx, buildLog); err != nil {
		return errors.WithStack(err)
	}

	var opts docker.RuntimeOptions
	if exists(path.Join(workDir, ".duci/config.yml")) {
		content, err := ioutil.ReadFile(path.Join(workDir, ".duci/config.yml"))
		if err != nil {
			return errors.WithStack(err)
		}
		content = []byte(os.ExpandEnv(string(content)))
		if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&opts); err != nil {
			return errors.WithStack(err)
		}
	}

	containerId, runLog, err := r.Docker.Run(ctx, opts, tagName, command...)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := r.logAppend(ctx, runLog); err != nil {
		return errors.WithStack(err)
	}

	code, err := r.Docker.ExitCode(ctx, containerId)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := r.Docker.Rm(ctx, containerId); err != nil {
		return errors.WithStack(err)
	}
	if code != 0 {
		return Failure
	}

	return err
}

func (r *DockerRunner) logAppend(ctx context.Context, log docker.Log) error {
	for {
		line, err := log.ReadLine()
		if err != nil && err != io.EOF {
			logger.Debugf(ctx.UUID(), "skip read line with error: %s", err.Error())
			continue
		}
		logger.Info(ctx.UUID(), string(line.Message))
		if err := r.LogStore.Append(ctx.UUID(), model.Message{Time: line.Timestamp, Text: string(line.Message)}); err != nil {
			return errors.WithStack(err)
		}
		if err == io.EOF {
			return nil
		}
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
