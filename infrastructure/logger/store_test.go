package logger_test

import (
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestOpenMemDb(t *testing.T) {
	// when
	actual, err := logger.OpenMemDb()

	// then
	if _, ok := actual.(*leveldb.DB); !ok {
		t.Error("must be a *leveldb.DB, but not.")
	}

	if err != nil {
		t.Errorf("error must not occur, but got %+v", err)
	}
}
