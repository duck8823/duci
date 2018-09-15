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
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
)

var Failure = errors.New("Task Failure")

type TargetSource struct {
	Repo github.Repository
	Ref  string
	SHA  plumbing.Hash
}

type Runner interface {
	Run(ctx context.Context, src TargetSource, command ...string) error
}

type DockerRunner struct {
	Git         git.Service
	GitHub      github.Service
	Docker      docker.Client
	LogStore    logstore.Service
	BaseWorkDir string
}

func (r *DockerRunner) Run(ctx context.Context, src TargetSource, command ...string) error {
	if err := r.LogStore.Start(ctx.UUID()); err != nil {
		r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.ERROR, err.Error())
		return errors.WithStack(err)
	}

	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		errs <- r.run(timeout, src, command...)
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		r.timeout(timeout, src)
		return timeout.Err()
	case err := <-errs:
		r.finish(ctx, src, err)
		return err
	}
}

func (r *DockerRunner) run(ctx context.Context, src TargetSource, command ...string) error {
	workDir := path.Join(r.BaseWorkDir, random.String(36, random.Alphanumeric))

	if err := r.Git.Clone(ctx, workDir, src.Repo.GetSSHURL(), src.Ref, src.SHA); err != nil {
		return errors.WithStack(err)
	}

	r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.PENDING, "started job")

	if err := r.dockerBuild(ctx, workDir, src.Repo); err != nil {
		return errors.WithStack(err)
	}

	conID, err := r.dockerRun(ctx, workDir, src.Repo, command...)
	if err != nil {
		return errors.WithStack(err)
	}

	code, err := r.Docker.ExitCode(ctx, conID)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := r.Docker.Rm(ctx, conID); err != nil {
		return errors.WithStack(err)
	}
	if code != 0 {
		return Failure
	}

	return err
}

func (r *DockerRunner) dockerBuild(ctx context.Context, dir string, repo github.Repository) error {
	tarball, err := createTarball(dir)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tarball.Close()

	buildLog, err := r.Docker.Build(ctx, tarball, repo.GetFullName(), dockerfilePath(dir))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := r.logAppend(ctx, buildLog); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func createTarball(workDir string) (*os.File, error) {
	tarFilePath := path.Join(workDir, "duci.tar")
	writeFile, err := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer writeFile.Close()

	if err := tar.Create(workDir, writeFile); err != nil {
		return nil, errors.WithStack(err)
	}

	readFile, _ := os.Open(tarFilePath)
	return readFile, nil
}

func dockerfilePath(workDir string) string {
	dockerfile := "./Dockerfile"
	if exists(path.Join(workDir, ".duci/Dockerfile")) {
		dockerfile = ".duci/Dockerfile"
	}
	return dockerfile
}

func (r *DockerRunner) dockerRun(ctx context.Context, dir string, repo github.Repository, cmd ...string) (string, error) {
	opts, err := runtimeOpts(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	conID, runLog, err := r.Docker.Run(ctx, opts, repo.GetFullName(), cmd...)
	if err != nil {
		return conID, errors.WithStack(err)
	}
	if err := r.logAppend(ctx, runLog); err != nil {
		return conID, errors.WithStack(err)
	}
	return conID, nil
}

func runtimeOpts(workDir string) (docker.RuntimeOptions, error) {
	var opts docker.RuntimeOptions

	if !exists(path.Join(workDir, ".duci/config.yml")) {
		return opts, nil
	}
	content, err := ioutil.ReadFile(path.Join(workDir, ".duci/config.yml"))
	if err != nil {
		return opts, errors.WithStack(err)
	}
	content = []byte(os.ExpandEnv(string(content)))
	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&opts); err != nil {
		return opts, errors.WithStack(err)
	}
	return opts, nil
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

func (r *DockerRunner) timeout(ctx context.Context, src TargetSource) {
	if ctx.Err() != nil {
		logger.Errorf(ctx.UUID(), "%+v", ctx.Err())
		r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.ERROR, ctx.Err().Error())
	}
	r.LogStore.Finish(ctx.UUID())
}

func (r *DockerRunner) finish(ctx context.Context, src TargetSource, err error) {
	if err == Failure {
		logger.Error(ctx.UUID(), err.Error())
		r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.FAILURE, "failure job")
	} else if err != nil {
		logger.Errorf(ctx.UUID(), "%+v", err)
		r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.ERROR, err.Error())
	} else {
		r.GitHub.CreateCommitStatus(ctx, src.Repo, src.SHA, github.SUCCESS, "success")
	}
	r.LogStore.Finish(ctx.UUID())
}
