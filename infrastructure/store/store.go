package store

import (
	leveldb_errors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var (
	// NotFoundError is a leveldb/errors.ErrNotFound
	NotFoundError = leveldb_errors.ErrNotFound
)

// ReadOptions is a type alias of leveldb/opt.ReadOptions
type ReadOptions = opt.ReadOptions

// WriteOptions is a type alias of leveldb/opt.WriteOptions
type WriteOptions = opt.WriteOptions

// Store is a interface represents key-value store.
type Store interface {
	Get(key []byte, ro *ReadOptions) (value []byte, err error)
	Has(key []byte, ro *ReadOptions) (ret bool, err error)
	Put(key, value []byte, wo *WriteOptions) error
	Close() error
}
