package logger_store_test

import (
	"encoding/json"
	"github.com/duck8823/duci/infrastructure/logger/store"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	t.Run("with correct directory", func(t *testing.T) {
		// setup
		tempTestDir := path.Join(os.TempDir(), ".duci_test")
		dir := path.Join(tempTestDir, strconv.FormatInt(time.Now().Unix(), 10))

		// and
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("error cooured: %+v", err)
		}

		// expect
		if _, err := os.Open(dir); err != nil {
			t.Errorf("must not occur error: %+v", err)
		}

		// cleanup
		os.RemoveAll(tempTestDir)
	})

	t.Run("when directory not found", func(t *testing.T) {
		// given
		dir := "/path/to/not/found"

		// expect
		if _, err := os.Open(dir); err == nil {
			t.Error("must occur error")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("when database not open", func(t *testing.T) {
		// setup
		logger_store.Close()

		// when
		job, err := logger_store.Get(uuid.New())

		// then
		if job != nil {
			t.Errorf("job must be nil, but got %+v", job)
		}

		if err == nil {
			t.Error("must occur error")
		}
	})

	t.Run("with exists data", func(t *testing.T) {
		// setup
		tempTestDir := path.Join(os.TempDir(), ".duci_test")
		dir := path.Join(tempTestDir, strconv.FormatInt(time.Now().Unix(), 10))

		// and
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("error cooured: %+v", err)
		}

		// given
		existId, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("errour occured: %+v", err)
		}

		// and
		expected := &logger_store.Job{
			Finished: true,
			Stream:   []logger_store.Message{{Level: "INFO", Time: "time", Text: "text"}},
		}

		data, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		if err := db.Put([]byte(existId.String()), data, nil); err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		db.Close()

		// and
		logger_store.Open(dir)

		t.Run("when get stored data", func(t *testing.T) {
			// when
			actual, err := logger_store.Get(existId)

			// then
			if err != nil {
				t.Errorf("error must be nil, but got %+v", err)
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("wont %+v, but got %+v", expected, actual)
			}
		})

		t.Run("when try to get not stored data", func(t *testing.T) {
			// when
			actual, err := logger_store.Get(uuid.New())

			// then
			if err.Error() != logger_store.NotFound.Error() {
				t.Errorf("error must be %+v, but got %+v", logger_store.NotFound, err)
			}

			if actual != nil {
				t.Errorf("job must be nil, buft got %+v", err)
			}
		})

		// cleanup
		logger_store.Close()

		os.RemoveAll(tempTestDir)
	})
}

func TestAppend(t *testing.T) {
	t.Run("when database not open", func(t *testing.T) {
		// setup
		logger_store.Close()

		// expected
		if err := logger_store.Append(uuid.New(), "level", "hello world"); err == nil {
			t.Error("must occur error")
		}
	})

	t.Run("with exists data", func(t *testing.T) {
		// setup
		tempTestDir := path.Join(os.TempDir(), ".duci_test")
		dir := path.Join(tempTestDir, strconv.FormatInt(time.Now().Unix(), 10))

		// and
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("error cooured: %+v", err)
		}

		// given
		existId, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("errour occured: %+v", err)
		}

		// and
		job := &logger_store.Job{
			Finished: true,
			Stream:   []logger_store.Message{{Level: "INFO", Time: "time", Text: "text"}},
		}

		data, err := json.Marshal(job)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		if err := db.Put([]byte(existId.String()), data, nil); err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		db.Close()

		// and
		logger_store.Open(dir)

		t.Run("with stored data", func(t *testing.T) {
			// when
			logger_store.Append(existId, "Error", "Append Message")

			// and
			actual, _ := logger_store.Get(existId)

			// then
			if len(actual.Stream) != 2 {
				t.Errorf("length of stream must be 2, but got %d", len(actual.Stream))
			}

			if actual.Stream[1].Level != "Error" {
				t.Errorf("appended level wont `Error`, but got `%s`", actual.Stream[1].Level)
			}

			if actual.Stream[1].Text != "Append Message" {
				t.Errorf("appended level wont `Append Message`, but got `%s`", actual.Stream[1].Text)
			}
		})

		t.Run("with not stored data", func(t *testing.T) {
			// given
			newId := uuid.New()

			// when
			logger_store.Append(newId, "Error", "Append Message")

			// and
			actual, _ := logger_store.Get(newId)

			// then
			if len(actual.Stream) != 1 {
				t.Errorf("length of stream must be 2, but got %d", len(actual.Stream))
			}

			if actual.Stream[0].Level != "Error" {
				t.Errorf("appended level wont `Error`, but got `%s`", actual.Stream[0].Level)
			}

			if actual.Stream[0].Text != "Append Message" {
				t.Errorf("appended level wont `Append Message`, but got `%s`", actual.Stream[0].Text)
			}
		})

		// cleanup
		logger_store.Close()

		os.RemoveAll(tempTestDir)
	})
}

func TestFinish(t *testing.T) {
	t.Run("when database not open", func(t *testing.T) {
		// setup
		logger_store.Close()

		// expected
		if err := logger_store.Finish(uuid.New()); err == nil {
			t.Error("must occur error")
		}
	})

	t.Run("with exists data", func(t *testing.T) {
		// setup
		tempTestDir := path.Join(os.TempDir(), ".duci_test")
		dir := path.Join(tempTestDir, strconv.FormatInt(time.Now().Unix(), 10))

		// and
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("error cooured: %+v", err)
		}

		// given
		existId, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("errour occured: %+v", err)
		}

		// and
		job := &logger_store.Job{
			Finished: true,
			Stream:   []logger_store.Message{{Level: "INFO", Time: "time", Text: "text"}},
		}

		data, err := json.Marshal(job)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		db, err := leveldb.OpenFile(dir, nil)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		if err := db.Put([]byte(existId.String()), data, nil); err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		db.Close()

		// and
		logger_store.Open(dir)

		t.Run("with stored data", func(t *testing.T) {
			// when
			err := logger_store.Finish(existId)

			// and
			actual, _ := logger_store.Get(existId)

			// then
			if err != nil {
				t.Errorf("error must not be nil, but got %+v", err)
			}

			// then
			if actual.Finished != true {
				t.Error("job must be finished")
			}
		})

		t.Run("when try to get not stored data", func(t *testing.T) {
			// when
			err := logger_store.Finish(uuid.New())

			// then
			if err == nil {
				t.Errorf("error must be %+v, but got %+v", logger_store.NotFound, err)
			}
		})

		// cleanup
		logger_store.Close()

		os.RemoveAll(tempTestDir)
	})
}

func TestClose(t *testing.T) {
	t.Run("when not opened", func(t *testing.T) {
		// setup
		logger_store.Close()

		// expect
		if err := logger_store.Close(); err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when opened", func(t *testing.T) {
		// setup
		tempTestDir := path.Join(os.TempDir(), ".duci_test")
		dir := path.Join(tempTestDir, strconv.FormatInt(time.Now().Unix(), 10))

		// and
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Fatalf("error cooured: %+v", err)
		}

		// and
		if err := logger_store.Open(dir); err != nil {
			t.Fatalf("error occurred.: %+v", err)
		}

		// expect
		if err := logger_store.Close(); err != nil {
			t.Errorf("error must not be occurred, but got %+v", err)
		}

		// cleanup
		os.RemoveAll(tempTestDir)
	})
}
