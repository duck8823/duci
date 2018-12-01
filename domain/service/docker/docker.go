package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	moby "github.com/docker/docker/client"
	. "github.com/duck8823/duci/domain/model/docker"
	. "github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
	"io"
)

// Docker is a interface describe docker service.
type Docker interface {
	Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (Log, error)
	Run(ctx context.Context, opts RuntimeOptions, tag Tag, cmd Command) (ContainerID, Log, error)
	RemoveContainer(ctx context.Context, containerID ContainerID) error
	RemoveImage(ctx context.Context, tag Tag) error
	ExitCode(ctx context.Context, containerID ContainerID) (ExitCode, error)
	Status() error
}

type dockerService struct {
	moby Moby
}

// New returns instance of docker service
func New() (Docker, error) {
	cli, err := moby.NewClientWithOpts(moby.FromEnv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dockerService{moby: cli}, nil
}

// Build a docker image.
func (s *dockerService) Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (Log, error) {
	opts := types.ImageBuildOptions{
		Tags:       []string{tag.ToString()},
		Dockerfile: dockerfile.ToString(),
		Remove:     true,
	}
	resp, err := s.moby.ImageBuild(ctx, file, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewBuildLog(resp.Body), nil
}

// Run docker container with command.
func (s *dockerService) Run(ctx context.Context, opts RuntimeOptions, tag Tag, cmd Command) (ContainerID, Log, error) {
	con, err := s.moby.ContainerCreate(ctx, &container.Config{
		Image:   tag.ToString(),
		Env:     opts.Environments.ToArray(),
		Volumes: opts.Volumes.ToMap(),
		Cmd:     cmd.ToSlice(),
	}, &container.HostConfig{
		Binds: opts.Volumes,
	}, nil, "")
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if err := s.moby.ContainerStart(ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return ContainerID(con.ID), nil, errors.WithStack(err)
	}

	logs, err := s.moby.ContainerLogs(ctx, con.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return ContainerID(con.ID), nil, errors.WithStack(err)
	}

	return ContainerID(con.ID), NewRunLog(logs), nil
}

// RemoveContainer remove docker container.
func (s *dockerService) RemoveContainer(ctx context.Context, conID ContainerID) error {
	if err := s.moby.ContainerRemove(ctx, conID.ToString(), types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RemoveImage remove docker image.
func (s *dockerService) RemoveImage(ctx context.Context, tag Tag) error {
	if _, err := s.moby.ImageRemove(ctx, tag.ToString(), types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode returns exit code specific container id.
func (s *dockerService) ExitCode(ctx context.Context, conID ContainerID) (ExitCode, error) {
	body, err := s.moby.ContainerWait(ctx, conID.ToString(), container.WaitConditionNotRunning)
	select {
	case b := <-body:
		return ExitCode(b.StatusCode), nil
	case e := <-err:
		return -1, errors.WithStack(e)
	}
}

// Status returns error of docker daemon status.
func (s *dockerService) Status() error {
	if _, err := s.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
