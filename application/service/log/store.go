package log

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/domain/model"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

type Level = string

type StoreService interface {
	Get(uuid uuid.UUID) (*model.Job, error)
	Append(uuid uuid.UUID, message model.Message) error
	Finish(uuid uuid.UUID) error
	Close() error
}

type storeServiceImpl struct {
	db logger.Store
}

func NewStoreService() (StoreService, error) {
	database, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &storeServiceImpl{database}, nil
}

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
	if err == logger.NotFoundError {
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

func (s *storeServiceImpl) Close() error {
	if err := s.db.Close(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
