package job_service

import (
	. "github.com/duck8823/duci/domain/model/job"
	. "github.com/duck8823/duci/infrastructure/job"
	"github.com/pkg/errors"
)

// Service represents job service
type Service interface {
	FindBy(id ID) (*Job, error)
	Start(id ID) error
	Append(id ID, line LogLine) error
	Finish(id ID) error
}

var instance *serviceImpl

type serviceImpl struct {
	repo Repository
}

// Initialize implementation of job service
func Initialize(path string) error {
	if instance != nil {
		return errors.New("instance already initialized.")
	}

	dataSource, err := NewDataSource(path)
	if err != nil {
		return errors.WithStack(err)
	}

	instance = &serviceImpl{repo: dataSource}
	return nil
}

// GetInstance returns job service
func GetInstance() (Service, error) {
	if instance == nil {
		return nil, errors.New("instance still not initialized.")
	}
	return instance, nil
}

// FindBy returns job is found by ID
func (s *serviceImpl) FindBy(id ID) (*Job, error) {
	job, err := s.repo.FindBy(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// Start store empty job
func (s *serviceImpl) Start(id ID) error {
	job := Job{ID: id, Finished: false}
	if err := s.repo.Save(job); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Append log to job
func (s *serviceImpl) Append(id ID, line LogLine) error {
	job, err := s.findOrInitialize(id)
	if err != nil {
		return errors.WithStack(err)
	}
	job.AppendLog(line)

	if err := s.repo.Save(*job); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *serviceImpl) findOrInitialize(id ID) (*Job, error) {
	job, err := s.repo.FindBy(id)
	if err == NotFound {
		return &Job{ID: id, Finished: false}, nil
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return job, nil
}

// Finish store finished job
func (s *serviceImpl) Finish(id ID) error {
	job, err := s.repo.FindBy(id)
	if err != nil {
		return errors.WithStack(err)
	}
	job.Finish()

	if err := s.repo.Save(*job); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
