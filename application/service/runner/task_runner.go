package runner

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/semaphore"
	"github.com/duck8823/duci/application/service/docker"
	"github.com/pkg/errors"
)

// Tag is a docker tag.
type Tag string

// Command is docker command.
type Command []string

// RunOptions is options for docker.
type RunOptions struct {
	Tag
	Command
}

// TaskRunner is a interface describes task runner.
type TaskRunner interface {
	Run(ctx context.Context, dir string, opts RunOptions) error
}

// DockerTaskRunner is a implement of TaskRunner with Docker
type DockerTaskRunner struct {
	Docker    docker.Service
	StartFunc []func(context.Context)
	LogFunc   []func(context.Context, docker.Log)
	EndFunc   []func(context.Context, error)
}

// Run task in docker container.
func (r *DockerTaskRunner) Run(ctx context.Context, dir string, opts RunOptions) error {
	for _, f := range r.StartFunc {
		go f(ctx)
	}

	errs := make(chan error, 1)

	timeout, cancel := context.WithTimeout(ctx, application.Config.Timeout())
	defer cancel()

	go func() {
		semaphore.Acquire()
		errs <- r.run(timeout, dir, opts)
		semaphore.Release()
	}()

	select {
	case <-timeout.Done():
		for _, f := range r.EndFunc {
			go f(ctx, timeout.Err())
		}
		return timeout.Err()
	case err := <-errs:
		for _, f := range r.EndFunc {
			go f(ctx, err)
		}
		return err
	}
}

func (r *DockerTaskRunner) run(ctx context.Context, dir string, opts RunOptions) error {
	if err := r.dockerBuild(ctx, dir, opts.Tag); err != nil {
		return errors.WithStack(err)
	}

	conID, err := r.dockerRun(ctx, dir, opts)
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

func (r *DockerTaskRunner) dockerBuild(ctx context.Context, dir string, tag Tag) error {
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

func (r *DockerTaskRunner) dockerRun(ctx context.Context, dir string, opts RunOptions) (docker.ContainerID, error) {
	dockerOpts, err := runtimeOpts(dir)
	if err != nil {
		return "", errors.WithStack(err)
	}

	conID, runLog, err := r.Docker.Run(ctx, dockerOpts, docker.Tag(opts.Tag), docker.Command(opts.Command))
	if err != nil {
		return conID, errors.WithStack(err)
	}
	for _, f := range r.LogFunc {
		go f(ctx, runLog)
	}
	return conID, nil
}
