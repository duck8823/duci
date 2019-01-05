package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	moby "github.com/docker/docker/client"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
	"io"
)

type dockerImpl struct {
	moby Moby
}

// New returns instance of docker dockerImpl
func New() (Docker, error) {
	cli, err := moby.NewClientWithOpts(moby.FromEnv)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dockerImpl{moby: cli}, nil
}

// Build a docker image.
func (c *dockerImpl) Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (job.Log, error) {
	opts := types.ImageBuildOptions{
		Tags:       []string{tag.String()},
		Dockerfile: dockerfile.String(),
		Remove:     true,
	}
	resp, err := c.moby.ImageBuild(ctx, file, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewBuildLog(resp.Body), nil
}

// Run docker container with command.
func (c *dockerImpl) Run(ctx context.Context, opts RuntimeOptions, tag Tag, cmd Command) (ContainerID, job.Log, error) {
	con, err := c.moby.ContainerCreate(ctx, &container.Config{
		Image:   tag.String(),
		Env:     opts.Environments.Array(),
		Volumes: opts.Volumes.Map(),
		Cmd:     cmd.Slice(),
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
func (c *dockerImpl) RemoveContainer(ctx context.Context, conID ContainerID) error {
	if err := c.moby.ContainerRemove(ctx, conID.String(), types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// RemoveImage remove docker image.
func (c *dockerImpl) RemoveImage(ctx context.Context, tag Tag) error {
	if _, err := c.moby.ImageRemove(ctx, tag.String(), types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode returns exit code specific container id.
func (c *dockerImpl) ExitCode(ctx context.Context, conID ContainerID) (ExitCode, error) {
	body, err := c.moby.ContainerWait(ctx, conID.String(), container.WaitConditionNotRunning)
	select {
	case b := <-body:
		return ExitCode(b.StatusCode), nil
	case e := <-err:
		return -1, errors.WithStack(e)
	}
}

// Status returns error of docker daemon status.
func (c *dockerImpl) Status() error {
	if _, err := c.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
