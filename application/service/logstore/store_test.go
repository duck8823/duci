package logstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/duck8823/duci/data/model"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/duck8823/duci/infrastructure/logger/mock_logger"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestNewStoreService(t *testing.T) {
	// when
	actual, err := New()

	// then
	if _, ok := actual.(*storeServiceImpl); !ok {
		t.Error("must be a Service, but not.")
	}

	if err != nil {
		t.Errorf("error must not occur, but got %+v", err)
	}
}

func TestStoreServiceImpl_Append(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("when store returns correct data", func(t *testing.T) {
		// given
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		date1 := time.Date(2020, time.April, 1, 12, 3, 00, 00, jst)
		date2 := time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
		job := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: date1, Text: "Hello World."}},
		}
		storedData, err := json.Marshal(job)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		clock.Now = func() time.Time {
			return date2
		}

		// and
		expected := &model.Job{
			Finished: false,
			Stream: []model.Message{
				{Time: date1, Text: "Hello World."},
				{Time: date2, Text: "Hello Testing."},
			},
		}
		expectedData, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Eq(expectedData), gomock.Nil()).
			Times(1).
			Return(nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Not(expectedData), gomock.Nil()).
			Do(func(_ interface{}, data []byte, _ interface{}) {
				t.Logf("wont: %s", string(expectedData))
				t.Logf("got:  %s", string(data))
			}).
			Return(errors.New("must not call this"))

		// expect
		if err := service.Append(id, model.Message{Time: date2, Text: "Hello Testing."}); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// cleanup
		clock.Adjust()
	})

	t.Run("when value not found", func(t *testing.T) {
		// given
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		clock.Now = func() time.Time {
			return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
		}

		// and
		clock.Now = func() time.Time {
			return time.Time{}
		}

		expected := &model.Job{
			Finished: false,
			Stream: []model.Message{
				{Time: clock.Now(), Text: "Hello Testing."},
			},
		}
		expectedData, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("error occurred: %+v", err)
		}

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(nil, logger.NotFoundError)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Eq(expectedData), gomock.Nil()).
			Times(1).
			Return(nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Not(expectedData), gomock.Nil()).
			Do(func(_ interface{}, data []byte, _ interface{}) {
				t.Logf("wont: %s", string(expectedData))
				t.Logf("got:  %s", string(data))
			}).
			Return(errors.New("must not call this"))

		// expect
		if err := service.Append(id, model.Message{Time: time.Time{}, Text: "Hello Testing."}); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// cleanup
		clock.Adjust()
	})

	t.Run("when store.Get returns error", func(t *testing.T) {
		// given
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(nil, errors.New("hello testing"))

		// expect
		if err := service.Append(id, model.Message{Text: "Hello Testing."}); err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("when store.Get returns invalid data", func(t *testing.T) {
		// given
		storedData := []byte("invalid data")

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Any(), gomock.Nil()).
			Times(1).
			Do(func(_, _, _ interface{}) {
				t.Fatalf("must not call this.")
			})

		// expect
		if err := service.Append(id, model.Message{Text: "Hello Testing."}); err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("when store.Put returns invalid data", func(t *testing.T) {
		// given
		job := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: clock.Now(), Text: "Hello World."}},
		}
		storedData, err := json.Marshal(job)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		clock.Now = func() time.Time {
			return time.Date(1987, time.March, 27, 19, 19, 00, 00, jst)
		}

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Any(), gomock.Nil()).
			Times(1).
			Return(errors.New("hello error"))

		// expect
		if err := service.Append(id, model.Message{Text: "Hello Testing."}); err == nil {
			t.Error("error must occur, but got nil")
		}

		// cleanup
		clock.Adjust()
	})
}

