package docker

import (
	"context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/pkg/errors"
	"io"
)

// Dockerfile represents a path to dockerfile
type Dockerfile string

// ContainerID describes a container id of docker
type ContainerID string

// Tag describes a docker tag
type Tag string

// ExitCode describes a exit code
type ExitCode int64

// Command describes a docker CMD
type Command []string

// Service is a interface describe docker service.
type Service interface {
	Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (docker.Log, error)
	Run(ctx context.Context, opts docker.RuntimeOptions, tag Tag, cmd Command) (ContainerID, docker.Log, error)
	Rm(ctx context.Context, containerID ContainerID) error
	Rmi(ctx context.Context, tag Tag) error
	ExitCode(ctx context.Context, containerID ContainerID) (ExitCode, error)
	Status() error
}

type serviceImpl struct {
	moby docker.Client
}

// New returns instance of docker service
func New(moby docker.Client) Service {
	return &serviceImpl{moby}
}

// Build a docker image.
func (s *serviceImpl) Build(ctx context.Context, file io.Reader, tag Tag, dockerfile Dockerfile) (docker.Log, error) {
	log, err := s.moby.Build(ctx, file, string(tag), string(dockerfile))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return log, nil
}

// Run docker container with command.
func (s *serviceImpl) Run(ctx context.Context, opts docker.RuntimeOptions, tag Tag, cmd Command) (ContainerID, docker.Log, error) {
	conID, log, err := s.moby.Run(ctx, opts, string(tag), cmd...)
	if err != nil {
		return ContainerID(conID), nil, errors.WithStack(err)
	}
	return ContainerID(conID), log, nil
}

// Rm remove docker container.
func (s *serviceImpl) Rm(ctx context.Context, conID ContainerID) error {
	if err := s.moby.Rm(ctx, string(conID)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rmi remove docker image.
func (s *serviceImpl) Rmi(ctx context.Context, tag Tag) error {
	if err := s.moby.Rmi(ctx, string(tag)); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ExitCode returns exit code specific container id.
func (s *serviceImpl) ExitCode(ctx context.Context, conID ContainerID) (ExitCode, error) {
	code, err := s.moby.ExitCode(ctx, string(conID))
	if err != nil {
		return ExitCode(code), errors.WithStack(err)
	}
	return ExitCode(code), nil
}

// Status returns error of docker daemon status.
func (s *serviceImpl) Status() error {
	if _, err := s.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
