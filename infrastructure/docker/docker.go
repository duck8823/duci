package docker

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	moby "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"io"
	"strings"
)

// RuntimeOptions is a docker options.
type RuntimeOptions struct {
	Environments Environments
	Volumes      Volumes
}

// Environments represents a docker `-e` option.
type Environments map[string]interface{}

func (e Environments) ToArray() []string {
	var a []string
	for key, val := range e {
		a = append(a, fmt.Sprintf("%s=%v", key, val))
	}
	return a
}

// Environments represents a docker `-v` option.
type Volumes []string

func (v Volumes) ToMap() map[string]struct{} {
	m := make(map[string]struct{})
	for _, volume := range v {
		key := strings.Split(volume, ":")[0]
		m[key] = struct{}{}
	}
	return m
}

// Client is a interface of docker client
type Client interface {
	Build(ctx context.Context, file io.Reader, tag string, dockerfile string) (Log, error)
	Run(ctx context.Context, opts RuntimeOptions, tag string, cmd ...string) (string, Log, error)
	Rm(ctx context.Context, containerID string) error
	Rmi(ctx context.Context, tag string) error
	ExitCode(ctx context.Context, containerID string) (int64, error)
}

type clientImpl struct {
	moby *moby.Client
}

// New returns docker client.
func New() (Client, error) {
	cli, err := moby.NewEnvClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &clientImpl{moby: cli}, nil
}

// Build docker image.
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

// Run id a function create, start container.
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

// Rm remove docker container.
func (c *clientImpl) Rm(ctx context.Context, containerID string) error {
	if err := c.moby.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rmi remove docker image.
func (c *clientImpl) Rmi(ctx context.Context, tag string) error {
	if _, err := c.moby.ImageRemove(ctx, tag, types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode wait container until exit and returns exit code.
func (c *clientImpl) ExitCode(ctx context.Context, containerID string) (int64, error) {
	body, err := c.moby.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case b := <-body:
		return b.StatusCode, nil
	case e := <-err:
		return -1, errors.WithStack(e)
	}
}
