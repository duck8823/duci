package job

import (
	"github.com/duck8823/duci/domain/model/job"
	jobDataSource "github.com/duck8823/duci/infrastructure/job"
	"github.com/duck8823/duci/internal/container"
	"github.com/pkg/errors"
)

type serviceImpl struct {
	repo job.Repository
}

// Initialize implementation of job service
func Initialize(path string) error {
	dataSource, err := jobDataSource.NewDataSource(path)
	if err != nil {
		return errors.WithStack(err)
	}

	service := new(Service)
	*service = &serviceImpl{repo: dataSource}
	if err := container.Submit(service); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetInstance returns job service
func GetInstance() (Service, error) {
	ins := new(Service)
	if err := container.Get(ins); err != nil {
		return nil, errors.WithStack(err)
	}
	return *ins, nil
}

// FindBy returns job is found by ID
func (s *serviceImpl) FindBy(id job.ID) (*job.Job, error) {
	job, err := s.repo.FindBy(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// Start store empty job
func (s *serviceImpl) Start(id job.ID) error {
	job := job.Job{ID: id, Finished: false}
	if err := s.repo.Save(job); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Append log to job
func (s *serviceImpl) Append(id job.ID, line job.LogLine) error {
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

func (s *serviceImpl) findOrInitialize(id job.ID) (*job.Job, error) {
	j, err := s.repo.FindBy(id)
	if err == job.ErrNotFound {
		return &job.Job{ID: id, Finished: false}, nil
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	return j, nil
}

// Finish store finished job
func (s *serviceImpl) Finish(id job.ID) error {
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
