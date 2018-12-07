package job

import (
	"bytes"
	"encoding/json"
	. "github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type dataSource struct {
	db LevelDB
}

func NewDataSource(path string) (*dataSource, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dataSource{db}, nil
}

func (d *dataSource) FindBy(id ID) (*Job, error) {
	data, err := d.db.Get(id.ToSlice(), nil)
	if err == leveldb.ErrNotFound {
		return nil, NotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	job := &Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	job.ID = id
	return job, nil
}

func (d *dataSource) Save(job Job) error {
	data, err := job.ToBytes()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.db.Put(job.ID.ToSlice(), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
