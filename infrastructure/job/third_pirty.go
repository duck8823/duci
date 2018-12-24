package job

import (
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// LevelDB is a interface represents key-value store.
type LevelDB interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Has(key []byte, ro *opt.ReadOptions) (ret bool, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
	Close() error
}
