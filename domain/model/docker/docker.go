package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	moby "github.com/docker/docker/client"
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

type client struct {
	moby Moby
}

// New returns instance of docker client
func New() (Docker, error) {
	cli, err := moby.NewClientWithOpts(moby.FromEnv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &client{moby: cli}, nil
}

// Build a docker image.
func (c *client) Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (Log, error) {
	opts := types.ImageBuildOptions{
		Tags:       []string{tag.ToString()},
		Dockerfile: dockerfile.ToString(),
		Remove:     true,
	}
	resp, err := c.moby.ImageBuild(ctx, file, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewBuildLog(resp.Body), nil
}

// Run docker container with command.
func (c *client) Run(ctx context.Context, opts RuntimeOptions, tag Tag, cmd Command) (ContainerID, Log, error) {
	con, err := c.moby.ContainerCreate(ctx, &container.Config{
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

	if err := c.moby.ContainerStart(ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return ContainerID(con.ID), nil, errors.WithStack(err)
	}

	logs, err := c.moby.ContainerLogs(ctx, con.ID, types.ContainerLogsOptions{
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
func (c *client) RemoveContainer(ctx context.Context, conID ContainerID) error {
	if err := c.moby.ContainerRemove(ctx, conID.ToString(), types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RemoveImage remove docker image.
func (c *client) RemoveImage(ctx context.Context, tag Tag) error {
	if _, err := c.moby.ImageRemove(ctx, tag.ToString(), types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode returns exit code specific container id.
func (c *client) ExitCode(ctx context.Context, conID ContainerID) (ExitCode, error) {
	body, err := c.moby.ContainerWait(ctx, conID.ToString(), container.WaitConditionNotRunning)
	select {
	case b := <-body:
		return ExitCode(b.StatusCode), nil
	case e := <-err:
		return -1, errors.WithStack(e)
	}
}

// Status returns error of docker daemon status.
func (c *client) Status() error {
	if _, err := c.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
