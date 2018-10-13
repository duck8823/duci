package docker

import (
	"context"
	"github.com/duck8823/duci/infrastructure/docker"
	"github.com/pkg/errors"
)

type Service interface {
	Status() error
}

type serviceImpl struct {
	moby docker.Client
}

func New(moby docker.Client) Service {
	return &serviceImpl{moby}
}

func (s *serviceImpl) Status() error {
	if _, err := s.moby.Info(context.Background()); err != nil {
		return errors.Wrap(err, "Couldn't connect to Docker daemon.")
	}
	return nil
}