func TestStoreServiceImpl_Get(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("with error", func(t *testing.T) {
		// setup
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// given
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(nil, errors.New("hello testing"))

		// when
		actual, err := service.Get(id)

		// then
		if actual != nil {
			t.Errorf("job must be nil, but got %+v", actual)
		}

		if err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("with invalid data", func(t *testing.T) {
		// given
		storedData := []byte("invalid data")

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)

		// when
		actual, err := service.Get(id)

		// then
		if err == nil {
			t.Error("error must occur, but got nil")
		}

		if actual != nil {
			t.Errorf("job must be nil, but got %+v", err)
		}
	})

	t.Run("with stored data", func(t *testing.T) {
		// given
		clock.Now = func() time.Time {
			return time.Unix(0, 0)
		}

		expected := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: clock.Now(), Text: "Hello World."}},
		}
		storedData, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)

		// when
		actual, err := service.Get(id)

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		if !cmp.Equal(actual.Stream[0].Time, expected.Stream[0].Time) {
			t.Errorf("wont %+v, but got %+v", expected, actual)
		}

		// cleanup
		clock.Adjust()
	})
}

func TestStoreServiceImpl_Start(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("when put success", func(t *testing.T) {
		// given
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		expected, err := json.Marshal(&model.Job{Finished: false})
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Eq(expected), gomock.Nil()).
			Times(1).
			Return(nil)

		// when
		err = service.Start(id)

		// then
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("when put fail", func(t *testing.T) {
		// given
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		mockStore.EXPECT().
			Put(gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(errors.New("test error"))

		// when
		err = service.Start(id)

		// then
		if err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestStoreServiceImpl_Finish(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("with error", func(t *testing.T) {
		// setup
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// given
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(nil, errors.New("hello testing"))

		// expect
		if err := service.Finish(id); err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("with invalid data", func(t *testing.T) {
		// given
		storedData := []byte("invalid data")

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)

		// expect
		if err := service.Finish(id); err == nil {
			t.Error("error must occur, but got nil")
		}
	})

	t.Run("with stored data", func(t *testing.T) {
		// given
		given := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: time.Unix(0, 0), Text: "Hello World."}},
		}
		storedData, err := json.Marshal(given)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		expected := &model.Job{
			Finished: true,
			Stream:   []model.Message{{Time: time.Unix(0, 0), Text: "Hello World."}},
		}
		expectedData, err := json.Marshal(expected)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Eq(expectedData), gomock.Nil()).
			Times(1)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Not(expectedData), gomock.Nil()).
			Do(func(_, arg, _ interface{}) {
				actual := &model.Job{}
				json.NewDecoder(bytes.NewReader(arg.([]byte))).Decode(actual)
				t.Fatalf("invalid argument: wont %+v, but got %+v", expected, actual)
			})

		// when
		err = service.Finish(id)

		// and
		if err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("when failed put", func(t *testing.T) {
		// given
		given := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: clock.Now(), Text: "Hello World."}},
		}
		storedData, err := json.Marshal(given)
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}

		// and
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatalf("error occured: %+v", err)
		}
		storedId := []byte(id.String())

		// and
		mockStore.EXPECT().
			Get(gomock.Eq(storedId), gomock.Nil()).
			Times(1).
			Return(storedData, nil)
		mockStore.EXPECT().
			Put(gomock.Eq(storedId), gomock.Any(), gomock.Nil()).
			Times(1).
			Return(errors.New("hello testing"))

		// expect
		if err := service.Finish(id); err == nil {
			t.Error("error must occur, but got nil")
		}
	})
}

func TestStoreServiceImpl_Close(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("with error", func(t *testing.T) {
		// given
		mockStore.EXPECT().
			Close().
			Return(errors.New("hello testing"))

		// expect
		if err := service.Close(); err == nil {
			t.Errorf("error must not occur, but got %+v", err)
		}
	})

	t.Run("without error", func(t *testing.T) {
		// given
		mockStore.EXPECT().
			Close().
			Return(nil)

		// expect
		if err := service.Close(); err != nil {
			t.Error("error must occur, but got nil")
		}
	})
}
