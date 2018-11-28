package runner_test

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/log/mock_log"
	"github.com/duck8823/duci/domain/service/docker/mock_docker"
	"github.com/duck8823/duci/domain/service/runner"
	. "github.com/golang/mock/gomock"
	"github.com/labstack/gommon/random"
	"os"
	"path"
	"testing"
)

func TestDockerTaskRunner_Run(t *testing.T) {
	// given
	dir, rmDir := tmpDir(t)
	defer rmDir()

	// and
	ctrl := NewController(t)
	defer ctrl.Finish()

	// and
	mockDocker := mock_docker.NewMockDocker(ctrl)
	mockDocker.EXPECT().
		Build(Any(), Any(), Any(), Any()).
		Return(mock_log.NewMockLog(ctrl), nil).
		Times(1)
	mockDocker.EXPECT().
		Run(Any(), Any(), Any(), Any()).
		Return(docker.ContainerID(""), mock_log.NewMockLog(ctrl), nil).
		Times(1)
	mockDocker.EXPECT().
		ExitCode(Any(), Any()).
		Return(docker.ExitCode(0), nil).
		Times(1)
	mockDocker.EXPECT().
		Rm(Any(), Any()).
		Return(nil).
		Times(1)

	// and
	sut := &runner.DockerRunnerImpl{
		Docker: mockDocker,
	}

	// when
	err := sut.Run(context.Background(), dir, docker.Tag("duck8823/duci:test"), []string{})

	// then
	if err != nil {
		t.Errorf("must not error. but: %+v", err)
	}
}

func tmpDir(t *testing.T) (string, func()) {
	tmpDir := path.Join(os.TempDir(), random.String(16, random.Alphanumeric))
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		t.Fatalf("error occured: %+v", err)
	}
	return tmpDir, func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Fatalf("error occured: %+v", err)
		}
	}
}
