package logger

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

var (
	NotFoundError = leveldb_errors.ErrNotFound
)

type ReadOptions = opt.ReadOptions
type WriteOptions = opt.WriteOptions

type Store interface {
	Get(key []byte, ro *ReadOptions) (value []byte, err error)
	Has(key []byte, ro *ReadOptions) (ret bool, err error)
	Put(key, value []byte, wo *WriteOptions) error
	Close() error
}

func OpenMemDb() (Store, error) {
	database, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return database, nil
}
