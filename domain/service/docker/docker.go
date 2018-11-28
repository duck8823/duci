package docker

import (
	"context"
	. "github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/log"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/pkg/errors"
	"io"
)

// Docker is a interface describe docker service.
type Docker interface {
	Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (log.Log, error)
	Run(ctx context.Context, opts docker.RuntimeOptions, tag Tag, cmd Command) (ContainerID, log.Log, error)
	Rm(ctx context.Context, containerID ContainerID) error
	Rmi(ctx context.Context, tag Tag) error
	ExitCode(ctx context.Context, containerID ContainerID) (ExitCode, error)
	Status() error
}

type dockerService struct {
	moby docker.Client
}

// New returns instance of docker service
func New() (Docker, error) {
	cli, err := docker.New()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dockerService{moby: cli}, nil
}

// Build a docker image.
func (s *dockerService) Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (log.Log, error) {
	buildLog, err := s.moby.Build(ctx, file, string(tag), string(dockerfile))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return buildLog.(log.Log), nil
}

// Run docker container with command.
func (s *dockerService) Run(ctx context.Context, opts docker.RuntimeOptions, tag Tag, cmd Command) (ContainerID, log.Log, error) {
	conID, runLog, err := s.moby.Run(ctx, opts, string(tag), cmd...)
	if err != nil {
		return ContainerID(conID), nil, errors.WithStack(err)
	}
	return ContainerID(conID), runLog.(log.Log), nil
}

// Rm remove docker container.
func (s *dockerService) Rm(ctx context.Context, conID ContainerID) error {
	if err := s.moby.Rm(ctx, string(conID)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rmi remove docker image.
func (s *dockerService) Rmi(ctx context.Context, tag Tag) error {
	if err := s.moby.Rmi(ctx, string(tag)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode returns exit code specific container id.
func (s *dockerService) ExitCode(ctx context.Context, conID ContainerID) (ExitCode, error) {
	code, err := s.moby.ExitCode(ctx, string(conID))
	if err != nil {
		return ExitCode(code), errors.WithStack(err)
	}
	return ExitCode(code), nil
}

// Status returns error of docker daemon status.
func (s *dockerService) Status() error {
	if _, err := s.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
