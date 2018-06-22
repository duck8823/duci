package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/google/logger"
	moby "github.com/moby/moby/client"
	"github.com/pkg/errors"
	"io"
)

type Environments map[string]interface{}

func (e Environments) ToArray() []string {
	var a []string
	for key, val := range e {
		a = append(a, fmt.Sprintf("%s=%v", key, val))
	}
	return a
}

type TaskFailure error

type client struct {
	moby *moby.Client
}

func New() (*client, error) {
	cli, err := moby.NewEnvClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &client{moby: cli}, nil
}

func (c *client) Build(ctx context.Context, file io.Reader, tag string) error {
	resp, err := c.moby.ImageBuild(ctx, file, types.ImageBuildOptions{Tags: []string{tag}})
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 128)
	for {
		_, err := resp.Body.Read(buf)
		logger.Info(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (c *client) Run(ctx context.Context, env Environments, tag string, cmd ...string) (string, error) {
	con, err := c.moby.ContainerCreate(ctx, &container.Config{
		Image: tag,
		Env:   env.ToArray(),
		Cmd:   cmd,
	}, nil, nil, "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err := c.moby.ContainerStart(ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.WithStack(err)
	}

	log, err := c.moby.ContainerLogs(ctx, con.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	go func() {
		buf := make([]byte, 128)
		for {
			_, err := log.Read(buf)
			logger.Info(buf)
			if err == io.EOF {
				break
			}
		}
	}()

	if code, err := c.moby.ContainerWait(ctx, con.ID); err != nil {
		return "", errors.WithStack(err)
	} else if code != 0 {
		return con.ID, new(TaskFailure)
	}

	return con.ID, nil
}

func (c *client) Rm(ctx context.Context, containerId string) error {
	if err := c.moby.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *client) Rmi(ctx context.Context, tag string) error {
	if _, err := c.moby.ImageRemove(context.Background(), tag, types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
