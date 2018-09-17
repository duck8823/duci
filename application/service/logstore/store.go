package logstore

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/data/model"
	"github.com/duck8823/duci/infrastructure/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

// Level describes a log level.
type Level = string

// Service is a interface describe store for log.
type Service interface {
	Get(uuid uuid.UUID) (*model.Job, error)
	Append(uuid uuid.UUID, message model.Message) error
	Start(uuid uuid.UUID) error
	Finish(uuid uuid.UUID) error
	Close() error
}

type storeServiceImpl struct {
	db store.Store
}

// New returns a implementation of Service interface.
func New() (Service, error) {
	database, err := leveldb.OpenFile(application.Config.Server.DatabasePath, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &storeServiceImpl{database}, nil
}

// Append a message to store.
func (s *storeServiceImpl) Append(uuid uuid.UUID, message model.Message) error {
	job, err := s.findOrInitialize(uuid)
	if err != nil {
		return errors.WithStack(err)
	}

	job.Stream = append(job.Stream, message)

	data, err := json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := s.db.Put([]byte(uuid.String()), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *storeServiceImpl) findOrInitialize(uuid uuid.UUID) (*model.Job, error) {
	job := &model.Job{}

	data, err := s.db.Get([]byte(uuid.String()), nil)
	if err == store.NotFoundError {
		return job, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// Get a job from store.
func (s *storeServiceImpl) Get(uuid uuid.UUID) (*model.Job, error) {
	data, err := s.db.Get([]byte(uuid.String()), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	job := &model.Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

// Start stores initialized job to store.
func (s *storeServiceImpl) Start(uuid uuid.UUID) error {
	started, _ := json.Marshal(&model.Job{Finished: false})
	if err := s.db.Put([]byte(uuid.String()), started, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Finish stores with finished flag.
func (s *storeServiceImpl) Finish(uuid uuid.UUID) error {
	data, err := s.db.Get([]byte(uuid.String()), nil)
	if err != nil {
		return errors.WithStack(err)
	}

	job := &model.Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return errors.WithStack(err)
	}

	job.Finished = true

	finished, err := json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := s.db.Put([]byte(uuid.String()), finished, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Close a data store.
func (s *storeServiceImpl) Close() error {
	if err := s.db.Close(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
