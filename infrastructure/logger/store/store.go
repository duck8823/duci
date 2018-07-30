package logger_store

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

type Job struct {
	Finished bool      `json:"finished"`
	Stream   []Message `json:"stream"`
}

type Message struct {
	Level string `json:"level"`
	Time  string `json:"time"`
	Text  string `json:"message"`
}

var (
	db       *leveldb.DB = nil
	NotFound             = errors.New("leveldb: not found")
)

func Open(path string) error {
	database, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	db = database

	return nil
}

func Append(uuid uuid.UUID, level, message string) error {
	if db == nil {
		return errors.New("database not open")
	}
	data, err := db.Get([]byte(uuid.String()), nil)
	if err != nil && err.Error() != NotFound.Error() {
		return errors.WithStack(err)
	}
	job := &Job{}
	if data != nil {
		json.NewDecoder(bytes.NewReader(data)).Decode(job)
	}

	msg := Message{
		Level: level,
		Time:  time.Now().Format("2006-01-02 15:04:05.000"),
		Text:  message,
	}
	job.Stream = append(job.Stream, msg)

	data, err = json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := db.Put([]byte(uuid.String()), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Get(uuid uuid.UUID) (*Job, error) {
	if db == nil {
		return nil, errors.New("database not open")
	}
	data, err := db.Get([]byte(uuid.String()), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	job := &Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

func Finish(uuid uuid.UUID) error {
	if db == nil {
		return errors.New("database not open")
	}
	data, _ := db.Get([]byte(uuid.String()), nil)
	job := &Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return errors.WithStack(err)
	}

	job.Finished = true

	data, err := json.Marshal(job)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := db.Put([]byte(uuid.String()), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Close() error {
	if db == nil {
		return errors.New("database not open")
	}
	if err := db.Close(); err != nil {
		return errors.WithStack(err)
	}
	db = nil
	return nil
}
