package job

import (
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

func (d *dataSource) Get(id ID) (*Job, error) {
	data, err := d.db.Get(id.ToSlice(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	job, err := NewJob(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

func (d *dataSource) Start(id ID) error {
	data, _ := (&Job{Finished: false}).ToBytes()
	if err := d.db.Put(id.ToSlice(), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *dataSource) Append(id ID, line LogLine) error {
	job, err := d.findOrInitialize(id)
	if err != nil {
		return errors.WithStack(err)
	}

	job.AppendLog(line)

	data, err := job.ToBytes()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.db.Put(id.ToSlice(), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (d *dataSource) findOrInitialize(id ID) (*Job, error) {
	data, err := d.db.Get(id.ToSlice(), nil)
	if err == leveldb.ErrNotFound {
		return &Job{}, nil
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	job, err := NewJob(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return job, nil
}

func (d *dataSource) Finish(id ID) error {
	data, err := d.db.Get(id.ToSlice(), nil)
	if err != nil {
		return errors.WithStack(err)
	}

	job, err := NewJob(data)
	if err != nil {
		return errors.WithStack(err)
	}

	job.Finish()

	finished, err := job.ToBytes()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.db.Put(id.ToSlice(), finished, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
