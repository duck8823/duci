package executor_test

import (
	"context"
	"github.com/duck8823/duci/application"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner/mock_runner"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"os"
	"path"
	"testing"
	"time"
)

func TestJobExecutor_Execute(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		ctx := context.Background()
		target := &executor.StubTarget{
			Dir:     job.WorkDir(path.Join(os.TempDir(), random.String(16))),
			Cleanup: func() error { return nil },
			Err:     nil,
		}

		// and
		var calledStartFunc, calledEndFunc bool

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		runner := mock_runner.NewMockDockerRunner(ctrl)
		runner.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil)

		// and
		sut := &executor.JobExecutor{}
		defer sut.SetDockerRunner(runner)()
		defer sut.SetStartFunc(func(context.Context) {
			calledStartFunc = true
		})()
		defer sut.SetEndFunc(func(context.Context, error) {
			calledEndFunc = true
		})()

		// when
		err := sut.Execute(ctx, target)

		// then
		if err != nil {
			t.Errorf("must be nil, but got %+v", err)
		}

		// and
		if !calledStartFunc {
			t.Errorf("must be called startFunc")
		}

		// and
		if !calledEndFunc {
			t.Errorf("must be called endFunc")
		}
	})

	t.Run("with error", func(t *testing.T) {
		// given
		ctx := context.Background()
		target := &executor.StubTarget{
			Dir:     job.WorkDir(path.Join(os.TempDir(), random.String(16))),
			Cleanup: func() error { return nil },
			Err:     nil,
		}

		// and
		wantErr := errors.New("test error")

		// and
		var calledStartFunc, calledEndFunc bool

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		runner := mock_runner.NewMockDockerRunner(ctrl)
		runner.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(wantErr)

		// and
		sut := &executor.JobExecutor{}
		defer sut.SetDockerRunner(runner)()
		defer sut.SetStartFunc(func(context.Context) {
			calledStartFunc = true
		})()
		defer sut.SetEndFunc(func(context.Context, error) {
			calledEndFunc = true
		})()

		// when
		err := sut.Execute(ctx, target)

		// then
		if err != wantErr {
			t.Errorf("must be equal. want %+v, but got %+v", wantErr, err)
		}

		// and
		if !calledStartFunc {
			t.Errorf("must be called startFunc")
		}

		// and
		if !calledEndFunc {
			t.Errorf("must be called endFunc")
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		// given
		timeout := application.Config.Timeout()
		application.Config.Job.Timeout = 1
		defer func() {
			application.Config.Job.Timeout = timeout.Nanoseconds() * 1000 * 1000
		}()

		// and
		ctx := context.Background()
		target := &executor.StubTarget{
			Dir:     job.WorkDir(path.Join(os.TempDir(), random.String(16))),
			Cleanup: func() error { return nil },
			Err:     nil,
		}

		// and
		var calledStartFunc, calledEndFunc bool

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		runner := mock_runner.NewMockDockerRunner(ctrl)
		runner.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Do(func(_, _, _, _ interface{}) {
				time.Sleep(5 * time.Second)
			}).
			Return(nil)

		// and
		sut := &executor.JobExecutor{}
		defer sut.SetDockerRunner(runner)()
		defer sut.SetStartFunc(func(context.Context) {
			calledStartFunc = true
		})()
		defer sut.SetEndFunc(func(context.Context, error) {
			calledEndFunc = true
		})()

		// when
		err := sut.Execute(ctx, target)

		// then
		if err != context.DeadlineExceeded {
			t.Errorf("must be equal. want %+v, but got %+v", context.DeadlineExceeded, err)
		}

		// and
		if !calledStartFunc {
			t.Errorf("must be called startFunc")
		}

		// and
		if !calledEndFunc {
			t.Errorf("must be called endFunc")
		}
	})

	t.Run("when prepare returns error", func(t *testing.T) {
		// given
		ctx := context.Background()
		target := &executor.StubTarget{
			Dir:     job.WorkDir(path.Join(os.TempDir(), random.String(16))),
			Cleanup: func() error { return nil },
			Err:     errors.New("test error"),
		}

		// and
		var calledStartFunc, calledEndFunc bool

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		runner := mock_runner.NewMockDockerRunner(ctrl)
		runner.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := &executor.JobExecutor{}
		defer sut.SetDockerRunner(runner)()
		defer sut.SetStartFunc(func(context.Context) {
			calledStartFunc = true
		})()
		defer sut.SetEndFunc(func(context.Context, error) {
			calledEndFunc = true
		})()

		// when
		err := sut.Execute(ctx, target)

		// then
		if err == nil {
			t.Error("must not be nil")
		}

		// and
		if !calledStartFunc {
			t.Errorf("must be called startFunc")
		}

		// and
		if calledEndFunc {
			t.Errorf("must not be called endFunc")
		}
	})
}
