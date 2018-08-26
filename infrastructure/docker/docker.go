package docker

import (
	"bufio"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	moby "github.com/docker/docker/client"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/pkg/errors"
	"io"
	"strings"
)

type RuntimeOptions struct {
	Environments Environments
	Volumes      Volumes
}

type Environments map[string]interface{}

func (e Environments) ToArray() []string {
	var a []string
	for key, val := range e {
		a = append(a, fmt.Sprintf("%s=%v", key, val))
	}
	return a
}

type Volumes []string

func (v Volumes) ToMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, volume := range v {
		key := strings.Split(volume, ":")[0]
		m[key] = struct{}{}
	}
	return m
}

type Client interface {
	Build(ctx context.Context, file io.Reader, tag string, dockerfile string) (Log, error)
	Run(ctx context.Context, opts RuntimeOptions, tag string, cmd ...string) (string, Log, error)
	Rm(ctx context.Context, containerId string) error
	Rmi(ctx context.Context, tag string) error
	ExitCode(ctx context.Context, containerId string) (int64, error)
}

type clientImpl struct {
	moby *moby.Client
}

func New() (Client, error) {
	cli, err := moby.NewEnvClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &clientImpl{moby: cli}, nil
}

func (c *clientImpl) Build(ctx context.Context, file io.Reader, tag string, dockerfile string) (Log, error) {
	opts := types.ImageBuildOptions{
		Tags:       []string{tag},
		Dockerfile: dockerfile,
		Remove:     true,
	}
	resp, err := c.moby.ImageBuild(ctx, file, opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &buildLogger{bufio.NewReader(resp.Body)}, nil
}

func (c *clientImpl) Run(ctx context.Context, opts RuntimeOptions, tag string, cmd ...string) (string, Log, error) {
	con, err := c.moby.ContainerCreate(ctx, &container.Config{
		Image:   tag,
		Env:     opts.Environments.ToArray(),
		Volumes: opts.Volumes.ToMap(),
		Cmd:     cmd,
	}, &container.HostConfig{
		Binds: opts.Volumes,
	}, nil, "")
	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	if err := c.moby.ContainerStart(ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return con.ID, nil, errors.WithStack(err)
	}

	log, err := c.moby.ContainerLogs(ctx, con.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return con.ID, nil, errors.WithStack(err)
	}

	return con.ID, &runLogger{bufio.NewReader(log)}, nil
}

func (c *clientImpl) Rm(ctx context.Context, containerId string) error {
	if err := c.moby.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *clientImpl) Rmi(ctx context.Context, tag string) error {
	if _, err := c.moby.ImageRemove(ctx, tag, types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *clientImpl) ExitCode(ctx context.Context, containerId string) (int64, error) {
	body, err := c.moby.ContainerWait(ctx, containerId, container.WaitConditionNotRunning)
	select {
	case b := <-body:
		return b.StatusCode, nil
	case e := <-err:
		return -1, errors.WithStack(e)
	}
}
