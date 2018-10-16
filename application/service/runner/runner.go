package runner

import (
	"bytes"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/context"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/docker"
	"github.com/duck8823/duci/application/service/git"
	"github.com/duck8823/duci/application/service/github"
	"github.com/duck8823/duci/application/service/logstore"
	"github.com/duck8823/duci/data/model"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// ErrFailure is a error describes task failure.
var ErrFailure = errors.New("Task Failure")

// Runner is a interface describes task runner.
type Runner interface {
	Run(ctx context.Context, src *github.TargetSource, command ...string) error
}

// DockerRunner represents a runner implement for docker.
type DockerRunner struct {
	Git         git.Service
	GitHub      github.Service
	Docker      docker.Service
	LogStore    logstore.Service
	BaseWorkDir string
}

// Run task in docker container.
func (r *DockerRunner) Run(ctx context.Context, src *github.TargetSource, command ...string) error {
	if err := r.LogStore.Start(ctx.UUID()); err != nil {
		r.GitHub.CreateCommitStatus(ctx, src, github.ERROR, err.Error())
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

func (r *DockerRunner) run(ctx context.Context, src *github.TargetSource, command ...string) error {
	workDir := path.Join(r.BaseWorkDir, random.String(36, random.Alphanumeric))

	if err := r.Git.Clone(ctx, workDir, src); err != nil {
		return errors.WithStack(err)
	}

	r.GitHub.CreateCommitStatus(ctx, src, github.PENDING, "started job")

	if err := r.dockerBuild(ctx, workDir, src.Repo); err != nil {
		return errors.WithStack(err)
	}

	conID, err := r.dockerRun(ctx, workDir, src.Repo, command)
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
		return ErrFailure
	}

	return err
}

func (r *DockerRunner) dockerBuild(ctx context.Context, dir string, repo github.Repository) error {
	tarball, err := createTarball(dir)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tarball.Close()

	tag := docker.Tag(repo.GetFullName())
	buildLog, err := r.Docker.Build(ctx, tarball, tag, dockerfilePath(dir))
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

func dockerfilePath(workDir string) docker.Dockerfile {
	dockerfile := "./Dockerfile"
	if exists(path.Join(workDir, ".duci/Dockerfile")) {
		dockerfile = ".duci/Dockerfile"
	}
	return docker.Dockerfile(dockerfile)
}

func (r *DockerRunner) dockerRun(ctx context.Context, dir string, repo github.Repository, cmd docker.Command) (docker.ContainerID, error) {
	opts, err := runtimeOpts(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	tag := docker.Tag(repo.GetFullName())
	conID, runLog, err := r.Docker.Run(ctx, opts, tag, cmd)
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

func (r *DockerRunner) timeout(ctx context.Context, src *github.TargetSource) {
	if ctx.Err() != nil {
		logger.Errorf(ctx.UUID(), "%+v", ctx.Err())
		r.GitHub.CreateCommitStatus(ctx, src, github.ERROR, ctx.Err().Error())
	}
	r.LogStore.Finish(ctx.UUID())
}

func (r *DockerRunner) finish(ctx context.Context, src *github.TargetSource, err error) {
	if err == ErrFailure {
		logger.Error(ctx.UUID(), err.Error())
		r.GitHub.CreateCommitStatus(ctx, src, github.FAILURE, "failure job")
	} else if err != nil {
		logger.Errorf(ctx.UUID(), "%+v", err)
		r.GitHub.CreateCommitStatus(ctx, src, github.ERROR, err.Error())
	} else {
		r.GitHub.CreateCommitStatus(ctx, src, github.SUCCESS, "success")
	}
	r.LogStore.Finish(ctx.UUID())
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
