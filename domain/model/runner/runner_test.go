package runner_test

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/docker/mock_docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/mock_job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDockerRunnerImpl_Run(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, log, nil)
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(docker.ExitCode(0), nil)
		mockDocker.EXPECT().
			RemoveContainer(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(nil)
		mockDocker.EXPECT().
			RemoveImage(gomock.Any(), gomock.Eq(tag)).
			Times(1).
			Return(nil)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err != nil {
			t.Errorf("error must be nil, but got %+v", err)
		}
	})

	t.Run("when failure create tarball", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		if err := os.MkdirAll(filepath.Join(dir.String(), "duci.tar"), 0700); err != nil {
			t.Fatalf("error occur: %+v", err)
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Errorf("error must not be nil")
		}
	})

	t.Run("when failure load runtime options", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		if err := os.MkdirAll(filepath.Join(dir.String(), ".duci", "config.yml"), 0700); err != nil {
			t.Fatalf("error occur: %+v", err)
		}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// expect
		if err := sut.Run(context.Background(), dir, tag, cmd); err == nil {
			t.Errorf("error must not be nil")
		}
	})

	t.Run("when failure docker build", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, errors.New("error test"))
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(0)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Errorf("error must not be nil")
		}
	})

	t.Run("when failure docker run", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, nil, errors.New("test error"))
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(0)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when failure to get exit code", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, log, nil)
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(docker.ExitCode(0), errors.New("test error"))
		mockDocker.EXPECT().
			RemoveContainer(gomock.Any(), gomock.Eq(conID)).
			Times(0).
			Return(nil)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when exit code is not zero", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, log, nil)
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(docker.ExitCode(-1), nil)
		mockDocker.EXPECT().
			RemoveContainer(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(nil)
		mockDocker.EXPECT().
			RemoveImage(gomock.Any(), gomock.Eq(tag)).
			Times(1).
			Return(nil)

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err != runner.ErrFailure {
			t.Errorf("error must be ErrFailure, but got %+v", err)
		}
	})

	t.Run("when failure docker remove container", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, log, nil)
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(docker.ExitCode(0), nil)
		mockDocker.EXPECT().
			RemoveContainer(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})

	t.Run("when failure docker remove image", func(t *testing.T) {
		// given
		dir, cleanup := tmpDir(t)
		defer cleanup()

		tag := docker.Tag(fmt.Sprintf("duci/test:%s", random.String(8)))
		cmd := docker.Command{"echo", "test"}

		// and
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// and
		log := stubLog(t, ctrl)
		conID := docker.ContainerID(random.String(16, random.Alphanumeric))

		mockDocker := mock_docker.NewMockDocker(ctrl)
		mockDocker.EXPECT().
			Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(log, nil)
		mockDocker.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Times(1).
			Return(conID, log, nil)
		mockDocker.EXPECT().
			ExitCode(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(docker.ExitCode(0), nil)
		mockDocker.EXPECT().
			RemoveContainer(gomock.Any(), gomock.Eq(conID)).
			Times(1).
			Return(nil)
		mockDocker.EXPECT().
			RemoveImage(gomock.Any(), gomock.Eq(tag)).
			Times(1).
			Return(errors.New("test error"))

		// and
		sut := runner.DockerRunnerImpl{}
		defer sut.SetDocker(mockDocker)()
		defer sut.SetLogFunc(runner.NothingToDo)()

		// when
		err := sut.Run(context.Background(), dir, tag, cmd)

		// then
		if err == nil {
			t.Error("error must not be nil")
		}
	})
}

func tmpDir(t *testing.T) (workDir job.WorkDir, clean func()) {
	t.Helper()

	dir := job.WorkDir(filepath.Join(os.TempDir(), random.String(16)))
	if err := os.MkdirAll(dir.String(), 0700); err != nil {
		t.Fatalf("error occur: %+v", err)
	}
	return dir, func() {
		_ = os.RemoveAll(dir.String())
	}
}

func stubLog(t *testing.T, ctrl *gomock.Controller) *mock_job.MockLog {
	t.Helper()

	log := mock_job.NewMockLog(ctrl)
	log.EXPECT().
		ReadLine().
		AnyTimes().
		Return(&job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}, io.EOF)
	return log
}
