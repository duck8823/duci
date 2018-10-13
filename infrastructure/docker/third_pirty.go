package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"io"
)

// Moby is a interface of docker client
// see also github.com/moby/moby/client
type Moby interface {
	ImageBuild(
		ctx context.Context,
		buildContext io.Reader,
		options types.ImageBuildOptions,
	) (types.ImageBuildResponse, error)
	ContainerCreate(
		ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		containerName string,
	) (container.ContainerCreateCreatedBody, error)
	ContainerStart(
		ctx context.Context,
		containerID string,
		options types.ContainerStartOptions,
	) error
	ContainerLogs(
		ctx context.Context,
		container string,
		options types.ContainerLogsOptions,
	) (io.ReadCloser, error)
	ContainerRemove(
		ctx context.Context,
		containerID string,
		options types.ContainerRemoveOptions,
	) error
	ImageRemove(
		ctx context.Context,
		imageID string,
		options types.ImageRemoveOptions,
	) ([]types.ImageDeleteResponseItem, error)
	ContainerWait(
		ctx context.Context,
		containerID string,
		condition container.WaitCondition,
	) (<-chan container.ContainerWaitOKBody, <-chan error)
	Info(
		ctx context.Context,
	) (types.Info, error)
}
