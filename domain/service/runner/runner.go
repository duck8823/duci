package runner

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/log"
	. "github.com/duck8823/duci/domain/service/docker"
	"github.com/pkg/errors"
)

// DockerRunner is a interface describes task runner.
type DockerRunner interface {
	Run(ctx context.Context, dir string, tag docker.Tag, cmd docker.Command) error
}

// dockerRunnerImpl is a implement of DockerRunner
type dockerRunnerImpl struct {
	Docker
	LogFunc []func(context.Context, log.Log)
}

// Run task in docker container
func (r *dockerRunnerImpl) Run(ctx context.Context, dir string, tag docker.Tag, cmd docker.Command) error {
	if err := r.dockerBuild(ctx, dir, tag); err != nil {
		return errors.WithStack(err)
	}

	conID, err := r.dockerRun(ctx, dir, tag, cmd)
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

	return nil
}

// dockerBuild build a docker image
func (r *dockerRunnerImpl) dockerBuild(ctx context.Context, dir string, tag docker.Tag) error {
	tarball, err := createTarball(dir)
	if err != nil {
		return errors.WithStack(err)
	}
	defer tarball.Close()

	buildLog, err := r.Docker.Build(ctx, tarball, docker.Tag(tag), dockerfilePath(dir))
	if err != nil {
		return errors.WithStack(err)
	}
	for _, f := range r.LogFunc {
		go f(ctx, buildLog)
	}
	return nil
}

// dockerRun run docker container
func (r *dockerRunnerImpl) dockerRun(ctx context.Context, dir string, tag docker.Tag, cmd docker.Command) (docker.ContainerID, error) {
	opts, err := runtimeOptions(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	conID, runLog, err := r.Docker.Run(ctx, opts, tag, cmd)
	if err != nil {
		return conID, errors.WithStack(err)
	}
	for _, f := range r.LogFunc {
		go f(ctx, runLog)
	}
	return conID, nil
}
