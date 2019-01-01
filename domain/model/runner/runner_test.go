package runner_test

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/docker/mock_docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/mock_job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	"io"
	"os"
	"path"
	"testing"
	"time"
)

func TestDockerRunnerImpl_Run(t *testing.T) {
	// given
	dir := job.WorkDir(path.Join(os.TempDir(), random.String(16, random.Alphanumeric)))
	if err := os.MkdirAll(dir.String(), 0700); err != nil {
		t.Fatalf("error occur: %+v", err)
	}
	defer func() {
		_ = os.RemoveAll(dir.String())
	}()

	// and
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// and
	log := mock_job.NewMockLog(ctrl)
	log.EXPECT().
		ReadLine().
		AnyTimes().
		Return(&job.LogLine{Timestamp: time.Now(), Message: "Hello Test"}, io.EOF)

	// and
	containerID := docker.ContainerID(random.String(16, random.Alphanumeric))

	mockDocker := mock_docker.NewMockDocker(ctrl)
	mockDocker.EXPECT().
		Build(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(log, nil)
	mockDocker.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(containerID, log, nil)
	mockDocker.EXPECT().
		ExitCode(gomock.Any(), gomock.Eq(containerID)).
		Times(1).
		Return(docker.ExitCode(0), nil)
	mockDocker.EXPECT().
		RemoveContainer(gomock.Any(), gomock.Eq(containerID)).
		Times(1).
		Return(nil)

	// and
	sut := runner.DockerRunnerImpl{}
	defer sut.SetDocker(mockDocker)()
	defer sut.SetLogFunc(runner.NothingToDo)()

	// when
	err := sut.Run(context.Background(), dir, "", []string{""})

	// then
	if err != nil {
		t.Errorf("error must be nil, but got %+v", err)
	}
}
