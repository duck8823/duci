package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/duck8823/minimal-ci/infrastructure/logger"
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

var Failure = errors.New("Failure")

type Client struct {
	Moby *moby.Client
}

func New() (*Client, error) {
	cli, err := moby.NewEnvClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Client{Moby: cli}, nil
}

func (c *Client) Build(ctx context.Context, file io.Reader, tag string) error {
	resp, err := c.Moby.ImageBuild(ctx, file, types.ImageBuildOptions{Tags: []string{tag}})
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	logStream(resp.Body)
	return nil
}

func (c *Client) Run(ctx context.Context, env Environments, tag string, cmd ...string) (string, error) {
	con, err := c.Moby.ContainerCreate(ctx, &container.Config{
		Image: tag,
		Env:   env.ToArray(),
		Cmd:   cmd,
	}, nil, nil, "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err := c.Moby.ContainerStart(ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return "", errors.WithStack(err)
	}

	log, err := c.Moby.ContainerLogs(ctx, con.ID, types.ContainerLogsOptions{
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
				logger.Info(string(progress[0]))
			}
			if err == io.EOF {
				break
			}
		}
	}()

	if code, err := c.Moby.ContainerWait(ctx, con.ID); err != nil {
		return "", errors.WithStack(err)
	} else if code != 0 {
		return con.ID, Failure
	}

	return con.ID, nil
}

func (c *Client) Rm(ctx context.Context, containerId string) error {
	if err := c.Moby.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) Rmi(ctx context.Context, tag string) error {
	if _, err := c.Moby.ImageRemove(context.Background(), tag, types.ImageRemoveOptions{}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func logStream(log io.Reader) error {
	reader := bufio.NewReaderSize(log, 1024)
	for {
		line, _, err := reader.ReadLine()
		stream := &struct {
			Stream string `json:"stream"`
		}{}
		json.Unmarshal(line, stream)
		if len(stream.Stream) > 0 {
			logger.Info(stream.Stream)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
