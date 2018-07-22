package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/duck8823/duci/infrastructure/context"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	moby "github.com/moby/moby/client"
	"github.com/pkg/errors"
	"io"
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
		m[volume] = struct{}{}
	}
	return m
}

var Failure = errors.New("Task Failure")

type Client interface {
	Build(ctx context.Context, file io.Reader, tag string, dockerfile string) error
	Run(ctx context.Context, opts RuntimeOptions, tag string, cmd ...string) (string, error)
	Rm(ctx context.Context, containerId string) error
	Rmi(ctx context.Context, tag string) error
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

func (c *clientImpl) Build(ctx context.Context, file io.Reader, tag string, dockerfile string) error {
	opts := types.ImageBuildOptions{
		Tags:       []string{tag},
		Dockerfile: dockerfile,
	}
	resp, err := c.moby.ImageBuild(ctx, file, opts)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	logStream(ctx.UUID(), resp.Body)
	return nil
}

func (c *clientImpl) Run(ctx context.Context, opts RuntimeOptions, tag string, cmd ...string) (string, error) {
	con, err := c.moby.ContainerCreate(ctx, &container.Config{
		Image:   tag,
		Env:     opts.Environments.ToArray(),
		Volumes: opts.Volumes.ToMap(),
		Cmd:     cmd,
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
	defer log.Close()

	go func() {
		reader := bufio.NewReaderSize(log, 1024)
		for {
			line, _, err := reader.ReadLine()
			if len(line) > 8 {
				// detect log prefix
				// see https://godoc.org/github.com/docker/docker/client#Client.ContainerLogs
				if !((line[0] == 1 || line[0] == 2) && (line[1] == 0 && line[2] == 0 && line[3] == 0)) {
					continue
				}
				messages := line[8:]

				// prevent to CR
				progress := bytes.Split(messages, []byte{'\r'})
				logger.Info(ctx.UUID(), string(progress[0]))
			}
			if err == io.EOF {
				break
			}
		}
	}()

	if code, err := c.moby.ContainerWait(ctx, con.ID); err != nil {
		return "", errors.WithStack(err)
	} else if code != 0 {
		return con.ID, Failure
	}

	return con.ID, nil
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

func logStream(uuid uuid.UUID, log io.Reader) error {
	reader := bufio.NewReaderSize(log, 1024)
	for {
		line, _, err := reader.ReadLine()
		stream := &struct {
			Stream string `json:"stream"`
		}{}
		json.Unmarshal(line, stream)
		if len(stream.Stream) > 0 {
			logger.Info(uuid, stream.Stream)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
