package job_test

import (
	"encoding/json"
	"github.com/duck8823/duci/domain/model/job"
	. "github.com/duck8823/duci/infrastructure/job"
	"github.com/duck8823/duci/infrastructure/job/mock_job"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path"
	"testing"
	"time"
)

func TestNewDataSource(t *testing.T) {
	t.Run("with temporary path", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// when
		got, err := NewDataSource(tmpDir)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if got == nil {
			t.Error("must not be nil")
		}
	})

	t.Run("with wrong path", func(t *testing.T) {
		// given
		tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		// and
		db, err := leveldb.OpenFile(tmpDir, nil)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}
		defer db.Close()

		// when
		got, err := NewDataSource(tmpDir)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})
}

func TestDataSource_FindBy(t *testing.T) {
	t.Run("when returns data", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		want := &job.Job{
			ID:       id,
			Finished: false,
			Stream:   []job.LogLine{{Timestamp: time.Now(), Message: "Hello Test"}},
		}
		data, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Get(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Nil()).
			Times(1).
			Return(data, nil)

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}

		// and
		if !cmp.Equal(got, want) {
			t.Errorf("must be equal, but %+v", cmp.Diff(got, want))
		}
	})

	t.Run("when returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Get(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Nil()).
			Times(1).
			Return(nil, errors.New("test error"))

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})

	t.Run("when leveldb.ErrNotFound", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Get(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Nil()).
			Times(1).
			Return(nil, leveldb.ErrNotFound)

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})

	t.Run("when stored data is invalid format", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Get(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Nil()).
			Times(1).
			Return([]byte("invalid format"), nil)

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// when
		got, err := sut.FindBy(id)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}

		// and
		if got != nil {
			t.Errorf("must be nil, but got %+v", got)
		}
	})
}

func TestDataSource_Save(t *testing.T) {
	t.Run("when returns no error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		j := &job.Job{
			ID:       id,
			Finished: false,
			Stream:   []job.LogLine{{Timestamp: time.Now(), Message: "Hello Test"}},
		}
		data, err := json.Marshal(j)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Put(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Eq(data), gomock.Nil()).
			Times(1).
			Return(nil)

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// expect
		if err := sut.Save(*j); err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when returns error", func(t *testing.T) {
		// given
		id := job.ID(uuid.New())

		// and
		j := &job.Job{
			ID:       id,
			Finished: false,
			Stream:   []job.LogLine{{Timestamp: time.Now(), Message: "Hello Test"}},
		}
		data, err := json.Marshal(j)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db := mock_job.NewMockLevelDB(ctrl)
		db.EXPECT().
			Put(gomock.Eq([]byte(uuid.UUID(id).String())), gomock.Eq(data), gomock.Nil()).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := &DataSource{}
		defer sut.SetDB(db)()

		// expect
		if err := sut.Save(*j); err == nil {
			t.Error("error must not be nil")
		}
	})

}
