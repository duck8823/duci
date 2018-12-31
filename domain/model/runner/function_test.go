package runner_test

import (
	"context"
	"github.com/duck8823/duci/domain/model/job/mock_job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNothingToDo(t *testing.T) {
	// given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock_job.NewMockLog(ctrl)
	log.EXPECT().
		ReadLine().
		Times(0).
		Do(func() {
			t.Error("must not call this.")
		})

	// and
	sut := runner.NothingToDo

	// expect
	sut(context.Background(), log)
}
