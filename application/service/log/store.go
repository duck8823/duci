package log

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/domain/model"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"time"
)

type Level = string

type StoreService interface {
	Get(uuid uuid.UUID) (*model.Job, error)
	Append(uuid uuid.UUID, level Level, message string) error
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

func (s *storeServiceImpl) Append(uuid uuid.UUID, level, message string) error {
	data, err := s.db.Get([]byte(uuid.String()), nil)
	if err != nil && err != leveldb_errors.ErrNotFound {
		return errors.WithStack(err)
	}
	job := &model.Job{}
	if data != nil {
		json.NewDecoder(bytes.NewReader(data)).Decode(job)
	}

	msg := model.Message{
		Level: level,
		Time:  time.Now().Format("2006-01-02 15:04:05.000"),
		Text:  message,
	}
	job.Stream = append(job.Stream, msg)

	data, err = json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := s.db.Put([]byte(uuid.String()), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
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
	data, _ := s.db.Get([]byte(uuid.String()), nil)
	job := &model.Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return errors.WithStack(err)
	}

	job.Finished = true

	data, err := json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := s.db.Put([]byte(uuid.String()), data, nil); err != nil {
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
