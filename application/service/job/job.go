package job_service

import (
	. "github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
)

type Service interface {
	FindBy(id ID) (*Job, error)
}

type serviceImpl struct {
	repo Repository
}

func New() (Service, error) {
	return &serviceImpl{}, nil
}

func (s *serviceImpl) FindBy(id ID) (*Job, error) {
	job, err := s.repo.Get(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}
