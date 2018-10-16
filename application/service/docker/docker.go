package docker

import (
	"context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/pkg/errors"
	"io"
)

type Service interface {
	Build(ctx context.Context, file io.Reader, tag string, dockerfile string) (docker.Log, error)
	Run(ctx context.Context, opts docker.RuntimeOptions, tag string, cmd ...string) (string, docker.Log, error)
	Rm(ctx context.Context, containerID string) error
	Rmi(ctx context.Context, tag string) error
	ExitCode(ctx context.Context, containerID string) (int64, error)
	Status() error
}

type serviceImpl struct {
	moby docker.Client
}

func New(moby docker.Client) Service {
	return &serviceImpl{moby}
}

func (s *serviceImpl) Build(ctx context.Context, file io.Reader, tag string, dockerfile string) (docker.Log, error) {
	log, err := s.moby.Build(ctx, file, tag, dockerfile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return log, nil
}

func (s *serviceImpl) Run(ctx context.Context, opts docker.RuntimeOptions, tag string, cmd ...string) (string, docker.Log, error) {
	conID, log, err := s.moby.Run(ctx, opts, tag, cmd...)
	if err != nil {
		return conID, nil, errors.WithStack(err)
	}
	return conID, log, nil
}

func (s *serviceImpl) Rm(ctx context.Context, conID string) error {
	if err := s.moby.Rm(ctx, conID); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *serviceImpl) Rmi(ctx context.Context, tag string) error {
	if err := s.moby.Rmi(ctx, tag); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *serviceImpl) ExitCode(ctx context.Context, conID string) (int64, error) {
	code, err := s.moby.ExitCode(ctx, conID)
	if err != nil {
		return code, errors.WithStack(err)
	}
	return code, nil
}

func (s *serviceImpl) Status() error {
	if _, err := s.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
