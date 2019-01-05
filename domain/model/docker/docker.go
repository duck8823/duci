package docker

import (
	"context"
	"github.com/duck8823/duci/domain/model/job"
	"io"
)

// Docker is a interface describe docker service.
type Docker interface {
	Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (job.Log, error)
	Run(ctx context.Context, opts RuntimeOptions, tag Tag, cmd Command) (ContainerID, job.Log, error)
	RemoveContainer(ctx context.Context, containerID ContainerID) error
	RemoveImage(ctx context.Context, tag Tag) error
	ExitCode(ctx context.Context, containerID ContainerID) (ExitCode, error)
	Status() error
}
