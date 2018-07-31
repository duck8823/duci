package log

import (
	"encoding/json"
	"errors"
	"github.com/duck8823/duci/domain/model"
	"github.com/duck8823/duci/infrastructure/clock"
	"github.com/duck8823/duci/infrastructure/logger/mock_logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestStoreServiceImpl_Append(t *testing.T) {
	// setup
	ctrl := gomock.NewController(t)
	mockStore := mock_logger.NewMockStore(ctrl)

	service := &storeServiceImpl{mockStore}
	t.Run("when store returns correct data", func(t *testing.T) {
		// given
		job := &model.Job{
			Finished: false,
			Stream:   []model.Message{{Time: "stored time", Level: "INFO", Text: "Hello World."}},
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
		expected := &model.Job{
			Finished: false,
			Stream: []model.Message{
				{Time: "stored time", Level: "INFO", Text: "Hello World."},
				{Time: clock.Now().String(), Level: "INFO", Text: "Hello Testing."},
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
		if err := service.Append(id, "INFO", "Hello Testing."); err != nil {
			t.Errorf("error must not occur, but got %+v", err)
		}

		// cleanup
		clock.Adjust()
	})
}
